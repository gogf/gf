package gkvdb

import (
    "g/encoding/gcrc32"
    "os"
    "g/os/gfile"
    "strings"
    "fmt"
    "g/encoding/gbinary"
    "g/os/gfilepool"
    "errors"
    "encoding/json"
    "g/encoding/gcompress"
    "g/encoding/gjson"
)

const (
    gREAD_CACHE_TIMEOUT      = 60  // bucket的读取缓存过期时间(秒)
    gFILE_POOL_CACHE_TIMEOUT = 60  // 指针连接池缓存时间(秒)
    gBUCKET_SIZE             = 100 // 每个数据块的数据集大小（约数）
)

// KV数据库
type DB struct {
    path   string          // 数据文件存放目录路径
    prefix string          // 数据文件名前缀
    ixfp   *gfilepool.Pool // 索引0文件打开指针池
    dbfp   *gfilepool.Pool // 数据文件打开指针池
}

// KV数据记录
type DBRecord struct {
    code    uint32
    offset  int64
    dbstart int64
    dbend   int64
    dbcap   int32
    dbmap   map[string][]byte
}

func New(path, prefix string) (*DB, error) {
    path = strings.TrimRight(path, gfile.Separator)
    if prefix == "" {
        prefix = "gkvdb"
    }
    if !gfile.Exists(path) {
        if err := gfile.Mkdir(path); err != nil {
            return nil, err
        }
    }
    // 目录权限检测
    if !gfile.IsWritable(path) {
        return nil, errors.New(path + " is not writable")
    }
    // 索引/数据文件权限检测
    ixpath := path + gfile.Separator + prefix + ".ix"
    dbpath := path + gfile.Separator + prefix + ".db"
    if gfile.Exists(ixpath) && (!gfile.IsWritable(ixpath) || !gfile.IsReadable(ixpath)){
        return nil, errors.New("permission denied to index file: " + ixpath)
    }
    if gfile.Exists(dbpath) && (!gfile.IsWritable(dbpath) || !gfile.IsReadable(dbpath)){
        return nil, errors.New("permission denied to data file: " + dbpath)
    }
    // 创建文件指针池
    ixfp := gfilepool.New(ixpath, os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
    dbfp := gfilepool.New(dbpath, os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
    return &DB {
        path   : path,
        prefix : prefix,
        ixfp   : ixfp,
        dbfp   : dbfp,
    }, nil
}

// 计算关键字的校验码
func (db *DB) getHashCode(key string) uint32 {
    return gcrc32.EncodeString(key)
}

// 计算键名对应到索引文件中的偏移量
func (db *DB) getOffsetByCode(code uint32) int64 {
    offset := int64(code/gBUCKET_SIZE)
    if offset > 0 {
        offset = (offset - 1)*20
    }
    return offset
}

// 查询数据
func (db *DB) getDBRecordByKey(key string) (*DBRecord, error) {
    code    := db.getHashCode(key)
    offset  := db.getOffsetByCode(code)
    pf, err := db.ixfp.File()
    if err != nil {
        return nil, err
    }
    defer pf.Close()
    record  := DBRecord {
        code    : code,
        offset  : offset,
        dbmap   : make(map[string][]byte),
    }
    buffer := gfile.GetBinContentByTwoOffsets(pf.File(), offset, offset + 20)
    if buffer != nil && len(buffer) == 20 {
        record.dbstart,_ = gbinary.DecodeToInt64(buffer[0:8])
        record.dbend,_   = gbinary.DecodeToInt64(buffer[8:16])
        record.dbcap,_   = gbinary.DecodeToInt32(buffer[16:])
        if record.dbcap > 0 {
            if pf, err := db.dbfp.File(); err == nil {
                defer pf.Close()
                buffer := gfile.GetBinContentByTwoOffsets(pf.File(), record.dbstart, record.dbend)
                if buffer != nil {
                    json.Unmarshal(gcompress.UnZlib(buffer), &record.dbmap)
                }
            }
        }
    }
    return &record, nil
}

//// 获取数据缓存键名
//func (db *DB) getCacheKey(offset int64) string {
//    return fmt.Sprintf("gkvdb_%d", offset)
//}

// 查询KV数据
func (db *DB) Get(key string) ([]byte, error) {
    record, err := db.getDBRecordByKey(key)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("data index failed for key: %s, error:%s", key, err.Error()))
    }
    if v, ok := record.dbmap[key]; ok {
        return v, nil
    }
    return nil, nil
}

// 设置KV数据
func (db *DB) Set(key string, value []byte) error {
    record, err := db.getDBRecordByKey(key)
    if err != nil {
        return errors.New(fmt.Sprintf("data index failed for key: %s, error:%s", key, err.Error()))
    }
    // 数据整合
    record.dbmap[key] = value
    // 数据结构化处理及压缩
    data := gcompress.Zlib([]byte(gjson.Encode(record.dbmap)))
    // 创建文件操作指针
    pfix, err := db.ixfp.File()
    if err != nil {
        return err
    }
    defer pfix.Close()
    pfdb, err := db.dbfp.File()
    if err != nil {
        return err
    }
    defer pfdb.Close()
    // 写入数据文件
    var dbstart int64 = record.dbstart
    var dbend   int64 = record.dbend
    var dbcap   int32 = record.dbcap
    length := int32(len(data))
    if record.dbcap >= length {
        // 如果原本的空间大小满足本次数据，那么直接覆写
        if _, err = pfdb.File().WriteAt(data, record.dbstart); err != nil {
            return err
        }
        // 更新数据索引信息
        dbend = record.dbstart + int64(length)
    } else {
        // 否则从文件末尾重新写入数据，原本的空间依靠另外的清理线程回收处理
        pos, err := pfdb.File().Seek(0, 2)
        if err != nil {
            return err
        }
        if _, err = pfdb.File().WriteAt(data, pos); err != nil {
            return err
        }
        dbstart = pos
        dbend   = dbstart + int64(length)
        dbcap   = length
    }
    // 写入索引文件
    buffer, _ := gbinary.Encode(dbstart, dbend, dbcap)
    if _, err := pfix.File().WriteAt(buffer, record.offset); err != nil {
        return err
    }
    return nil
}






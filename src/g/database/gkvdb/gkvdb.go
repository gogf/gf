package gkvdb

import (
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
    "g/encoding/ghash"
)

const (
    gINDEX1_CACHE_TIMEOUT    = 60       // 二级索引缓存时间(秒)
    gFILE_POOL_CACHE_TIMEOUT = 60       // 文件指针池缓存时间(秒)
    gPARTITION_SIZE          = 4194304  // 哈希表分区大小(文件大小最大约为64MB)
    gINDEX1_SIZE             = 520000   // 二级索引允许的最大长度，超过此长度则重新
)

// KV数据库
type DB struct {
    path   string          // 数据文件存放目录路径
    prefix string          // 数据文件名前缀
    ix0fp  *gfilepool.Pool // 一级索引文件打开指针池(用以高并发下的IO复用)
    ix1fp  *gfilepool.Pool // 二级索引文件打开指针池
    dbfp   *gfilepool.Pool // 数据文件打开指针池
}

// KV数据记录
type DBRecord struct {
    hash    uint64 // 64位的hash code
    offset0 int64  // 一级索引地址
    offset1 int64  // 二级索引地址
    keysize uint16 // 关键字长度，用以切分数据
    dbstart int64  // 数据文件中的开始地址
    dbend   int64  // 数据文件中的结束地址
    dbcap   int32  // 数据允许存放的的最大长度（用以修改对比）
}

// 创建一个KV数据库
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
    ix0path := path + gfile.Separator + prefix + ".0.ix"
    ix1path := path + gfile.Separator + prefix + ".1.ix"
    dbpath  := path + gfile.Separator + prefix + ".db"
    if gfile.Exists(ix0path) && (!gfile.IsWritable(ix0path) || !gfile.IsReadable(ix0path)){
        return nil, errors.New("permission denied to 0 index file: " + ix0path)
    }
    if gfile.Exists(ix1path) && (!gfile.IsWritable(ix1path) || !gfile.IsReadable(ix1path)){
        return nil, errors.New("permission denied to 1 index file: " + ix1path)
    }
    if gfile.Exists(dbpath) && (!gfile.IsWritable(dbpath) || !gfile.IsReadable(dbpath)){
        return nil, errors.New("permission denied to data file: " + dbpath)
    }
    // 创建文件指针池
    ix0fp := gfilepool.New(ix0path, os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
    ix1fp := gfilepool.New(ix1path, os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
    dbfp := gfilepool.New(dbpath, os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
    return &DB {
        path   : path,
        prefix : prefix,
        ix0fp  : ix0fp,
        ix1fp  : ix1fp,
        dbfp   : dbfp,
    }, nil
}

// 计算关键字的hash code
func (db *DB) getHash(key []byte) uint64 {
    return ghash.BKDRHash64(key)
}

// 计算关键字再一级索引文件中的偏移量
func (db *DB) getOffset0ByHash(hash uint64) int64 {
    return int64(hash%gPARTITION_SIZE)
}

// 计算关键字在二级索引文件中的位置
func (db *DB) getOffset1ByHash(hash uint64) (int64, int64) {
    offset0 := db.getOffset0ByHash(hash)
    pf, err := db.ix0fp.File()
    if err != nil {
        return -1, -1
    }
    defer pf.Close()
    buffer := gfile.GetBinContentByTwoOffsets(pf.File(), offset0, offset0 + 16)
    if buffer != nil && len(buffer) == 16 {

    }
    record.dbstart,_ = gbinary.DecodeToInt64(buffer[0:8])
    record.dbend,_   = gbinary.DecodeToInt64(buffer[8:16])
    return -1, -1
}

// 查询数据信息
func (db *DB) getDBRecordByKey(key []byte) (*DBRecord, error) {
    hash    := db.getHash(key)
    offset0 := db.getOffset0ByHash(code)
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

// 查询KV数据
//func (db *DB) Get(key string) ([]byte, error) {
//    record, err := db.getDBRecordByKey(key)
//    if err != nil {
//        return nil, errors.New(fmt.Sprintf("data index failed for key: %s, error:%s", key, err.Error()))
//    }
//    if v, ok := record.dbmap[key]; ok {
//        return v, nil
//    }
//    return nil, nil
//}








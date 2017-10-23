// 基于哈希分区的KV嵌入式数据库
// 一级索引结构：二级索引文件开始位置(8) 二级索引文件结束位置(8)
// 二级索引结构：键名64位哈希值(8) 数据键名长度(2) 数据分配的存储长度(4) 数据文件索引开始位置(8) 数据文件索引结束位置(8)
// 数据索引结构：键名键值

package gkvdb

import (
    "os"
    "g/os/gfile"
    "strings"
    "g/encoding/gbinary"
    "g/os/gfilepool"
    "errors"
    "g/encoding/ghash"
)

const (
    gINDEX1_BUCKET_SIZE      = 30*10     // 二级索引分块大小
    gINDEX1_CACHE_TIMEOUT    = 60        // 二级索引缓存时间(秒)
    gFILE_POOL_CACHE_TIMEOUT = 60        // 文件指针池缓存时间(秒)
    gPARTITION_SIZE          = 100000    // 哈希表分区大小
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
    hash    uint64     // 64位的hash code
    part    int64      // 分区位置
    offset0 struct {
        start int64    // 一级索引开始位置(关键字列表起始位置)
        end   int64    // 一级索引结束位置
    }
    offset1 struct {
        start int64    // 二级索引开始位置(关键字列表中匹配关键字的准确的起始位置)
        end   int64    // 二级索引结束位置
    }
    dbinfo  struct {
        start   int64  // 数据文件中的开始地址
        end     int64  // 数据文件中的结束地址
        cap     uint32 // 数据允许存放的的最大长度（用以修改对比）
        keysize uint16 // 关键字长度，用以切分数据
    }
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
    dbfp  := gfilepool.New(dbpath,  os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
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
func (db *DB) getPartitionByHash(hash uint64) int64 {
    return int64(hash%gPARTITION_SIZE)
}

// 获得一级索引信息
func (db *DB) getOffset0ByPart(part int64) (int64, int64, error) {
    pf, err := db.ix0fp.File()
    if err != nil {
        return -1, -1, err
    }
    defer pf.Close()
    start  := part*16
    buffer := gfile.GetBinContentByTwoOffsets(pf.File(), start, start + 16)
    if buffer != nil {
        return gbinary.DecodeToInt64(buffer[0:8]), gbinary.DecodeToInt64(buffer[8:16]), nil
    }
    return -1, -1, nil
}

// 查询索引信息
func (db *DB) getRecordByKey(key []byte) (*DBRecord, error) {
    hash    := db.getHash(key)
    part    := db.getPartitionByHash(hash)
    record  := &DBRecord {
        hash    : hash,
        part    : part,
    }
    // 查询一级索引信息
    offset0start, offset0end, err := db.getOffset0ByPart(part)
    record.offset0.start = offset0start
    record.offset0.end   = offset0end
    if err != nil {
        return record, err
    }
    // 查询二级索引信息
    if offset0end > 0 {
        pf, err := db.ix1fp.File()
        if err != nil {
            return record, err
        }
        defer pf.Close()
        buffer := gfile.GetBinContentByTwoOffsets(pf.File(), offset0start, offset0end)
        if buffer != nil {
            for i := 0; i < len(buffer); i += 30 {
                hash64 := gbinary.DecodeToUint64(buffer[i: i + 8])
                if hash == hash64 {
                    if uint16(len(key)) == gbinary.DecodeToUint16(buffer[i + 8: i + 8 + 2]) {
                        record.offset1.start  = offset0start + int64(i)
                        record.offset1.end    = record.offset1.start + 30
                        record.dbinfo.cap     = gbinary.DecodeToUint32(buffer[i + 10: i + 14])
                        record.dbinfo.start   = gbinary.DecodeToInt64(buffer[i + 14: i + 22])
                        record.dbinfo.end     = gbinary.DecodeToInt64(buffer[i + 22: i + 30])
                        record.dbinfo.keysize = uint16(len(key))
                        return record, nil
                    }
                }
            }
        }

    }
    return record, nil
}

// 查询数据信息键值
func (db *DB) getValueByKey(key []byte) ([]byte, error) {
    record, err := db.getRecordByKey(key)
    if err != nil {
        return nil, err
    }
    if record.dbinfo.end > 0 {
        pf, err := db.dbfp.File()
        if err != nil {
            return nil, err
        }
        defer pf.Close()
        buffer := gfile.GetBinContentByTwoOffsets(pf.File(), record.dbinfo.start + int64(record.dbinfo.keysize), record.dbinfo.end)
        if buffer != nil {
            return buffer, nil
        }
    }
    return nil, nil
}

// 查询KV数据
func (db *DB) Get(key []byte) ([]byte, error) {
    value, err := db.getValueByKey(key)
    if err != nil {
        return nil, err
    }
    return value, nil
}

// 设置KV数据
func (db *DB) Set(key []byte, value []byte) error {
    record, err := db.getRecordByKey(key)
    if err != nil {
        return err
    }
    // 写入数据文件，并更新record信息
    if err := db.insertDataByRecord(key, value, record); err != nil {
        return err
    }
    // 根据record信息更新索引文件
    if err := db.createIndexByRecord(record); err != nil {
        return err
    }
    return nil
}

// 插入一条KV数据
func (db *DB) insertDataByRecord(key []byte, value []byte, record *DBRecord) error {
    dbpf, err := db.dbfp.File()
    if err != nil {
        return err
    }
    defer dbpf.Close()
    dbcap   := record.dbinfo.cap
    dbstart := record.dbinfo.start
    length  := uint32(len(key) + len(value))
    if record.dbinfo.end <= 0 || record.dbinfo.cap < length {
        pos, err := dbpf.File().Seek(0, 2)
        if err != nil {
            return err
        }
        dbcap   = length
        dbstart = pos
    }
    data := make([]byte, 0)
    data  = append(data, key...)
    data  = append(data, value...)
    if _, err = dbpf.File().WriteAt(data, dbstart); err != nil {
        return err
    }
    record.dbinfo.start   = dbstart
    record.dbinfo.end     = dbstart + int64(length)
    record.dbinfo.cap     = dbcap
    record.dbinfo.keysize = uint16(len(key))
    return nil
}

// 根据record重新创建索引信息
func (db *DB) createIndexByRecord(record *DBRecord) error {
    // 创建二级索引信息
    ix1pf, err := db.ix1fp.File()
    if err != nil {
        return err
    }
    defer ix1pf.Close()
    // 如果一级索引都不存在，那么需要同时更新一级和二级索引信息，这里先获取索引信息
    if record.offset0.end <= 0 {
        pos, err := ix1pf.File().Seek(0, 2)
        if err != nil {
            return err
        }
        // 每次分配必须为gINDEX1_BUCKET_SIZE
        r := pos%gINDEX1_BUCKET_SIZE
        if r != 0 {
            pos += gINDEX1_BUCKET_SIZE - r
        }
        record.offset0.start = pos
        record.offset0.end   = pos + 30
        record.offset1.start = pos
        record.offset1.end   = pos + 30
    }
    // 如果一级索引存在，那么写入到二级索引数据信息列表末尾
    if record.offset1.end <= 0 {
        record.offset1.start = record.offset0.end
    }
    data := make([]byte, 0)
    data  = append(data, gbinary.EncodeUint64(record.hash)...)
    data  = append(data, gbinary.EncodeUint16(record.dbinfo.keysize)...)
    data  = append(data, gbinary.EncodeUint32(record.dbinfo.cap)...)
    data  = append(data, gbinary.EncodeInt64(record.dbinfo.start)...)
    data  = append(data, gbinary.EncodeInt64(record.dbinfo.end)...)
    if _, err = ix1pf.File().WriteAt(data, record.offset1.start); err != nil {
        return err
    }
    // 创建一级索引信息
    ix0pf, err := db.ix0fp.File()
    if err != nil {
        return err
    }
    defer ix0pf.Close()
    data = make([]byte, 0)
    data = append(data, gbinary.EncodeInt64(record.offset0.start)...)
    data = append(data, gbinary.EncodeInt64(record.offset0.end)...)
    if _, err = ix0pf.File().WriteAt(data, record.part*16); err != nil {
        return err
    }
    return nil
}



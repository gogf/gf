// 基于哈希分区的KV嵌入式数据库
// 一级索引结构：二级索引文件偏移量(8) 索引列表分配长度(4) 索引列表真实长度(4)
// 一级索引结构：[键名32位哈希值(4) 数据分配长度(4) 数据真实长度(4) 数据文件偏移量(8)](变长，按照键名32位哈希值升序排序)
// 数据文件结构：键名长度(2) 键名键值(变长)

package gkvdb

import (
    "os"
    "g/os/gfile"
    "strings"
    "g/encoding/gbinary"
    "g/os/gfilepool"
    "errors"
    "g/encoding/ghash"
    "fmt"
)

const (
    //gPARTITION_SIZE          = 419430   // 哈希表分区大小
    gPARTITION_SIZE          = 1   // 哈希表分区大小
    gINDEX1_BUCKET_SIZE      = 1       // 二级索引索引文件列表分块大小(值越大，初始化时占用的空间越大)
    gFILE_POOL_CACHE_TIMEOUT = 60       // 文件指针池缓存时间(秒)
)

// KV数据库
type DB struct {
    path   string          // 数据文件存放目录路径
    prefix string          // 数据文件名前缀
    ix0fp  *gfilepool.Pool // 一级索引文件打开指针池(用以高并发下的IO复用)
    ix1fp  *gfilepool.Pool // 二级索引文件打开指针池
    dbfp   *gfilepool.Pool // 数据文件打开指针池
}

// KV数据检索记录
type Record struct {
    hash32    uint32 // 32位的hash code
    hash64    uint64 // 64位的hash code
    part      int64  // 分区位置
    index0 struct {
        start int64  // 一级索引开始位置
        end   int64  // 一级索引结束位置
    }
    index1 struct {
        start  int64  // 二级索引开始位置
        end    int64  // 二级索引结束位置
        cap    int32  // 二级索引分配大小(条数)
        size   int32  // 二级索引项大小(条数)
        buffer []byte // 索引列表([]byte)
        match  int32  // list中匹配到的索引位置
        near   int32  // list中未匹配到的相邻索引位置
        cmp    int8   // 当near存在时有效，判断给定的key比near大还是小
    }
    index2 struct {
        start int64   // 数据文件中的开始地址
        end   int64   // 数据文件中的结束地址
        cap   int32   // 数据允许存放的的最大长度（用以修改对比）
        klen  uint16  // 关键字长度，用以切分数据
        size  int32   // 数据总长度，用以计算结束位置
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
    ix0path := path + gfile.Separator + prefix + ".ix0"
    ix1path := path + gfile.Separator + prefix + ".ix1"
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
func (db *DB) getHash(key []byte) (uint32, uint64) {
    return ghash.BKDRHash(key), ghash.BKDRHash64(key)
}

// 计算关键字再一级索引文件中的偏移量
func (db *DB) getPartitionByHash64(hash uint64) int64 {
    return int64(hash%gPARTITION_SIZE)
}

// 获得一级索引信息
func (db *DB) getIndexInfoByRecord(record *Record) error {
    pf, err := db.ix0fp.File()
    if err != nil {
        return err
    }
    defer pf.Close()
    record.index0.start = record.part*16
    record.index0.end   = record.index0.start + 16
    if buffer := gfile.GetBinContentByTwoOffsets(pf.File(), record.index0.start, record.index0.end); buffer != nil {
        record.index1.start = gbinary.DecodeToInt64(buffer[0:8])
        record.index1.cap   = gbinary.DecodeToInt32(buffer[8:12])
        record.index1.size  = gbinary.DecodeToInt32(buffer[12:16])
        record.index1.end   = record.index1.start + int64(record.index1.size*20)
        return nil
    }
    return nil
}

// 获得二级级索引信息
func (db *DB) getDataInfoByRecord(record *Record) error {
    pf, err := db.ix1fp.File()
    if err != nil {
        return err
    }
    defer pf.Close()
    record.index1.buffer = gfile.GetBinContentByTwoOffsets(pf.File(), record.index1.start, record.index1.end)
    if record.index1.buffer != nil {
        //fmt.Println("get record", record)
        // 获取到二级索引数据后，进行二分查找
        record.index1.match = -1
        min := int32(0)
        max := record.index1.size - 1
        for {
            if record.index1.match != -1 || min > max {
                break
            }
            for {
                mid    := int32((min + max) / 2)
                hash32 := gbinary.DecodeToUint32(record.index1.buffer[mid*20 : mid*20 + 4])
                cmp    := 0
                //fmt.Println("mid:", mid, record.hash32, "VS", hash32)
                if record.hash32 < hash32 {
                    max = mid - 1
                    cmp = -1
                } else if record.hash32 > hash32 {
                    min = mid + 1
                    cmp = 1
                } else {
                    record.index1.match = mid
                    break
                }
                if min > max {
                    record.index1.near  = mid
                    record.index1.cmp   = int8(cmp)
                    break
                }
            }
        }
        if record.index1.match != -1 {
            match                := record.index1.match*20
            record.index2.cap     = gbinary.DecodeToInt32(record.index1.buffer[match +  4 : match + 8])
            record.index2.size    = gbinary.DecodeToInt32(record.index1.buffer[match +  8 : match + 12])
            record.index2.start   = gbinary.DecodeToInt64(record.index1.buffer[match + 12 : match + 20])
            record.index2.end     = record.index2.start + int64(record.index2.size)
        }
    }
    return nil
}

// 查询索引信息
func (db *DB) getRecordByKey(key []byte) (*Record, error) {
    hash32,hash64 := db.getHash(key)
    part          := db.getPartitionByHash64(hash64)
    record        := &Record {
        hash32  : hash32,
        hash64  : hash64,
        part    : part,
    }
    // 查询索引信息
    if err := db.getIndexInfoByRecord(record); err != nil {
        return record, err
    }
    // 查询数据信息
    if record.index1.end > 0 {
        record.index2.klen = uint16(len(key))
        if err := db.getDataInfoByRecord(record); err != nil {
            return record, err
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
    //fmt.Println(record)
    if record.index2.end > 0 {
        pf, err := db.dbfp.File()
        if err != nil {
            return nil, err
        }
        defer pf.Close()
        buffer := gfile.GetBinContentByTwoOffsets(pf.File(), record.index2.start + 2 + int64(record.index2.klen), record.index2.end)
        if buffer != nil {
            return buffer, nil
        }
    }
    return nil, nil
}

// 关闭数据库链接
func (db *DB) Close() {
    db.ix0fp.Close()
    db.ix1fp.Close()
    db.dbfp.Close()
}

// 删除数据库
func (db *DB) Remove(sure bool) {
    if sure {
        db.Close()
        gfile.Remove(db.path)
    }
}

// 查询KV数据
func (db *DB) Get(key []byte) []byte {
    value, _ := db.getValueByKey(key)
    return value
}


// 设置KV数据
func (db *DB) Set(key []byte, value []byte) error {
    record, err := db.getRecordByKey(key)
    if err != nil {
        return err
    }
    // fmt.Println(record)
    //return nil
    // 写入数据文件，并更新record信息
    if err := db.insertDataByRecord(key, value, record); err != nil {
        return err
    }
    oldcap := record.index1.cap
    // 根据record信息更新索引文件
    if err := db.createIndexByRecord(record); err != nil {
        return err
    }
    if record.index1.cap != oldcap {
        if record.index1.cap > gINDEX1_BUCKET_SIZE {
            fmt.Printf("new cap %d for string: %s\n", record.index1.cap, string(key))
        }
    }
    return nil
}

// 插入一条KV数据
func (db *DB) insertDataByRecord(key []byte, value []byte, record *Record) error {
    dbpf, err := db.dbfp.File()
    if err != nil {
        return err
    }
    defer dbpf.Close()
    dbcap   := record.index2.cap
    dbstart := record.index2.start
    length  := int32(len(key) + len(value)) + 2
    if record.index2.end <= 0 || record.index2.cap < length {
        pos, err := dbpf.File().Seek(0, 2)
        if err != nil {
            return err
        }
        dbcap   = length
        dbstart = pos
    }

    data := make([]byte, 0)
    data  = append(data, gbinary.EncodeUint16(uint16(len(key)))...)
    data  = append(data, key...)
    data  = append(data, value...)
    if _, err = dbpf.File().WriteAt(data, dbstart); err != nil {
        return err
    }
    record.index2.start   = dbstart
    record.index2.end     = dbstart + int64(length)
    record.index2.cap     = dbcap
    record.index2.size    = length
    if record.index2.klen <= 0 {
        record.index2.klen    = uint16(len(key))
    }
    return nil
}

// 根据record重新创建索引信息
func (db *DB) createIndexByRecord(record *Record) error {
    // 创建二级索引信息
    ix1pf, err := db.ix1fp.File()
    if err != nil {
        return err
    }
    defer ix1pf.Close()

    data := make([]byte, 0)
    data  = append(data, gbinary.EncodeUint32(record.hash32)...)
    data  = append(data, gbinary.EncodeInt32(record.index2.cap)...)
    data  = append(data, gbinary.EncodeInt32(record.index2.size)...)
    data  = append(data, gbinary.EncodeInt64(record.index2.start)...)
    fmt.Println(record)
    // 判断是否需要重新分配空间
    if record.index1.end <= 0 || (record.index1.match == -1 && (record.index1.cap < record.index1.size + 1)) {
        // 如果二级索引不存在，或者分配的空间大小不够，那么直接写入到二级索引列表末尾
        pos, err := ix1pf.File().Seek(0, 2)
        if err != nil {
            return err
        }
        // 每次分配必须为gINDEX1_BUCKET_SIZE
        t := int64(gINDEX1_BUCKET_SIZE*20)
        r := pos % t
        if r != 0 {
            pos += t - r
        }
        record.index1.start = pos
        record.index1.end   = pos + 20 + int64(record.index1.size*20)
        record.index1.cap   = (int32(record.index1.size/gINDEX1_BUCKET_SIZE) + 1)*gINDEX1_BUCKET_SIZE
    }
    // 写入数据处理
    buffer := make([]byte, 0)
    if record.index1.match != -1 {
        // 更新
        if record.index1.size > 0 {
            buffer = record.index1.buffer
            copy(buffer[record.index1.match*20 : ], data)
        } else {
            buffer = data
            record.index1.size++
        }
    } else {
        var length int32 = 0
        if record.index1.cmp > 0 {
            // 插入到near后面
            size := record.index1.near + 1
            if size > record.index1.size {
                size = record.index1.size
            }
            length = size*20
        } else {
            // 插入到near前面
            length = record.index1.near*20
        }
        buffer = record.index1.buffer[0 : length]
        buffer = append(buffer, data...)
        buffer = append(buffer, record.index1.buffer[length : ]...)
        record.index1.size++
        record.index1.end = record.index1.end + 20
        fmt.Println(length)
        fmt.Println(buffer)
    }

    //fmt.Println(record)

    if _, err = ix1pf.File().WriteAt(buffer, record.index1.start); err != nil {
        return err
    }
    // 创建一级索引信息
    ix0pf, err := db.ix0fp.File()
    if err != nil {
        return err
    }
    defer ix0pf.Close()
    buffer = make([]byte, 0)
    buffer = append(buffer, gbinary.EncodeInt64(record.index1.start)...)
    buffer = append(buffer, gbinary.EncodeInt32(record.index1.cap)...)
    buffer = append(buffer, gbinary.EncodeInt32(record.index1.size)...)
    if _, err = ix0pf.File().WriteAt(buffer, record.part*16); err != nil {
        return err
    }
    return nil
}



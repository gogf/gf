// 基于哈希分区的KV嵌入式数据库
// KV数据库其实只需要保存键值即可，但本数据库同时保存了键名，以便于后期遍历需要

// 数据库支持的范围：2546540(极端) <= n <= 229064922452(理论)

// 数据结构要点   ：数据的分配长度cap >= 数据真实长度len，且 cap - len <= bucket，
//               当数据存储内容发生改变时，依靠碎片管理器对碎片进行回收再利用，且碎片大小 >= bucket

// 索引文件结构  ：元数据文件偏移量(36bit) 元数据文件列表项大小(20bit)
// 元数据文件结构 :[键名哈希28(28bit) 键名长度(8bit) 键值长度(22bit,最大4MB) 数据文件偏移量(38bit)](变长,链表)
// 数据文件结构  ：键名(变长) 键值(变长)

// 数据项类型 :
// 0: 元数据文件中的数据项键值放在第5项中，最大长度为5byte
// 1: 元数据文件中的数据项键值放在数据文件中，第5项为数据文件中的索引位置

package gkvdb

import (
    "os"
    "g/os/gfile"
    "strings"
    "g/encoding/gbinary"
    "g/os/gfilepool"
    "errors"
    "g/encoding/ghash"
    "g/os/gfilespace"
    "sync"
)

const (
    gPARTITION_SIZE          = 1497965                  // 哈希表分区大小(大小约为10MB)
    //gPARTITION_SIZE          = 2
    gMAX_KEY_SIZE            = 0xFFFF                   // 键名最大长度(255byte)
    gMAX_VALUE_SIZE          = 0xFFFFFF >> 2            // 键值最大长度(4194303byte = 4MB)
    gINDEX_BUCKET_SIZE       = 7                        // 索引文件数据块大小
    gMETA_BUCKET_SIZE        = 12*5                     // 元数据数据分块大小(byte, 值越大，数据增长时占用的空间越大)
    gDATA_BUCKET_SIZE        = 32                       // 数据分块大小(byte, 值越大，数据增长时占用的空间越大)
    gFILE_POOL_CACHE_TIMEOUT = 60                       // 文件指针池缓存时间(秒)
)

// KV数据库
type DB struct {
    mu      sync.RWMutex
    path    string            // 数据文件存放目录路径
    name    string            // 数据文件名
    ixfp    *gfilepool.Pool   // 索引文件打开指针池(用以高并发下的IO复用)
    mtfp    *gfilepool.Pool   // 元数据文件打开指针池(元数据，包含索引信息和部分数据信息)
    dbfp    *gfilepool.Pool   // 数据文件打开指针池
    mtsp    *gfilespace.Space // 元数据文件碎片管理
    dbsp    *gfilespace.Space // 数据文件碎片管理器
}

// KV数据检索记录
type Record struct {
    hash64    uint    // 64位的hash code
    hash28    uint    // 28位的hash code
    part      int64   // 分区位置
    key       []byte  // 键名
    value     []byte  // 键值(当键值<=5时直接存放到mt文件中，检索时便能直接获取到值)
    ix struct {
        start int64   // 索引开始位置
        end   int64   // 索引结束位置
    }
    mt struct {
        start  int64  // 开始位置
        end    int64  // 结束位置
        cap    uint   // 列表分配长度(byte)
        size   uint   // 列表真实长度(byte)
        buffer []byte // 数据项列表([]byte)
        match  bool   // 是否在查找中准确匹配key
        index  int    // (匹配时有效, match=true)列表匹配的索引位置

    }
    db struct {
        start  int64  // 数据文件中的开始地址
        end    int64  // 数据文件中的结束地址
        cap    uint   // 数据允许存放的的最大长度（用以修改对比）
        size   uint   // klen + vlen
        klen   uint   // 键名大小
        vlen   uint   // 键值大小(byte)
    }
}

// 创建一个KV数据库
func New(path, name string) (*DB, error) {
    path = strings.TrimRight(path, gfile.Separator)
    if name == "" {
        name = "gkvdb"
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
    db := &DB {path : path, name : name}
    // 索引/数据文件权限检测
    ixpath := db.getIndexFilePath()
    mtpath := db.getMetaFilePath()
    dbpath := db.getDataFilePath()
    fspath := db.getSpaceFilePath()
    if gfile.Exists(ixpath) && (!gfile.IsWritable(ixpath) || !gfile.IsReadable(ixpath)){
        return nil, errors.New("permission denied to index file: " + ixpath)
    }
    if gfile.Exists(mtpath) && (!gfile.IsWritable(mtpath) || !gfile.IsReadable(mtpath)){
        return nil, errors.New("permission denied to meta file: " + mtpath)
    }
    if gfile.Exists(dbpath) && (!gfile.IsWritable(dbpath) || !gfile.IsReadable(dbpath)){
        return nil, errors.New("permission denied to data file: " + dbpath)
    }
    if gfile.Exists(fspath) && (!gfile.IsWritable(fspath) || !gfile.IsReadable(fspath)){
        return nil, errors.New("permission denied to space file: " + fspath)
    }
    // 创建文件指针池
    db.ixfp = gfilepool.New(ixpath, os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
    db.mtfp = gfilepool.New(mtpath, os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
    db.dbfp = gfilepool.New(dbpath, os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
    db.init()
    return db, nil
}

func (db *DB) getIndexFilePath() string {
    return db.path + gfile.Separator + db.name + ".ix"
}

func (db *DB) getMetaFilePath() string {
    return db.path + gfile.Separator + db.name + ".mt"
}

func (db *DB) getDataFilePath() string {
    return db.path + gfile.Separator + db.name + ".db"
}

func (db *DB) getSpaceFilePath() string {
    return db.path + gfile.Separator + db.name + ".fs"
}

// 数据库启动自检，整体快速扫描数据库，修复可能的异常
func (db *DB) init() {
    db.initFileSpace()
    db.restoreFileSpace()
}

// 根据元数据的size计算cap
func (db *DB) getMetaCapBySize(size uint) uint {
    if size > 0 && size%gMETA_BUCKET_SIZE != 0 {
        return size + gMETA_BUCKET_SIZE - size%gMETA_BUCKET_SIZE
    }
    return size
}

// 根据数据的size计算cap
func (db *DB) getDataCapBySize(size uint) uint {
    if size > 0 && size%gDATA_BUCKET_SIZE != 0 {
        return size + gDATA_BUCKET_SIZE - size%gDATA_BUCKET_SIZE
    }
    return size
}

// 初始化碎片管理器
func (db *DB) initFileSpace() {
    db.mtsp = gfilespace.New()
    db.dbsp = gfilespace.New()
}

// 计算关键字的hash code，使用64位哈希函数
func (db *DB) getHash64(key []byte) uint64 {
    return ghash.BKDRHash64(key)
}

// 计算关键字的hash code，使用32位不同于64位的哈希函数
func (db *DB) getHash32(key []byte) uint32 {
    return ghash.APHash(key)
}

// 计算关键字在索引文件中的偏移量
func (db *DB) getPartitionByHash64(hash uint64) int64 {
    return int64(hash%gPARTITION_SIZE)
}

// 获得索引信息
func (db *DB) getIndexInfoByRecord(record *Record) error {
    pf, err := db.ixfp.File()
    if err != nil {
        return err
    }
    defer pf.Close()
    record.ix.start = record.part*gINDEX_BUCKET_SIZE
    record.ix.end   = record.ix.start + gINDEX_BUCKET_SIZE
    if buffer := gfile.GetBinContentByTwoOffsets(pf.File(), record.ix.start, record.ix.end); buffer != nil {
        bits           := gbinary.DecodeBytesToBits(buffer)
        record.mt.start = int64(gbinary.DecodeBits(bits[0 : 36]))*gMETA_BUCKET_SIZE
        record.mt.size  = uint(gbinary.DecodeBits(bits[36 : ]))*12
        record.mt.cap   = db.getMetaCapBySize(record.mt.size)
        record.mt.end   = record.mt.start + int64(record.mt.size)
        return nil
    }
    return nil
}

// 获得数据检索信息
func (db *DB) getDataInfoByRecord(record *Record) error {
    pf, err := db.mtfp.File()
    if err != nil {
        return err
    }
    defer pf.Close()

    record.mt.buffer = gfile.GetBinContentByTwoOffsets(pf.File(), record.mt.start, record.mt.end)
    if record.mt.buffer != nil {
        //fmt.Println("get record", record)
        //fmt.Println("get:", record.mt.start, record.mt.size, record.mt.buffer)
        // 线性查找
        for i := 0; i < len(record.mt.buffer); i += 12 {
            buffer := record.mt.buffer[i : i + 12]
            bits   := gbinary.DecodeBytesToBits(buffer)
            hash28 := gbinary.DecodeBits(bits[0 : 28])
            if hash28 == record.hash28 {
                record.mt.index  = i
                record.mt.match  = true
                record.db.klen   = gbinary.DecodeBits(bits[28 : 36])
                record.db.vlen   = gbinary.DecodeBits(bits[36 : 58])
                record.db.size   = record.db.klen + record.db.vlen
                record.db.cap    = db.getDataCapBySize(record.db.size)
                record.db.start  = int64(gbinary.DecodeBits(bits[58 : 96]))*gDATA_BUCKET_SIZE
                record.db.end    = record.db.start + int64(record.db.size)
                break
            }
        }
    }
    return nil
}

// 查询检索信息
func (db *DB) getRecordByKey(key []byte) (*Record, error) {
    hash64 := db.getHash64(key)
    hash32 := db.getHash32(key)
    part   := db.getPartitionByHash64(hash64)
    record := &Record {
        hash64  : uint(hash64),
        hash28  : uint(hash32 & 0xFFFFFFF),
        part    : part,
        key     : key,
        value   : nil,
    }

    // 查询索引信息
    if err := db.getIndexInfoByRecord(record); err != nil {
        return record, err
    }

    // 查询数据信息
    if record.mt.end > 0 {
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

    if record == nil {
        return nil, nil
    }

    if record.value != nil {
        return record.value, nil
    }

    if record.db.end > 0 {
        pf, err := db.dbfp.File()
        if err != nil {
            return nil, err
        }
        defer pf.Close()
        buffer := gfile.GetBinContentByTwoOffsets(pf.File(), record.db.start + int64(record.db.klen), record.db.end)
        if buffer != nil {
            return buffer, nil
        }
    }
    return nil, nil
}

// 根据索引信息删除指定数据
func (db *DB) removeDataByRecord(record *Record) error {
    if err := db.removeDataFromDb(record); err != nil {
        return err
    }
    if err := db.removeDataFromMt(record); err != nil {
        return err
    }
    if err := db.removeDataFromIx(record); err != nil {
        return err
    }
    return nil
}

// 从数据文件中删除指定数据
func (db *DB) removeDataFromDb(record *Record) error {
    pf, err := db.dbfp.File()
    if err != nil {
        return err
    }
    defer pf.Close()
    // 内容空间必须执行清0
    if _, err = pf.File().WriteAt(make([]byte, record.db.vlen), record.db.start); err != nil {
        return err
    }
    // 添加碎片
    db.addDbFileSpace(int(record.db.start), record.db.cap)
    return nil
}

// 从元数据中删除指定数据
func (db *DB) removeDataFromMt(record *Record) error {
    pf, err := db.mtfp.File()
    if err != nil {
        return err
    }
    defer pf.Close()

    buffer   := make([]byte, 0)
    buffer    = append(buffer, record.mt.buffer[ : record.mt.index]...)
    endindex := record.mt.index + int(record.db.klen + 9)
    if endindex <= len(record.mt.buffer) - 1 {
        buffer  = append(buffer, record.mt.buffer[endindex : ]...)
    }
    record.mt.buffer = buffer
    record.mt.size   = uint(len(buffer))
    for i := 0; i < int(record.mt.cap - record.mt.size); i++ {
        buffer = append(buffer, byte(0))
    }
    if _, err = pf.File().WriteAt(buffer, record.mt.start); err != nil {
        return err
    }
    if record.mt.size == 0 {
        // 如果列表被清空，那么添加整块空间到碎片管理器
        db.addMtFileSpace(int(record.mt.start), record.mt.cap)
    } else {
        // 如果列表分配大小比较实际大小超过bucket，那么进行空间切分，添加多余的空间到碎片管理器
        db.checkAndResizeMtCap(record)
    }
    return nil
}

// 从索引中删除指定数据
func (db *DB) removeDataFromIx(record *Record) error {
    return db.updateIndexByRecord(record)
}

// 检查并更新元数据分配大小与实际大小，如果有多余的空间，交给碎片管理器
func (db *DB) checkAndResizeMtCap(record *Record) {
    if int(record.mt.cap - record.mt.size) >= gMETA_BUCKET_SIZE {
        realcap := db.getMetaCapBySize(record.mt.size)
        diffcap := int(record.mt.cap - realcap)
        if diffcap >= gMETA_BUCKET_SIZE {
            record.mt.cap = realcap
            db.addMtFileSpace(int(record.mt.start)+int(realcap), uint(diffcap))
        }
    }
}

// 检查并更新数据分配大小与实际大小，如果有多余的空间，交给碎片管理器
func (db *DB) checkAndResizeDbCap(record *Record) {
    if int(record.db.cap - record.db.size) >= gDATA_BUCKET_SIZE {
        realcap := db.getDataCapBySize(record.db.size)
        diffcap := int(record.db.cap - realcap)
        if diffcap >= gDATA_BUCKET_SIZE {
            record.db.cap = realcap
            db.addDbFileSpace(int(record.db.start)+int(realcap), uint(diffcap))
        }
    }
}

// 插入一条KV数据
func (db *DB) insertDataByRecord(key []byte, value []byte, record *Record) error {
    record.db.klen = uint(len(key))
    record.db.vlen = uint(len(value))
    record.db.size = record.db.klen + record.db.vlen

    // 写入数据文件
    if err := db.insertDataIntoDb(key, value, record); err != nil {
        return err
    }

    // 写入元数据
    if err := db.insertDataIntoMt(key, value, record); err != nil {
        return err
    }
    return nil
}

// 将数据写入到数据文件中，并更新信息到record
func (db *DB) insertDataIntoDb(key []byte, value []byte, record *Record) error {
    pf, err := db.dbfp.File()
    if err != nil {
        return err
    }
    defer pf.Close()
    // 判断是否额外分配键值存储空间
    if record.db.end <= 0 || record.db.cap < record.db.size {
        // 不用的空间添加到碎片管理器
        if record.db.end > 0 && record.db.cap > 0 {
            //fmt.Println("add db block", int(record.db.start), uint(record.db.cap))
            db.addDbFileSpace(int(record.db.start), record.db.cap)
        }
        // 重新计算所需空间
        if record.db.cap < record.db.size {
            for {
                record.db.cap += gDATA_BUCKET_SIZE
                if record.db.cap >= record.db.size {
                    break
                }
            }
        }
        // 首先从碎片管理器中获取，如果不够，那么再从文件末尾分配
        index, size := db.getDbFileSpace(record.db.cap)
        if index >= 0 {
            // 只能分配cap大小，多余的空间放回管理器继续分配
            extra := int(size - record.db.cap)
            if extra > 0 {
                //fmt.Println("readd db block", index + int(record.db.cap), extra)
                db.addDbFileSpace(index + int(record.db.cap), uint(extra))
            }
            record.db.start = int64(index)
            record.db.end   = int64(index) + int64(record.db.size)
        } else {
            start, err := pf.File().Seek(0, 2)
            if err != nil {
                return err
            }
            record.db.start = start
            record.db.end   = start + int64(record.db.size)
        }
    }
    // vlen不够vcap的对末尾进行补0占位
    buffer := make([]byte, 0)
    buffer  = append(buffer, key...)
    buffer  = append(buffer, value...)
    for i := 0; i < int(record.db.cap - record.db.size); i++ {
        buffer = append(buffer, byte(0))
    }
    if _, err = pf.File().WriteAt(buffer, record.db.start); err != nil {
        return err
    }
    db.checkAndResizeDbCap(record)
    return nil
}

// 将数据写入到元数据文件中，并更新信息到record
func (db *DB) insertDataIntoMt(key []byte, value []byte, record *Record) error {
    pf, err := db.mtfp.File()
    if err != nil {
        return err
    }
    defer pf.Close()

    bits   := make([]gbinary.Bit, 0)
    buffer := make([]byte, 0)
    // 二进制打包
    bits = gbinary.EncodeBits(bits, record.hash28,   28)
    bits = gbinary.EncodeBits(bits, record.db.klen,  8)
    bits = gbinary.EncodeBits(bits, record.db.vlen,  22)
    bits = gbinary.EncodeBits(bits, uint(record.db.start/gDATA_BUCKET_SIZE), 38)
    // 数据列表打包(判断位置进行覆盖或者插入)
    buffer = append(buffer, gbinary.EncodeBitsToBytes(bits)...)
    if len(record.mt.buffer) > 0 {
        if record.mt.match {
            buffer = append(buffer, record.mt.buffer[0 : record.mt.index]...)
            buffer = append(buffer, record.mt.buffer[record.mt.index + 12 :]...)
        } else {
            // @todo 二分排序处理
            buffer = append(buffer, record.mt.buffer...)
        }
    }
    // 判断数据列表空间是否足够
    record.mt.size = uint(len(buffer))
    if record.mt.end <= 0 || record.mt.cap < record.mt.size {
        // 不用的空间添加到碎片管理器
        if record.mt.end > 0 && record.mt.cap > 0 {
            //fmt.Println("add mt block", int(record.mt.start), uint(record.mt.cap))
            db.addMtFileSpace(int(record.mt.start), uint(record.mt.cap))
        }
        // 重新计算所需空间
        if record.mt.cap < record.mt.size {
            for {
                record.mt.cap += gMETA_BUCKET_SIZE
                if record.mt.cap >= record.mt.size {
                    break
                }
            }
        }
        // 首先从碎片管理器中获取，如果不够，那么再从文件末尾分配
        index, size := db.getMtFileSpace(record.mt.cap)
        if index >= 0 {
            //fmt.Println("get mt block:", index, size, record.mt.cap)
            // 只能分配cap大小，多余的空间放回管理器继续分配
            extra := int(size - record.mt.cap)
            if extra > 0 {
                //fmt.Println("readd mt block", index + int(record.mt.cap), extra)
                db.addMtFileSpace(index + int(record.mt.cap), uint(extra))
            }
            //fmt.Println(db.mtsp.GetAllBlocksByIndex())
            record.mt.start = int64(index)
            record.mt.end   = int64(index) + int64(record.mt.size)
        } else {
            //if db.mtsp.GetMaxSize() >= record.mt.cap {
            //    //fmt.Printf("get mt block failed, request: %d, max: %d\n", record.mt.cap, db.mtsp.GetMaxSize())
            //    os.Exit(1)
            //}
            start, err := pf.File().Seek(0, 2)
            if err != nil {
                return err
            }
            //fmt.Println("new mt offset:", start)
            record.mt.start = start
            record.mt.end   = start + int64(record.mt.cap)
        }
    }
    // size不够cap的对末尾进行补0占位
    for i := 0; i < int(record.mt.cap - record.mt.size); i++ {
        buffer = append(buffer, byte(0))
    }
    //fmt.Println("set:",record.mt.start, record.mt.size, buffer)
    if _, err = pf.File().WriteAt(buffer, record.mt.start); err != nil {
        return err
    }
    db.checkAndResizeMtCap(record)
    return nil
}

// 根据record更新索引信息
func (db *DB) updateIndexByRecord(record *Record) error {
    ixpf, err := db.ixfp.File()
    if err != nil {
        return err
    }
    defer ixpf.Close()

    bits := make([]gbinary.Bit, 0)
    if record.mt.size > 0 {
        // 添加/修改/部分删除
        bits = gbinary.EncodeBits(bits, uint(record.mt.start/gMETA_BUCKET_SIZE),   36)
        bits = gbinary.EncodeBits(bits, record.mt.size/12,                         20)
    } else {
        // 数据全部删除完
        bits = make([]gbinary.Bit, gINDEX_BUCKET_SIZE)
    }

    if _, err = ixpf.File().WriteAt(gbinary.EncodeBitsToBytes(bits), record.ix.start); err != nil {
        return err
    }
    return nil
}



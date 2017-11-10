// 基于哈希分区的KV嵌入式数据库
// KV数据库其实只需要保存键值即可，但本数据库同时保存了键名，以便于后期遍历需要

// 数据库支持的范围：2546540(极端) <= n <= 229064922452(理论)

// 数据结构要点   ：数据的分配长度cap >= 数据真实长度len，且 cap - len <= bucket，
//               当数据存储内容发生改变时，依靠碎片管理器对碎片进行回收再利用，且碎片大小 >= bucket

// 索引文件结构  ：元数据文件偏移量(36bit,64GB) 元数据文件列表项大小(20bit,1048575)
// 元数据文件结构 :[键名哈希28(28bit) 键名长度(8bit) 键值长度(22bit,4MB) 数据文件偏移量(38bit,256GB)](变长,链表)
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
    "sync/atomic"
    "strconv"
    "g/os/gcache"
)

const (
    gPARTITION_SIZE          = 1497965                  // 哈希表分区大小(大小约为10MB)
    //gPARTITION_SIZE          = 1
    gMAX_KEY_SIZE            = 0xFFFF                   // 键名最大长度(255byte)
    gMAX_VALUE_SIZE          = 0xFFFFFF >> 2            // 键值最大长度(4194303byte = 4MB)
    gINDEX_BUCKET_SIZE       = 7                        // 索引文件数据块大小
    gMETA_BUCKET_SIZE        = 12*5                     // 元数据数据分块大小(byte, 值越大，数据增长时占用的空间越大)
    gDATA_BUCKET_SIZE        = 32                       // 数据分块大小(byte, 值越大，数据增长时占用的空间越大)
    gFILE_POOL_CACHE_TIMEOUT = 60                       // 文件指针池缓存时间(秒)
    gCACHE_DEFAULT_TIMEOUT   = 60000                    // gcache默认缓存时间(毫秒)
    gAUTO_SAVING_TIMEOUT     = 1000                     // 自动同步到磁盘的时间(毫秒)
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
    memt    *MemTable         // MemTable
    cache   int32             // 是否开启缓存功能
}

// KV数据检索记录
type Record struct {
    hash64    uint    // 64位的hash code
    hash28    uint    // 28位的hash code
    part      int64   // 分区位置
    key       []byte  // 键名
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
        match  int    // 是否在查找中匹配结果(-2, -1, 0, 1)
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
    db := &DB {
        path  : path,
        name  : name,
        cache : 1,
    }
    db.memt = newMemTable(db)

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

    // 初始化相关服务及数据
    db.initFileSpace()
    db.restoreFileSpace()
    db.startAutoSavingLoop()
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

func (db *DB) getCache() bool {
    return atomic.LoadInt32(&db.cache) > 0
}

func (db *DB) setCache(v int32) {
    atomic.StoreInt32(&db.cache, v)
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
    var buffer []byte
    key            := "index_cache_" + strconv.Itoa(int(record.part))
    record.ix.start = record.part*gINDEX_BUCKET_SIZE
    record.ix.end   = record.ix.start + gINDEX_BUCKET_SIZE
    if v := gcache.Get(key); v != nil {
        buffer = v.([]byte)
    } else {
        pf, err := db.ixfp.File()
        if err != nil {
            return err
        }
        defer pf.Close()
        buffer = gfile.GetBinContentByTwoOffsets(pf.File(), record.ix.start, record.ix.end)
        gcache.Set(key, buffer, gCACHE_DEFAULT_TIMEOUT)
    }
    if buffer != nil {
        bits           := gbinary.DecodeBytesToBits(buffer)
        record.mt.start = int64(gbinary.DecodeBits(bits[0 : 36]))*gMETA_BUCKET_SIZE
        record.mt.size  = uint(gbinary.DecodeBits(bits[36 : ]))*12
        record.mt.cap   = db.getMetaCapBySize(record.mt.size)
        record.mt.end   = record.mt.start + int64(record.mt.size)
        return nil
    }
    return nil
}

// 获得元数据信息
func (db *DB) getMetaInfoByRecord(record *Record) error {
    pf, err := db.mtfp.File()
    if err != nil {
        return err
    }
    defer pf.Close()

    record.mt.buffer = gfile.GetBinContentByTwoOffsets(pf.File(), record.mt.start, record.mt.end)
    if record.mt.buffer != nil {
        // 二分查找
        min := 0
        max := len(record.mt.buffer)/12 - 1
        mid := 0
        cmp := -2
        for {
            if cmp == 0 || min > max {
                break
            }
            for {
                mid     = int((min + max) / 2)
                buffer := record.mt.buffer[mid*12 : mid*12 + 12]
                bits   := gbinary.DecodeBytesToBits(buffer)
                hash28 := gbinary.DecodeBits(bits[0 : 28])
                if record.hash28 < hash28 {
                    max = mid - 1
                    cmp = -1
                } else if record.hash28 > hash28 {
                    min = mid + 1
                    cmp = 1
                } else {
                    cmp = 0
                    record.db.klen   = gbinary.DecodeBits(bits[28 : 36])
                    record.db.vlen   = gbinary.DecodeBits(bits[36 : 58])
                    record.db.size   = record.db.klen + record.db.vlen
                    record.db.cap    = db.getDataCapBySize(record.db.size)
                    record.db.start  = int64(gbinary.DecodeBits(bits[58 : 96]))*gDATA_BUCKET_SIZE
                    record.db.end    = record.db.start + int64(record.db.size)
                    break
                }
                if cmp == 0 || min > max {
                    break
                }
            }
        }
        record.mt.index = mid*12
        record.mt.match = cmp
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
    }
    record.mt.match = -2

    // 查询索引信息
    if err := db.getIndexInfoByRecord(record); err != nil {
        return record, err
    }

    // 查询数据信息
    if record.mt.end > 0 {
        if err := db.getMetaInfoByRecord(record); err != nil {
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
    oldr := *record
    if err := db.removeDataFromDb(record); err != nil {
        return err
    }
    if err := db.removeDataFromMt(record); err != nil {
        return err
    }
    if oldr.mt.start != record.mt.start || oldr.mt.size != record.mt.size {
        if err := db.removeDataFromIx(record); err != nil {
            return err
        }
    }
    return nil
}

// 从数据文件中删除指定数据
func (db *DB) removeDataFromDb(record *Record) error {
    // 添加碎片
    db.addDbFileSpace(int(record.db.start), record.db.cap)
    return nil
}

// 从元数据中删除指定数据
func (db *DB) removeDataFromMt(record *Record) error {
    // 如果没有匹配到数据，那么也没必要执行删除了
    if record.mt.match != 0 {
        return nil
    }
    pf, err := db.mtfp.File()
    if err != nil {
        return err
    }
    defer pf.Close()

    record.mt.buffer = db.removeMeta(record.mt.buffer, record.mt.index)
    record.mt.size   = uint(len(record.mt.buffer))
    if record.mt.size == 0 {
        // 如果列表被清空，那么添加整块空间到碎片管理器
        db.addMtFileSpace(int(record.mt.start), record.mt.cap)
    } else {
        if _, err = pf.File().WriteAt(record.mt.buffer, record.mt.start); err != nil {
            return err
        }
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
    // vlen不够vcap的对末尾进行补0占位(便于文件末尾分配空间)
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

// 添加一项, cmp < 0往前插入，cmp >= 0往后插入
func (db *DB) saveMeta(slice []byte, buffer []byte, index int, cmp int) []byte {
    if cmp == 0 {
        copy(slice[index:], buffer)
        return slice
    }
    pos := index
    if cmp == -1 {
        // 添加到前面
    } else {
        // 添加到后面
        pos = index + 12
        if pos >= len(slice) {
            pos = len(slice)
        }
    }
    rear  := append([]byte{}, slice[pos : ]...)
    slice  = append(slice[0 : pos], buffer...)
    slice  = append(slice, rear...)
    return slice
}


// 删除一项
func (db *DB) removeMeta(slice []byte, index int) []byte {
    return append(slice[ : index], slice[index + 12 : ]...)
}


// 将数据写入到元数据文件中，并更新信息到record
func (db *DB) insertDataIntoMt(key []byte, value []byte, record *Record) error {
    pf, err := db.mtfp.File()
    if err != nil {
        return err
    }
    defer pf.Close()

    // 二进制打包
    bits := make([]gbinary.Bit, 0)
    bits  = gbinary.EncodeBits(bits, record.hash28,   28)
    bits  = gbinary.EncodeBits(bits, record.db.klen,  8)
    bits  = gbinary.EncodeBits(bits, record.db.vlen,  22)
    bits  = gbinary.EncodeBits(bits, uint(record.db.start/gDATA_BUCKET_SIZE), 38)
    // 数据列表打包(判断位置进行覆盖或者插入)
    record.mt.buffer = db.saveMeta(record.mt.buffer, gbinary.EncodeBitsToBytes(bits), record.mt.index, record.mt.match)
    record.mt.size   = uint(len(record.mt.buffer))
    if record.mt.end <= 0 || record.mt.cap < record.mt.size {
        // 不用的空间添加到碎片管理器
        if record.mt.end > 0 && record.mt.cap > 0 {
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
            // 只能分配cap大小，多余的空间放回管理器继续分配
            extra := int(size - record.mt.cap)
            if extra > 0 {
                db.addMtFileSpace(index + int(record.mt.cap), uint(extra))
            }
            record.mt.start = int64(index)
            record.mt.end   = int64(index) + int64(record.mt.size)
        } else {
            start, err := pf.File().Seek(0, 2)
            if err != nil {
                return err
            }
            record.mt.start = start
            record.mt.end   = start + int64(record.mt.cap)
        }
    }
    // size不够cap的对末尾进行补0占位(便于文件末尾分配空间)
    for i := 0; i < int(record.mt.cap - record.mt.size); i++ {
        record.mt.buffer = append(record.mt.buffer, byte(0))
    }
    if _, err = pf.File().WriteAt(record.mt.buffer, record.mt.start); err != nil {
        return err
    }
    db.checkAndResizeMtCap(record)
    return nil
}

// 根据record更新索引信息
func (db *DB) updateIndexByRecord(record *Record) error {
    defer gcache.Remove("index_cache_" + strconv.Itoa(int(record.part)))

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



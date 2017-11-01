// 基于哈希分区的KV嵌入式数据库

// 索引文件结构  ：数据0文件偏移量(5) 数据0文件列表分配大小(2 buckets) 数据0文件列表真实长度(3)
// 元数据文件结构1：[数据项长度(10bit) 键值分配长度(12bit buckets) 键值真实长度(2) 数据项类型(2bit - 0) 键值(变长,最大5) 键名(变长)](变长,链表)
// 元数据文件结构2：[数据项长度(10bit) 键值分配长度(12bit buckets) 键值真实长度(2) 数据项类型(2bit - 1|2|3) 数据文件偏移量(5) 键名(变长)](变长,链表)
// 数据文件结构  ：键值(变长)
// 数据项类型 :
// 0: 元数据文件中的数据项键值放在第5项中，最大长度为5byte
// 1: 元数据文件中的数据项键值放在数据文件中，第5项为数据文件中的索引位置，键值cap存放gBUCKET_SIZE的倍数
// 2: 保留
// 3: 保留

package gkvdb

import (
    "os"
    "g/os/gfile"
    "strings"
    "g/encoding/gbinary"
    "g/os/gfilepool"
    "errors"
    "g/encoding/ghash"
    "bytes"
    "strconv"
    "g/os/gfilespace"
    "sync"
)

const (
    gPARTITION_SIZE          = 1048576                    // 哈希表分区大小(大小约为10MB)
    gMAX_KEY_SIZE            = (0xFFFF >> 6) - 10         // 键名最大长度(1013)
    gMAX_VALUE_SIZE          = 0xFFFF                     // 键值最大长度(65535)
    gBUCKET_SIZE             = 64                         // 元数据文件文件列表分块大小(byte, 值越大，初始化时占用的空间越大)
    gFILE_POOL_CACHE_TIMEOUT = 60                         // 文件指针池缓存时间(秒)
)

// KV数据库
type DB struct {
    mu    sync.RWMutex
    path  string            // 数据文件存放目录路径
    name  string            // 数据文件名
    ixfp  *gfilepool.Pool   // 索引文件打开指针池(用以高并发下的IO复用)
    mtfp  *gfilepool.Pool   // 元数据文件打开指针池(元数据，包含索引信息和部分数据信息)
    dbfp  *gfilepool.Pool   // 数据文件打开指针池(纯键值存储)
    mtsp  *gfilespace.Space // 元数据文件碎片管理
    dbsp  *gfilespace.Space // 数据文件碎片管理器
}

// KV数据检索记录
type Record struct {
    hash64    uint64  // 64位的hash code
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
        index  int    // 列表匹配的索引位置
        match  bool   // 是否在查找中准确匹配key
    }
    db struct {
        start  int64  // 数据文件中的开始地址
        end    int64  // 数据文件中的结束地址
        vcap   uint   // 键值允许存放的的最大长度（用以修改对比）
        klen   uint   // 键名大小
        vlen   uint   // 键值大小(byte)
        vtype  uint   // 键值类型
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
    if gfile.Exists(ixpath) && (!gfile.IsWritable(ixpath) || !gfile.IsReadable(ixpath)){
        return nil, errors.New("permission denied to index file: " + ixpath)
    }
    if gfile.Exists(mtpath) && (!gfile.IsWritable(mtpath) || !gfile.IsReadable(mtpath)){
        return nil, errors.New("permission denied to meta file: " + mtpath)
    }
    if gfile.Exists(dbpath) && (!gfile.IsWritable(dbpath) || !gfile.IsReadable(dbpath)){
        return nil, errors.New("permission denied to data file: " + dbpath)
    }
    // 创建文件指针池
    db.ixfp = gfilepool.New(ixpath, os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
    db.mtfp = gfilepool.New(mtpath, os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
    db.dbfp = gfilepool.New(dbpath, os.O_RDWR|os.O_CREATE, gFILE_POOL_CACHE_TIMEOUT)
    db.initFileSpace()
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

// 初始化碎片管理器，计算当前数据库存在的碎片，保存到内存中用于运行时计算和分配
func (db *DB) initFileSpace() {
    if db.mtsp == nil {
        db.mtsp  = gfilespace.New()
    }
    if db.dbsp == nil {
        db.dbsp = gfilespace.New()
    }
    pf, _ := db.mtfp.File()
    defer pf.Close()

    usedmtsp := gfilespace.New()
    useddbsp := gfilespace.New()
    ixbuffer := gfile.GetBinContents(db.getIndexFilePath())
    for i := 0; i < len(ixbuffer); i += 10 {
        mtindex := gbinary.DecodeToInt64(ixbuffer[i : i + 5])
        mtcap   := int(gbinary.DecodeToUint16(ixbuffer[i + 5 : i + 7])*gBUCKET_SIZE)
        mtsize  := int(gbinary.DecodeToUint32(ixbuffer[i + 7 : i + 10]))
        if mtcap > 0 {
            usedmtsp.AddBlock(int(mtindex), uint(mtcap))
            // 获取数据列表
            if buffer := gfile.GetBinContentByTwoOffsets(pf.File(), mtindex, mtindex + int64(mtsize)); buffer != nil {
                bits   := gbinary.DecodeBytesToBits(buffer[0:5])
                dbcap  := gbinary.DecodeBits(bits[10 : 22])
                dbtype := gbinary.DecodeBits(bits[38 : 40])
                if dbtype > 0 {
                    dbindex := gbinary.DecodeToInt64(buffer[5 : 10])
                    useddbsp.AddBlock(int(dbindex), uint(dbcap)*gBUCKET_SIZE)
                }
            }
        }
    }
    start := 0
    for _, v := range usedmtsp.GetAllBlocksByIndex() {
        if v.Index() > start {
            db.dbsp.AddBlock(start, uint(v.Index() - start))
            start = v.Index() + int(v.Size())
        }
    }
    start  = 0
    for _, v := range useddbsp.GetAllBlocksByIndex() {
        if v.Index() > start {
            db.mtsp.AddBlock(start, uint(v.Index() - start))
            start = v.Index() + int(v.Size())
        }
    }
    //fmt.Println(db.mtsp.GetAllBlocksByIndex())
    //fmt.Println(db.dbsp.GetAllBlocksByIndex())
}

// 计算关键字的hash code
func (db *DB) getHash(key []byte) uint64 {
    return ghash.BKDRHash64(key)
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
    record.ix.start = record.part*10
    record.ix.end   = record.ix.start + 10
    if buffer := gfile.GetBinContentByTwoOffsets(pf.File(), record.ix.start, record.ix.end); buffer != nil {
        //fmt.Println("get index:",buffer)
        record.mt.start = gbinary.DecodeToInt64(buffer[0:5])
        record.mt.cap   = uint(gbinary.DecodeToUint16(buffer[5:7])*gBUCKET_SIZE)
        record.mt.size  = uint(gbinary.DecodeToUint32(buffer[7:10]))
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
        // 线性查找
        for i := 0; i < len(record.mt.buffer); {
            buffer := record.mt.buffer[i:]
            bits   := gbinary.DecodeBytesToBits(buffer[0:5])
            length := gbinary.DecodeBits(bits[0 : 10])
            key    := buffer[10 : length]
            if bytes.Compare(key, record.key) == 0 {
                record.mt.index  = i
                record.mt.match  = true
                record.db.klen   = length - 10
                record.db.vcap   = gbinary.DecodeBits(bits[10 : 22])
                record.db.vlen   = gbinary.DecodeBits(bits[22 : 38])
                record.db.vtype  = gbinary.DecodeBits(bits[38 : 40])
                if record.db.vtype == 0 {
                    record.value = buffer[5 : 5 + record.db.vlen]
                } else {
                    record.db.vcap  = record.db.vcap*gBUCKET_SIZE
                    record.db.start = gbinary.DecodeToInt64(buffer[5 : 10])
                    record.db.end   = record.db.start + int64(record.db.vlen)
                }
                break
            } else {
                i += int(length)
            }
        }
    }
    return nil
}

// 查询检索信息
func (db *DB) getRecordByKey(key []byte) (*Record, error) {
    hash64 := db.getHash(key)
    part   := db.getPartitionByHash64(hash64)
    record := &Record {
        hash64  : hash64,
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
//fmt.Println(record)
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
        buffer := gfile.GetBinContentByTwoOffsets(pf.File(), record.db.start, record.db.end)
        if buffer != nil {
            return buffer, nil
        }
    }
    return nil, nil
}

// 关闭数据库链接
func (db *DB) Close() {
    db.ixfp.Close()
    db.mtfp.Close()
    db.dbfp.Close()
}

// 查询KV数据
func (db *DB) Get(key []byte) []byte {
    value, _ := db.getValueByKey(key)
    return value
}


// 设置KV数据
func (db *DB) Set(key []byte, value []byte) error {
    if len(key) > gMAX_KEY_SIZE {
        return errors.New("too large key size, max allowed: " + strconv.Itoa(gMAX_KEY_SIZE) + " bytes")
    }
    if len(value) > gMAX_VALUE_SIZE {
        return errors.New("too large value size, max allowed: " + strconv.Itoa(gMAX_VALUE_SIZE) + " bytes")
    }

    // 查询索引信息
    record, err := db.getRecordByKey(key)
    if err != nil {
        return err
    }

    // 写入数据文件，并更新record信息
    if err := db.insertDataByRecord(key, value, record); err != nil {
        return errors.New("inserting data error: " + err.Error())
    }

    // 根据record信息更新索引文件
    if err := db.updateIndexByRecord(record); err != nil {
        return errors.New("creating index error: " + err.Error())
    }
    return nil
}

// 删除KV数据
func (db *DB) Remove(key []byte) error {
    // 查询索引信息
    record, err := db.getRecordByKey(key)
    if err != nil {
        return err
    }
    // 如果找到匹配才执行删除操作
    if record.mt.match {
        return db.removeDataByRecord(record)
    }
    return nil
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
    if record.db.vtype == 0 {
        return nil
    }

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
    db.dbsp.AddBlock(int(record.db.start), record.db.vcap)
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
    endindex := record.mt.index + int(record.db.klen + 10)
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
    // 如果列表被清空，那么添加到碎片管理器
    if record.mt.size == 0 {
        db.mtsp.AddBlock(int(record.mt.start), record.mt.cap)
    }
    return nil
}

// 从索引中删除指定数据
func (db *DB) removeDataFromIx(record *Record) error {
    return db.updateIndexByRecord(record)
}

// 插入一条KV数据
func (db *DB) insertDataByRecord(key []byte, value []byte, record *Record) error {
    // type0可以变为type1,但是一旦成为type1将无法逆转，否则会造成碎片管理问题
    if record.db.vtype == 0 && len(value) > 5 {
        record.db.vtype = 1
    }
    record.db.klen = uint(len(key))
    record.db.vlen = uint(len(value))
    if record.db.vtype > 0 {
        if err := db.insertDataIntoDb(key, value, record); err != nil {
            return err
        }
        // 改变value的值为db文件偏移地址
        value = gbinary.EncodeUint64(uint64(record.db.start))[0:5]
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
    if record.db.end <= 0 || record.db.vcap < record.db.vlen {
        // 不用的空间添加到碎片管理器
        if record.db.end > 0 && record.db.vcap > 0 {
            db.dbsp.AddBlock(int(record.db.start), record.db.vcap)
        }
        // 重新计算所需空间
        if record.db.vcap < record.db.vlen {
            for {
                record.db.vcap += gBUCKET_SIZE
                if record.db.vcap >= record.db.vlen {
                    break
                }
            }
        }
        // 首先从碎片管理器中获取，如果不够，那么再从文件末尾分配
        index, size := db.dbsp.GetBlock(record.db.vcap)
        if index >= 0 {
            // 如果还可以对分配的空间进行切分
            extra := size - record.db.vcap
            if extra > gBUCKET_SIZE {
                db.dbsp.AddBlock(index + int(record.db.vcap), extra)
            }
            record.db.start = int64(index)
            record.db.end   = int64(index) + int64(record.db.vlen)
        } else {
            start, err := pf.File().Seek(0, 2)
            if err != nil {
                return err
            }
            record.db.start = start
            record.db.end   = start + int64(record.db.vlen)
        }
    }
    // vlen不够vcap的对末尾进行补0占位
    vbuffer := make([]byte, 0)
    vbuffer  = append(vbuffer, value...)
    for i := 0; i < int(record.db.vcap - record.db.vlen); i++ {
        vbuffer = append(vbuffer, byte(0))
    }
    if _, err = pf.File().WriteAt(vbuffer, record.db.start); err != nil {
        return err
    }
    return nil
}

// 将数据写入到元数据文件中，并更新信息到record
func (db *DB) insertDataIntoMt(key []byte, value []byte, record *Record) error {
    pf, err := db.mtfp.File()
    if err != nil {
        return err
    }
    defer pf.Close()
    bits   := make([]uint8, 0)
    data   := make([]byte, 0)
    buffer := make([]byte, 0)
    // 二进制打包
    bits = gbinary.EncodeBits(bits, record.db.klen + 10, 10)
    bits = gbinary.EncodeBits(bits, record.db.vcap/gBUCKET_SIZE, 12)
    bits = gbinary.EncodeBits(bits, record.db.vlen, 16)
    bits = gbinary.EncodeBits(bits, record.db.vtype, 2)
    data = append(data, gbinary.EncodeBitsToBytes(bits)...)
    data = append(data, value...)
    // vlen不够vcap的对末尾进行补0占位
    for i := 0; i < 5 - len(value); i++ {
        data  = append(data, byte(0))
    }
    data = append(data, key...)
    // 数据列表打包
    buffer = append(buffer, data...)
    if len(record.mt.buffer) > 0 {
        if record.mt.match {
            buffer = append(buffer, record.mt.buffer[0 : record.mt.index]...)
            buffer = append(buffer, record.mt.buffer[record.mt.index + 10 + int(record.db.klen) :]...)
        } else {
            buffer = append(buffer, record.mt.buffer...)
        }
    }
    // 判断数据列表空间是否足够
    record.mt.size = uint(len(buffer))
    if record.mt.end <= 0 || record.mt.cap < record.mt.size {
        // 不用的空间添加到碎片管理器
        if record.mt.end > 0 && record.mt.cap > 0 {
            db.mtsp.AddBlock(int(record.mt.start), uint(record.mt.cap))
        }
        // 重新计算所需空间
        if record.mt.cap < record.mt.size {
            for {
                record.mt.cap += gBUCKET_SIZE
                if record.mt.cap >= record.mt.size {
                    break
                }
            }
        }
        // 首先从碎片管理器中获取，如果不够，那么再从文件末尾分配
        index, size := db.mtsp.GetBlock(record.mt.cap)
        if index >= 0 {
            // 如果还可以对分配的空间进行切分
            extra := size - record.mt.cap
            if extra > gBUCKET_SIZE {
                db.mtsp.AddBlock(index + int(record.mt.cap), extra)
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
    // vlen不够vcap的对末尾进行补0占位
    for i := 0; i < int(record.mt.cap - record.mt.size); i++ {
        buffer = append(buffer, byte(0))
    }
    if _, err = pf.File().WriteAt(buffer, record.mt.start); err != nil {
        return err
    }
    return nil
}

// 根据record更新索引信息
func (db *DB) updateIndexByRecord(record *Record) error {
    ixpf, err := db.ixfp.File()
    if err != nil {
        return err
    }
    defer ixpf.Close()

    buffer := make([]byte, 0)
    if record.mt.size > 0 {
        mtcap  := uint16(record.mt.cap/gBUCKET_SIZE)
        mtsize := uint32(record.mt.size)
        buffer  = append(buffer, gbinary.EncodeInt64(record.mt.start)[0:5]...)
        buffer  = append(buffer, gbinary.EncodeUint16(mtcap)...)
        buffer  = append(buffer, gbinary.EncodeUint32(mtsize)[0:3]...)
    } else {
        buffer = make([]byte, 10)
    }
    if _, err = ixpf.File().WriteAt(buffer, record.part*10); err != nil {
        return err
    }
    return nil
}



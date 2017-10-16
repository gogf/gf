package gkvdb

import (
    "g/encoding/gcrc32"
    "os"
    "g/os/gfile"
    "strings"
    "fmt"
    "g/os/gcache"
    "g/core/types/gmap"
    "encoding/json"
    "g/encoding/gcompress"
    "g/encoding/gbinary"
)

const (
    gCACHE_TIMEOUT = 60  // bucket的读取缓存过期时间(秒)
    gBUCKET_SIZE   = 100 // 每个数据块的数据集大小（约数）
)

type DB struct {
    path   string   // 数据文件存放目录路径
    prefix string   // 数据文件名前缀
    ixfile *os.File // 索引文件打开指针
    dbfile *os.File // 数据文件打开指针
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
    ixfile, err := gfile.OpenWithFlag(path + gfile.Separator + prefix + ".ix", os.O_RDWR|os.O_CREATE)
    if err != nil {
        return nil,err
    }
    dbfile, err := gfile.OpenWithFlag(path + gfile.Separator + prefix + ".db", os.O_RDWR|os.O_CREATE)
    if err != nil {
        return nil,err
    }
    return &DB{
        path   : path,
        prefix : prefix,
        ixfile : ixfile,
        dbfile : dbfile,
    }, nil
}

func (db *DB) getHashCode(k string) uint32 {
    return gcrc32.EncodeString(k)
}

// 计算键名对应到数据文件中的偏移量
func (db *DB) getOffset(k string) int64 {
    return 0
    code   := db.getHashCode(k)
    offset := int64(code/gBUCKET_SIZE)
    if offset > 0 {
        offset = (offset - 1)*16
    }
    return offset
}

// 获取数据缓存键名
func (db *DB) getCacheKey(offset int64) string {
    return fmt.Sprintf("gkvdb_%d", offset)
}

// 查询KV数据
func (db *DB) Get(key string) ([]byte, error) {
    // 首先查询缓存
    offset   := db.getOffset(key)
    cachekey := db.getCacheKey(offset)
    result   := gcache.Get(cachekey)
    if result != nil {
        gm := result.(*gmap.StringInterfaceMap)
        r := gm.Get(key)
        if r != nil {
            return r.([]byte), nil
        }
    }
    // 其次查询文件(只需要查询数据长度即可)
    position := make([]byte, 16)
    _, err   := db.ixfile.ReadAt(position, offset)
    if err != nil {
        return nil, err
    }

    start,_  := gbinary.DecodeToInt64(position[0:8])
    end,_    := gbinary.DecodeToInt64(position[8:])
    buffer   := make([]byte, end - start)
    if _, err  = db.dbfile.ReadAt(buffer, start); err != nil {
        return nil, err
    }
    // 解析压缩数据
    m := make(map[string][]byte)
    if err = json.Unmarshal(gcompress.UnZlib(buffer), &m); err != nil {
        return nil, err
    }
    gm := gmap.NewStringInterfaceMap()
    for k, v := range m {
        gm.Set(k, v)
    }
    gcache.Set(cachekey, gm, gCACHE_TIMEOUT*1000)
    return buffer, nil
}

// 设置KV数据
func (db *DB) Set(key string, value []byte) error {
    var start int64 = 0
    var end   int64 = 0
    var cap   int32 = 0
    offset := db.getOffset(key)
    buffer := make([]byte, 20)
    _, err := db.ixfile.ReadAt(buffer, offset)
    if err == nil {
        start,_  = gbinary.DecodeToInt64(buffer[0:8])
        end,_    = gbinary.DecodeToInt64(buffer[8:])
        cap,_    = gbinary.DecodeToInt32(buffer[16:])
    }
    // 写入数据文件
    length := int32(len(value))
    if cap >= length {
        // 如果原本的空间大小满足本次数据，那么直接覆写
        if _, err = db.dbfile.WriteAt(value, start); err != nil {
            return err
        }
    } else {
        // 否则从文件末尾重新写入数据，原本的空间依靠另外的清理线程回收处理
        pos, err := db.dbfile.Seek(0, 2)
        if err != nil {
            return err
        }
        if _, err = db.dbfile.WriteAt(value, pos); err != nil {
            return err
        }
        start = pos
        cap   = length
    }
    end = start + int64(length)
    // 写入索引文件
    buffer, err = gbinary.Encode(start, end, cap)
    if err != nil {
        return err
    }
    if _, err = db.ixfile.WriteAt(buffer, offset); err != nil {
        return err
    }
    // 删除缓存键名，以便下次读取的时候从文件更新
    gcache.Remove(db.getCacheKey(offset))
    return nil
}






package gkvdb

import (
    "g/core/types/gmap"
    "strconv"
    "errors"
    "g/encoding/ghash"
)

// 用于磁盘与接口之间的数据缓冲层，异步线程将会定期同步到磁盘
type MemTable struct {
    m  *gmap.UintInterfaceMap
    db *DB
}

// 数据项
type MemTableItem struct {
    key     []byte
    value   []byte
    deleted bool
}

// 创建一个MemTable
func newMemTable(db *DB) *MemTable {
    return &MemTable{
        m  : gmap.NewUintInterfaceMap(),
        db : db,
    }
}

// 计算哈希64值
func (table *MemTable) hash64(key []byte) uint {
    return uint(ghash.BKDRHash64(key))
}

// 保存
func (table *MemTable) set(key []byte, value []byte) error {
    if len(key) > gMAX_KEY_SIZE {
        return errors.New("too large key size, max allowed: " + strconv.Itoa(gMAX_KEY_SIZE) + " bytes")
    }
    if len(value) > gMAX_VALUE_SIZE {
        return errors.New("too large value size, max allowed: " + strconv.Itoa(gMAX_VALUE_SIZE) + " bytes")
    }
    table.m.Set(table.hash64(key), MemTableItem{key, value, false})
    return nil
}

// 获取
func (table *MemTable) get(key []byte) ([]byte, bool) {
    if v := table.m.Get(table.hash64(key)); v != nil {
        item := v.(MemTableItem)
        if item.deleted {
             return nil, true
        } else {
            return item.value, true
        }
    }
    return nil, false
}

// 删除
func (table *MemTable) remove(key []byte) error {
    table.m.Set(table.hash64(key), MemTableItem{key, nil, true})
    return nil
}

// 同步数据到磁盘
// 注意这里的map不会在delete之后马上释放，而是随着GC的逻辑缓慢进行
func (table *MemTable) sync() {
    if table.m.IsEmpty() {
        return
    }
    for _, k := range table.m.Keys() {
        v    := table.m.GetAndRemove(k)
        item := v.(MemTableItem)
        if item.deleted {
            table.db.remove(item.key)
        } else {
            table.db.set(item.key, item.value)
        }
    }
}

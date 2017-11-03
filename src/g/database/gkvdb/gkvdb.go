// 底层基于leveldb的KV数据库封装
// leveldb: https://github.com/syndtr/goleveldb

package gkvdb

import (
    "sync"
    "github.com/syndtr/goleveldb/leveldb"
    "g/os/gfile"
)

type DB struct {
    mu    sync.RWMutex
    path  string
    lvdb  *leveldb.DB
}

// 创建一个KV数据库
func New(path string) (*DB, error) {
    lvdb, err := leveldb.OpenFile(path, nil)
    if err != nil {
        return nil, err
    }
    return &DB{
        path : path,
        lvdb : lvdb,
    }, nil
}

// 关闭数据库链接
func (db *DB) Close() error {
    return db.lvdb.Close()
}

// 清空数据库数据
func (db *DB) Clear() error {
    if err := gfile.Remove(db.path); err != nil {
        return err
    }
    db.lvdb.Close()

    lvdb, err := leveldb.OpenFile(db.path, nil)
    if err != nil {
        return err
    }
    db.lvdb = lvdb
    return nil
}

// 查询KV数据
func (db *DB) Get(key []byte) []byte {
    value, _ := db.lvdb.Get(key, nil)
    return value
}


// 设置KV数据
func (db *DB) Set(key []byte, value []byte) error {
    if err := db.lvdb.Put(key, value, nil); err != nil {
        return err
    }
    return nil
}

// 是否含有键
func (db *DB) Contains(key []byte) bool {
    r, _ := db.lvdb.Has(key, nil)
    return r
}

// 删除KV数据
func (db *DB) Remove(key []byte) error {
    return db.lvdb.Delete(key, nil)
}

// 获得数据库大小
func (db *DB) Size() uint {
    size := uint(0)
    iter := db.lvdb.NewIterator(nil, nil)
    for iter.Next() {
        size++
    }
    iter.Release()
    return size
}

// 获得所有的键名
func (db *DB) Keys() [][]byte {
    keys := make([][]byte, 0)
    iter := db.lvdb.NewIterator(nil, nil)
    for iter.Next() {
        keys = append(keys, iter.Key())
    }
    iter.Release()
    return keys
}

// 遍历数据库，给定参数位遍历回调函数
func (db *DB) Iterate(f func(key, value []byte) bool) {
    iter := db.lvdb.NewIterator(nil, nil)
    for iter.Next() {
        if !f(iter.Key(), iter.Value()) {
            break
        }
    }
    iter.Release()
}

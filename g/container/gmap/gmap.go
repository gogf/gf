// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gmap provides kinds of concurrent-safe(alternative) maps.
//
// 并发安全MAP.
package gmap

import "gitee.com/johng/gf/g/internal/rwmutex"

// 注意:
// 1、这个Map是所有并发安全Map中效率最低的，如果对效率要求比较高的场合，请合理选择对应数据类型的Map；
// 2、这个Map的优点是使用简便，由于键值都是interface{}类型，因此对键值的数据类型要求不高；
// 3、底层实现比较类似于sync.Map；

type Map struct {
    mu *rwmutex.RWMutex
    m  map[interface{}]interface{}
}

func New(safe...bool) *Map {
    return NewMap(safe...)
}

func NewMap(unsafe...bool) *Map {
    return &Map{
        m  : make(map[interface{}]interface{}),
        mu : rwmutex.New(unsafe...),
    }
}

func NewFromArray(keys []interface{}, values []interface{}, unsafe...bool) *Map {
    m := make(map[interface{}]interface{})
    l := len(values)
    for i, k := range keys {
        if i < l {
            m[k] = values[i]
        } else {
            m[k] = interface{}(nil)
        }
    }
    return &Map{
        m  : m,
        mu : rwmutex.New(unsafe...),
    }
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (gm *Map) Iterator(f func (k interface{}, v interface{}) bool) {
    gm.mu.RLock()
    defer gm.mu.RUnlock()
    for k, v := range gm.m {
        if !f(k, v) {
            break
        }
    }
}

// 哈希表克隆
func (gm *Map) Clone() map[interface{}]interface{} {
    m := make(map[interface{}]interface{})
    gm.mu.RLock()
    for k, v := range gm.m {
        m[k] = v
    }
    gm.mu.RUnlock()
    return m
}

// 设置键值对
func (gm *Map) Set(key interface{}, val interface{}) {
    gm.mu.Lock()
    gm.m[key] = val
    gm.mu.Unlock()
}

// 批量设置键值对
func (gm *Map) BatchSet(m map[interface{}]interface{}) {
    gm.mu.Lock()
    for k, v := range m {
        gm.m[k] = v
    }
    gm.mu.Unlock()
}

// 获取键值
func (gm *Map) Get(key interface{}) interface{} {
    gm.mu.RLock()
    val, _ := gm.m[key]
    gm.mu.RUnlock()
    return val
}

// 设置kv缓存键值对，内部会对键名的存在性使用写锁进行二次检索确认，如果存在则不再写入；返回键名对应的键值。
// 在高并发下有用，防止数据写入的并发逻辑错误。
func (gm *Map) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
    gm.mu.Lock()
    defer gm.mu.Unlock()
    if v, ok := gm.m[key]; ok {
        return v
    }
    if f, ok := value.(func() interface {}); ok {
        value = f()
    }
    gm.m[key] = value
    return value
}

// 当键名存在时返回其键值，否则写入指定的键值
func (gm *Map) GetOrSet(key interface{}, value interface{}) interface{} {
    if v := gm.Get(key); v == nil {
        return gm.doSetWithLockCheck(key, value)
    } else {
        return v
    }
}

// 当键名存在时返回其键值，否则写入指定的键值，键值由指定的函数生成
func (gm *Map) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
    if v := gm.Get(key); v == nil {
        return gm.doSetWithLockCheck(key, f())
    } else {
        return v
    }
}

// 与GetOrSetFunc不同的是，f是在写锁机制内执行
func (gm *Map) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
    if v := gm.Get(key); v == nil {
        return gm.doSetWithLockCheck(key, f)
    } else {
        return v
    }
}

// 当键名不存在时写入，并返回true；否则返回false。
func (gm *Map) SetIfNotExist(key interface{}, value interface{}) bool {
    if !gm.Contains(key) {
        gm.doSetWithLockCheck(key, value)
        return true
    }
    return false
}

// 批量删除键值对
func (gm *Map) BatchRemove(keys []interface{}) {
    gm.mu.Lock()
    for _, key := range keys {
        delete(gm.m, key)
    }
    gm.mu.Unlock()
}

// 返回对应的键值，并删除该键值
func (gm *Map) Remove(key interface{}) interface{} {
    gm.mu.Lock()
    val, exists := gm.m[key]
    if exists {
        delete(gm.m, key)
    }
    gm.mu.Unlock()
    return val
}

// 返回键列表
func (gm *Map) Keys() []interface{} {
    gm.mu.RLock()
    keys := make([]interface{}, 0)
    for key, _ := range gm.m {
        keys = append(keys, key)
    }
    gm.mu.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (gm *Map) Values() []interface{} {
    gm.mu.RLock()
    vals := make([]interface{}, 0)
    for _, val := range gm.m {
        vals = append(vals, val)
    }
    gm.mu.RUnlock()
    return vals
}

// 是否存在某个键
func (gm *Map) Contains(key interface{}) bool {
    gm.mu.RLock()
    _, exists := gm.m[key]
    gm.mu.RUnlock()
    return exists
}

// 哈希表大小
func (gm *Map) Size() int {
    gm.mu.RLock()
    length := len(gm.m)
    gm.mu.RUnlock()
    return length
}

// 哈希表是否为空
func (gm *Map) IsEmpty() bool {
    gm.mu.RLock()
    empty := len(gm.m) == 0
    gm.mu.RUnlock()
    return empty
}

// 清空哈希表
func (gm *Map) Clear() {
    gm.mu.Lock()
    gm.m = make(map[interface{}]interface{})
    gm.mu.Unlock()
}

// 并发安全锁操作，使用自定义方法执行加锁修改操作
func (gm *Map) LockFunc(f func(m map[interface{}]interface{})) {
    gm.mu.Lock(true)
    defer gm.mu.Unlock(true)
    f(gm.m)
}

// 并发安全锁操作，使用自定义方法执行加锁读取操作
func (gm *Map) RLockFunc(f func(m map[interface{}]interface{})) {
    gm.mu.RLock(true)
    defer gm.mu.RUnlock(true)
    f(gm.m)
}

// 交换Map中的键和值.
func (gm *Map) Flip() {
    gm.mu.Lock()
    defer gm.mu.Unlock()
    n := make(map[interface{}]interface{}, len(gm.m))
    for i, v := range gm.m {
        n[v] = i
    }
    gm.m = n
}

// 合并两个Map.
func (gm *Map) Merge(m *Map) {
    gm.mu.Lock()
    defer gm.mu.Unlock()
    if m != gm {
        m.mu.RLock()
        defer m.mu.RUnlock()
    }
    for k, v := range m.m {
        gm.m[k] = v
    }
}
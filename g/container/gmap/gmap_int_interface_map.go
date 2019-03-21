// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.
//

package gmap

import (
    "github.com/gogf/gf/g/internal/rwmutex"
    "github.com/gogf/gf/g/util/gconv"
)

type IntInterfaceMap struct {
	mu *rwmutex.RWMutex
	m  map[int]interface{}
}

func NewIntInterfaceMap(unsafe...bool) *IntInterfaceMap {
	return &IntInterfaceMap{
        m  : make(map[int]interface{}),
        mu : rwmutex.New(unsafe...),
    }
}

func NewIntInterfaceMapFrom(m map[int]interface{}, unsafe...bool) *IntInterfaceMap {
    return &IntInterfaceMap{
        m  : m,
        mu : rwmutex.New(unsafe...),
    }
}

func NewIntInterfaceMapFromArray(keys []int, values []interface{}, unsafe...bool) *IntInterfaceMap {
    m := make(map[int]interface{})
    l := len(values)
    for i, k := range keys {
        if i < l {
            m[k] = values[i]
        } else {
            m[k] = interface{}(nil)
        }
    }
    return &IntInterfaceMap{
        m  : m,
        mu : rwmutex.New(unsafe...),
    }
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (gm *IntInterfaceMap) Iterator(f func (k int, v interface{}) bool) {
    gm.mu.RLock()
    defer gm.mu.RUnlock()
    for k, v := range gm.m {
        if !f(k, v) {
            break
        }
    }
}

// 哈希表克隆.
func (gm *IntInterfaceMap) Clone() *IntInterfaceMap {
    return NewIntInterfaceMapFrom(gm.Map(), !gm.mu.IsSafe())
}

// 返回当前哈希表的数据Map.
func (gm *IntInterfaceMap) Map() map[int]interface{} {
	m := make(map[int]interface{})
	gm.mu.RLock()
	for k, v := range gm.m {
		m[k] = v
	}
    gm.mu.RUnlock()
	return m
}

// 设置键值对
func (gm *IntInterfaceMap) Set(key int, val interface{}) {
	gm.mu.Lock()
	gm.m[key] = val
	gm.mu.Unlock()
}

// 批量设置键值对
func (gm *IntInterfaceMap) BatchSet(m map[int]interface{}) {
	gm.mu.Lock()
	for k, v := range m {
		gm.m[k] = v
	}
	gm.mu.Unlock()
}

// 获取键值
func (gm *IntInterfaceMap) Get(key int) (interface{}) {
	gm.mu.RLock()
	val, _ := gm.m[key]
	gm.mu.RUnlock()
	return val
}

// 设置kv缓存键值对，内部会对键名的存在性使用写锁进行二次检索确认，如果存在则不再写入；返回键名对应的键值。
// 在高并发下有用，防止数据写入的并发逻辑错误。
func (gm *IntInterfaceMap) doSetWithLockCheck(key int, value interface{}) interface{} {
    gm.mu.Lock()
    defer gm.mu.Unlock()
    if v, ok := gm.m[key]; ok {
        return v
    }
    if f, ok := value.(func() interface {}); ok {
        value = f()
    }
    if value != nil {
        gm.m[key] = value
    }
    return value
}

// 当键名存在时返回其键值，否则写入指定的键值
func (gm *IntInterfaceMap) GetOrSet(key int, value interface{}) interface{} {
    if v := gm.Get(key); v == nil {
        return gm.doSetWithLockCheck(key, value)
    } else {
        return v
    }
}

// 当键名存在时返回其键值，否则写入指定的键值，键值由指定的函数生成
func (gm *IntInterfaceMap) GetOrSetFunc(key int, f func() interface{}) interface{} {
    if v := gm.Get(key); v == nil {
        return gm.doSetWithLockCheck(key, f())
    } else {
        return v
    }
}

// 与GetOrSetFunc不同的是，f是在写锁机制内执行
func (gm *IntInterfaceMap) GetOrSetFuncLock(key int, f func() interface{}) interface{} {
    if v := gm.Get(key); v == nil {
        return gm.doSetWithLockCheck(key, f)
    } else {
        return v
    }
}

// 当键名不存在时写入，并返回true；否则返回false。
func (gm *IntInterfaceMap) SetIfNotExist(key int, value interface{}) bool {
    if !gm.Contains(key) {
        gm.doSetWithLockCheck(key, value)
        return true
    }
    return false
}

// 批量删除键值对
func (gm *IntInterfaceMap) BatchRemove(keys []int) {
    gm.mu.Lock()
    for _, key := range keys {
        delete(gm.m, key)
    }
    gm.mu.Unlock()
}

// 返回对应的键值，并删除该键值
func (gm *IntInterfaceMap) Remove(key int) interface{} {
    gm.mu.Lock()
    val, exists := gm.m[key]
    if exists {
        delete(gm.m, key)
    }
    gm.mu.Unlock()
    return val
}

// 返回键列表
func (gm *IntInterfaceMap) Keys() []int {
    gm.mu.RLock()
    keys := make([]int, 0)
    for key, _ := range gm.m {
        keys = append(keys, key)
    }
    gm.mu.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (gm *IntInterfaceMap) Values() []interface{} {
    gm.mu.RLock()
    vals := make([]interface{}, 0)
    for _, val := range gm.m {
        vals = append(vals, val)
    }
    gm.mu.RUnlock()
    return vals
}

// 是否存在某个键
func (gm *IntInterfaceMap) Contains(key int) bool {
    gm.mu.RLock()
    _, exists := gm.m[key]
    gm.mu.RUnlock()
    return exists
}

// 哈希表大小
func (gm *IntInterfaceMap) Size() int {
    gm.mu.RLock()
    length := len(gm.m)
    gm.mu.RUnlock()
    return length
}

// 哈希表是否为空
func (gm *IntInterfaceMap) IsEmpty() bool {
    gm.mu.RLock()
    empty := len(gm.m) == 0
    gm.mu.RUnlock()
    return empty
}

// 清空哈希表
func (gm *IntInterfaceMap) Clear() {
    gm.mu.Lock()
    gm.m = make(map[int]interface{})
    gm.mu.Unlock()
}

// 并发安全锁操作，使用自定义方法执行加锁修改操作
func (gm *IntInterfaceMap) LockFunc(f func(m map[int]interface{})) {
    gm.mu.Lock()
    defer gm.mu.Unlock()
    f(gm.m)
}

// 并发安全锁操作，使用自定义方法执行加锁读取操作
func (gm *IntInterfaceMap) RLockFunc(f func(m map[int]interface{})) {
    gm.mu.RLock()
    defer gm.mu.RUnlock()
    f(gm.m)
}

// 交换Map中的键和值.
func (gm *IntInterfaceMap) Flip() {
    gm.mu.Lock()
    defer gm.mu.Unlock()
    n := make(map[int]interface{}, len(gm.m))
    for k, v := range gm.m {
        n[gconv.Int(v)] = k
    }
    gm.m = n
}

// 合并两个Map.
func (gm *IntInterfaceMap) Merge(m *IntInterfaceMap) {
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
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

type StringInterfaceMap struct {
	mu *rwmutex.RWMutex
	m  map[string]interface{}
}

func NewStringInterfaceMap(unsafe...bool) *StringInterfaceMap {
	return &StringInterfaceMap{
		m  : make(map[string]interface{}),
		mu : rwmutex.New(unsafe...),
	}
}

func NewStringInterfaceMapFrom(m map[string]interface{}, unsafe...bool) *StringInterfaceMap {
    return &StringInterfaceMap{
        m  : m,
        mu : rwmutex.New(unsafe...),
    }
}

func NewStringInterfaceMapFromArray(keys []string, values []interface{}, unsafe...bool) *StringInterfaceMap {
    m := make(map[string]interface{})
    l := len(values)
    for i, k := range keys {
        if i < l {
            m[k] = values[i]
        } else {
            m[k] = interface{}(nil)
        }
    }
    return &StringInterfaceMap{
        m  : m,
        mu : rwmutex.New(unsafe...),
    }
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (gm *StringInterfaceMap) Iterator(f func (k string, v interface{}) bool) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	for k, v := range gm.m {
		if !f(k, v) {
			break
		}
	}
}

// 哈希表克隆.
func (gm *StringInterfaceMap) Clone() *StringInterfaceMap {
    return NewStringInterfaceMapFrom(gm.Map(), !gm.mu.IsSafe())
}

// 返回当前哈希表的数据Map.
func (gm *StringInterfaceMap) Map() map[string]interface{} {
    m := make(map[string]interface{})
    gm.mu.RLock()
    for k, v := range gm.m {
        m[k] = v
    }
    gm.mu.RUnlock()
    return m
}

// 设置键值对
func (gm *StringInterfaceMap) Set(key string, val interface{}) {
	gm.mu.Lock()
	gm.m[key] = val
	gm.mu.Unlock()
}

// 批量设置键值对
func (gm *StringInterfaceMap) BatchSet(m map[string]interface{}) {
    gm.mu.Lock()
    for k, v := range m {
        gm.m[k] = v
    }
    gm.mu.Unlock()
}

// 获取键值
func (gm *StringInterfaceMap) Get(key string) interface{} {
	gm.mu.RLock()
	val, _ := gm.m[key]
	gm.mu.RUnlock()
	return val
}

// 设置kv缓存键值对，内部会对键名的存在性使用写锁进行二次检索确认，如果存在则不再写入；返回键名对应的键值。
// 在高并发下有用，防止数据写入的并发逻辑错误。
func (gm *StringInterfaceMap) doSetWithLockCheck(key string, value interface{}) interface{} {
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
func (gm *StringInterfaceMap) GetOrSet(key string, value interface{}) interface{} {
    if v := gm.Get(key); v == nil {
        return gm.doSetWithLockCheck(key, value)
    } else {
        return v
    }
}

// 当键名存在时返回其键值，否则写入指定的键值，键值由指定的函数生成
func (gm *StringInterfaceMap) GetOrSetFunc(key string, f func() interface{}) interface{} {
    if v := gm.Get(key); v == nil {
        return gm.doSetWithLockCheck(key, f())
    } else {
        return v
    }
}

// 与GetOrSetFunc不同的是，f是在写锁机制内执行
func (gm *StringInterfaceMap) GetOrSetFuncLock(key string, f func() interface{}) interface{} {
    if v := gm.Get(key); v == nil {
        return gm.doSetWithLockCheck(key, f)
    } else {
        return v
    }
}

// 当键名不存在时写入，并返回true；否则返回false。
func (gm *StringInterfaceMap) SetIfNotExist(key string, value interface{}) bool {
    if !gm.Contains(key) {
        gm.doSetWithLockCheck(key, value)
        return true
    }
    return false
}

// 批量删除键值对
func (gm *StringInterfaceMap) BatchRemove(keys []string) {
    gm.mu.Lock()
    for _, key := range keys {
        delete(gm.m, key)
    }
    gm.mu.Unlock()
}

// 返回对应的键值，并删除该键值
func (gm *StringInterfaceMap) Remove(key string) interface{} {
	gm.mu.Lock()
	val, exists := gm.m[key]
	if exists {
		delete(gm.m, key)
	}
	gm.mu.Unlock()
	return val
}

// 返回键列表
func (gm *StringInterfaceMap) Keys() []string {
	gm.mu.RLock()
	keys := make([]string, 0)
	for key, _ := range gm.m {
		keys = append(keys, key)
	}
    gm.mu.RUnlock()
	return keys
}

// 返回值列表(注意是随机排序)
func (gm *StringInterfaceMap) Values() []interface{} {
	gm.mu.RLock()
	vals := make([]interface{}, 0)
	for _, val := range gm.m {
		vals = append(vals, val)
	}
	gm.mu.RUnlock()
	return vals
}

// 是否存在某个键
func (gm *StringInterfaceMap) Contains(key string) bool {
	gm.mu.RLock()
	_, exists := gm.m[key]
	gm.mu.RUnlock()
	return exists
}

// 哈希表大小
func (gm *StringInterfaceMap) Size() int {
	gm.mu.RLock()
	length := len(gm.m)
	gm.mu.RUnlock()
	return length
}

// 哈希表是否为空
func (gm *StringInterfaceMap) IsEmpty() bool {
	gm.mu.RLock()
	empty := len(gm.m) == 0
	gm.mu.RUnlock()
	return empty
}

// 清空哈希表
func (gm *StringInterfaceMap) Clear() {
    gm.mu.Lock()
    gm.m = make(map[string]interface{})
    gm.mu.Unlock()
}

// 并发安全写锁操作，使用自定义方法执行加锁修改操作
func (gm *StringInterfaceMap) LockFunc(f func(m map[string]interface{})) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	f(gm.m)
}

// 并发安全读锁操作，使用自定义方法执行加锁读取操作
func (gm *StringInterfaceMap) RLockFunc(f func(m map[string]interface{})) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	f(gm.m)
}

// 交换Map中的键和值.
func (gm *StringInterfaceMap) Flip() {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	n := make(map[string]interface{}, len(gm.m))
	for k, v := range gm.m {
		n[gconv.String(v)] = k
	}
	gm.m = n
}

// 合并两个Map.
func (gm *StringInterfaceMap) Merge(m *StringInterfaceMap) {
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
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package gmap

import (
	"sync"
    "gitee.com/johng/gf/g/util/gconv"
)

type StringInterfaceMap struct {
	mu sync.RWMutex
	m  map[string]interface{}
}

func NewStringInterfaceMap() *StringInterfaceMap {
	return &StringInterfaceMap{
		m: make(map[string]interface{}),
	}
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (this *StringInterfaceMap) Iterator(f func (k string, v interface{}) bool) {
	this.mu.RLock()
	for k, v := range this.m {
		if !f(k, v) {
			break
		}
	}
	this.mu.RUnlock()
}

// 哈希表克隆
func (this *StringInterfaceMap) Clone() *map[string]interface{} {
    m := make(map[string]interface{})
    this.mu.RLock()
    for k, v := range this.m {
        m[k] = v
    }
    this.mu.RUnlock()
    return &m
}

// 设置键值对
func (this *StringInterfaceMap) Set(key string, val interface{}) {
	this.mu.Lock()
	this.m[key] = val
	this.mu.Unlock()
}

// 批量设置键值对
func (this *StringInterfaceMap) BatchSet(m map[string]interface{}) {
    this.mu.Lock()
    for k, v := range m {
        this.m[k] = v
    }
    this.mu.Unlock()
}

// 获取键值
func (this *StringInterfaceMap) Get(key string) interface{} {
	this.mu.RLock()
	val, _ := this.m[key]
	this.mu.RUnlock()
	return val
}

// 获取键值，如果键值不存在则写入默认值
func (this *StringInterfaceMap) GetWithDefault(key string, value interface{}) interface{} {
	this.mu.Lock()
	val, ok := this.m[key]
	if !ok {
        this.m[key] = value
        val         = value
    }
	this.mu.Unlock()
	return val
}

func (this *StringInterfaceMap) GetBool(key string) bool {
    return gconv.Bool(this.Get(key))
}

func (this *StringInterfaceMap) GetInt(key string) int {
    return gconv.Int(this.Get(key))
}

func (this *StringInterfaceMap) GetUint (key string) uint {
    return gconv.Uint(this.Get(key))
}

func (this *StringInterfaceMap) GetFloat32 (key string) float32 {
    return gconv.Float32(this.Get(key))
}

func (this *StringInterfaceMap) GetFloat64 (key string) float64 {
    return gconv.Float64(this.Get(key))
}

func (this *StringInterfaceMap) GetString (key string) string {
    return gconv.String(this.Get(key))
}

// 删除键值对
func (this *StringInterfaceMap) Remove(key string) {
	this.mu.Lock()
	delete(this.m, key)
	this.mu.Unlock()
}

// 批量删除键值对
func (this *StringInterfaceMap) BatchRemove(keys []string) {
    this.mu.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.mu.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *StringInterfaceMap) GetAndRemove(key string) interface{} {
	this.mu.Lock()
	val, exists := this.m[key]
	if exists {
		delete(this.m, key)
	}
	this.mu.Unlock()
	return val
}

// 返回键列表
func (this *StringInterfaceMap) Keys() []string {
	this.mu.RLock()
	keys := make([]string, 0)
	for key, _ := range this.m {
		keys = append(keys, key)
	}
    this.mu.RUnlock()
	return keys
}

// 返回值列表(注意是随机排序)
func (this *StringInterfaceMap) Values() []interface{} {
	this.mu.RLock()
	vals := make([]interface{}, 0)
	for _, val := range this.m {
		vals = append(vals, val)
	}
	this.mu.RUnlock()
	return vals
}

// 是否存在某个键
func (this *StringInterfaceMap) Contains(key string) bool {
	this.mu.RLock()
	_, exists := this.m[key]
	this.mu.RUnlock()
	return exists
}

// 哈希表大小
func (this *StringInterfaceMap) Size() int {
	this.mu.RLock()
	length := len(this.m)
	this.mu.RUnlock()
	return length
}

// 哈希表是否为空
func (this *StringInterfaceMap) IsEmpty() bool {
	this.mu.RLock()
	empty := (len(this.m) == 0)
	this.mu.RUnlock()
	return empty
}

// 清空哈希表
func (this *StringInterfaceMap) Clear() {
    this.mu.Lock()
    this.m = make(map[string]interface{})
    this.mu.Unlock()
}

// 使用自定义方法执行加锁修改操作
func (this *StringInterfaceMap) LockFunc(f func(m map[string]interface{})) {
	this.mu.Lock()
	f(this.m)
	this.mu.Unlock()
}

// 使用自定义方法执行加锁读取操作
func (this *StringInterfaceMap) RLockFunc(f func(m map[string]interface{})) {
	this.mu.RLock()
	f(this.m)
	this.mu.RUnlock()
}

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

type IntInterfaceMap struct {
	mu sync.RWMutex
	m  map[int]interface{}
}

func NewIntInterfaceMap() *IntInterfaceMap {
	return &IntInterfaceMap{
        m: make(map[int]interface{}),
    }
}

// 哈希表克隆
func (this *IntInterfaceMap) Clone() *map[int]interface{} {
	m := make(map[int]interface{})
	this.mu.RLock()
	for k, v := range this.m {
		m[k] = v
	}
    this.mu.RUnlock()
	return &m
}

// 设置键值对
func (this *IntInterfaceMap) Set(key int, val interface{}) {
	this.mu.Lock()
	this.m[key] = val
	this.mu.Unlock()
}

// 批量设置键值对
func (this *IntInterfaceMap) BatchSet(m map[int]interface{}) {
	this.mu.Lock()
	for k, v := range m {
		this.m[k] = v
	}
	this.mu.Unlock()
}

// 获取键值
func (this *IntInterfaceMap) Get(key int) (interface{}) {
	this.mu.RLock()
	val, _ := this.m[key]
	this.mu.RUnlock()
	return val
}

func (this *IntInterfaceMap) GetBool(key int) bool {
    return gconv.Bool(this.Get(key))
}

func (this *IntInterfaceMap) GetInt(key int) int {
    return gconv.Int(this.Get(key))
}

func (this *IntInterfaceMap) GetUint (key int) uint {
    return gconv.Uint(this.Get(key))
}

func (this *IntInterfaceMap) GetFloat32 (key int) float32 {
    return gconv.Float32(this.Get(key))
}

func (this *IntInterfaceMap) GetFloat64 (key int) float64 {
    return gconv.Float64(this.Get(key))
}

func (this *IntInterfaceMap) GetString (key int) string {
    return gconv.String(this.Get(key))
}

// 删除键值对
func (this *IntInterfaceMap) Remove(key int) {
    this.mu.Lock()
    delete(this.m, key)
    this.mu.Unlock()
}

// 批量删除键值对
func (this *IntInterfaceMap) BatchRemove(keys []int) {
    this.mu.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.mu.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *IntInterfaceMap) GetAndRemove(key int) (interface{}) {
    this.mu.Lock()
    val, exists := this.m[key]
    if exists {
        delete(this.m, key)
    }
    this.mu.Unlock()
    return val
}

// 返回键列表
func (this *IntInterfaceMap) Keys() []int {
    this.mu.RLock()
    keys := make([]int, 0)
    for key, _ := range this.m {
        keys = append(keys, key)
    }
    this.mu.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *IntInterfaceMap) Values() []interface{} {
    this.mu.RLock()
    vals := make([]interface{}, 0)
    for _, val := range this.m {
        vals = append(vals, val)
    }
    this.mu.RUnlock()
    return vals
}

// 是否存在某个键
func (this *IntInterfaceMap) Contains(key int) bool {
    this.mu.RLock()
    _, exists := this.m[key]
    this.mu.RUnlock()
    return exists
}

// 哈希表大小
func (this *IntInterfaceMap) Size() int {
    this.mu.RLock()
    len := len(this.m)
    this.mu.RUnlock()
    return len
}

// 哈希表是否为空
func (this *IntInterfaceMap) IsEmpty() bool {
    this.mu.RLock()
    empty := (len(this.m) == 0)
    this.mu.RUnlock()
    return empty
}

// 清空哈希表
func (this *IntInterfaceMap) Clear() {
    this.mu.Lock()
    this.m = make(map[int]interface{})
    this.mu.Unlock()
}


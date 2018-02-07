// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package gmap

import (
	"sync"
)

type IntBoolMap struct {
	mu sync.RWMutex
	m  map[int]bool
}

func NewIntBoolMap() *IntBoolMap {
	return &IntBoolMap{
        m: make(map[int]bool),
    }
}

// 哈希表克隆
func (this *IntBoolMap) Clone() *map[int]bool {
	m := make(map[int]bool)
	this.mu.RLock()
	for k, v := range this.m {
		m[k] = v
	}
    this.mu.RUnlock()
	return &m
}

// 给定回调函数对原始内容进行遍历
func (this *IntBoolMap) Iterator(f func (k int, v bool)) {
    this.mu.RLock()
    for k, v := range this.m {
        f(k, v)
    }
    this.mu.RUnlock()
}

// 设置键值对
func (this *IntBoolMap) Set(key int, val bool) {
	this.mu.Lock()
	this.m[key] = val
	this.mu.Unlock()
}

// 批量设置键值对
func (this *IntBoolMap) BatchSet(m map[int]bool) {
	this.mu.Lock()
	for k, v := range m {
		this.m[k] = v
	}
	this.mu.Unlock()
}

// 获取键值
func (this *IntBoolMap) Get(key int) (bool) {
	this.mu.RLock()
	val, _ := this.m[key]
	this.mu.RUnlock()
	return val
}

// 删除键值对
func (this *IntBoolMap) Remove(key int) {
    this.mu.Lock()
    delete(this.m, key)
    this.mu.Unlock()
}

// 批量删除键值对
func (this *IntBoolMap) BatchRemove(keys []int) {
    this.mu.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.mu.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *IntBoolMap) GetAndRemove(key int) (bool) {
    this.mu.Lock()
    val, exists := this.m[key]
    if exists {
        delete(this.m, key)
    }
    this.mu.Unlock()
    return val
}

// 返回键列表
func (this *IntBoolMap) Keys() []int {
    this.mu.RLock()
    keys := make([]int, 0)
    for key, _ := range this.m {
        keys = append(keys, key)
    }
    this.mu.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
//func (this *IntBoolMap) Values() []bool {
//    this.mu.RLock()
//    vals := make([]bool, 0)
//    for _, val := range this.m {
//        vals = append(vals, val)
//    }
//    this.mu.RUnlock()
//    return vals
//}

// 是否存在某个键
func (this *IntBoolMap) Contains(key int) bool {
	this.mu.RLock()
	_, exists := this.m[key]
	this.mu.RUnlock()
	return exists
}

// 哈希表大小
func (this *IntBoolMap) Size() int {
    this.mu.RLock()
    len := len(this.m)
    this.mu.RUnlock()
    return len
}

// 哈希表是否为空
func (this *IntBoolMap) IsEmpty() bool {
    this.mu.RLock()
    empty := (len(this.m) == 0)
    this.mu.RUnlock()
    return empty
}

// 清空哈希表
func (this *IntBoolMap) Clear() {
    this.mu.Lock()
    this.m = make(map[int]bool)
    this.mu.Unlock()
}


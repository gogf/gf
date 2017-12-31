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

type StringStringMap struct {
	sync.RWMutex
	m map[string]string
}

func NewStringStringMap() *StringStringMap {
	return &StringStringMap{
		m: make(map[string]string),
	}
}

// 哈希表克隆
func (this *StringStringMap) Clone() *map[string]string {
    m := make(map[string]string)
    this.RLock()
    for k, v := range this.m {
        m[k] = v
    }
    this.RUnlock()
    return &m
}

// 设置键值对
func (this *StringStringMap) Set(key string, val string) {
	this.Lock()
	this.m[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *StringStringMap) BatchSet(m map[string]string) {
    this.Lock()
    for k, v := range m {
        this.m[k] = v
    }
    this.Unlock()
}

// 获取键值
func (this *StringStringMap) Get(key string) string {
	this.RLock()
	val, _ := this.m[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *StringStringMap) Remove(key string) {
	this.Lock()
	delete(this.m, key)
	this.Unlock()
}

// 批量删除键值对
func (this *StringStringMap) BatchRemove(keys []string) {
    this.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *StringStringMap) GetAndRemove(key string) string {
	this.Lock()
	val, exists := this.m[key]
	if exists {
		delete(this.m, key)
	}
	this.Unlock()
	return val
}

// 返回键列表
func (this *StringStringMap) Keys() []string {
	this.RLock()
	keys := make([]string, 0)
	for key, _ := range this.m {
		keys = append(keys, key)
	}
    this.RUnlock()
	return keys
}

// 返回值列表(注意是随机排序)
func (this *StringStringMap) Values() []string {
	this.RLock()
	vals := make([]string, 0)
	for _, val := range this.m {
		vals = append(vals, val)
	}
	this.RUnlock()
	return vals
}

// 是否存在某个键
func (this *StringStringMap) Contains(key string) bool {
	this.RLock()
	_, exists := this.m[key]
	this.RUnlock()
	return exists
}

// 哈希表大小
func (this *StringStringMap) Size() int {
	this.RLock()
	len := len(this.m)
	this.RUnlock()
	return len
}

// 哈希表是否为空
func (this *StringStringMap) IsEmpty() bool {
	this.RLock()
	empty := (len(this.m) == 0)
	this.RUnlock()
	return empty
}

// 清空哈希表
func (this *StringStringMap) Clear() {
    this.Lock()
    this.m = make(map[string]string)
    this.Unlock()
}

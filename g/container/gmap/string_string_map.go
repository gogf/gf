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
	mu sync.RWMutex
	m  map[string]string
}

func NewStringStringMap() *StringStringMap {
	return &StringStringMap{
		m: make(map[string]string),
	}
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (this *StringStringMap) Iterator(f func (k string, v string) bool) {
	this.mu.RLock()
	for k, v := range this.m {
		if !f(k, v) {
			break
		}
	}
	this.mu.RUnlock()
}

// 哈希表克隆
func (this *StringStringMap) Clone() *map[string]string {
    m := make(map[string]string)
    this.mu.RLock()
    for k, v := range this.m {
        m[k] = v
    }
    this.mu.RUnlock()
    return &m
}

// 设置键值对
func (this *StringStringMap) Set(key string, val string) {
	this.mu.Lock()
	this.m[key] = val
	this.mu.Unlock()
}

// 批量设置键值对
func (this *StringStringMap) BatchSet(m map[string]string) {
    this.mu.Lock()
    for k, v := range m {
        this.m[k] = v
    }
    this.mu.Unlock()
}

// 获取键值
func (this *StringStringMap) Get(key string) string {
	this.mu.RLock()
	val, _ := this.m[key]
	this.mu.RUnlock()
	return val
}

// 获取键值，如果键值不存在则写入默认值
func (this *StringStringMap) GetWithDefault(key string, value string) string {
	this.mu.Lock()
	val, ok := this.m[key]
	if !ok {
		this.m[key] = value
		val         = value
	}
	this.mu.Unlock()
	return val
}

// 删除键值对
func (this *StringStringMap) Remove(key string) {
	this.mu.Lock()
	delete(this.m, key)
	this.mu.Unlock()
}

// 批量删除键值对
func (this *StringStringMap) BatchRemove(keys []string) {
    this.mu.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.mu.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *StringStringMap) GetAndRemove(key string) string {
	this.mu.Lock()
	val, exists := this.m[key]
	if exists {
		delete(this.m, key)
	}
	this.mu.Unlock()
	return val
}

// 返回键列表
func (this *StringStringMap) Keys() []string {
	this.mu.RLock()
	keys := make([]string, 0)
	for key, _ := range this.m {
		keys = append(keys, key)
	}
    this.mu.RUnlock()
	return keys
}

// 返回值列表(注意是随机排序)
func (this *StringStringMap) Values() []string {
	this.mu.RLock()
	vals := make([]string, 0)
	for _, val := range this.m {
		vals = append(vals, val)
	}
	this.mu.RUnlock()
	return vals
}

// 是否存在某个键
func (this *StringStringMap) Contains(key string) bool {
	this.mu.RLock()
	_, exists := this.m[key]
	this.mu.RUnlock()
	return exists
}

// 哈希表大小
func (this *StringStringMap) Size() int {
	this.mu.RLock()
	length := len(this.m)
	this.mu.RUnlock()
	return length
}

// 哈希表是否为空
func (this *StringStringMap) IsEmpty() bool {
	this.mu.RLock()
	empty := (len(this.m) == 0)
	this.mu.RUnlock()
	return empty
}

// 清空哈希表
func (this *StringStringMap) Clear() {
    this.mu.Lock()
    this.m = make(map[string]string)
    this.mu.Unlock()
}

// 使用自定义方法执行加锁修改操作
func (this *StringStringMap) LockFunc(f func(m map[string]string)) {
	this.mu.Lock()
	f(this.m)
	this.mu.Unlock()
}

// 使用自定义方法执行加锁读取操作
func (this *StringStringMap) RLockFunc(f func(m map[string]string)) {
	this.mu.RLock()
	f(this.m)
	this.mu.RUnlock()
}

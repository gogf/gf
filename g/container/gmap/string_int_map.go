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

type StringIntMap struct {
	mu sync.RWMutex
	m  map[string]int
}

func NewStringIntMap() *StringIntMap {
	return &StringIntMap{
        m: make(map[string]int),
    }
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (this *StringIntMap) Iterator(f func (k string, v int) bool) {
    this.mu.RLock()
    for k, v := range this.m {
        if !f(k, v) {
            break
        }
    }
    this.mu.RUnlock()
}

// 哈希表克隆
func (this *StringIntMap) Clone() *map[string]int {
    m := make(map[string]int)
    this.mu.RLock()
    for k, v := range this.m {
        m[k] = v
    }
    this.mu.RUnlock()
    return &m
}

// 设置键值对
func (this *StringIntMap) Set(key string, val int) {
	this.mu.Lock()
	this.m[key] = val
	this.mu.Unlock()
}

// 批量设置键值对
func (this *StringIntMap) BatchSet(m map[string]int) {
	this.mu.Lock()
	for k, v := range m {
		this.m[k] = v
	}
	this.mu.Unlock()
}

// 获取键值
func (this *StringIntMap) Get(key string) int {
	this.mu.RLock()
	val, _ := this.m[key]
    this.mu.RUnlock()
	return val
}

// 获取键值，如果键值不存在则写入默认值
func (this *StringIntMap) GetWithDefault(key string, value int) int {
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
func (this *StringIntMap) Remove(key string) {
    this.mu.Lock()
    delete(this.m, key)
    this.mu.Unlock()
}

// 批量删除键值对
func (this *StringIntMap) BatchRemove(keys []string) {
    this.mu.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.mu.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *StringIntMap) GetAndRemove(key string) (int) {
    this.mu.Lock()
    val, exists := this.m[key]
    if exists {
        delete(this.m, key)
    }
    this.mu.Unlock()
    return val
}

// 返回键列表
func (this *StringIntMap) Keys() []string {
    this.mu.RLock()
    keys := make([]string, 0)
    for key, _ := range this.m {
        keys = append(keys, key)
    }
    this.mu.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *StringIntMap) Values() []int {
    this.mu.RLock()
    vals := make([]int, 0)
    for _, val := range this.m {
        vals = append(vals, val)
    }
    this.mu.RUnlock()
    return vals
}

// 是否存在某个键
func (this *StringIntMap) Contains(key string) bool {
    this.mu.RLock()
    _, exists := this.m[key]
    this.mu.RUnlock()
    return exists
}

// 哈希表大小
func (this *StringIntMap) Size() int {
    this.mu.RLock()
    length := len(this.m)
    this.mu.RUnlock()
    return length
}

// 哈希表是否为空
func (this *StringIntMap) IsEmpty() bool {
    this.mu.RLock()
    empty := (len(this.m) == 0)
    this.mu.RUnlock()
    return empty
}

// 清空哈希表
func (this *StringIntMap) Clear() {
    this.mu.Lock()
    this.m = make(map[string]int)
    this.mu.Unlock()
}

// 使用自定义方法执行加锁修改操作
func (this *StringIntMap) LockFunc(f func(m map[string]int)) {
    this.mu.Lock()
    f(this.m)
    this.mu.Unlock()
}

// 使用自定义方法执行加锁读取操作
func (this *StringIntMap) RLockFunc(f func(m map[string]int)) {
    this.mu.RLock()
    f(this.m)
    this.mu.RUnlock()
}

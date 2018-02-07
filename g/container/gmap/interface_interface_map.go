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

type InterfaceInterfaceMap struct {
	mu sync.RWMutex
	m  map[interface{}]interface{}
}

func NewInterfaceInterfaceMap() *InterfaceInterfaceMap {
	return &InterfaceInterfaceMap{
		m: make(map[interface{}]interface{}),
	}
}

// 给定回调函数对原始内容进行遍历
func (this *InterfaceInterfaceMap) Iterator(f func (k interface{}, v interface{})) {
	this.mu.RLock()
	for k, v := range this.m {
		f(k, v)
	}
	this.mu.RUnlock()
}

// 哈希表克隆
func (this *InterfaceInterfaceMap) Clone() *map[interface{}]interface{} {
    m := make(map[interface{}]interface{})
    this.mu.RLock()
    for k, v := range this.m {
        m[k] = v
    }
    this.mu.RUnlock()
    return &m
}

// 设置键值对
func (this *InterfaceInterfaceMap) Set(key interface{}, val interface{}) {
	this.mu.Lock()
	this.m[key] = val
	this.mu.Unlock()
}

// 批量设置键值对
func (this *InterfaceInterfaceMap) BatchSet(m map[interface{}]interface{}) {
    this.mu.Lock()
    for k, v := range m {
        this.m[k] = v
    }
    this.mu.Unlock()
}

// 获取键值
func (this *InterfaceInterfaceMap) Get(key interface{}) (interface{}) {
	this.mu.RLock()
	val, _ := this.m[key]
	this.mu.RUnlock()
	return val
}

func (this *InterfaceInterfaceMap) GetBool(key interface{}) bool {
    return gconv.Bool(this.Get(key))
}

func (this *InterfaceInterfaceMap) GetInt(key interface{}) int {
    return gconv.Int(this.Get(key))
}

func (this *InterfaceInterfaceMap) GetUint (key interface{}) uint {
    return gconv.Uint(this.Get(key))
}

func (this *InterfaceInterfaceMap) GetFloat32 (key interface{}) float32 {
    return gconv.Float32(this.Get(key))
}

func (this *InterfaceInterfaceMap) GetFloat64 (key interface{}) float64 {
    return gconv.Float64(this.Get(key))
}

func (this *InterfaceInterfaceMap) GetString (key interface{}) string {
    return gconv.String(this.Get(key))
}

// 删除键值对
func (this *InterfaceInterfaceMap) Remove(key interface{}) {
	this.mu.Lock()
	delete(this.m, key)
	this.mu.Unlock()
}

// 批量删除键值对
func (this *InterfaceInterfaceMap) BatchRemove(keys []interface{}) {
    this.mu.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.mu.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *InterfaceInterfaceMap) GetAndRemove(key interface{}) (interface{}) {
	this.mu.Lock()
	val, exists := this.m[key]
	if exists {
		delete(this.m, key)
	}
	this.mu.Unlock()
	return val
}

// 返回键列表
func (this *InterfaceInterfaceMap) Keys() []interface{} {
	this.mu.RLock()
	keys := make([]interface{}, 0)
	for key, _ := range this.m {
		keys = append(keys, key)
	}
    this.mu.RUnlock()
	return keys
}

// 返回值列表(注意是随机排序)
func (this *InterfaceInterfaceMap) Values() []interface{} {
	this.mu.RLock()
	vals := make([]interface{}, 0)
	for _, val := range this.m {
		vals = append(vals, val)
	}
	this.mu.RUnlock()
	return vals
}

// 是否存在某个键
func (this *InterfaceInterfaceMap) Contains(key interface{}) bool {
	this.mu.RLock()
	_, exists := this.m[key]
	this.mu.RUnlock()
	return exists
}

// 哈希表大小
func (this *InterfaceInterfaceMap) Size() int {
	this.mu.RLock()
	len := len(this.m)
	this.mu.RUnlock()
	return len
}

// 哈希表是否为空
func (this *InterfaceInterfaceMap) IsEmpty() bool {
	this.mu.RLock()
	empty := (len(this.m) == 0)
	this.mu.RUnlock()
	return empty
}

// 清空哈希表
func (this *InterfaceInterfaceMap) Clear() {
    this.mu.Lock()
    this.m = make(map[interface{}]interface{})
    this.mu.Unlock()
}

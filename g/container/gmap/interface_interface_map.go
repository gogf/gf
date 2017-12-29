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

type InterfaceInterfaceMap struct {
	sync.RWMutex
	m map[interface{}]interface{}
}

func NewInterfaceInterfaceMap() *InterfaceInterfaceMap {
	return &InterfaceInterfaceMap{
		m: make(map[interface{}]interface{}),
	}
}

// 哈希表克隆
func (this *InterfaceInterfaceMap) Clone() *map[interface{}]interface{} {
    m := make(map[interface{}]interface{})
    this.RLock()
    for k, v := range this.m {
        m[k] = v
    }
    this.RUnlock()
    return &m
}

// 设置键值对
func (this *InterfaceInterfaceMap) Set(key interface{}, val interface{}) {
	this.Lock()
	this.m[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *InterfaceInterfaceMap) BatchSet(m map[interface{}]interface{}) {
    this.Lock()
    for k, v := range m {
        this.m[k] = v
    }
    this.Unlock()
}

// 获取键值
func (this *InterfaceInterfaceMap) Get(key interface{}) (interface{}) {
	this.RLock()
	val, _ := this.m[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *InterfaceInterfaceMap) Remove(key interface{}) {
	this.Lock()
	delete(this.m, key)
	this.Unlock()
}

// 批量删除键值对
func (this *InterfaceInterfaceMap) BatchRemove(keys []interface{}) {
    this.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *InterfaceInterfaceMap) GetAndRemove(key interface{}) (interface{}) {
	this.Lock()
	val, exists := this.m[key]
	if exists {
		delete(this.m, key)
	}
	this.Unlock()
	return val
}

// 返回键列表
func (this *InterfaceInterfaceMap) Keys() []interface{} {
	this.RLock()
	keys := make([]interface{}, 0)
	for key, _ := range this.m {
		keys = append(keys, key)
	}
    this.RUnlock()
	return keys
}

// 返回值列表(注意是随机排序)
func (this *InterfaceInterfaceMap) Values() []interface{} {
	this.RLock()
	vals := make([]interface{}, 0)
	for _, val := range this.m {
		vals = append(vals, val)
	}
	this.RUnlock()
	return vals
}

// 是否存在某个键
func (this *InterfaceInterfaceMap) Contains(key interface{}) bool {
	this.RLock()
	_, exists := this.m[key]
	this.RUnlock()
	return exists
}

// 哈希表大小
func (this *InterfaceInterfaceMap) Size() int {
	this.RLock()
	len := len(this.m)
	this.RUnlock()
	return len
}

// 哈希表是否为空
func (this *InterfaceInterfaceMap) IsEmpty() bool {
	this.RLock()
	empty := (len(this.m) == 0)
	this.RUnlock()
	return empty
}

// 清空哈希表
func (this *InterfaceInterfaceMap) Clear() {
    this.Lock()
    this.m = make(map[interface{}]interface{})
    this.Unlock()
}

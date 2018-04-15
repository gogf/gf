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

type UintInterfaceMap struct {
	mu sync.RWMutex
	m  map[uint]interface{}
}

func NewUintInterfaceMap() *UintInterfaceMap {
	return &UintInterfaceMap{
        m: make(map[uint]interface{}),
    }
}

func (this *UintInterfaceMap) Iterator(f func (k uint, v interface{}) bool) {
    this.mu.RLock()
    for k, v := range this.m {
        if !f(k, v) {
            break
        }
    }
    this.mu.RUnlock()
}

// 哈希表克隆
func (this *UintInterfaceMap) Clone() *map[uint]interface{} {
	m := make(map[uint]interface{})
	this.mu.RLock()
	for k, v := range this.m {
		m[k] = v
	}
    this.mu.RUnlock()
	return &m
}

// 设置键值对
func (this *UintInterfaceMap) Set(key uint, val interface{}) {
	this.mu.Lock()
	this.m[key] = val
	this.mu.Unlock()
}

// 批量设置键值对
func (this *UintInterfaceMap) BatchSet(m map[uint]interface{}) {
	this.mu.Lock()
	for k, v := range m {
		this.m[k] = v
	}
	this.mu.Unlock()
}

// 获取键值
func (this *UintInterfaceMap) Get(key uint) interface{} {
	this.mu.RLock()
	val, _ := this.m[key]
	this.mu.RUnlock()
	return val
}

// 获取键值，如果键值不存在则写入默认值
func (this *UintInterfaceMap) GetWithDefault(key uint, value interface{}) interface{} {
    this.mu.Lock()
    val, ok := this.m[key]
    if !ok {
        this.m[key] = value
        val         = value
    }
    this.mu.Unlock()
    return val
}

func (this *UintInterfaceMap) GetBool(key uint) bool {
    return gconv.Bool(this.Get(key))
}

func (this *UintInterfaceMap) GetInt(key uint) int {
    return gconv.Int(this.Get(key))
}

func (this *UintInterfaceMap) GetUint (key uint) uint {
    return gconv.Uint(this.Get(key))
}

func (this *UintInterfaceMap) GetFloat32 (key uint) float32 {
    return gconv.Float32(this.Get(key))
}

func (this *UintInterfaceMap) GetFloat64 (key uint) float64 {
    return gconv.Float64(this.Get(key))
}

func (this *UintInterfaceMap) GetString (key uint) string {
    return gconv.String(this.Get(key))
}

// 删除键值对
func (this *UintInterfaceMap) Remove(key uint) {
    this.mu.Lock()
    delete(this.m, key)
    this.mu.Unlock()
}

// 批量删除键值对
func (this *UintInterfaceMap) BatchRemove(keys []uint) {
    this.mu.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.mu.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *UintInterfaceMap) GetAndRemove(key uint) (interface{}) {
    this.mu.Lock()
    val, exists := this.m[key]
    if exists {
        delete(this.m, key)
    }
    this.mu.Unlock()
    return val
}

// 返回键列表
func (this *UintInterfaceMap) Keys() []uint {
    this.mu.RLock()
    keys := make([]uint, 0)
    for key, _ := range this.m {
        keys = append(keys, key)
    }
    this.mu.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *UintInterfaceMap) Values() []interface{} {
    this.mu.RLock()
    vals := make([]interface{}, 0)
    for _, val := range this.m {
        vals = append(vals, val)
    }
    this.mu.RUnlock()
    return vals
}

// 是否存在某个键
func (this *UintInterfaceMap) Contains(key uint) bool {
    this.mu.RLock()
    _, exists := this.m[key]
    this.mu.RUnlock()
    return exists
}

// 哈希表大小
func (this *UintInterfaceMap) Size() int {
    this.mu.RLock()
    length := len(this.m)
    this.mu.RUnlock()
    return length
}

// 哈希表是否为空
func (this *UintInterfaceMap) IsEmpty() bool {
    this.mu.RLock()
    empty := (len(this.m) == 0)
    this.mu.RUnlock()
    return empty
}

// 清空哈希表
func (this *UintInterfaceMap) Clear() {
    this.mu.Lock()
    this.m = make(map[uint]interface{})
    this.mu.Unlock()
}

// 使用自定义方法执行加锁修改操作
func (this *UintInterfaceMap) LockFunc(f func(m map[uint]interface{})) {
    this.mu.Lock()
    f(this.m)
    this.mu.Unlock()
}

// 使用自定义方法执行加锁读取操作
func (this *UintInterfaceMap) RLockFunc(f func(m map[uint]interface{})) {
    this.mu.RLock()
    f(this.m)
    this.mu.RUnlock()
}

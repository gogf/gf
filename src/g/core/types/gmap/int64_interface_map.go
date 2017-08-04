package gmap

import (
	"sync"
)

type Int64InterfaceMap struct {
	m sync.RWMutex
	M map[int64]interface{}
}

func NewInt64InterfaceMap() *Int64InterfaceMap {
	return &Int64InterfaceMap{
        M: make(map[int64]interface{}),
    }
}

// 哈希表克隆
func (this *Int64InterfaceMap) Clone() *map[int64]interface{} {
	m := make(map[int64]interface{})
	this.m.RLock()
	for k, v := range this.M {
		m[k] = v
	}
    this.m.RUnlock()
	return &m
}

// 设置键值对
func (this *Int64InterfaceMap) Set(key int64, val interface{}) {
	this.m.Lock()
	this.M[key] = val
	this.m.Unlock()
}

// 批量设置键值对
func (this *Int64InterfaceMap) BatchSet(m map[int64]interface{}) {
	todo := make(map[int64]interface{})
	this.m.RLock()
	for k, v := range m {
		old, exists := this.M[k]
		if exists && v == old {
			continue
		}
		todo[k] = v
	}
	this.m.RUnlock()

	if len(todo) == 0 {
		return
	}

	this.m.Lock()
	for k, v := range todo {
		this.M[k] = v
	}
	this.m.Unlock()
}

// 获取键值
func (this *Int64InterfaceMap) Get(key int64) (interface{}) {
	this.m.RLock()
	val, _ := this.M[key]
	this.m.RUnlock()
	return val
}

// 删除键值对
func (this *Int64InterfaceMap) Remove(key int64) {
    this.m.Lock()
    delete(this.M, key)
    this.m.Unlock()
}

// 批量删除键值对
func (this *Int64InterfaceMap) BatchRemove(keys []int64) {
    this.m.Lock()
    for _, key := range keys {
        delete(this.M, key)
    }
    this.m.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *Int64InterfaceMap) GetAndRemove(key int64) (interface{}) {
    this.m.Lock()
    val, exists := this.M[key]
    if exists {
        delete(this.M, key)
    }
    this.m.Unlock()
    return val
}

// 返回键列表
func (this *Int64InterfaceMap) Keys() []int64 {
    this.m.RLock()
    keys := make([]int64, 0)
    for key, _ := range this.M {
        keys = append(keys, key)
    }
    this.m.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *Int64InterfaceMap) Values() []interface{} {
    this.m.RLock()
    vals := make([]interface{}, 0)
    for _, val := range this.M {
        vals = append(vals, val)
    }
    this.m.RUnlock()
    return vals
}

// 是否存在某个键
func (this *Int64InterfaceMap) Contains(key int64) bool {
    this.m.RLock()
    _, exists := this.M[key]
    this.m.RUnlock()
    return exists
}

// 哈希表大小
func (this *Int64InterfaceMap) Size() int {
    this.m.RLock()
    len := len(this.M)
    this.m.RUnlock()
    return len
}

// 哈希表是否为空
func (this *Int64InterfaceMap) IsEmpty() bool {
    this.m.RLock()
    empty := (len(this.M) == 0)
    this.m.RUnlock()
    return empty
}

// 清空哈希表
func (this *Int64InterfaceMap) Clear() {
    this.m.Lock()
    this.M = make(map[int64]interface{})
    this.m.Unlock()
}


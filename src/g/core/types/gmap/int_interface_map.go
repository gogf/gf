package gmap

import (
	"sync"
)

type IntInterfaceMap struct {
	m sync.RWMutex
	M map[int]interface{}
}

func NewIntInterfaceMap() *IntInterfaceMap {
	return &IntInterfaceMap{
        M: make(map[int]interface{}),
    }
}

// 哈希表克隆
func (this *IntInterfaceMap) Clone() *map[int]interface{} {
	m := make(map[int]interface{})
	this.m.RLock()
	for k, v := range this.M {
		m[k] = v
	}
    this.m.RUnlock()
	return &m
}

// 设置键值对
func (this *IntInterfaceMap) Set(key int, val interface{}) {
	this.m.Lock()
	this.M[key] = val
	this.m.Unlock()
}

// 批量设置键值对
func (this *IntInterfaceMap) BatchSet(m map[int]interface{}) {
	todo := make(map[int]interface{})
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
func (this *IntInterfaceMap) Get(key int) (interface{}, bool) {
	this.m.RLock()
	val, exists := this.M[key]
	this.m.RUnlock()
	return val, exists
}

// 删除键值对
func (this *IntInterfaceMap) Remove(key int) {
    this.m.Lock()
    delete(this.M, key)
    this.m.Unlock()
}

// 批量删除键值对
func (this *IntInterfaceMap) BatchRemove(keys []int) {
    this.m.Lock()
    for _, key := range keys {
        delete(this.M, key)
    }
    this.m.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *IntInterfaceMap) GetAndRemove(key int) (interface{}, bool) {
    this.m.Lock()
    val, exists := this.M[key]
    if exists {
        delete(this.M, key)
    }
    this.m.Unlock()
    return val, exists
}

// 返回键列表
func (this *IntInterfaceMap) Keys() []int {
    this.m.RLock()
    keys := make([]int, 0)
    for key, _ := range this.M {
        keys = append(keys, key)
    }
    this.m.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *IntInterfaceMap) Values() []interface{} {
    this.m.RLock()
    vals := make([]interface{}, 0)
    for _, val := range this.M {
        vals = append(vals, val)
    }
    this.m.RUnlock()
    return vals
}

// 是否存在某个键
func (this *IntInterfaceMap) Contains(key int) bool {
	_, exists := this.Get(key)
	return exists
}

// 哈希表大小
func (this *IntInterfaceMap) Size() int {
    this.m.RLock()
    len := len(this.M)
    this.m.RUnlock()
    return len
}

// 哈希表是否为空
func (this *IntInterfaceMap) IsEmpty() bool {
    this.m.RLock()
    empty := (len(this.M) == 0)
    this.m.RUnlock()
    return empty
}

// 清空哈希表
func (this *IntInterfaceMap) Clear() {
    this.m.Lock()
    this.M = make(map[int]interface{})
    this.m.Unlock()
}


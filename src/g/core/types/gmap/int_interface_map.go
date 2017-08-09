package gmap

import (
	"sync"
)

type IntInterfaceMap struct {
	sync.RWMutex
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
	this.RLock()
	for k, v := range this.M {
		m[k] = v
	}
    this.RUnlock()
	return &m
}

// 设置键值对
func (this *IntInterfaceMap) Set(key int, val interface{}) {
	this.Lock()
	this.M[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *IntInterfaceMap) BatchSet(m map[int]interface{}) {
	todo := make(map[int]interface{})
	this.RLock()
	for k, v := range m {
		old, exists := this.M[k]
		if exists && v == old {
			continue
		}
		todo[k] = v
	}
	this.RUnlock()

	if len(todo) == 0 {
		return
	}

	this.Lock()
	for k, v := range todo {
		this.M[k] = v
	}
	this.Unlock()
}

// 获取键值
func (this *IntInterfaceMap) Get(key int) (interface{}) {
	this.RLock()
	val, _ := this.M[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *IntInterfaceMap) Remove(key int) {
    this.Lock()
    delete(this.M, key)
    this.Unlock()
}

// 批量删除键值对
func (this *IntInterfaceMap) BatchRemove(keys []int) {
    this.Lock()
    for _, key := range keys {
        delete(this.M, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *IntInterfaceMap) GetAndRemove(key int) (interface{}) {
    this.Lock()
    val, exists := this.M[key]
    if exists {
        delete(this.M, key)
    }
    this.Unlock()
    return val
}

// 返回键列表
func (this *IntInterfaceMap) Keys() []int {
    this.RLock()
    keys := make([]int, 0)
    for key, _ := range this.M {
        keys = append(keys, key)
    }
    this.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *IntInterfaceMap) Values() []interface{} {
    this.RLock()
    vals := make([]interface{}, 0)
    for _, val := range this.M {
        vals = append(vals, val)
    }
    this.RUnlock()
    return vals
}

// 是否存在某个键
func (this *IntInterfaceMap) Contains(key int) bool {
    this.RLock()
    _, exists := this.M[key]
    this.RUnlock()
    return exists
}

// 哈希表大小
func (this *IntInterfaceMap) Size() int {
    this.RLock()
    len := len(this.M)
    this.RUnlock()
    return len
}

// 哈希表是否为空
func (this *IntInterfaceMap) IsEmpty() bool {
    this.RLock()
    empty := (len(this.M) == 0)
    this.RUnlock()
    return empty
}

// 清空哈希表
func (this *IntInterfaceMap) Clear() {
    this.Lock()
    this.M = make(map[int]interface{})
    this.Unlock()
}


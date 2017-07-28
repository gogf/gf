package gmap

import (
	"sync"
)

type IntStringMap struct {
	m sync.RWMutex
	M map[int]string
}

func NewIntStringMap() *IntStringMap {
	return &IntStringMap{
        M: make(map[int]string),
    }
}

// 哈希表克隆
func (this *IntStringMap) Clone() *map[int]string {
	m := make(map[int]string)
	this.m.RLock()
	for k, v := range this.M {
		m[k] = v
	}
    this.m.RUnlock()
	return &m
}

// 设置键值对
func (this *IntStringMap) Set(key int, val string) {
	this.m.Lock()
	this.M[key] = val
	this.m.Unlock()
}

// 批量设置键值对
func (this *IntStringMap) Sets(m map[int]string) {
	todo := make(map[int]string)
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
func (this *IntStringMap) Get(key int) (string, bool) {
	this.m.RLock()
	val, exists := this.M[key]
	this.m.RUnlock()
	return val, exists
}

// 删除键值对
func (this *IntStringMap) Remove(key int) {
    this.m.Lock()
    delete(this.M, key)
    this.m.Unlock()
}

// 批量删除键值对
func (this *IntStringMap) BatchRemove(keys []int) {
    this.m.Lock()
    for _, key := range keys {
        delete(this.M, key)
    }
    this.m.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *IntStringMap) GetAndRemove(key int) (string, bool) {
    this.m.Lock()
    val, exists := this.M[key]
    if exists {
        delete(this.M, key)
    }
    this.m.Unlock()
    return val, exists
}

// 返回键列表
func (this *IntStringMap) Keys() []int {
    this.m.RLock()
    keys := make([]int, 0)
    for key, _ := range this.M {
        keys = append(keys, key)
    }
    this.m.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *IntStringMap) Values() []string {
    this.m.RLock()
    vals := make([]string, 0)
    for _, val := range this.M {
        vals = append(vals, val)
    }
    this.m.RUnlock()
    return vals
}

// 是否存在某个键
func (this *IntStringMap) Contains(key int) bool {
	_, exists := this.Get(key)
	return exists
}

// 哈希表大小
func (this *IntStringMap) Size() int {
    this.m.RLock()
    len := len(this.M)
    this.m.RUnlock()
    return len
}

// 哈希表是否为空
func (this *IntStringMap) IsEmpty() bool {
    this.m.RLock()
    empty := (len(this.M) == 0)
    this.m.RUnlock()
    return empty
}

// 清空哈希表
func (this *IntStringMap) Clear() {
    this.m.Lock()
    this.M = make(map[int]string)
    this.m.Unlock()
}


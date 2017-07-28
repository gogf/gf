package gmap

import (
	"sync"
)

type IntInt64Map struct {
	m sync.RWMutex
	M map[int]int64
}

func NewIntInt64Map() *IntInt64Map {
	return &IntInt64Map{
        M: make(map[int]int64),
    }
}

// 哈希表克隆
func (this *IntInt64Map) Clone() *map[int]int64 {
	m := make(map[int]int64)
	this.m.RLock()
	for k, v := range this.M {
		m[k] = v
	}
    this.m.RUnlock()
	return &m
}

// 设置键值对
func (this *IntInt64Map) Set(key int, val int64) {
	this.m.Lock()
	this.M[key] = val
	this.m.Unlock()
}

// 批量设置键值对
func (this *IntInt64Map) BatchSet(m map[int]int64) {
	todo := make(map[int]int64)
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
func (this *IntInt64Map) Get(key int) (int64, bool) {
	this.m.RLock()
	val, exists := this.M[key]
	this.m.RUnlock()
	return val, exists
}

// 删除键值对
func (this *IntInt64Map) Remove(key int) {
    this.m.Lock()
    delete(this.M, key)
    this.m.Unlock()
}

// 批量删除键值对
func (this *IntInt64Map) BatchRemove(keys []int) {
    this.m.Lock()
    for _, key := range keys {
        delete(this.M, key)
    }
    this.m.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *IntInt64Map) GetAndRemove(key int) (int64, bool) {
    this.m.Lock()
    val, exists := this.M[key]
    if exists {
        delete(this.M, key)
    }
    this.m.Unlock()
    return val, exists
}

// 返回键列表
func (this *IntInt64Map) Keys() []int {
    this.m.RLock()
    keys := make([]int, 0)
    for key, _ := range this.M {
        keys = append(keys, key)
    }
    this.m.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *IntInt64Map) Values() []int64 {
    this.m.RLock()
    vals := make([]int64, 0)
    for _, val := range this.M {
        vals = append(vals, val)
    }
    this.m.RUnlock()
    return vals
}

// 是否存在某个键
func (this *IntInt64Map) Contains(key int) bool {
	_, exists := this.Get(key)
	return exists
}

// 哈希表大小
func (this *IntInt64Map) Size() int {
    this.m.RLock()
    len := len(this.M)
    this.m.RUnlock()
    return len
}

// 哈希表是否为空
func (this *IntInt64Map) IsEmpty() bool {
    this.m.RLock()
    empty := (len(this.M) == 0)
    this.m.RUnlock()
    return empty
}

// 清空哈希表
func (this *IntInt64Map) Clear() {
    this.m.Lock()
    this.M = make(map[int]int64)
    this.m.Unlock()
}


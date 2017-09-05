package gmap

import (
	"sync"
)

type IntIntMap struct {
	sync.RWMutex
	m map[int]int
}

func NewIntIntMap() *IntIntMap {
	return &IntIntMap{
        m: make(map[int]int),
    }
}

// 哈希表克隆
func (this *IntIntMap) Clone() *map[int]int {
	m := make(map[int]int)
	this.RLock()
	for k, v := range this.m {
		m[k] = v
	}
    this.RUnlock()
	return &m
}

// 设置键值对
func (this *IntIntMap) Set(key int, val int) {
	this.Lock()
	this.m[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *IntIntMap) BatchSet(m map[int]int) {
	this.Lock()
	for k, v := range m {
		this.m[k] = v
	}
	this.Unlock()
}

// 获取键值
func (this *IntIntMap) Get(key int) (int) {
	this.RLock()
	val, _ := this.m[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *IntIntMap) Remove(key int) {
    this.Lock()
    delete(this.m, key)
    this.Unlock()
}

// 批量删除键值对
func (this *IntIntMap) BatchRemove(keys []int) {
    this.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *IntIntMap) GetAndRemove(key int) (int) {
    this.Lock()
    val, exists := this.m[key]
    if exists {
        delete(this.m, key)
    }
    this.Unlock()
    return val
}

// 返回键列表
func (this *IntIntMap) Keys() []int {
    this.RLock()
    keys := make([]int, 0)
    for key, _ := range this.m {
        keys = append(keys, key)
    }
    this.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *IntIntMap) Values() []int {
    this.RLock()
    vals := make([]int, 0)
    for _, val := range this.m {
        vals = append(vals, val)
    }
    this.RUnlock()
    return vals
}

// 是否存在某个键
func (this *IntIntMap) Contains(key int) bool {
    this.RLock()
    _, exists := this.m[key]
    this.RUnlock()
    return exists
}

// 哈希表大小
func (this *IntIntMap) Size() int {
    this.RLock()
    len := len(this.m)
    this.RUnlock()
    return len
}

// 哈希表是否为空
func (this *IntIntMap) IsEmpty() bool {
    this.RLock()
    empty := (len(this.m) == 0)
    this.RUnlock()
    return empty
}

// 清空哈希表
func (this *IntIntMap) Clear() {
    this.Lock()
    this.m = make(map[int]int)
    this.Unlock()
}


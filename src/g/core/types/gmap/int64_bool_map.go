package gmap

import (
	"sync"
)

type Int64BoolMap struct {
	sync.RWMutex
	m map[int64]bool
}

func NewInt64BoolMap() *Int64BoolMap {
	return &Int64BoolMap{
        m: make(map[int64]bool),
    }
}

// 哈希表克隆
func (this *Int64BoolMap) Clone() *map[int64]bool {
	m := make(map[int64]bool)
	this.RLock()
	for k, v := range this.m {
		m[k] = v
	}
    this.RUnlock()
	return &m
}

// 设置键值对
func (this *Int64BoolMap) Set(key int64, val bool) {
	this.Lock()
	this.m[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *Int64BoolMap) BatchSet(m map[int64]bool) {
	this.Lock()
	for k, v := range m {
		this.m[k] = v
	}
	this.Unlock()
}

// 获取键值
func (this *Int64BoolMap) Get(key int64) (bool) {
	this.RLock()
	val, _ := this.m[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *Int64BoolMap) Remove(key int64) {
    this.Lock()
    delete(this.m, key)
    this.Unlock()
}

// 批量删除键值对
func (this *Int64BoolMap) BatchRemove(keys []int64) {
    this.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *Int64BoolMap) GetAndRemove(key int64) (bool) {
    this.Lock()
    val, exists := this.m[key]
    if exists {
        delete(this.m, key)
    }
    this.Unlock()
    return val
}

// 返回键列表
func (this *Int64BoolMap) Keys() []int64 {
    this.RLock()
    keys := make([]int64, 0)
    for key, _ := range this.m {
        keys = append(keys, key)
    }
    this.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
//func (this *Int64BoolMap) Values() []bool {
//    this.RLock()
//    vals := make([]bool, 0)
//    for _, val := range this.m {
//        vals = append(vals, val)
//    }
//    this.RUnlock()
//    return vals
//}

// 是否存在某个键
func (this *Int64BoolMap) Contains(key int64) bool {
	this.RLock()
	_, exists := this.m[key]
	this.RUnlock()
	return exists
}

// 哈希表大小
func (this *Int64BoolMap) Size() int {
    this.RLock()
    len := len(this.m)
    this.RUnlock()
    return len
}

// 哈希表是否为空
func (this *Int64BoolMap) IsEmpty() bool {
    this.RLock()
    empty := (len(this.m) == 0)
    this.RUnlock()
    return empty
}

// 清空哈希表
func (this *Int64BoolMap) Clear() {
    this.Lock()
    this.m = make(map[int64]bool)
    this.Unlock()
}


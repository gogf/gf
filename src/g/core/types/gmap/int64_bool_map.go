package gmap

import (
	"sync"
)

type Int64BoolMap struct {
	m sync.RWMutex
	M map[int64]bool
}

func NewInt64BoolMap() *Int64BoolMap {
	return &Int64BoolMap{
        M: make(map[int64]bool),
    }
}

// 哈希表克隆
func (this *Int64BoolMap) Clone() *map[int64]bool {
	m := make(map[int64]bool)
	this.m.RLock()
	for k, v := range this.M {
		m[k] = v
	}
    this.m.RUnlock()
	return &m
}

// 设置键值对
func (this *Int64BoolMap) Set(key int64, val bool) {
	this.m.Lock()
	this.M[key] = val
	this.m.Unlock()
}

// 批量设置键值对
func (this *Int64BoolMap) BatchSet(m map[int64]bool) {
	todo := make(map[int64]bool)
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
func (this *Int64BoolMap) Get(key int64) (bool) {
	this.m.RLock()
	val, _ := this.M[key]
	this.m.RUnlock()
	return val
}

// 删除键值对
func (this *Int64BoolMap) Remove(key int64) {
    this.m.Lock()
    delete(this.M, key)
    this.m.Unlock()
}

// 批量删除键值对
func (this *Int64BoolMap) BatchRemove(keys []int64) {
    this.m.Lock()
    for _, key := range keys {
        delete(this.M, key)
    }
    this.m.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *Int64BoolMap) GetAndRemove(key int64) (bool) {
    this.m.Lock()
    val, exists := this.M[key]
    if exists {
        delete(this.M, key)
    }
    this.m.Unlock()
    return val
}

// 返回键列表
func (this *Int64BoolMap) Keys() []int64 {
    this.m.RLock()
    keys := make([]int64, 0)
    for key, _ := range this.M {
        keys = append(keys, key)
    }
    this.m.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
//func (this *Int64BoolMap) Values() []bool {
//    this.m.RLock()
//    vals := make([]bool, 0)
//    for _, val := range this.M {
//        vals = append(vals, val)
//    }
//    this.m.RUnlock()
//    return vals
//}

// 是否存在某个键
func (this *Int64BoolMap) Contains(key int64) bool {
	this.m.RLock()
	_, exists := this.M[key]
	this.m.RUnlock()
	return exists
}

// 哈希表大小
func (this *Int64BoolMap) Size() int {
    this.m.RLock()
    len := len(this.M)
    this.m.RUnlock()
    return len
}

// 哈希表是否为空
func (this *Int64BoolMap) IsEmpty() bool {
    this.m.RLock()
    empty := (len(this.M) == 0)
    this.m.RUnlock()
    return empty
}

// 清空哈希表
func (this *Int64BoolMap) Clear() {
    this.m.Lock()
    this.M = make(map[int64]bool)
    this.m.Unlock()
}


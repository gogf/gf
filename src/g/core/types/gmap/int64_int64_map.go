package gmap

import (
	"sync"
)

type Int64Int64Map struct {
	sync.RWMutex
	M map[int64]int64
}

func NewInt64Int64Map() *Int64Int64Map {
	return &Int64Int64Map{
        M: make(map[int64]int64),
    }
}

// 哈希表克隆
func (this *Int64Int64Map) Clone() *map[int64]int64 {
	m := make(map[int64]int64)
	this.RLock()
	for k, v := range this.M {
		m[k] = v
	}
    this.RUnlock()
	return &m
}

// 设置键值对
func (this *Int64Int64Map) Set(key int64, val int64) {
	this.Lock()
	this.M[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *Int64Int64Map) BatchSet(m map[int64]int64) {
	todo := make(map[int64]int64)
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
func (this *Int64Int64Map) Get(key int64) (int64) {
	this.RLock()
	val, _ := this.M[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *Int64Int64Map) Remove(key int64) {
    this.Lock()
    delete(this.M, key)
    this.Unlock()
}

// 批量删除键值对
func (this *Int64Int64Map) BatchRemove(keys []int64) {
    this.Lock()
    for _, key := range keys {
        delete(this.M, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *Int64Int64Map) GetAndRemove(key int64) (int64) {
    this.Lock()
    val, exists := this.M[key]
    if exists {
        delete(this.M, key)
    }
    this.Unlock()
    return val
}

// 返回键列表
func (this *Int64Int64Map) Keys() []int64 {
    this.RLock()
    keys := make([]int64, 0)
    for key, _ := range this.M {
        keys = append(keys, key)
    }
    this.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *Int64Int64Map) Values() []int64 {
    this.RLock()
    vals := make([]int64, 0)
    for _, val := range this.M {
        vals = append(vals, val)
    }
    this.RUnlock()
    return vals
}

// 是否存在某个键
func (this *Int64Int64Map) Contains(key int64) bool {
    this.RLock()
    _, exists := this.M[key]
    this.RUnlock()
    return exists
}

// 哈希表大小
func (this *Int64Int64Map) Size() int {
    this.RLock()
    len := len(this.M)
    this.RUnlock()
    return len
}

// 哈希表是否为空
func (this *Int64Int64Map) IsEmpty() bool {
    this.RLock()
    empty := (len(this.M) == 0)
    this.RUnlock()
    return empty
}

// 清空哈希表
func (this *Int64Int64Map) Clear() {
    this.Lock()
    this.M = make(map[int64]int64)
    this.Unlock()
}


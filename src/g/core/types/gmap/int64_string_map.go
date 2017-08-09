package gmap

import (
	"sync"
)

type Int64StringMap struct {
	sync.RWMutex
	M map[int64]string
}

func NewInt64StringMap() *Int64StringMap {
	return &Int64StringMap{
        M: make(map[int64]string),
    }
}

// 哈希表克隆
func (this *Int64StringMap) Clone() *map[int64]string {
	m := make(map[int64]string)
	this.RLock()
	for k, v := range this.M {
		m[k] = v
	}
    this.RUnlock()
	return &m
}

// 设置键值对
func (this *Int64StringMap) Set(key int64, val string) {
	this.Lock()
	this.M[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *Int64StringMap) BatchSet(m map[int64]string) {
	todo := make(map[int64]string)
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
func (this *Int64StringMap) Get(key int64) (string) {
	this.RLock()
	val, _ := this.M[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *Int64StringMap) Remove(key int64) {
    this.Lock()
    delete(this.M, key)
    this.Unlock()
}

// 批量删除键值对
func (this *Int64StringMap) BatchRemove(keys []int64) {
    this.Lock()
    for _, key := range keys {
        delete(this.M, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *Int64StringMap) GetAndRemove(key int64) (string) {
    this.Lock()
    val, exists := this.M[key]
    if exists {
        delete(this.M, key)
    }
    this.Unlock()
    return val
}

// 返回键列表
func (this *Int64StringMap) Keys() []int64 {
    this.RLock()
    keys := make([]int64, 0)
    for key, _ := range this.M {
        keys = append(keys, key)
    }
    this.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *Int64StringMap) Values() []string {
    this.RLock()
    vals := make([]string, 0)
    for _, val := range this.M {
        vals = append(vals, val)
    }
    this.RUnlock()
    return vals
}

// 是否存在某个键
func (this *Int64StringMap) Contains(key int64) bool {
    this.RLock()
    _, exists := this.M[key]
    this.RUnlock()
    return exists
}

// 哈希表大小
func (this *Int64StringMap) Size() int {
    this.RLock()
    len := len(this.M)
    this.RUnlock()
    return len
}

// 哈希表是否为空
func (this *Int64StringMap) IsEmpty() bool {
    this.RLock()
    empty := (len(this.M) == 0)
    this.RUnlock()
    return empty
}

// 清空哈希表
func (this *Int64StringMap) Clear() {
    this.Lock()
    this.M = make(map[int64]string)
    this.Unlock()
}


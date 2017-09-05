package gmap

import (
	"sync"
)

type Int64StringMap struct {
	sync.RWMutex
	m map[int64]string
}

func NewInt64StringMap() *Int64StringMap {
	return &Int64StringMap{
        m: make(map[int64]string),
    }
}

// 哈希表克隆
func (this *Int64StringMap) Clone() *map[int64]string {
	m := make(map[int64]string)
	this.RLock()
	for k, v := range this.m {
		m[k] = v
	}
    this.RUnlock()
	return &m
}

// 设置键值对
func (this *Int64StringMap) Set(key int64, val string) {
	this.Lock()
	this.m[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *Int64StringMap) BatchSet(m map[int64]string) {
	this.Lock()
	for k, v := range m {
		this.m[k] = v
	}
	this.Unlock()
}

// 获取键值
func (this *Int64StringMap) Get(key int64) (string) {
	this.RLock()
	val, _ := this.m[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *Int64StringMap) Remove(key int64) {
    this.Lock()
    delete(this.m, key)
    this.Unlock()
}

// 批量删除键值对
func (this *Int64StringMap) BatchRemove(keys []int64) {
    this.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *Int64StringMap) GetAndRemove(key int64) (string) {
    this.Lock()
    val, exists := this.m[key]
    if exists {
        delete(this.m, key)
    }
    this.Unlock()
    return val
}

// 返回键列表
func (this *Int64StringMap) Keys() []int64 {
    this.RLock()
    keys := make([]int64, 0)
    for key, _ := range this.m {
        keys = append(keys, key)
    }
    this.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *Int64StringMap) Values() []string {
    this.RLock()
    vals := make([]string, 0)
    for _, val := range this.m {
        vals = append(vals, val)
    }
    this.RUnlock()
    return vals
}

// 是否存在某个键
func (this *Int64StringMap) Contains(key int64) bool {
    this.RLock()
    _, exists := this.m[key]
    this.RUnlock()
    return exists
}

// 哈希表大小
func (this *Int64StringMap) Size() int {
    this.RLock()
    len := len(this.m)
    this.RUnlock()
    return len
}

// 哈希表是否为空
func (this *Int64StringMap) IsEmpty() bool {
    this.RLock()
    empty := (len(this.m) == 0)
    this.RUnlock()
    return empty
}

// 清空哈希表
func (this *Int64StringMap) Clear() {
    this.Lock()
    this.m = make(map[int64]string)
    this.Unlock()
}


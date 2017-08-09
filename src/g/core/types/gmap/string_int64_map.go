package gmap

import (
	"sync"
)

type StringInt64Map struct {
	sync.RWMutex
	M map[string]int64
}

func NewStringInt64Map() *StringInt64Map {
	return &StringInt64Map{
        M: make(map[string]int64),
    }
}

// 哈希表克隆
func (this *StringInt64Map) Clone() *map[string]int64 {
    m := make(map[string]int64)
    this.RLock()
    for k, v := range this.M {
        m[k] = v
    }
    this.RUnlock()
    return &m
}

// 设置键值对
func (this *StringInt64Map) Set(key string, val int64) {
	this.Lock()
	this.M[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *StringInt64Map) BatchSet(m map[string]int64) {
	todo := make(map[string]int64)
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
func (this *StringInt64Map) Get(key string) (int64) {
	this.RLock()
	val, _ := this.M[key]
    this.RUnlock()
	return val
}

// 删除键值对
func (this *StringInt64Map) Remove(key string) {
	this.Lock()
	delete(this.M, key)
    this.Unlock()
}

// 批量删除键值对
func (this *StringInt64Map) BatchRemove(keys []string) {
    this.Lock()
    for _, key := range keys {
        delete(this.M, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *StringInt64Map) GetAndRemove(key string) (int64) {
    this.Lock()
    val, exists := this.M[key]
    if exists {
        delete(this.M, key)
    }
    this.Unlock()
    return val
}

// 返回键列表
func (this *StringInt64Map) Keys() []string {
    this.RLock()
    keys := make([]string, 0)
    for key, _ := range this.M {
        keys = append(keys, key)
    }
    this.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *StringInt64Map) Values() []int64 {
    this.RLock()
    vals := make([]int64, 0)
    for _, val := range this.M {
        vals = append(vals, val)
    }
    this.RUnlock()
    return vals
}

// 是否存在某个键
func (this *StringInt64Map) Contains(key string) bool {
    this.RLock()
    _, exists := this.M[key]
    this.RUnlock()
    return exists
}

// 哈希表大小
func (this *StringInt64Map) Size() int {
    this.RLock()
    len := len(this.M)
    this.RUnlock()
    return len
}

// 哈希表是否为空
func (this *StringInt64Map) IsEmpty() bool {
    this.RLock()
    empty := (len(this.M) == 0)
    this.RUnlock()
    return empty
}

// 清空哈希表
func (this *StringInt64Map) Clear() {
    this.Lock()
    this.M = make(map[string]int64)
    this.Unlock()
}


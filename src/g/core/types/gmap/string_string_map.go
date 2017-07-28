package gmap

import (
	"sync"
)

type StringStringMap struct {
	m sync.RWMutex
	M map[string]string
}

func NewStringStringMap() *StringStringMap {
	return &StringStringMap{
		M: make(map[string]string),
	}
}

// 哈希表克隆
func (this *StringStringMap) Clone() *map[string]string {
    m := make(map[string]string)
    this.m.RLock()
    for k, v := range this.M {
        m[k] = v
    }
    this.m.RUnlock()
    return &m
}

// 设置键值对
func (this *StringStringMap) Set(key string, val string) {
	this.m.Lock()
	this.M[key] = val
	this.m.Unlock()
}

// 批量设置键值对
func (this *StringStringMap) BatchSet(m map[string]string) {
    todo := make(map[string]string)
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
func (this *StringStringMap) Get(key string) (string, bool) {
	this.m.RLock()
	val, exists := this.M[key]
	this.m.RUnlock()
	return val, exists
}

// 删除键值对
func (this *StringStringMap) Remove(key string) {
	this.m.Lock()
	delete(this.M, key)
	this.m.Unlock()
}

// 批量删除键值对
func (this *StringStringMap) BatchRemove(keys []string) {
    this.m.Lock()
    for _, key := range keys {
        delete(this.M, key)
    }
    this.m.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *StringStringMap) GetAndRemove(key string) (string, bool) {
	this.m.Lock()
	val, exists := this.M[key]
	if exists {
		delete(this.M, key)
	}
	this.m.Unlock()
	return val, exists
}

// 返回键列表
func (this *StringStringMap) Keys() []string {
	this.m.RLock()
	keys := make([]string, 0)
	for key, _ := range this.M {
		keys = append(keys, key)
	}
    this.m.RUnlock()
	return keys
}

// 返回值列表(注意是随机排序)
func (this *StringStringMap) Values() []string {
	this.m.RLock()
	vals := make([]string, 0)
	for _, val := range this.M {
		vals = append(vals, val)
	}
	this.m.RUnlock()
	return vals
}

// 是否存在某个键
func (this *StringStringMap) Contains(key string) bool {
	this.m.RLock()
	_, exists := this.M[key]
	this.m.RUnlock()
	return exists
}

// 哈希表大小
func (this *StringStringMap) Size() int {
	this.m.RLock()
	len := len(this.M)
	this.m.RUnlock()
	return len
}

// 哈希表是否为空
func (this *StringStringMap) IsEmpty() bool {
	this.m.RLock()
	empty := (len(this.M) == 0)
	this.m.RUnlock()
	return empty
}

// 清空哈希表
func (this *StringStringMap) Clear() {
    this.m.Lock()
    this.M = make(map[string]string)
    this.m.Unlock()
}

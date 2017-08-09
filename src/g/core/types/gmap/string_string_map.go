package gmap

import (
	"sync"
)

type StringStringMap struct {
	sync.RWMutex
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
    this.RLock()
    for k, v := range this.M {
        m[k] = v
    }
    this.RUnlock()
    return &m
}

// 设置键值对
func (this *StringStringMap) Set(key string, val string) {
	this.Lock()
	this.M[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *StringStringMap) BatchSet(m map[string]string) {
    todo := make(map[string]string)
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
func (this *StringStringMap) Get(key string) string {
	this.RLock()
	val, _ := this.M[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *StringStringMap) Remove(key string) {
	this.Lock()
	delete(this.M, key)
	this.Unlock()
}

// 批量删除键值对
func (this *StringStringMap) BatchRemove(keys []string) {
    this.Lock()
    for _, key := range keys {
        delete(this.M, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *StringStringMap) GetAndRemove(key string) string {
	this.Lock()
	val, exists := this.M[key]
	if exists {
		delete(this.M, key)
	}
	this.Unlock()
	return val
}

// 返回键列表
func (this *StringStringMap) Keys() []string {
	this.RLock()
	keys := make([]string, 0)
	for key, _ := range this.M {
		keys = append(keys, key)
	}
    this.RUnlock()
	return keys
}

// 返回值列表(注意是随机排序)
func (this *StringStringMap) Values() []string {
	this.RLock()
	vals := make([]string, 0)
	for _, val := range this.M {
		vals = append(vals, val)
	}
	this.RUnlock()
	return vals
}

// 是否存在某个键
func (this *StringStringMap) Contains(key string) bool {
	this.RLock()
	_, exists := this.M[key]
	this.RUnlock()
	return exists
}

// 哈希表大小
func (this *StringStringMap) Size() int {
	this.RLock()
	len := len(this.M)
	this.RUnlock()
	return len
}

// 哈希表是否为空
func (this *StringStringMap) IsEmpty() bool {
	this.RLock()
	empty := (len(this.M) == 0)
	this.RUnlock()
	return empty
}

// 清空哈希表
func (this *StringStringMap) Clear() {
    this.Lock()
    this.M = make(map[string]string)
    this.Unlock()
}

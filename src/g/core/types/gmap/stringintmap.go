package gmap

import (
	"sync"
)

type StringIntMap struct {
	m sync.RWMutex
	M map[string]int
}

func NewStringIntMap() *StringIntMap {
	return &StringIntMap{
        M: make(map[string]int),
    }
}

func (this *StringIntMap) Set(key string, val int) {
	this.m.Lock()
	this.M[key] = val
	this.m.Unlock()
}

func (this *StringIntMap) Sets(m map[string]int) {
	todo := make(map[string]int)
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

func (this *StringIntMap) Get(key string) (int, bool) {
	this.m.RLock()
	val, exists := this.M[key]
    this.m.RUnlock()
	return val, exists
}

func (this *StringIntMap) Contains(key string) bool {
	_, exists := this.Get(key)
	return exists
}

func (this *StringIntMap) Remove(key string) {
	this.m.Lock()
	delete(this.M, key)
    this.m.Unlock()
}

// 返回键列表
func (this *StringIntMap) Keys() []string {
    this.m.RLock()
    keys := make([]string, 0)
    for key, _ := range this.M {
        keys = append(keys, key)
    }
    this.m.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *StringIntMap) Values() []int {
    this.m.RLock()
    vals := make([]int, 0)
    for _, val := range this.M {
        vals = append(vals, val)
    }
    this.m.RUnlock()
    return vals
}

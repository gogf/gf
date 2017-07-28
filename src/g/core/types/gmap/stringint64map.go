package gmap

import (
	"sync"
)

type StringInt64Map struct {
	m sync.RWMutex
	M map[string]int64
}

func NewStringInt64Map() *StringInt64Map {
	return &StringInt64Map{
        M: make(map[string]int64),
    }
}

func (this *StringInt64Map) Set(key string, val int64) {
	this.m.Lock()
	this.M[key] = val
	this.m.Unlock()
}

func (this *StringInt64Map) Sets(m map[string]int64) {
	todo := make(map[string]int64)
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

func (this *StringInt64Map) Get(key string) (int64, bool) {
	this.m.RLock()
	val, exists := this.M[key]
    this.m.RUnlock()
	return val, exists
}

func (this *StringInt64Map) Contains(key string) bool {
	_, exists := this.Get(key)
	return exists
}

func (this *StringInt64Map) Remove(key string) {
	this.m.Lock()
	delete(this.M, key)
    this.m.Unlock()
}

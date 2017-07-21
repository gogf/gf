package maps

import (
	"sync"
)

type StringInt64Map struct {
	sync.RWMutex
	M map[string]int64
}

func NewStringInt64Map() *StringInt64Map {
	return &StringInt64Map{M: make(map[string]int64)}
}

func (this *StringInt64Map) Put(key string, val int64) {
	this.Lock()
	defer this.Unlock()
	this.M[key] = val
}

func (this *StringInt64Map) Puts(m map[string]int64) {
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

func (this *StringInt64Map) Get(key string) (int64, bool) {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[key]
	return val, exists
}

func (this *StringInt64Map) Exists(key string) bool {
	_, exists := this.Get(key)
	return exists
}

func (this *StringInt64Map) Remove(key string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, key)
}

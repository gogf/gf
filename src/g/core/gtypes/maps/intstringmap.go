package maps

import (
	"sync"
)

type IntStringMap struct {
	sync.RWMutex
	M map[int]string
}

func NewIntStringMap() *IntStringMap {
	return &IntStringMap{M: make(map[int]string)}
}

func (this *IntStringMap) Clone() map[int]string {
	m := make(map[int]string)
	this.RLock()
	defer this.RUnlock()
	for k, v := range this.M {
		m[k] = v
	}
	return m
}

func (this *IntStringMap) Put(key int, val string) {
	this.Lock()
	defer this.Unlock()
	this.M[key] = val
}

func (this *IntStringMap) Puts(m map[int]string) {
	todo := make(map[int]string)
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

func (this *IntStringMap) Get(key int) (string, bool) {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[key]
	return val, exists
}

func (this *IntStringMap) Exists(key int) bool {
	_, exists := this.Get(key)
	return exists
}

func (this *IntStringMap) Remove(key int) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, key)
}

func (this *IntStringMap) RemoveBatch(keys []int) {
	this.Lock()
	defer this.Unlock()
	for _, key := range keys {
		delete(this.M, key)
	}
}

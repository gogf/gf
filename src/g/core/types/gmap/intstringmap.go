package gmap

import (
	"sync"
)

type IntStringMap struct {
	m sync.RWMutex
	M map[int]string
}

func NewIntStringMap() *IntStringMap {
	return &IntStringMap{
        M: make(map[int]string),
    }
}

func (this *IntStringMap) Clone() *map[int]string {
	m := make(map[int]string)
	this.m.RLock()
	for k, v := range this.M {
		m[k] = v
	}
    this.m.RUnlock()
	return &m
}

func (this *IntStringMap) Set(key int, val string) {
	this.m.Lock()
	this.M[key] = val
	this.m.Unlock()
}

func (this *IntStringMap) Sets(m map[int]string) {
	todo := make(map[int]string)
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

func (this *IntStringMap) Get(key int) (string, bool) {
	this.m.RLock()
	val, exists := this.M[key]
	this.m.RUnlock()
	return val, exists
}

func (this *IntStringMap) Contains(key int) bool {
	_, exists := this.Get(key)
	return exists
}

func (this *IntStringMap) Remove(key int) {
	this.m.Lock()
	delete(this.M, key)
	this.m.Unlock()
}

func (this *IntStringMap) RemoveBatch(keys []int) {
	this.m.Lock()
	for _, key := range keys {
		delete(this.M, key)
	}
	this.m.Unlock()
}

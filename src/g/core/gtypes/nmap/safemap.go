package nmap

import (
	"sync"
)

type SafeMap struct {
	sync.RWMutex
	M map[string]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		M: make(map[string]interface{}),
	}
}

func (this *SafeMap) Put(key string, val interface{}) {
	this.Lock()
	this.M[key] = val
	this.Unlock()
}

func (this *SafeMap) Get(key string) (interface{}, bool) {
	this.RLock()
	val, exists := this.M[key]
	this.RUnlock()
	return val, exists
}

func (this *SafeMap) Remove(key string) {
	this.Lock()
	delete(this.M, key)
	this.Unlock()
}

func (this *SafeMap) GetAndRemove(key string) (interface{}, bool) {
	this.Lock()
	val, exists := this.M[key]
	if exists {
		delete(this.M, key)
	}
	this.Unlock()
	return val, exists
}

func (this *SafeMap) Clear() {
	this.Lock()
	this.M = make(map[string]interface{})
	this.Unlock()
}

func (this *SafeMap) Keys() []string {
	this.RLock()
	defer this.RUnlock()

	keys := make([]string, 0)
	for key, _ := range this.M {
		keys = append(keys, key)
	}
	return keys
}

func (this *SafeMap) Slice() []interface{} {
	this.RLock()
	defer this.RUnlock()

	vals := make([]interface{}, 0)
	for _, val := range this.M {
		vals = append(vals, val)
	}
	return vals
}

func (this *SafeMap) ContainsKey(key string) bool {
	this.RLock()
	_, exists := this.M[key]
	this.RUnlock()
	return exists
}

func (this *SafeMap) Size() int {
	this.RLock()
	len := len(this.M)
	this.RUnlock()
	return len
}

func (this *SafeMap) IsEmpty() bool {
	this.RLock()
	empty := (len(this.M) == 0)
	this.RUnlock()
	return empty
}

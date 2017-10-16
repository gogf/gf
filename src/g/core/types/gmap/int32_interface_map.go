package gmap

import (
	"sync"
)

type Int32InterfaceMap struct {
	sync.RWMutex
	m map[int32]interface{}
}

func NewInt32InterfaceMap() *Int32InterfaceMap {
	return &Int32InterfaceMap{
        m: make(map[int32]interface{}),
    }
}

// 哈希表克隆
func (this *Int32InterfaceMap) Clone() *map[int32]interface{} {
	m := make(map[int32]interface{})
	this.RLock()
	for k, v := range this.m {
		m[k] = v
	}
    this.RUnlock()
	return &m
}

// 设置键值对
func (this *Int32InterfaceMap) Set(key int32, val interface{}) {
	this.Lock()
	this.m[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *Int32InterfaceMap) BatchSet(m map[int32]interface{}) {
	this.Lock()
	for k, v := range m {
		this.m[k] = v
	}
	this.Unlock()
}

// 获取键值
func (this *Int32InterfaceMap) Get(key int32) (interface{}) {
	this.RLock()
	val, _ := this.m[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *Int32InterfaceMap) Remove(key int32) {
    this.Lock()
    delete(this.m, key)
    this.Unlock()
}

// 批量删除键值对
func (this *Int32InterfaceMap) BatchRemove(keys []int32) {
    this.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *Int32InterfaceMap) GetAndRemove(key int32) (interface{}) {
    this.Lock()
    val, exists := this.m[key]
    if exists {
        delete(this.m, key)
    }
    this.Unlock()
    return val
}

// 返回键列表
func (this *Int32InterfaceMap) Keys() []int32 {
    this.RLock()
    keys := make([]int32, 0)
    for key, _ := range this.m {
        keys = append(keys, key)
    }
    this.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *Int32InterfaceMap) Values() []interface{} {
    this.RLock()
    vals := make([]interface{}, 0)
    for _, val := range this.m {
        vals = append(vals, val)
    }
    this.RUnlock()
    return vals
}

// 是否存在某个键
func (this *Int32InterfaceMap) Contains(key int32) bool {
    this.RLock()
    _, exists := this.m[key]
    this.RUnlock()
    return exists
}

// 哈希表大小
func (this *Int32InterfaceMap) Size() int {
    this.RLock()
    len := len(this.m)
    this.RUnlock()
    return len
}

// 哈希表是否为空
func (this *Int32InterfaceMap) IsEmpty() bool {
    this.RLock()
    empty := (len(this.m) == 0)
    this.RUnlock()
    return empty
}

// 清空哈希表
func (this *Int32InterfaceMap) Clear() {
    this.Lock()
    this.m = make(map[int32]interface{})
    this.Unlock()
}


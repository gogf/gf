package gmap

import (
	"sync"
)

type UintInterfaceMap struct {
	sync.RWMutex
	m map[uint]interface{}
}

func NewUintInterfaceMap() *UintInterfaceMap {
	return &UintInterfaceMap{
        m: make(map[uint]interface{}),
    }
}

// 哈希表克隆
func (this *UintInterfaceMap) Clone() *map[uint]interface{} {
	m := make(map[uint]interface{})
	this.RLock()
	for k, v := range this.m {
		m[k] = v
	}
    this.RUnlock()
	return &m
}

// 设置键值对
func (this *UintInterfaceMap) Set(key uint, val interface{}) {
	this.Lock()
	this.m[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *UintInterfaceMap) BatchSet(m map[uint]interface{}) {
	this.Lock()
	for k, v := range m {
		this.m[k] = v
	}
	this.Unlock()
}

// 获取键值
func (this *UintInterfaceMap) Get(key uint) (interface{}) {
	this.RLock()
	val, _ := this.m[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *UintInterfaceMap) Remove(key uint) {
    this.Lock()
    delete(this.m, key)
    this.Unlock()
}

// 批量删除键值对
func (this *UintInterfaceMap) BatchRemove(keys []uint) {
    this.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *UintInterfaceMap) GetAndRemove(key uint) (interface{}) {
    this.Lock()
    val, exists := this.m[key]
    if exists {
        delete(this.m, key)
    }
    this.Unlock()
    return val
}

// 返回键列表
func (this *UintInterfaceMap) Keys() []uint {
    this.RLock()
    keys := make([]uint, 0)
    for key, _ := range this.m {
        keys = append(keys, key)
    }
    this.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *UintInterfaceMap) Values() []interface{} {
    this.RLock()
    vals := make([]interface{}, 0)
    for _, val := range this.m {
        vals = append(vals, val)
    }
    this.RUnlock()
    return vals
}

// 是否存在某个键
func (this *UintInterfaceMap) Contains(key uint) bool {
    this.RLock()
    _, exists := this.m[key]
    this.RUnlock()
    return exists
}

// 哈希表大小
func (this *UintInterfaceMap) Size() int {
    this.RLock()
    len := len(this.m)
    this.RUnlock()
    return len
}

// 哈希表是否为空
func (this *UintInterfaceMap) IsEmpty() bool {
    this.RLock()
    empty := (len(this.m) == 0)
    this.RUnlock()
    return empty
}

// 清空哈希表
func (this *UintInterfaceMap) Clear() {
    this.Lock()
    this.m = make(map[uint]interface{})
    this.Unlock()
}


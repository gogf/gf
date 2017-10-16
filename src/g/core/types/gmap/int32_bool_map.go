package gmap

import (
	"sync"
)

type Int32BoolMap struct {
	sync.RWMutex
	m map[int32]bool
}

func NewInt32BoolMap() *Int32BoolMap {
	return &Int32BoolMap{
        m: make(map[int32]bool),
    }
}

// 哈希表克隆
func (this *Int32BoolMap) Clone() *map[int32]bool {
	m := make(map[int32]bool)
	this.RLock()
	for k, v := range this.m {
		m[k] = v
	}
    this.RUnlock()
	return &m
}

// 设置键值对
func (this *Int32BoolMap) Set(key int32, val bool) {
	this.Lock()
	this.m[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *Int32BoolMap) BatchSet(m map[int32]bool) {
	this.Lock()
	for k, v := range m {
		this.m[k] = v
	}
	this.Unlock()
}

// 获取键值
func (this *Int32BoolMap) Get(key int32) (bool) {
	this.RLock()
	val, _ := this.m[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *Int32BoolMap) Remove(key int32) {
    this.Lock()
    delete(this.m, key)
    this.Unlock()
}

// 批量删除键值对
func (this *Int32BoolMap) BatchRemove(keys []int32) {
    this.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *Int32BoolMap) GetAndRemove(key int32) (bool) {
    this.Lock()
    val, exists := this.m[key]
    if exists {
        delete(this.m, key)
    }
    this.Unlock()
    return val
}

// 返回键列表
func (this *Int32BoolMap) Keys() []int32 {
    this.RLock()
    keys := make([]int32, 0)
    for key, _ := range this.m {
        keys = append(keys, key)
    }
    this.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
//func (this *Int32BoolMap) Values() []bool {
//    this.RLock()
//    vals := make([]bool, 0)
//    for _, val := range this.m {
//        vals = append(vals, val)
//    }
//    this.RUnlock()
//    return vals
//}

// 是否存在某个键
func (this *Int32BoolMap) Contains(key int32) bool {
	this.RLock()
	_, exists := this.m[key]
	this.RUnlock()
	return exists
}

// 哈希表大小
func (this *Int32BoolMap) Size() int {
    this.RLock()
    len := len(this.m)
    this.RUnlock()
    return len
}

// 哈希表是否为空
func (this *Int32BoolMap) IsEmpty() bool {
    this.RLock()
    empty := (len(this.m) == 0)
    this.RUnlock()
    return empty
}

// 清空哈希表
func (this *Int32BoolMap) Clear() {
    this.Lock()
    this.m = make(map[int32]bool)
    this.Unlock()
}


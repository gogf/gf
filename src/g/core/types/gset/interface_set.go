package gset

import (
	"fmt"
	"sync"
)

type InterfaceSet struct {
	m sync.RWMutex
	M map[interface{}]struct{}
}

func NewInterfaceSet() *InterfaceSet {
	return &InterfaceSet{M: make(map[interface{}]struct{})}
}

// 设置键
func (this *InterfaceSet) Add(item interface{}) *InterfaceSet {
	if this.Contains(item) {
		return this
	}
	this.m.Lock()
	this.M[item] = struct{}{}
	this.m.Unlock()
	return this
}

// 批量添加设置键
func (this *InterfaceSet) BatchAdd(items []interface{}) *InterfaceSet {
    count := len(items)
    if count == 0 {
        return this
    }

    todo := make([]interface{}, 0, count)
    this.m.RLock()
    for i := 0; i < count; i++ {
        _, exists := this.M[items[i]]
        if exists {
            continue
        }

        todo = append(todo, items[i])
    }
    this.m.RUnlock()

    count = len(todo)
    if count == 0 {
        return this
    }

    this.m.Lock()
    for i := 0; i < count; i++ {
        this.M[todo[i]] = struct{}{}
    }
    this.m.Unlock()
    return this
}

// 键是否存在
func (this *InterfaceSet) Contains(item interface{}) bool {
	this.m.RLock()
	_, exists := this.M[item]
	this.m.RUnlock()
	return exists
}

// 删除键值对
func (this *InterfaceSet) Remove(key interface{}) {
	this.m.Lock()
	delete(this.M, key)
	this.m.Unlock()
}

// 大小
func (this *InterfaceSet) Size() int {
	this.m.RLock()
	l := len(this.M)
	this.m.RUnlock()
	return l
}

// 清空set
func (this *InterfaceSet) Clear() {
	this.m.Lock()
	this.M = make(map[interface{}]struct{})
	this.m.Unlock()
}

// 转换为数组
func (this *InterfaceSet) Slice() []interface{} {
	this.m.RLock()
	ret := make([]interface{}, len(this.M))
	i := 0
	for item := range this.M {
		ret[i] = item
		i++
	}

	this.m.RUnlock()
	return ret
}

// 转换为字符串
func (this *InterfaceSet) String() string {
	s := this.Slice()
	return fmt.Sprint(s)
}

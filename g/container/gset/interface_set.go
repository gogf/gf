// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package gset

import (
	"fmt"
	"sync"
)

type InterfaceSet struct {
	sync.RWMutex
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
	this.Lock()
	this.M[item] = struct{}{}
	this.Unlock()
	return this
}

// 批量添加设置键
func (this *InterfaceSet) BatchAdd(items []interface{}) *InterfaceSet {
    count := len(items)
    if count == 0 {
        return this
    }

    todo := make([]interface{}, 0, count)
    this.RLock()
    for i := 0; i < count; i++ {
        _, exists := this.M[items[i]]
        if exists {
            continue
        }

        todo = append(todo, items[i])
    }
    this.RUnlock()

    count = len(todo)
    if count == 0 {
        return this
    }

    this.Lock()
    for i := 0; i < count; i++ {
        this.M[todo[i]] = struct{}{}
    }
    this.Unlock()
    return this
}

// 键是否存在
func (this *InterfaceSet) Contains(item interface{}) bool {
	this.RLock()
	_, exists := this.M[item]
	this.RUnlock()
	return exists
}

// 删除键值对
func (this *InterfaceSet) Remove(key interface{}) {
	this.Lock()
	delete(this.M, key)
	this.Unlock()
}

// 大小
func (this *InterfaceSet) Size() int {
	this.RLock()
	l := len(this.M)
	this.RUnlock()
	return l
}

// 清空set
func (this *InterfaceSet) Clear() {
	this.Lock()
	this.M = make(map[interface{}]struct{})
	this.Unlock()
}

// 转换为数组
func (this *InterfaceSet) Slice() []interface{} {
	this.RLock()
	ret := make([]interface{}, len(this.M))
	i := 0
	for item := range this.M {
		ret[i] = item
		i++
	}

	this.RUnlock()
	return ret
}

// 转换为字符串
func (this *InterfaceSet) String() string {
	s := this.Slice()
	return fmt.Sprint(s)
}

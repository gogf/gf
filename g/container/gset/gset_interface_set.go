// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package gset

import (
	"fmt"
	"gitee.com/johng/gf/g/internal/rwmutex"
)

type InterfaceSet struct {
	mu *rwmutex.RWMutex
	m  map[interface{}]struct{}
}

func NewInterfaceSet(unsafe...bool) *InterfaceSet {
	return &InterfaceSet{
		m  : make(map[interface{}]struct{}),
		mu : rwmutex.New(unsafe...),
    }
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (set *InterfaceSet) Iterator(f func (v interface{}) bool) {
    set.mu.RLock()
    defer set.mu.RUnlock()
    for k, _ := range set.m {
		if !f(k) {
			break
		}
    }
}

// 添加
func (set *InterfaceSet) Add(item interface{}) *InterfaceSet {
	set.mu.Lock()
	set.m[item] = struct{}{}
	set.mu.Unlock()
	return set
}

// 批量添加
func (set *InterfaceSet) BatchAdd(items []interface{}) *InterfaceSet {
	set.mu.Lock()
	for _, item := range items {
		set.m[item] = struct{}{}
	}
	set.mu.Unlock()
    return set
}

// 键是否存在
func (set *InterfaceSet) Contains(item interface{}) bool {
	set.mu.RLock()
	_, exists := set.m[item]
	set.mu.RUnlock()
	return exists
}

// 删除键值对
func (set *InterfaceSet) Remove(key interface{}) {
	set.mu.Lock()
	delete(set.m, key)
	set.mu.Unlock()
}

// 大小
func (set *InterfaceSet) Size() int {
	set.mu.RLock()
	l := len(set.m)
	set.mu.RUnlock()
	return l
}

// 清空set
func (set *InterfaceSet) Clear() {
	set.mu.Lock()
	set.m = make(map[interface{}]struct{})
	set.mu.Unlock()
}

// 转换为数组
func (set *InterfaceSet) Slice() []interface{} {
	set.mu.RLock()
	i   := 0
	ret := make([]interface{}, len(set.m))
	for item := range set.m {
		ret[i] = item
		i++
	}
	set.mu.RUnlock()
	return ret
}

// 转换为字符串
func (set *InterfaceSet) String() string {
	return fmt.Sprint(set.Slice())
}

func (set *InterfaceSet) LockFunc(f func(m map[interface{}]struct{})) {
	set.mu.Lock(true)
	defer set.mu.Unlock(true)
	f(set.m)
}

func (set *InterfaceSet) RLockFunc(f func(m map[interface{}]struct{})) {
	set.mu.RLock(true)
	defer set.mu.RUnlock(true)
	f(set.m)
}

// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//
//

package gset

import (
	"fmt"
	"gitee.com/johng/gf/g/internal/rwmutex"
)

type IntSet struct {
	mu *rwmutex.RWMutex
	m  map[int]struct{}
}

func NewIntSet(unsafe...bool) *IntSet {
	return &IntSet{
		m  : make(map[int]struct{}),
		mu : rwmutex.New(unsafe...),
	}
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (set *IntSet) Iterator(f func (v int) bool) {
    set.mu.RLock()
    defer set.mu.RUnlock()
	for k, _ := range set.m {
		if !f(k) {
			break
		}
	}
}

// 设置键
func (set *IntSet) Add(item int) *IntSet {
	set.mu.Lock()
	set.m[item] = struct{}{}
	set.mu.Unlock()
	return set
}

// 批量添加设置键
func (set *IntSet) BatchAdd(items []int) *IntSet {
	set.mu.Lock()
	for _, item := range items {
		set.m[item] = struct{}{}
	}
	set.mu.Unlock()
    return set
}

// 键是否存在
func (set *IntSet) Contains(item int) bool {
	set.mu.RLock()
	_, exists := set.m[item]
	set.mu.RUnlock()
	return exists
}

// 删除键值对
func (set *IntSet) Remove(key int) {
	set.mu.Lock()
	delete(set.m, key)
	set.mu.Unlock()
}

// 大小
func (set *IntSet) Size() int {
	set.mu.RLock()
	l := len(set.m)
	set.mu.RUnlock()
	return l
}

// 清空set
func (set *IntSet) Clear() {
	set.mu.Lock()
	set.m = make(map[int]struct{})
	set.mu.Unlock()
}

// 转换为数组
func (set *IntSet) Slice() []int {
	set.mu.RLock()
	ret := make([]int, len(set.m))
	i := 0
	for item := range set.m {
		ret[i] = item
		i++
	}

	set.mu.RUnlock()
	return ret
}

// 转换为字符串
func (set *IntSet) String() string {
	return fmt.Sprint(set.Slice())
}

func (set *IntSet) LockFunc(f func(m map[int]struct{})) {
	set.mu.Lock(true)
	defer set.mu.Unlock(true)
	f(set.m)
}

func (set *IntSet) RLockFunc(f func(m map[int]struct{})) {
    set.mu.RLock(true)
    defer set.mu.RUnlock(true)
    f(set.m)
}
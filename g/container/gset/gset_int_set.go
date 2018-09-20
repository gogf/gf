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
	"gitee.com/johng/gf/g/container/internal/rwmutex"
)

type IntSet struct {
	mu *rwmutex.RWMutex
	m  map[int]struct{}
}

func NewIntSet(safe...bool) *IntSet {
	return &IntSet{
		m  : make(map[int]struct{}),
		mu : rwmutex.New(safe...),
	}
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (this *IntSet) Iterator(f func (v int) bool) {
    this.mu.RLock()
    defer this.mu.RUnlock()
	for k, _ := range this.m {
		if !f(k) {
			break
		}
	}
}

// 设置键
func (this *IntSet) Add(item int) *IntSet {
	this.mu.Lock()
	this.m[item] = struct{}{}
	this.mu.Unlock()
	return this
}

// 批量添加设置键
func (this *IntSet) BatchAdd(items []int) *IntSet {
	this.mu.Lock()
	for _, item := range items {
		this.m[item] = struct{}{}
	}
	this.mu.Unlock()
    return this
}

// 键是否存在
func (this *IntSet) Contains(item int) bool {
	this.mu.RLock()
	_, exists := this.m[item]
	this.mu.RUnlock()
	return exists
}

// 删除键值对
func (this *IntSet) Remove(key int) {
	this.mu.Lock()
	delete(this.m, key)
	this.mu.Unlock()
}

// 大小
func (this *IntSet) Size() int {
	this.mu.RLock()
	l := len(this.m)
	this.mu.RUnlock()
	return l
}

// 清空set
func (this *IntSet) Clear() {
	this.mu.Lock()
	this.m = make(map[int]struct{})
	this.mu.Unlock()
}

// 转换为数组
func (this *IntSet) Slice() []int {
	this.mu.RLock()
	ret := make([]int, len(this.m))
	i := 0
	for item := range this.m {
		ret[i] = item
		i++
	}

	this.mu.RUnlock()
	return ret
}

// 转换为字符串
func (this *IntSet) String() string {
	return fmt.Sprint(this.Slice())
}

func (this *IntSet) LockFunc(f func(m map[int]struct{})) {
	this.mu.Lock(true)
	defer this.mu.Unlock(true)
	f(this.m)
}

func (this *IntSet) RLockFunc(f func(m map[int]struct{})) {
    this.mu.RLock(true)
    defer this.mu.RUnlock(true)
    f(this.m)
}
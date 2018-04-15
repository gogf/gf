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

type UintSet struct {
	mu sync.RWMutex
	m  map[uint]struct{}
}

func NewUintSet() *UintSet {
	return &UintSet{m: make(map[uint]struct{})}
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (this *UintSet) Iterator(f func (v uint) bool) {
    this.mu.RLock()
    for k, _ := range this.m {
		if !f(k) {
			break
		}
    }
    this.mu.RUnlock()
}

// 添加
func (this *UintSet) Add(item uint) *UintSet {
	this.mu.Lock()
	this.m[item] = struct{}{}
	this.mu.Unlock()
	return this
}

// 批量添加
func (this *UintSet) BatchAdd(items []uint) *UintSet {
	this.mu.Lock()
	for _, item := range items {
		this.m[item] = struct{}{}
	}
	this.mu.Unlock()
    return this
}

// 键是否存在
func (this *UintSet) Contains(item uint) bool {
	this.mu.RLock()
	_, exists := this.m[item]
	this.mu.RUnlock()
	return exists
}

// 删除键值对
func (this *UintSet) Remove(key uint) {
	this.mu.Lock()
	delete(this.m, key)
	this.mu.Unlock()
}

// 大小
func (this *UintSet) Size() int {
	this.mu.RLock()
	l := len(this.m)
	this.mu.RUnlock()
	return l
}

// 清空set
func (this *UintSet) Clear() {
	this.mu.Lock()
	this.m = make(map[uint]struct{})
	this.mu.Unlock()
}

// 转换为数组
func (this *UintSet) Slice() []uint {
	this.mu.RLock()
	i   := 0
	ret := make([]uint, len(this.m))
	for item := range this.m {
		ret[i] = item
		i++
	}
	this.mu.RUnlock()
	return ret
}

// 转换为字符串
func (this *UintSet) String() string {
	return fmt.Sprint(this.Slice())
}

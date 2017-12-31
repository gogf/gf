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

type Int64Set struct {
	sync.RWMutex
	M map[int64]struct{}
}

func NewInt64Set() *Int64Set {
	return &Int64Set{M: make(map[int64]struct{})}
}

// 设置键
func (this *Int64Set) Add(item int64) *Int64Set {
	if this.Contains(item) {
		return this
	}
	this.Lock()
	this.M[item] = struct{}{}
	this.Unlock()
	return this
}

// 批量添加设置键
func (this *Int64Set) BatchAdd(items []int64) *Int64Set {
    count := len(items)
    if count == 0 {
        return this
    }

    todo := make([]int64, 0, count)
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
func (this *Int64Set) Contains(item int64) bool {
	this.RLock()
	_, exists := this.M[item]
	this.RUnlock()
	return exists
}

// 删除键值对
func (this *Int64Set) Remove(key int64) {
	this.Lock()
	delete(this.M, key)
	this.Unlock()
}

// 大小
func (this *Int64Set) Size() int {
	this.RLock()
	l := len(this.M)
	this.RUnlock()
	return l
}

// 清空set
func (this *Int64Set) Clear() {
	this.Lock()
	this.M = make(map[int64]struct{})
	this.Unlock()
}

// 转换为数组
func (this *Int64Set) Slice() []int64 {
	this.RLock()
	ret := make([]int64, len(this.M))
	i := 0
	for item := range this.M {
		ret[i] = item
		i++
	}

	this.RUnlock()
	return ret
}

// 转换为字符串
func (this *Int64Set) String() string {
	s := this.Slice()
	return fmt.Sprint(s)
}

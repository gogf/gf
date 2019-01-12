// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package garray

import (
	"gitee.com/johng/gf/g/internal/rwmutex"
)

type IntArray struct {
	mu    *rwmutex.RWMutex // 互斥锁
	cap   int              // 初始化设置的数组容量
	size  int              // 初始化设置的数组大小
	array []int            // 底层数组
}

func NewIntArray(size int, cap int, unsafe...bool) *IntArray {
	a := &IntArray{
		mu : rwmutex.New(unsafe...),
	}
	a.size = size
	if cap > 0 {
		a.cap   = cap
		a.array = make([]int, size, cap)
	} else {
		a.array = make([]int, size)
	}
	return a
}

// 获取指定索引的数据项, 调用方注意判断数组边界
func (a *IntArray) Get(index int) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	value := a.array[index]
	return value
}

// 设置指定索引的数据项, 调用方注意判断数组边界
func (a *IntArray) Set(index int, value int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array[index] = value
}

// 在当前索引位置前插入一个数据项, 调用方注意判断数组边界
func (a *IntArray) InsertBefore(index int, value int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	rear   := append([]int{}, a.array[index : ]...)
	a.array = append(a.array[0 : index], value)
	a.array = append(a.array, rear...)
}

// 在当前索引位置后插入一个数据项, 调用方注意判断数组边界
func (a *IntArray) InsertAfter(index int, value int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	rear   := append([]int{}, a.array[index + 1:]...)
	a.array = append(a.array[0 : index + 1], value)
	a.array = append(a.array, rear...)

}

// 删除指定索引的数据项, 调用方注意判断数组边界
func (a *IntArray) Remove(index int) int {
	a.mu.Lock()
	defer a.mu.Unlock()
	// 边界删除判断，以提高删除效率
	if index == 0 {
		value  := a.array[0]
		a.array = a.array[1 : ]
		return value
	} else if index == len(a.array) - 1 {
		value  := a.array[index]
		a.array = a.array[: index]
		return value
	}
	// 如果非边界删除，会涉及到数组创建，那么删除的效率差一些
	value  := a.array[index]
	a.array = append(a.array[ : index], a.array[index + 1 : ]...)
	return value
}

// 追加数据项
func (a *IntArray) Append(value...int) {
	a.mu.Lock()
	a.array = append(a.array, value...)
    a.mu.Unlock()
}

// 数组长度
func (a *IntArray) Len() int {
	a.mu.RLock()
	length := len(a.array)
	a.mu.RUnlock()
	return length
}

// 返回原始数据数组
func (a *IntArray) Slice() []int {
	array := ([]int)(nil)
	if a.mu.IsSafe() {
		a.mu.RLock()
		array = make([]int, len(a.array))
		for k, v := range a.array {
			array[k] = v
		}
		a.mu.RUnlock()
	} else {
		array = a.array
	}
	return array
}

// 清空数据数组
func (a *IntArray) Clear() {
	a.mu.Lock()
	if len(a.array) > 0 {
		if a.cap > 0 {
			a.array = make([]int, a.size, a.cap)
		} else {
			a.array = make([]int, a.size)
		}
	}
    a.mu.Unlock()
}

// 查找指定数值的索引位置，返回索引位置，如果查找不到则返回-1
func (a *IntArray) Search(value int) int {
	if len(a.array) == 0 {
		return -1
	}
	a.mu.RLock()
	result := -1
	for index, v := range a.array {
		if v == value {
			result = index
			break
		}
	}
	a.mu.RUnlock()

	return result
}

// 清理数组中重复的元素项
func (a *IntArray) Unique() *IntArray {
	a.mu.Lock()
	for i := 0; i < len(a.array) - 1; i++ {
		for j := i + 1; j < len(a.array); j++ {
			if a.array[i] == a.array[j] {
				a.array = append(a.array[ : j], a.array[j + 1 : ]...)
			}
		}
	}
	a.mu.Unlock()
	return a
}

// 使用自定义方法执行加锁修改操作
func (a *IntArray) LockFunc(f func(array []int)) {
	a.mu.Lock(true)
	defer a.mu.Unlock(true)
	f(a.array)
}

// 使用自定义方法执行加锁读取操作
func (a *IntArray) RLockFunc(f func(array []int)) {
	a.mu.RLock(true)
	defer a.mu.RUnlock(true)
	f(a.array)
}

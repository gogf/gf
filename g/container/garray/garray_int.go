// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package garray

import "sync"

type IntArray struct {
    mu           sync.RWMutex        // 互斥锁
    cap          int                 // 初始化设置的数组容量
    size         int                 // 初始化设置的数组大小
    array        []int               // 底层数组
    compareFunc func(v1, v2 int) int // 比较函数，返回值 -1: v1 < v2；0: v1 == v2；1: v1 > v2
}

func NewIntArray(size int, cap ... int) *IntArray {
    a     := &IntArray{}
    a.size = size
    if len(cap) > 0 {
        a.cap   = cap[0]
        a.array = make([]int, size, cap[0])
    } else {
        a.array = make([]int, size)
    }
    a.compareFunc = func(v1, v2 int) int {
        if v1 < v2 {
            return -1
        }
        if v1 > v2 {
            return 1
        }
        return 0
    }
    return a
}

// 获取指定索引的数据项, 调用方注意判断数组边界
func (a *IntArray) Get(index int) int {
    a.mu.RLock()
    value := a.array[index]
    a.mu.RUnlock()
    return value
}

// 设置指定索引的数据项, 调用方注意判断数组边界
func (a *IntArray) Set(index int, value int) {
    a.mu.Lock()
    a.array[index] = value
    a.mu.Unlock()
}

// 在当前索引位置前插入一个数据项, 调用方注意判断数组边界
func (a *IntArray) Insert(index int, value int) {
    a.mu.Lock()
    rear   := append([]int{}, a.array[index : ]...)
    a.array = append(a.array[0 : index], value)
    a.array = append(a.array, rear...)
    a.mu.Unlock()
}

// 删除指定索引的数据项, 调用方注意判断数组边界
func (a *IntArray) Remove(index int) {
    a.mu.Lock()
    a.array = append(a.array[ : index], a.array[index + 1 : ]...)
    a.mu.Unlock()
}

// 追加数据项
func (a *IntArray) Append(value int) {
    a.mu.Lock()
    a.array = append(a.array, value)
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
    a.mu.RLock()
    array := a.array
    a.mu.RUnlock()
    return array
}

// 清空数据数组
func (a *IntArray) Clear() {
    a.mu.Lock()
    if a.cap > 0 {
        a.array = make([]int, a.size, a.cap)
    } else {
        a.array = make([]int, a.size)
    }
    a.mu.Unlock()
}

// 查找指定数值的索引位置，返回索引位置，如果查找不到则返回-1
func (a *IntArray) Search(value int) int {
    if len(a.array) == 0 {
        return -1
    }
    a.mu.RLock()
    min := 0
    max := len(a.array) - 1
    mid := 0
    cmp := -2
    for {
        if cmp == 0 || min > max {
            break
        }
        for {
            mid = int((min + max) / 2)
            cmp = a.compareFunc(value, a.array[mid])
            switch cmp {
                case -1 : max = mid - 1
                case  0 :
                case  1 : min = mid + 1
            }
            if cmp == 0 || min > max {
                break
            }
        }
    }
    a.mu.RUnlock()

    if cmp == 0 {
        return mid
    }
    return -1
}

// 使用自定义方法执行加锁修改操作
func (a *IntArray) LockFunc(f func(array []int)) {
    a.mu.Lock()
    f(a.array)
    a.mu.Unlock()
}

// 使用自定义方法执行加锁读取操作
func (a *IntArray) RLockFunc(f func(array []int)) {
    a.mu.RLock()
    f(a.array)
    a.mu.RUnlock()
}

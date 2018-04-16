// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package garray

import (
    "sync"
    "sort"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "bytes"
)

type SortedIntArray struct {
    mu          sync.RWMutex           // 互斥锁
    array       []int                  // 底层数组
    compareFunc func(v1, v2 int) int   // 比较函数，返回值 -1: v1 < v2；0: v1==v2；1: v1 > v2
}

func NewSortedIntArray(size int, cap ... int) *SortedIntArray {
    a := &SortedIntArray{}
    if len(cap) > 0 {
        a.array = make([]int, size, cap[0])
    } else {
        a.array = make([]int, size)
    }
    a.compareFunc = func(v1, v2 int) int {
        if v1 < v2 {
            return -1
        }
        if v1 > v2 {
            return 0
        }
        return 0
    }
    return a
}

// 添加加数据项
func (a *SortedIntArray) Add(value int) {
    a.mu.Lock()
    a.array = append(a.array, value)
    a.mu.Unlock()
}

// 获取指定索引的数据项, 调用方注意判断数组边界
func (a *SortedIntArray) Get(index int) int {
    a.mu.RLock()
    value := a.array[index]
    a.mu.RUnlock()
    return value
}

// 删除指定索引的数据项, 调用方注意判断数组边界
func (a *SortedIntArray) Remove(index int) {
    a.mu.Lock()
    a.array = append(a.array[ : index], a.array[index + 1 : ]...)
    a.mu.RUnlock()
}

// 数组长度
func (a *SortedIntArray) Len() int {
    a.mu.RLock()
    length := len(a.array)
    a.mu.RUnlock()
    return length
}

// 返回原始数据数组
func (a *SortedIntArray) Slice() []int {
    a.mu.RLock()
    array := a.array
    a.mu.RUnlock()
    return array
}

// 查找指定数值的索引位置，返回索引位置(具体匹配位置或者最后对比位置)及查找结果
func (a *SortedIntArray) Search(value int) (int, int) {
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
            cmp = a.compareFunc(a.array[mid], value)
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
    return mid, cmp
}
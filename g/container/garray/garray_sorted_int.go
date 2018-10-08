// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package garray

import (
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/container/internal/rwmutex"
)

// 默认按照从低到高进行排序
type SortedIntArray struct {
    mu          *rwmutex.RWMutex     // 互斥锁
    cap         int                  // 初始化设置的数组容量
    array       []int                // 底层数组
    unique      *gtype.Bool          // 是否要求不能重复
    compareFunc func(v1, v2 int) int // 比较函数，返回值 -1: v1 < v2；0: v1 == v2；1: v1 > v2
}

// 创建一个排序的int数组
func NewSortedIntArray(size int, cap int, safe...bool) *SortedIntArray {
    return &SortedIntArray {
        mu          : rwmutex.New(safe...),
        array       : make([]int, 0, cap),
        unique      : gtype.NewBool(),
        compareFunc : func(v1, v2 int) int {
            if v1 < v2 {
                return -1
            }
            if v1 > v2 {
                return 1
            }
            return 0
        },
    }
}

// 添加加数据项
func (a *SortedIntArray) Add(values...int) {
    if len(values) == 0 {
        return
    }
    for _, value := range values {
        index, cmp := a.Search(value)
        if a.unique.Val() && cmp == 0 {
            continue
        }
        a.mu.Lock()
        if index < 0 {
            a.array = append(a.array, value)
            a.mu.Unlock()
            continue
        }
        // 加到指定索引后面
        if cmp > 0 {
            index++
        }
        rear   := append([]int{}, a.array[index : ]...)
        a.array = append(a.array[0 : index], value)
        a.array = append(a.array, rear...)
        a.mu.Unlock()
    }
}

// 获取指定索引的数据项, 调用方注意判断数组边界
func (a *SortedIntArray) Get(index int) int {
    a.mu.RLock()
    defer a.mu.RUnlock()
    value := a.array[index]
    return value
}

// 删除指定索引的数据项, 调用方注意判断数组边界
func (a *SortedIntArray) Remove(index int) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.array = append(a.array[ : index], a.array[index + 1 : ]...)
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

// 查找指定数值的索引位置，返回索引位置(具体匹配位置或者最后对比位置)及查找结果
// 返回值: 最后比较位置, 比较结果
func (a *SortedIntArray) Search(value int) (int, int) {
    if len(a.array) == 0 {
        return -1, -2
    }
    a.mu.RLock()
    min := 0
    max := len(a.array) - 1
    mid := 0
    cmp := -2
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
    a.mu.RUnlock()
    return mid, cmp
}

// 设置是否允许数组唯一
func (a *SortedIntArray) SetUnique(unique bool) {
    oldUnique := a.unique.Val()
    a.unique.Set(unique)
    if unique && oldUnique != unique {
        a.doUnique()
    }
}

// 清理数组中重复的元素项
func (a *SortedIntArray) doUnique() {
    a.mu.Lock()
    i := 0
    for {
        if i == len(a.array) - 1 {
            break
        }
        if a.compareFunc(a.array[i], a.array[i + 1]) == 0 {
            a.array = append(a.array[ : i + 1], a.array[i + 1 + 1 : ]...)
        } else {
            i++
        }
    }
    a.mu.Unlock()
}

// 清空数据数组
func (a *SortedIntArray) Clear() {
    a.mu.Lock()
    a.array = make([]int, 0, a.cap)
    a.mu.Unlock()
}

// 使用自定义方法执行加锁修改操作
func (a *SortedIntArray) LockFunc(f func(array []int)) {
    a.mu.Lock(true)
    defer a.mu.Unlock(true)
    f(a.array)
}

// 使用自定义方法执行加锁读取操作
func (a *SortedIntArray) RLockFunc(f func(array []int)) {
    a.mu.RLock(true)
    defer a.mu.RUnlock(true)
    f(a.array)
}
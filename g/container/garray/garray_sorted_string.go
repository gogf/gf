// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package garray

import (
    "sync"
    "gitee.com/johng/gf/g/container/gtype"
    "strings"
)

// 默认按照从低到高进行排序
type SortedStringArray struct {
    mu          sync.RWMutex            // 互斥锁
    cap         int                     // 初始化设置的数组容量
    size        int                     // 初始化设置的数组大小
    array       []string                // 底层数组
    unique      *gtype.Bool             // 是否要求不能重复
    compareFunc func(v1, v2 string) int // 比较函数，返回值 -1: v1 < v2；0: v1 == v2；1: v1 > v2
}

func NewSortedStringArray(size int, cap ... int) *SortedStringArray {
    a := &SortedStringArray {
        unique      : gtype.NewBool(),
        compareFunc : func(v1, v2 string) int {
            return strings.Compare(v1, v2)
        },
    }
    a.size = size
    if len(cap) > 0 {
        a.cap   = cap[0]
        a.array = make([]string, size, cap[0])
    } else {
        a.array = make([]string, size)
    }
    return a
}

// 添加加数据项
func (a *SortedStringArray) Add(value string) {
    index, cmp := a.Search(value)
    if a.unique.Val() && cmp == 0 {
        return
    }
    if index < 0 {
        a.mu.Lock()
        a.array = append(a.array, value)
        a.mu.Unlock()
        return
    }
    // 加到指定索引后面
    if cmp > 0 {
        index++
    }
    a.mu.Lock()
    rear   := append([]string{}, a.array[index : ]...)
    a.array = append(a.array[0 : index], value)
    a.array = append(a.array, rear...)
    a.mu.Unlock()
}

// 获取指定索引的数据项, 调用方注意判断数组边界
func (a *SortedStringArray) Get(index int) string {
    a.mu.RLock()
    value := a.array[index]
    a.mu.RUnlock()
    return value
}

// 删除指定索引的数据项, 调用方注意判断数组边界
func (a *SortedStringArray) Remove(index int) {
    a.mu.Lock()
    a.array = append(a.array[ : index], a.array[index + 1 : ]...)
    a.mu.Unlock()
}

// 数组长度
func (a *SortedStringArray) Len() int {
    a.mu.RLock()
    length := len(a.array)
    a.mu.RUnlock()
    return length
}

// 返回原始数据数组
func (a *SortedStringArray) Slice() []string {
    a.mu.RLock()
    array := a.array
    a.mu.RUnlock()
    return array
}

// 查找指定数值的索引位置，返回索引位置(具体匹配位置或者最后对比位置)及查找结果
func (a *SortedStringArray) Search(value string) (int, int) {
    if len(a.array) == 0 {
        return -1, -2
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
    return mid, cmp
}

// 设置是否允许数组唯一
func (a *SortedStringArray) SetUnique(unique bool) {
    oldUnique := a.unique.Val()
    a.unique.Set(unique)
    if unique && oldUnique != unique {
        a.doUnique()
    }
}

// 清理数组中重复的元素项
func (a *SortedStringArray) doUnique() {
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
func (a *SortedStringArray) Clear() {
    a.mu.Lock()
    if a.cap > 0 {
        a.array = make([]string, a.size, a.cap)
    } else {
        a.array = make([]string, a.size)
    }
    a.mu.Unlock()
}

// 使用自定义方法执行加锁修改操作
func (a *SortedStringArray) LockFunc(f func(array []string)) {
    a.mu.Lock()
    f(a.array)
    a.mu.Unlock()
}

// 使用自定义方法执行加锁读取操作
func (a *SortedStringArray) RLockFunc(f func(array []string)) {
    a.mu.RLock()
    f(a.array)
    a.mu.RUnlock()
}
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package garray

import (
    "gitee.com/johng/gf/g/internal/rwmutex"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/grand"
    "math"
    "sort"
    "strings"
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

func NewIntArrayEmpty(unsafe...bool) *IntArray {
    return NewIntArray(0, 0, unsafe...)
}

func NewIntArrayFrom(array []int, unsafe...bool) *IntArray {
	return &IntArray{
        mu    : rwmutex.New(unsafe...),
        array : array,
    }
}

// 获取指定索引的数据项, 调用方注意判断数组边界
func (a *IntArray) Get(index int) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	value := a.array[index]
	return value
}

// 设置指定索引的数据项, 调用方注意判断数组边界
func (a *IntArray) Set(index int, value int) *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array[index] = value
	return a
}

// 设置底层数组变量.
func (a *IntArray) SetArray(array []int) *IntArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.array = array
    return a
}

// 使用指定数组替换到对应的索引元素值.
func (a *IntArray) Replace(array []int) *IntArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    max := len(array)
    if max > len(a.array) {
        max = len(a.array)
    }
    for i := 0; i < max; i++ {
        a.array[i] = array[i]
    }
    return a
}

// Calculate the sum of values in an array.
//
// 对数组中的元素项求和。
func (a *IntArray) Sum() (sum int) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    for _, v := range a.array {
        sum += v
    }
    return
}

// 将数组重新排序(从小到大).
func (a *IntArray) Sort() *IntArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    sort.Ints(a.array)
    return a
}

// 在当前索引位置前插入一个数据项, 调用方注意判断数组边界
func (a *IntArray) InsertBefore(index int, value int) *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	rear   := append([]int{}, a.array[index : ]...)
	a.array = append(a.array[0 : index], value)
	a.array = append(a.array, rear...)
    return a
}

// 在当前索引位置后插入一个数据项, 调用方注意判断数组边界
func (a *IntArray) InsertAfter(index int, value int) *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	rear   := append([]int{}, a.array[index + 1:]...)
	a.array = append(a.array[0 : index + 1], value)
	a.array = append(a.array, rear...)
    return a
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

// 将数据项添加到数组的最左端(索引为0)
func (a *IntArray) PushLeft(value...int) *IntArray {
    a.mu.Lock()
    a.array = append(value, a.array...)
    a.mu.Unlock()
    return a
}

// 将数据项添加到数组的最右端(索引为length - 1), 等于: Append
func (a *IntArray) PushRight(value...int) *IntArray {
    a.mu.Lock()
    a.array = append(a.array, value...)
    a.mu.Unlock()
    return a
}

// 将最左端(索引为0)的数据项移出数组，并返回该数据项
func (a *IntArray) PopLeft() int {
    a.mu.Lock()
    defer a.mu.Unlock()
    value  := a.array[0]
    a.array = a.array[1 : ]
    return value
}

// 将最右端(索引为length - 1)的数据项移出数组，并返回该数据项
func (a *IntArray) PopRight() int {
    a.mu.Lock()
    defer a.mu.Unlock()
    index  := len(a.array) - 1
    value  := a.array[index]
    a.array = a.array[: index]
    return value
}

// 随机将一个数据项移出数组，并返回该数据项
func (a *IntArray) PopRand() int {
    return a.Remove(grand.Intn(len(a.array)))
}

// 追加数据项
func (a *IntArray) Append(value...int) *IntArray {
	a.mu.Lock()
	a.array = append(a.array, value...)
    a.mu.Unlock()
    return a
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
		defer a.mu.RUnlock()
		array = make([]int, len(a.array))
		copy(array, a.array)
	} else {
		array = a.array
	}
	return array
}

// 清空数据数组
func (a *IntArray) Clear() *IntArray {
	a.mu.Lock()
	if len(a.array) > 0 {
		if a.cap > 0 {
			a.array = make([]int, a.size, a.cap)
		} else {
			a.array = make([]int, a.size)
		}
	}
    a.mu.Unlock()
    return a
}

// 查找指定数值是否存在
func (a *IntArray) Contains(value int) bool {
    return a.Search(value) != -1
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
func (a *IntArray) LockFunc(f func(array []int)) *IntArray {
	a.mu.Lock(true)
	defer a.mu.Unlock(true)
	f(a.array)
    return a
}

// 使用自定义方法执行加锁读取操作
func (a *IntArray) RLockFunc(f func(array []int)) *IntArray {
	a.mu.RLock(true)
	defer a.mu.RUnlock(true)
	f(a.array)
    return a
}

// 合并两个数组.
func (a *IntArray) Merge(array *IntArray) *IntArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    if a != array {
        array.mu.RLock()
        defer array.mu.RUnlock()
    }
    a.array = append(a.array, array.array...)
    return a
}

// Fills an array with num entries of the value of the value parameter, keys starting at the startIndex parameter.
func (a *IntArray) Fill(startIndex int, num int, value int) *IntArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    if startIndex < 0 {
        startIndex = 0
    }
    for i := startIndex; i < startIndex + num; i++ {
        if i > len(a.array) - 1 {
            a.array = append(a.array, value)
        } else {
            a.array[i] = value
        }
    }
    return a
}

// Chunks an array into arrays with size elements. The last chunk may contain less than size elements.
func (a *IntArray) Chunk(size int) [][]int {
    if size < 1 {
        panic("size: cannot be less than 1")
    }
    a.mu.RLock()
    defer a.mu.RUnlock()
    length := len(a.array)
    chunks := int(math.Ceil(float64(length) / float64(size)))
    var n [][]int
    for i, end := 0, 0; chunks > 0; chunks-- {
        end = (i + 1) * size
        if end > length {
            end = length
        }
        n = append(n, a.array[i*size : end])
        i++
    }
    return n
}

// Pad array to the specified length with a value.
// If size is positive then the array is padded on the right, or negative on the left.
// If the absolute value of size is less than or equal to the length of the array
// then no padding takes place.
func (a *IntArray) Pad(size int, value int) *IntArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    if size == 0 || (size > 0 && size < len(a.array)) || (size < 0 && size > -len(a.array)) {
        return a
    }
    n := size
    if size < 0 {
        n = -size
    }
    n   -= len(a.array)
    tmp := make([]int, n)
    for i := 0; i < n; i++ {
        tmp[i] = value
    }
    if size > 0 {
        a.array = append(a.array, tmp...)
    } else {
        a.array = append(tmp, a.array...)
    }
    return a
}

// Extract a slice of the array(If in concurrent safe usage, it returns a copy of the slice; else a pointer).
// It returns the sequence of elements from the array array as specified by the offset and length parameters.
func (a *IntArray) SubSlice(offset, size int) []int {
    a.mu.RLock()
    defer a.mu.RUnlock()
    if offset > len(a.array) {
        return nil
    }
    if offset + size > len(a.array) {
        size = len(a.array) - offset
    }
    if a.mu.IsSafe() {
        s := make([]int, size)
        copy(s, a.array[offset:])
        return s
    } else {
        return a.array[offset:]
    }
}

// Picks one or more random entries out of an array(a copy), and returns the key (or keys) of the random entries.
func (a *IntArray) Rand(size int) []int {
    a.mu.RLock()
    defer a.mu.RUnlock()
    if size > len(a.array) {
        size = len(a.array)
    }
    n := make([]int, size)
    for i, v := range grand.Perm(len(a.array)) {
        n[i] = a.array[v]
        if i == size - 1 {
            break
        }
    }
    return n
}

// Randomly shuffles the array.
func (a *IntArray) Shuffle() *IntArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    for i, v := range grand.Perm(len(a.array)) {
        a.array[i], a.array[v] = a.array[v], a.array[i]
    }
    return a
}

// Make array with elements in reverse order.
func (a *IntArray) Reverse() *IntArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    for i, j := 0, len(a.array) - 1; i < j; i, j = i + 1, j - 1 {
        a.array[i], a.array[j] = a.array[j], a.array[i]
    }
    return a
}

// Join array elements with a string.
func (a *IntArray) Join(glue string) string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return strings.Join(gconv.Strings(a.array), glue)
}
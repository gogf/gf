// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
    "bytes"
    "github.com/gogf/gf/g/container/gtype"
    "github.com/gogf/gf/g/internal/rwmutex"
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/util/grand"
    "math"
    "sort"
    "strings"
)

// 默认按照从小到大进行排序
type SortedStringArray struct {
    mu          *rwmutex.RWMutex        // 互斥锁
    array       []string                // 底层数组
    unique      *gtype.Bool             // 是否要求不能重复
    compareFunc func(v1, v2 string) int // 比较函数，返回值 -1: v1 < v2；0: v1 == v2；1: v1 > v2
}

// Create an empty sorted array.
// The param <unsafe> used to specify whether using array with un-concurrent-safety,
// which is false in default, means concurrent-safe in default.
//
// 创建一个空的排序数组对象，参数unsafe用于指定是否用于非并发安全场景，默认为false，表示并发安全。
func NewSortedStringArray(unsafe...bool) *SortedStringArray {
    return NewSortedStringArraySize(0, unsafe...)
}

// Create a sorted array with given size and cap.
// The param <unsafe> used to specify whether using array with un-concurrent-safety,
// which is false in default, means concurrent-safe in default.
//
// 创建一个指定大小的排序数组对象，参数unsafe用于指定是否用于非并发安全场景，默认为false，表示并发安全。
func NewSortedStringArraySize(cap int, unsafe...bool) *SortedStringArray {
    return &SortedStringArray {
        mu          : rwmutex.New(unsafe...),
        array       : make([]string, 0, cap),
        unique      : gtype.NewBool(),
        compareFunc : func(v1, v2 string) int {
            return strings.Compare(v1, v2)
        },
    }
}

// Create an array with given slice <array>.
// The param <unsafe> used to specify whether using array with un-concurrent-safety,
// which is false in default, means concurrent-safe in default.
//
// 通过给定的slice变量创建排序数组对象，参数unsafe用于指定是否用于非并发安全场景，默认为false，表示并发安全。
func NewSortedStringArrayFrom(array []string, unsafe...bool) *SortedStringArray {
    a := NewSortedStringArraySize(0, unsafe...)
    a.array = array
    sort.Strings(a.array)
    return a
}

// Create an array from a copy of given slice <array>.
// The param <unsafe> used to specify whether using array with un-concurrent-safety,
// which is false in default, means concurrent-safe in default.
//
// 通过给定的slice拷贝创建数组对象，参数unsafe用于指定是否用于非并发安全场景，默认为false，表示并发安全。
func NewSortedStringArrayFromCopy(array []string, unsafe...bool) *SortedStringArray {
    newArray := make([]string, len(array))
    copy(newArray, array)
    return &SortedStringArray{
        mu    : rwmutex.New(unsafe...),
        array : newArray,
    }
}

// Set the underlying slice array with the given <array> param.
//
// 设置底层数组变量.
func (a *SortedStringArray) SetArray(array []string) *SortedStringArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.array = array
    sort.Strings(a.array)
    return a
}

// Sort the array in increasing order.
//
// 将数组排序(默认从低到高).
func (a *SortedStringArray) Sort() *SortedStringArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    sort.Strings(a.array)
    return a
}

// And values to sorted array, the array always keeps sorted.
//
// 添加数据项.
func (a *SortedStringArray) Add(values...string) *SortedStringArray {
    if len(values) == 0 {
        return a
    }
    a.mu.Lock()
    defer a.mu.Unlock()
    for _, value := range values {
        index, cmp := a.binSearch(value, false)
        if a.unique.Val() && cmp == 0 {
            continue
        }
        if index < 0 {
            a.array = append(a.array, value)
            continue
        }
        // 加到指定索引后面
        if cmp > 0 {
            index++
        }
        rear   := append([]string{}, a.array[index : ]...)
        a.array = append(a.array[0 : index], value)
        a.array = append(a.array, rear...)
    }
    return a
}

// Get value by index.
//
// 获取指定索引的数据项, 调用方注意判断数组边界。
func (a *SortedStringArray) Get(index int) string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    value := a.array[index]
    return value
}

// Remove an item by index.
//
// 删除指定索引的数据项, 调用方注意判断数组边界。
func (a *SortedStringArray) Remove(index int) string {
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

// Push new items to the beginning of array.
//
// 将数据项添加到数组的最左端(索引为0)。
func (a *SortedStringArray) PopLeft() string {
    a.mu.Lock()
    defer a.mu.Unlock()
    value  := a.array[0]
    a.array = a.array[1 : ]
    return value
}

// Push new items to the end of array.
//
// 将数据项添加到数组的最右端(索引为length - 1)。
func (a *SortedStringArray) PopRight() string {
    a.mu.Lock()
    defer a.mu.Unlock()
    index  := len(a.array) - 1
    value  := a.array[index]
    a.array = a.array[: index]
    return value
}

// PopRand picks an random item out of array.
//
// 随机将一个数据项移出数组，并返回该数据项。
func (a *SortedStringArray) PopRand() string {
    return a.Remove(grand.Intn(len(a.array)))
}

// PopRands picks <size> items out of array.
//
// 随机将size个数据项移出数组，并返回该数据项。
func (a *SortedStringArray) PopRands(size int) []string {
    a.mu.Lock()
    defer a.mu.Unlock()
    if size > len(a.array) {
        size = len(a.array)
    }
    array := make([]string, size)
    for i := 0; i < size; i++ {
        index   := grand.Intn(len(a.array))
        array[i] = a.array[index]
        a.array  = append(a.array[ : index], a.array[index + 1 : ]...)
    }
    return array
}

// Pop <size> items from the beginning of array.
//
// 将最左端(首部)的size个数据项移出数组，并返回该数据项
func (a *SortedStringArray) PopLefts(size int) []string {
    a.mu.Lock()
    defer a.mu.Unlock()
    length := len(a.array)
    if size > length {
        size = length
    }
    value  := a.array[0 : size]
    a.array = a.array[size : ]
    return value
}

// Pop <size> items from the end of array.
//
// 将最右端(尾部)的size个数据项移出数组，并返回该数据项
func (a *SortedStringArray) PopRights(size int) []string {
    a.mu.Lock()
    defer a.mu.Unlock()
    index := len(a.array) - size
    if index < 0 {
        index = 0
    }
    value  := a.array[index :]
    a.array = a.array[ : index]
    return value
}

// Get items by range, returns array[start:end].
// Be aware that, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
//
// 将最右端(尾部)的size个数据项移出数组，并返回该数据项
func (a *SortedStringArray) Range(start, end int) []string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    length := len(a.array)
    if start > length || start > end {
        return nil
    }
    if start < 0 {
        start = 0
    }
    if end > length {
        end = length
    }
    array  := ([]string)(nil)
    if a.mu.IsSafe() {
        a.mu.RLock()
        defer a.mu.RUnlock()
        array = make([]string, end - start)
        copy(array, a.array[start : end])
    } else {
        array = a.array[start : end]
    }
    return array
}

// Calculate the sum of values in an array.
//
// 对数组中的元素项求和(将元素值转换为int类型后叠加)。
func (a *SortedStringArray) Sum() (sum int) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    for _, v := range a.array {
        sum += gconv.Int(v)
    }
    return
}

// Get the length of array.
//
// 数组长度。
func (a *SortedStringArray) Len() int {
    a.mu.RLock()
    length := len(a.array)
    a.mu.RUnlock()
    return length
}

// Get the underlying data of array.
// Be aware that, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
//
// 返回原始数据数组.
func (a *SortedStringArray) Slice() []string {
    array := ([]string)(nil)
    if a.mu.IsSafe() {
        a.mu.RLock()
        defer a.mu.RUnlock()
        array = make([]string, len(a.array))
        copy(array, a.array)
    } else {
        array = a.array
    }
    return array
}

// Check whether a value exists in the array.
//
// 查找指定数值是否存在。
func (a *SortedStringArray) Contains(value string) bool {
    return a.Search(value) == 0
}

// Search array by <value>, returns the index of <value>, returns -1 if not exists.
//
// 查找指定数值的索引位置，返回索引位置，如果查找不到则返回-1。
func (a *SortedStringArray) Search(value string) (index int) {
    index, _ = a.binSearch(value, true)
    return
}

// Binary search.
//
// 二分查找.
func (a *SortedStringArray) binSearch(value string, lock bool) (index int, result int) {
    if len(a.array) == 0 {
        return -1, -2
    }
    if lock {
        a.mu.RLock()
        defer a.mu.RUnlock()
    }
    min := 0
    max := len(a.array) - 1
    mid := 0
    cmp := -2
    for min <= max {
        mid = int((min + max) / 2)
        cmp = a.compareFunc(value, a.array[mid])
        switch {
            case cmp < 0 : max = mid - 1
            case cmp > 0 : min = mid + 1
            default :
                return mid, cmp
        }
    }
    return mid, cmp
}

// Set unique mark to the array,
// which means it does not contain any repeated items.
// It also do unique check, remove all repeated items.
//
// 设置是否允许数组唯一.
func (a *SortedStringArray) SetUnique(unique bool) *SortedStringArray {
    oldUnique := a.unique.Val()
    a.unique.Set(unique)
    if unique && oldUnique != unique {
        a.Unique()
    }
    return a
}

// Do unique check, remove all repeated items.
//
// 清理数组中重复的元素项.
func (a *SortedStringArray) Unique() *SortedStringArray {
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
    return a
}

// Return a new array, which is a copy of current array.
//
// 克隆当前数组，返回当前数组的一个拷贝。
func (a *SortedStringArray) Clone() (newArray *SortedStringArray) {
    a.mu.RLock()
    array := make([]string, len(a.array))
    copy(array, a.array)
    a.mu.RUnlock()
    return NewSortedStringArrayFrom(array, !a.mu.IsSafe())
}

// Clear array.
//
// 清空数据数组。
func (a *SortedStringArray) Clear() *SortedStringArray {
    a.mu.Lock()
    if len(a.array) > 0 {
        a.array = make([]string, 0)
    }
    a.mu.Unlock()
    return a
}

// Lock writing by callback function f.
//
// 使用自定义方法执行加锁修改操作。
func (a *SortedStringArray) LockFunc(f func(array []string)) *SortedStringArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    f(a.array)
    return a
}

// Lock reading by callback function f.
//
// 使用自定义方法执行加锁读取操作。
func (a *SortedStringArray) RLockFunc(f func(array []string)) *SortedStringArray {
    a.mu.RLock()
    defer a.mu.RUnlock()
    f(a.array)
    return a
}

// Merge two arrays. The parameter <array> can be any garray type or slice type.
// The difference between Merge and Add is Add supports only specified slice type,
// but Merge supports more variable types.
//
// 合并两个数组, 支持任意的garray数组类型及slice类型.
func (a *SortedStringArray) Merge(array interface{}) *SortedStringArray {
    switch v := array.(type) {
        case *Array:             a.Add(gconv.Strings(v.Slice())...)
        case *IntArray:          a.Add(gconv.Strings(v.Slice())...)
        case *StringArray:       a.Add(gconv.Strings(v.Slice())...)
        case *SortedArray:       a.Add(gconv.Strings(v.Slice())...)
        case *SortedIntArray:    a.Add(gconv.Strings(v.Slice())...)
        case *SortedStringArray: a.Add(gconv.Strings(v.Slice())...)
        default:
            a.Add(gconv.Strings(array)...)
    }
    return a
}

// Chunks an array into arrays with size elements.
// The last chunk may contain less than size elements.
//
// 将一个数组分割成多个数组，其中每个数组的单元数目由size决定。最后一个数组的单元数目可能会少于size个。
func (a *SortedStringArray) Chunk(size int) [][]string {
    if size < 1 {
        return nil
    }
    a.mu.RLock()
    defer a.mu.RUnlock()
    length := len(a.array)
    chunks := int(math.Ceil(float64(length) / float64(size)))
    var n [][]string
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

// Extract a slice of the array(If in concurrent safe usage,
// it returns a copy of the slice; else a pointer).
// It returns the sequence of elements from the array array as specified
// by the offset and length parameters.
//
// 返回根据offset和size参数所指定的数组中的一段序列。
func (a *SortedStringArray) SubSlice(offset, size int) []string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    if offset > len(a.array) {
        return nil
    }
    if offset + size > len(a.array) {
        size = len(a.array) - offset
    }
    if a.mu.IsSafe() {
        s := make([]string, size)
        copy(s, a.array[offset:])
        return s
    } else {
        return a.array[offset:]
    }
}

// Rand gets one random entry from array.
//
// 从数组中随机获得1个元素项(不删除)。
func (a *SortedStringArray) Rand() string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return a.array[grand.Intn(len(a.array))]
}

// Rands gets one or more random entries from array(a copy).
//
// 从数组中随机拷贝size个元素项，构成slice返回。
func (a *SortedStringArray) Rands(size int) []string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    if size > len(a.array) {
        size = len(a.array)
    }
    n := make([]string, size)
    for i, v := range grand.Perm(len(a.array)) {
        n[i] = a.array[v]
        if i == size - 1 {
            break
        }
    }
    return n
}

// Join array elements with a string.
//
// 使用glue字符串串连当前数组的元素项，构造成新的字符串返回。
func (a *SortedStringArray) Join(glue string) string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    buffer := bytes.NewBuffer(nil)
    for k, v := range a.array {
        buffer.WriteString(gconv.String(v))
        if k != len(a.array) - 1 {
            buffer.WriteString(glue)
        }
    }
    return buffer.String()
}
<<<<<<< HEAD
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
=======
// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
>>>>>>> upstream/master

package garray

import (
<<<<<<< HEAD
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
=======
    "bytes"
	"fmt"
	"github.com/gogf/gf/g/container/gtype"
    "github.com/gogf/gf/g/internal/rwmutex"
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/util/grand"
    "math"
    "sort"
    "strings"
)

// It's using increasing order in default.
type SortedStringArray struct {
    mu          *rwmutex.RWMutex
    array       []string
    unique      *gtype.Bool             // Whether enable unique feature(false)
    comparator func(v1, v2 string) int // Comparison function(it returns -1: v1 < v2; 0: v1 == v2; 1: v1 > v2)
}

// NewSortedStringArray creates and returns an empty sorted array.
// The param <unsafe> used to specify whether using array in un-concurrent-safety,
// which is false in default.
func NewSortedStringArray(unsafe...bool) *SortedStringArray {
    return NewSortedStringArraySize(0, unsafe...)
}

// NewSortedStringArraySize create and returns an sorted array with given size and cap.
// The param <unsafe> used to specify whether using array in un-concurrent-safety,
// which is false in default.
func NewSortedStringArraySize(cap int, unsafe...bool) *SortedStringArray {
    return &SortedStringArray {
        mu          : rwmutex.New(unsafe...),
        array       : make([]string, 0, cap),
        unique      : gtype.NewBool(),
        comparator : func(v1, v2 string) int {
            return strings.Compare(v1, v2)
        },
    }
}

// NewSortedStringArrayFrom creates and returns an sorted array with given slice <array>.
// The param <unsafe> used to specify whether using array in un-concurrent-safety,
// which is false in default.
func NewSortedStringArrayFrom(array []string, unsafe...bool) *SortedStringArray {
    a := NewSortedStringArraySize(0, unsafe...)
    a.array = array
    sort.Strings(a.array)
    return a
}

// NewSortedStringArrayFromCopy creates and returns an sorted array from a copy of given slice <array>.
// The param <unsafe> used to specify whether using array in un-concurrent-safety,
// which is false in default.
func NewSortedStringArrayFromCopy(array []string, unsafe...bool) *SortedStringArray {
    newArray := make([]string, len(array))
    copy(newArray, array)
    return &SortedStringArray{
        mu    : rwmutex.New(unsafe...),
        array : newArray,
    }
}

// SetArray sets the underlying slice array with the given <array>.
func (a *SortedStringArray) SetArray(array []string) *SortedStringArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.array = array
    sort.Strings(a.array)
    return a
}

// Sort sorts the array in increasing order.
// The param <reverse> controls whether sort
// in increasing order(default) or decreasing order.
func (a *SortedStringArray) Sort() *SortedStringArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    sort.Strings(a.array)
    return a
}

// Add adds one or multiple values to sorted array, the array always keeps sorted.
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
        if cmp > 0 {
            index++
        }
        rear   := append([]string{}, a.array[index : ]...)
        a.array = append(a.array[0 : index], value)
        a.array = append(a.array, rear...)
    }
    return a
}

// Get returns the value of the specified index,
// the caller should notice the boundary of the array.
func (a *SortedStringArray) Get(index int) string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    value := a.array[index]
    return value
}

// Remove removes an item by index.
func (a *SortedStringArray) Remove(index int) string {
    a.mu.Lock()
    defer a.mu.Unlock()
    // Determine array boundaries when deleting to improve deletion efficiency.
    if index == 0 {
        value  := a.array[0]
        a.array = a.array[1 : ]
        return value
    } else if index == len(a.array) - 1 {
        value  := a.array[index]
        a.array = a.array[: index]
        return value
    }
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
    value  := a.array[index]
    a.array = append(a.array[ : index], a.array[index + 1 : ]...)
    return value
}

// PopLeft pops and returns an item from the beginning of array.
func (a *SortedStringArray) PopLeft() string {
    a.mu.Lock()
    defer a.mu.Unlock()
    value  := a.array[0]
    a.array = a.array[1 : ]
    return value
}

// PopRight pops and returns an item from the end of array.
func (a *SortedStringArray) PopRight() string {
    a.mu.Lock()
    defer a.mu.Unlock()
    index  := len(a.array) - 1
    value  := a.array[index]
    a.array = a.array[: index]
    return value
}

// PopRand randomly pops and return an item out of array.
func (a *SortedStringArray) PopRand() string {
    return a.Remove(grand.Intn(len(a.array)))
}

// PopRands randomly pops and returns <size> items out of array.
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

// PopLefts pops and returns <size> items from the beginning of array.
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

// PopRights pops and returns <size> items from the end of array.
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

// Range picks and returns items by range, like array[start:end].
// Notice, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
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

// Sum returns the sum of values in an array.
func (a *SortedStringArray) Sum() (sum int) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    for _, v := range a.array {
        sum += gconv.Int(v)
    }
    return
}

// Len returns the length of array.
>>>>>>> upstream/master
func (a *SortedStringArray) Len() int {
    a.mu.RLock()
    length := len(a.array)
    a.mu.RUnlock()
    return length
}

<<<<<<< HEAD
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
=======
// Slice returns the underlying data of array.
// Notice, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
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

// Contains checks whether a value exists in the array.
func (a *SortedStringArray) Contains(value string) bool {
    return a.Search(value) == 0
}

// Search searches array by <value>, returns the index of <value>,
// or returns -1 if not exists.
func (a *SortedStringArray) Search(value string) (index int) {
    index, _ = a.binSearch(value, true)
    return
}

// Binary search.
func (a *SortedStringArray) binSearch(value string, lock bool) (index int, result int) {
    if len(a.array) == 0 {
        return -1, -2
    }
    if lock {
        a.mu.RLock()
        defer a.mu.RUnlock()
    }
>>>>>>> upstream/master
    min := 0
    max := len(a.array) - 1
    mid := 0
    cmp := -2
<<<<<<< HEAD
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
=======
    for min <= max {
        mid = int((min + max) / 2)
        cmp = a.comparator(value, a.array[mid])
        switch {
            case cmp < 0 : max = mid - 1
            case cmp > 0 : min = mid + 1
            default :
                return mid, cmp
        }
    }
    return mid, cmp
}

// SetUnique sets unique mark to the array,
// which means it does not contain any repeated items.
// It also do unique check, remove all repeated items.
func (a *SortedStringArray) SetUnique(unique bool) *SortedStringArray {
    oldUnique := a.unique.Val()
    a.unique.Set(unique)
    if unique && oldUnique != unique {
        a.Unique()
    }
    return a
}

// Unique uniques the array, clear repeated items.
func (a *SortedStringArray) Unique() *SortedStringArray {
>>>>>>> upstream/master
    a.mu.Lock()
    i := 0
    for {
        if i == len(a.array) - 1 {
            break
        }
<<<<<<< HEAD
        if a.compareFunc(a.array[i], a.array[i + 1]) == 0 {
=======
        if a.comparator(a.array[i], a.array[i + 1]) == 0 {
>>>>>>> upstream/master
            a.array = append(a.array[ : i + 1], a.array[i + 1 + 1 : ]...)
        } else {
            i++
        }
    }
    a.mu.Unlock()
<<<<<<< HEAD
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
=======
    return a
}

// Clone returns a new array, which is a copy of current array.
func (a *SortedStringArray) Clone() (newArray *SortedStringArray) {
    a.mu.RLock()
    array := make([]string, len(a.array))
    copy(array, a.array)
    a.mu.RUnlock()
    return NewSortedStringArrayFrom(array, !a.mu.IsSafe())
}

// Clear deletes all items of current array.
func (a *SortedStringArray) Clear() *SortedStringArray {
    a.mu.Lock()
    if len(a.array) > 0 {
        a.array = make([]string, 0)
    }
    a.mu.Unlock()
    return a
}

// LockFunc locks writing by callback function <f>.
func (a *SortedStringArray) LockFunc(f func(array []string)) *SortedStringArray {
    a.mu.Lock()
    defer a.mu.Unlock()
    f(a.array)
    return a
}

// RLockFunc locks reading by callback function <f>.
func (a *SortedStringArray) RLockFunc(f func(array []string)) *SortedStringArray {
    a.mu.RLock()
    defer a.mu.RUnlock()
    f(a.array)
    return a
}

// Merge merges <array> into current array.
// The parameter <array> can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
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

// Chunk splits an array into multiple arrays,
// the size of each array is determined by <size>.
// The last chunk may contain less than size elements.
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

// SubSlice returns a slice of elements from the array as specified
// by the <offset> and <size> parameters.
// If in concurrent safe usage, it returns a copy of the slice; else a pointer.
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

// Rand randomly returns one item from array(no deleting).
func (a *SortedStringArray) Rand() string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return a.array[grand.Intn(len(a.array))]
}

// Rands randomly returns <size> items from array(no deleting).
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

// Join joins array elements with a string <glue>.
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

// CountValues counts the number of occurrences of all values in the array.
func (a *SortedStringArray) CountValues() map[string]int {
	m := make(map[string]int)
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		m[v]++
	}
	return m
}

// String returns current array as a string.
func (a *SortedStringArray) String() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return fmt.Sprint(a.array)
>>>>>>> upstream/master
}
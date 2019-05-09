// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
    "bytes"
    "fmt"
    "github.com/gogf/gf/g/internal/rwmutex"
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/util/grand"
    "math"
    "sort"
)

type Array struct {
    mu    *rwmutex.RWMutex
    array []interface{}
}

// New creates and returns an empty array.
// The param <unsafe> used to specify whether using array in un-concurrent-safety,
// which is false in default.
func New(unsafe...bool) *Array {
    return NewArraySize(0, 0, unsafe...)
}

// See New.
func NewArray(unsafe...bool) *Array {
    return NewArraySize(0, 0, unsafe...)
}

// NewArraySize create and returns an array with given size and cap.
// The param <unsafe> used to specify whether using array in un-concurrent-safety,
// which is false in default.
func NewArraySize(size int, cap int, unsafe...bool) *Array {
    return &Array{
        mu    : rwmutex.New(unsafe...),
        array : make([]interface{}, size, cap),
    }
}

// See NewArrayFrom.
func NewFrom(array []interface{}, unsafe...bool) *Array {
    return NewArrayFrom(array, unsafe...)
}

// See NewArrayFromCopy.
func NewFromCopy(array []interface{}, unsafe...bool) *Array {
    return NewArrayFromCopy(array, unsafe...)
}

// NewArrayFrom creates and returns an array with given slice <array>.
// The param <unsafe> used to specify whether using array in un-concurrent-safety,
// which is false in default.
func NewArrayFrom(array []interface{}, unsafe...bool) *Array {
    return &Array{
        mu    : rwmutex.New(unsafe...),
        array : array,
    }
}

// NewArrayFromCopy creates and returns an array from a copy of given slice <array>.
// The param <unsafe> used to specify whether using array in un-concurrent-safety,
// which is false in default.
func NewArrayFromCopy(array []interface{}, unsafe...bool) *Array {
    newArray := make([]interface{}, len(array))
    copy(newArray, array)
    return &Array{
        mu    : rwmutex.New(unsafe...),
        array : newArray,
    }
}

// Get returns the value of the specified index,
// the caller should notice the boundary of the array.
func (a *Array) Get(index int) interface{} {
    a.mu.RLock()
    defer a.mu.RUnlock()
    value := a.array[index]
    return value
}

// Set sets value to specified index.
func (a *Array) Set(index int, value interface{}) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.array[index] = value
    return a
}

// SetArray sets the underlying slice array with the given <array>.
func (a *Array) SetArray(array []interface{}) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.array = array
    return a
}

// Replace replaces the array items by given <array> from the beginning of array.
func (a *Array) Replace(array []interface{}) *Array {
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

// Sum returns the sum of values in an array.
func (a *Array) Sum() (sum int) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    for _, v := range a.array {
        sum += gconv.Int(v)
    }
    return
}

// SortFunc sorts the array by custom function <less>.
func (a *Array) SortFunc(less func(v1, v2 interface{}) bool) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    sort.Slice(a.array, func(i, j int) bool {
        return less(a.array[i], a.array[j])
    })
    return a
}

// InsertBefore inserts the <value> to the front of <index>.
func (a *Array) InsertBefore(index int, value interface{}) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    rear   := append([]interface{}{}, a.array[index : ]...)
    a.array = append(a.array[0 : index], value)
    a.array = append(a.array, rear...)
    return a
}

// InsertAfter inserts the <value> to the back of <index>.
func (a *Array) InsertAfter(index int, value interface{}) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    rear   := append([]interface{}{}, a.array[index + 1 : ]...)
    a.array = append(a.array[0 : index + 1], value)
    a.array = append(a.array, rear...)
    return a
}

// Remove removes an item by index.
func (a *Array) Remove(index int) interface{} {
    a.mu.Lock()
    defer a.mu.Unlock()
	// Determine array boundaries when deleting to improve deletion efficiency。
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

// PushLeft pushes one or multiple items to the beginning of array.
func (a *Array) PushLeft(value...interface{}) *Array {
    a.mu.Lock()
    a.array = append(value, a.array...)
    a.mu.Unlock()
    return a
}

// PushRight pushes one or multiple items to the end of array.
// It equals to Append.
func (a *Array) PushRight(value...interface{}) *Array {
    a.mu.Lock()
    a.array = append(a.array, value...)
    a.mu.Unlock()
    return a
}

// PopRand randomly pops and return an item out of array.
func (a *Array) PopRand() interface{} {
    return a.Remove(grand.Intn(len(a.array)))
}

// PopRands randomly pops and returns <size> items out of array.
func (a *Array) PopRands(size int) []interface{} {
    a.mu.Lock()
    defer a.mu.Unlock()
    if size > len(a.array) {
        size = len(a.array)
    }
    array := make([]interface{}, size)
    for i := 0; i < size; i++ {
        index   := grand.Intn(len(a.array))
        array[i] = a.array[index]
        a.array  = append(a.array[ : index], a.array[index + 1 : ]...)
    }
    return array
}

// PopLeft pops and returns an item from the beginning of array.
func (a *Array) PopLeft() interface{} {
    a.mu.Lock()
    defer a.mu.Unlock()
    value  := a.array[0]
    a.array = a.array[1 : ]
    return value
}

// PopRight pops and returns an item from the end of array.
func (a *Array) PopRight() interface{} {
    a.mu.Lock()
    defer a.mu.Unlock()
    index  := len(a.array) - 1
    value  := a.array[index]
    a.array = a.array[: index]
    return value
}

// PopLefts pops and returns <size> items from the beginning of array.
func (a *Array) PopLefts(size int) []interface{} {
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
func (a *Array) PopRights(size int) []interface{} {
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
func (a *Array) Range(start, end int) []interface{} {
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
    array  := ([]interface{})(nil)
    if a.mu.IsSafe() {
        a.mu.RLock()
        defer a.mu.RUnlock()
        array = make([]interface{}, end - start)
        copy(array, a.array[start : end])
    } else {
        array = a.array[start : end]
    }
    return array
}

// See PushRight.
func (a *Array) Append(value...interface{}) *Array {
    a.PushRight(value...)
    return a
}

// Len returns the length of array.
func (a *Array) Len() int {
    a.mu.RLock()
    length := len(a.array)
    a.mu.RUnlock()
    return length
}

// Slice returns the underlying data of array.
// Notice, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
func (a *Array) Slice() []interface{} {
    array := ([]interface{})(nil)
    if a.mu.IsSafe() {
        a.mu.RLock()
        defer a.mu.RUnlock()
        array = make([]interface{}, len(a.array))
        copy(array, a.array)
    } else {
        array = a.array
    }
    return array
}

// Clone returns a new array, which is a copy of current array.
func (a *Array) Clone() (newArray *Array) {
    a.mu.RLock()
    array := make([]interface{}, len(a.array))
    copy(array, a.array)
    a.mu.RUnlock()
    return NewArrayFrom(array, !a.mu.IsSafe())
}

// Clear deletes all items of current array.
func (a *Array) Clear() *Array {
    a.mu.Lock()
    if len(a.array) > 0 {
        a.array = make([]interface{}, 0)
    }
    a.mu.Unlock()
    return a
}

// Contains checks whether a value exists in the array.
func (a *Array) Contains(value interface{}) bool {
    return a.Search(value) != -1
}

// Search searches array by <value>, returns the index of <value>,
// or returns -1 if not exists.
func (a *Array) Search(value interface{}) int {
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

// Unique uniques the array, clear repeated items.
func (a *Array) Unique() *Array {
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

// LockFunc locks writing by callback function <f>.
func (a *Array) LockFunc(f func(array []interface{})) *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    f(a.array)
    return a
}

// RLockFunc locks reading by callback function <f>.
func (a *Array) RLockFunc(f func(array []interface{})) *Array {
    a.mu.RLock()
    defer a.mu.RUnlock()
    f(a.array)
    return a
}

// Merge merges <array> into current array.
// The parameter <array> can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (a *Array) Merge(array interface{}) *Array {
    switch v := array.(type) {
        case *Array:             a.Append(gconv.Interfaces(v.Slice())...)
        case *IntArray:          a.Append(gconv.Interfaces(v.Slice())...)
        case *StringArray:       a.Append(gconv.Interfaces(v.Slice())...)
        case *SortedArray:       a.Append(gconv.Interfaces(v.Slice())...)
        case *SortedIntArray:    a.Append(gconv.Interfaces(v.Slice())...)
        case *SortedStringArray: a.Append(gconv.Interfaces(v.Slice())...)
        default:
            a.Append(gconv.Interfaces(array)...)
    }
    return a
}

// Fill fills an array with num entries of the value <value>,
// keys starting at the <startIndex> parameter.
func (a *Array) Fill(startIndex int, num int, value interface{}) *Array {
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

// Chunk splits an array into multiple arrays,
// the size of each array is determined by <size>.
// The last chunk may contain less than size elements.
func (a *Array) Chunk(size int) [][]interface{} {
    if size < 1 {
        return nil
    }
    a.mu.RLock()
    defer a.mu.RUnlock()
    length := len(a.array)
    chunks := int(math.Ceil(float64(length) / float64(size)))
    var n [][]interface{}
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

// Pad pads array to the specified length with <value>.
// If size is positive then the array is padded on the right, or negative on the left.
// If the absolute value of <size> is less than or equal to the length of the array
// then no padding takes place.
func (a *Array) Pad(size int, val interface{}) *Array {
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
    tmp := make([]interface{}, n)
    for i := 0; i < n; i++ {
        tmp[i] = val
    }
    if size > 0 {
        a.array = append(a.array, tmp...)
    } else {
        a.array = append(tmp, a.array...)
    }
    return a
}

// SubSlice returns a slice of elements from the array as specified
// by the <offset> and <size> parameters.
// If in concurrent safe usage, it returns a copy of the slice; else a pointer.
func (a *Array) SubSlice(offset, size int) []interface{} {
    a.mu.RLock()
    defer a.mu.RUnlock()
    if offset > len(a.array) {
        return nil
    }
    if offset + size > len(a.array) {
        size = len(a.array) - offset
    }
    if a.mu.IsSafe() {
        s := make([]interface{}, size)
        copy(s, a.array[offset:])
        return s
    } else {
        return a.array[offset:]
    }
}

// Rand randomly returns one item from array(no deleting).
func (a *Array) Rand() interface{} {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return a.array[grand.Intn(len(a.array))]
}

// Rands randomly returns <size> items from array(no deleting).
func (a *Array) Rands(size int) []interface{} {
    a.mu.RLock()
    defer a.mu.RUnlock()
    if size > len(a.array) {
        size = len(a.array)
    }
    n := make([]interface{}, size)
    for i, v := range grand.Perm(len(a.array)) {
        n[i] = a.array[v]
        if i == size - 1 {
            break
        }
    }
    return n
}

// Shuffle randomly shuffles the array.
func (a *Array) Shuffle() *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    for i, v := range grand.Perm(len(a.array)) {
        a.array[i], a.array[v] = a.array[v], a.array[i]
    }
    return a
}

// Reverse makes array with elements in reverse order.
func (a *Array) Reverse() *Array {
    a.mu.Lock()
    defer a.mu.Unlock()
    for i, j := 0, len(a.array) - 1; i < j; i, j = i + 1, j - 1 {
        a.array[i], a.array[j] = a.array[j], a.array[i]
    }
    return a
}

// Join joins array elements with a string <glue>.
func (a *Array) Join(glue string) string {
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
func (a *Array) CountValues() map[interface{}]int {
    m := make(map[interface{}]int)
    a.mu.RLock()
    defer a.mu.RUnlock()
    for _, v := range a.array {
        m[v]++
    }
    return m
}

// String returns current array as a string.
func (a *Array) String() string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return fmt.Sprint(a.array)
}
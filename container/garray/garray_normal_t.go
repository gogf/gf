// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
	"github.com/gogf/gf/v2/util/gconv"
)

// TArray is a golang array with rich features.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
// TArray is a wrapper of Array. It is designed to make using Array more convenient.
type TArray[T comparable] struct {
	Array
}

// NewTArray creates and returns an empty array.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewTArray[T comparable](safe ...bool) *TArray[T] {
	return &TArray[T]{
		Array: *NewArray(safe...),
	}
}

// NewTArraySize create and returns an array with given size and cap.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewTArraySize[T comparable](size int, cap int, safe ...bool) *TArray[T] {
	arr := NewArraySize(size, cap, safe...)
	ret := &TArray[T]{
		Array: *arr,
	}
	return ret
}

// NewTArrayRange creates and returns an array by a range from `start` to `end`
// with step value `step`.
func NewTArrayRange[T comparable](start, end, step int, safe ...bool) *TArray[T] {
	return &TArray[T]{
		Array: *NewArrayRange(start, end, step, safe...),
	}
}

// NewTArrayFrom creates and returns an array with given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewTArrayFrom[T comparable](array []T, safe ...bool) *TArray[T] {
	return &TArray[T]{
		Array: *NewArrayFrom(tToAnySlice(array), safe...),
	}
}

// NewTArrayFromCopy creates and returns an array from a copy of given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewTArrayFromCopy[T comparable](array []T, safe ...bool) *TArray[T] {
	return &TArray[T]{
		Array: *NewArrayFromCopy(tToAnySlice(array), safe...),
	}
}

// At returns the value by the specified index.
// If the given `index` is out of range of the array, it returns `nil`.
func (a *TArray[T]) At(index int) (value T) {
	value, _ = a.Array.At(index).(T)
	return
}

// Get returns the value by the specified index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *TArray[T]) Get(index int) (value T, found bool) {
	val, found := a.Array.Get(index)
	if !found {
		return
	}
	value, _ = val.(T)
	return
}

// Set sets value to specified index.
func (a *TArray[T]) Set(index int, value T) error {
	return a.Array.Set(index, value)
}

// SetArray sets the underlying slice array with the given `array`.
func (a *TArray[T]) SetArray(array []T) *TArray[T] {
	a.Array.SetArray(tToAnySlice(array))
	return a
}

// Replace replaces the array items by given `array` from the beginning of array.
func (a *TArray[T]) Replace(array []T) *TArray[T] {
	a.Array.Replace(tToAnySlice(array))
	return a
}

// Sum returns the sum of values in an array.
func (a *TArray[T]) Sum() int {
	return a.Array.Sum()
}

// SortFunc sorts the array by custom function `less`.
func (a *TArray[T]) SortFunc(less func(v1, v2 T) bool) *TArray[T] {
	a.Array.SortFunc(func(v1, v2 any) bool {
		v1t, _ := v1.(T)
		v2t, _ := v2.(T)
		return less(v1t, v2t)
	})
	return a
}

// InsertBefore inserts the `values` to the front of `index`.
func (a *TArray[T]) InsertBefore(index int, values ...T) error {
	return a.Array.InsertBefore(index, tToAnySlice(values)...)
}

// InsertAfter inserts the `values` to the back of `index`.
func (a *TArray[T]) InsertAfter(index int, values ...T) error {
	return a.Array.InsertAfter(index, tToAnySlice(values)...)
}

// Remove removes an item by index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *TArray[T]) Remove(index int) (value T, found bool) {
	val, found := a.Array.Remove(index)
	if !found {
		return
	}
	value, _ = val.(T)
	return
}

// RemoveValue removes an item by value.
// It returns true if value is found in the array, or else false if not found.
func (a *TArray[T]) RemoveValue(value T) bool {
	return a.Array.RemoveValue(value)
}

// RemoveValues removes multiple items by `values`.
func (a TArray[T]) RemoveValues(values ...T) {
	a.Array.RemoveValues(tToAnySlice(values)...)
}

// PushLeft pushes one or multiple items to the beginning of array.
func (a *TArray[T]) PushLeft(value ...T) *TArray[T] {
	a.Array.PushLeft(tToAnySlice(value)...)
	return a
}

// PushRight pushes one or multiple items to the end of array.
// It equals to Append.
func (a *TArray[T]) PushRight(value ...T) *TArray[T] {
	a.Array.PushRight(tToAnySlice(value)...)
	return a
}

// PopRand randomly pops and return an item out of array.
// Note that if the array is empty, the `found` is false.
func (a *TArray[T]) PopRand() (value T, found bool) {
	val, found := a.Array.PopRand()
	if !found {
		return
	}
	value, _ = val.(T)
	return
}

// PopRands randomly pops and returns `size` items out of array.
func (a *TArray[T]) PopRands(size int) []T {
	return anyToTSlice[T](a.Array.PopRands(size))
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the `found` is false.
func (a *TArray[T]) PopLeft() (value T, found bool) {
	val, found := a.Array.PopLeft()
	if !found {
		return
	}
	value, _ = val.(T)
	return
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the `found` is false.
func (a *TArray[T]) PopRight() (value T, found bool) {
	val, found := a.Array.PopRight()
	if !found {
		return
	}
	value, _ = val.(T)
	return
}

// PopLefts pops and returns `size` items from the beginning of array.
func (a *TArray[T]) PopLefts(size int) []T {
	return anyToTSlice[T](a.Array.PopLefts(size))
}

// PopRights pops and returns `size` items from the end of array.
func (a *TArray[T]) PopRights(size int) []T {
	return anyToTSlice[T](a.Array.PopRights(size))
}

// Range picks and returns items by range, like array[start:end].
// Notice, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
//
// If `end` is negative, then the offset will start from the end of array.
// If `end` is omitted, then the sequence will have everything from start up
// until the end of the array.
func (a *TArray[T]) Range(start int, end ...int) []T {
	return anyToTSlice[T](a.Array.Range(start, end...))
}

// SubSlice returns a slice of elements from the array as specified
// by the `offset` and `size` parameters.
// If in concurrent safe usage, it returns a copy of the slice; else a pointer.
//
// If offset is non-negative, the sequence will start at that offset in the array.
// If offset is negative, the sequence will start that far from the end of the array.
//
// If length is given and is positive, then the sequence will have up to that many elements in it.
// If the array is shorter than the length, then only the available array elements will be present.
// If length is given and is negative then the sequence will stop that many elements from the end of the array.
// If it is omitted, then the sequence will have everything from offset up until the end of the array.
//
// Any possibility crossing the left border of array, it will fail.
func (a *TArray[T]) SubSlice(offset int, length ...int) []T {
	return anyToTSlice[T](a.Array.SubSlice(offset, length...))
}

// Append is alias of PushRight, please See PushRight.
func (a *TArray[T]) Append(value ...T) *TArray[T] {
	a.Array.Append(tToAnySlice(value)...)
	return a
}

// Len returns the length of array.
func (a *TArray[T]) Len() int {
	return a.Array.Len()
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *TArray[T]) Slice() []T {
	return anyToTSlice[T](a.Array.Slice())
}

// Interfaces returns current array as []any.
func (a *TArray[T]) Interfaces() []any {
	return a.Array.Interfaces()
}

// Clone returns a new array, which is a copy of current array.
func (a *TArray[T]) Clone() *TArray[T] {
	return &TArray[T]{
		Array: *a.Array.Clone(),
	}
}

// Clear deletes all items of current array.
func (a *TArray[T]) Clear() *TArray[T] {
	a.Array.Clear()
	return a
}

// Contains checks whether a value exists in the array.
func (a *TArray[T]) Contains(value T) bool {
	return a.Array.Contains(value)
}

// Search searches array by `value`, returns the index of `value`,
// or returns -1 if not exists.
func (a *TArray[T]) Search(value T) int {
	return a.Array.Search(value)
}

// Unique uniques the array, clear repeated items.
// Example: [1,1,2,3,2] -> [1,2,3]
func (a *TArray[T]) Unique() *TArray[T] {
	a.Array.Unique()
	return a
}

// LockFunc locks writing by callback function `f`.
func (a *TArray[T]) LockFunc(f func(array []T)) *TArray[T] {
	a.Array.LockFunc(func(array []any) {
		vals := anyToTSlice[T](array)
		f(vals)
		for k, v := range vals {
			array[k] = v
		}
	})
	return a
}

// RLockFunc locks reading by callback function `f`.
func (a *TArray[T]) RLockFunc(f func(array []T)) *TArray[T] {
	a.Array.RLockFunc(func(array []any) {
		f(anyToTSlice[T](array))
	})
	return a
}

// Merge merges `array` into current array.
// The parameter `array` can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (a *TArray[T]) Merge(array any) *TArray[T] {
	switch v := array.(type) {
	case *Array:
		return a.Merge(v.Slice())
	case *StrArray:
		return a.Merge(v.Slice())
	case *IntArray:
		return a.Merge(v.Slice())
	case *TArray[T]:
		a.Array.Merge(&v.Array)
	case []T:
		a.Array.Merge(v)
	case T:
		a.Append(v)
	case TArray[T]:
		a.Array.Merge(&v.Array)
	default:
		var vals []T
		if err := gconv.Scan(v, &vals); err != nil {
			panic(err)
		}
		a.Append(vals...)
	}
	return a
}

// Fill fills an array with num entries of the value `value`,
// keys starting at the `startIndex` parameter.
func (a *TArray[T]) Fill(startIndex int, num int, value T) error {
	return a.Array.Fill(startIndex, num, value)
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by `size`.
// The last chunk may contain less than size elements.
func (a *TArray[T]) Chunk(size int) (values [][]T) {
	return anyToTSlices[T](a.Array.Chunk(size))
}

// Pad pads array to the specified length with `value`.
// If size is positive then the array is padded on the right, or negative on the left.
// If the absolute value of `size` is less than or equal to the length of the array
// then no padding takes place.
func (a *TArray[T]) Pad(size int, val T) *TArray[T] {
	a.Array.Pad(size, val)
	return a
}

// Rand randomly returns one item from array(no deleting).
func (a *TArray[T]) Rand() (value T, found bool) {
	val, found := a.Array.Rand()
	if !found {
		return
	}
	value, _ = val.(T)
	return
}

// Rands randomly returns `size` items from array(no deleting).
func (a *TArray[T]) Rands(size int) []T {
	return anyToTSlice[T](a.Array.Rands(size))
}

// Shuffle randomly shuffles the array.
func (a *TArray[T]) Shuffle() *TArray[T] {
	a.Array.Shuffle()
	return a
}

// Reverse makes array with elements in reverse order.
func (a *TArray[T]) Reverse() *TArray[T] {
	a.Array.Reverse()
	return a
}

// Join joins array elements with a string `glue`.
func (a *TArray[T]) Join(glue string) string {
	return a.Array.Join(glue)
}

// CountValues counts the number of occurrences of all values in the array.
func (a *TArray[T]) CountValues() (valueCnt map[T]int) {
	valueCnt = map[T]int{}
	for k, v := range a.Array.CountValues() {
		k0, _ := k.(T)
		valueCnt[k0] = v
	}
	return
}

// Iterator is alias of IteratorAsc.
func (a *TArray[T]) Iterator(f func(k int, v T) bool) {
	a.Array.Iterator(func(k int, v any) bool {
		v0, _ := v.(T)
		return f(k, v0)
	})
}

// IteratorAsc iterates the array readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *TArray[T]) IteratorAsc(f func(k int, v T) bool) {
	a.Array.IteratorAsc(func(k int, v any) bool {
		v0, _ := v.(T)
		return f(k, v0)
	})
}

// IteratorDesc iterates the array readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *TArray[T]) IteratorDesc(f func(k int, v T) bool) {
	a.Array.IteratorDesc(func(k int, v any) bool {
		v0, _ := v.(T)
		return f(k, v0)
	})
}

// String returns current array as a string, which implements like json.Marshal does.
func (a *TArray[T]) String() string {
	return a.Array.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (a TArray[T]) MarshalJSON() ([]byte, error) {
	return a.Array.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (a *TArray[T]) UnmarshalJSON(b []byte) error {
	return a.Array.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for array.
func (a *TArray[T]) UnmarshalValue(value any) error {
	return a.Array.UnmarshalValue(value)
}

// Filter iterates array and filters elements using custom callback function.
// It removes the element from array if callback function `filter` returns true,
// it or else does nothing and continues iterating.
func (a *TArray[T]) Filter(filter func(index int, value T) bool) *TArray[T] {
	a.Array.Filter(func(index int, value any) bool {
		val, _ := value.(T)
		return filter(index, val)
	})
	return a
}

// FilterNil removes all nil value of the array.
func (a *TArray[T]) FilterNil() *TArray[T] {
	a.Array.FilterNil()
	return a
}

// FilterEmpty removes all empty value of the array.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (a *TArray[T]) FilterEmpty() *TArray[T] {
	a.Array.FilterEmpty()
	return a
}

// Walk applies a user supplied function `f` to every item of array.
func (a *TArray[T]) Walk(f func(value T) T) *TArray[T] {
	a.Array.Walk(func(value any) any {
		val, _ := value.(T)
		return f(val)
	})
	return a
}

// IsEmpty checks whether the array is empty.
func (a *TArray[T]) IsEmpty() bool {
	return a.Array.IsEmpty()
}

// DeepCopy implements interface for deep copy of current type.
func (a *TArray[T]) DeepCopy() any {
	if a == nil {
		return nil
	}
	arr := a.Array.DeepCopy().(*Array)
	return &TArray[T]{
		Array: *arr,
	}
}

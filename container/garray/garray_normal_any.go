// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
	"fmt"

	"github.com/gogf/gf/v2/util/gconv"
)

// Array is a golang array with rich features.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type Array struct {
	*TArray[any]
}

// New creates and returns an empty array.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func New(safe ...bool) *Array {
	return NewArraySize(0, 0, safe...)
}

// NewArray is alias of New, please see New.
func NewArray(safe ...bool) *Array {
	return NewArraySize(0, 0, safe...)
}

// NewArraySize create and returns an array with given size and cap.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewArraySize(size int, cap int, safe ...bool) *Array {
	return &Array{
		TArray: NewTArraySize[any](size, cap, safe...),
	}
}

// NewArrayRange creates and returns an array by a range from `start` to `end`
// with step value `step`.
func NewArrayRange(start, end, step int, safe ...bool) *Array {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]any, 0)
	index := 0
	for i := start; i <= end; i += step {
		slice = append(slice, i)
		index++
	}
	return NewArrayFrom(slice, safe...)
}

// NewFrom is alias of NewArrayFrom.
// See NewArrayFrom.
func NewFrom(array []any, safe ...bool) *Array {
	return NewArrayFrom(array, safe...)
}

// NewFromCopy is alias of NewArrayFromCopy.
// See NewArrayFromCopy.
func NewFromCopy(array []any, safe ...bool) *Array {
	return NewArrayFromCopy(array, safe...)
}

// NewArrayFrom creates and returns an array with given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewArrayFrom(array []any, safe ...bool) *Array {
	return &Array{
		TArray: NewTArrayFrom(array, safe...),
	}
}

// NewArrayFromCopy creates and returns an array from a copy of given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewArrayFromCopy(array []any, safe ...bool) *Array {
	newArray := make([]any, len(array))
	copy(newArray, array)
	return NewArrayFrom(newArray, safe...)
}

// lazyInit lazily initializes the array.
func (a *Array) lazyInit() {
	if a.TArray == nil {
		a.TArray = NewTArray[any](false)
	}
}

// At returns the value by the specified index.
// If the given `index` is out of range of the array, it returns `nil`.
func (a *Array) At(index int) (value any) {
	a.lazyInit()
	return a.TArray.At(index)
}

// Get returns the value by the specified index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *Array) Get(index int) (value any, found bool) {
	a.lazyInit()
	return a.TArray.Get(index)
}

// Set sets value to specified index.
func (a *Array) Set(index int, value any) error {
	a.lazyInit()
	return a.TArray.Set(index, value)
}

// SetArray sets the underlying slice array with the given `array`.
func (a *Array) SetArray(array []any) *Array {
	a.lazyInit()
	a.TArray.SetArray(array)
	return a
}

// Replace replaces the array items by given `array` from the beginning of array.
func (a *Array) Replace(array []any) *Array {
	a.lazyInit()
	a.TArray.Replace(array)
	return a
}

// Sum returns the sum of values in an array.
func (a *Array) Sum() (sum int) {
	a.lazyInit()
	return a.TArray.Sum()
}

// SortFunc sorts the array by custom function `less`.
func (a *Array) SortFunc(less func(v1, v2 any) bool) *Array {
	a.lazyInit()
	a.TArray.SortFunc(less)
	return a
}

// InsertBefore inserts the `values` to the front of `index`.
func (a *Array) InsertBefore(index int, values ...any) error {
	a.lazyInit()
	return a.TArray.InsertBefore(index, values...)
}

// InsertAfter inserts the `values` to the back of `index`.
func (a *Array) InsertAfter(index int, values ...any) error {
	a.lazyInit()
	return a.TArray.InsertAfter(index, values...)
}

// Remove removes an item by index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *Array) Remove(index int) (value any, found bool) {
	a.lazyInit()
	return a.TArray.Remove(index)
}

// RemoveValue removes an item by value.
// It returns true if value is found in the array, or else false if not found.
func (a *Array) RemoveValue(value any) bool {
	a.lazyInit()
	return a.TArray.RemoveValue(value)
}

// RemoveValues removes multiple items by `values`.
func (a *Array) RemoveValues(values ...any) {
	a.lazyInit()
	a.TArray.RemoveValues(values...)
}

// PushLeft pushes one or multiple items to the beginning of array.
func (a *Array) PushLeft(value ...any) *Array {
	a.lazyInit()
	a.TArray.PushLeft(value...)
	return a
}

// PushRight pushes one or multiple items to the end of array.
// It equals to Append.
func (a *Array) PushRight(value ...any) *Array {
	a.lazyInit()
	a.TArray.PushRight(value...)
	return a
}

// PopRand randomly pops and return an item out of array.
// Note that if the array is empty, the `found` is false.
func (a *Array) PopRand() (value any, found bool) {
	a.lazyInit()
	return a.TArray.PopRand()
}

// PopRands randomly pops and returns `size` items out of array.
func (a *Array) PopRands(size int) []any {
	a.lazyInit()
	return a.TArray.PopRands(size)
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the `found` is false.
func (a *Array) PopLeft() (value any, found bool) {
	a.lazyInit()
	return a.TArray.PopLeft()
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the `found` is false.
func (a *Array) PopRight() (value any, found bool) {
	a.lazyInit()
	return a.TArray.PopRight()
}

// PopLefts pops and returns `size` items from the beginning of array.
func (a *Array) PopLefts(size int) []any {
	a.lazyInit()
	return a.TArray.PopLefts(size)
}

// PopRights pops and returns `size` items from the end of array.
func (a *Array) PopRights(size int) []any {
	a.lazyInit()
	return a.TArray.PopRights(size)
}

// Range picks and returns items by range, like array[start:end].
// Notice, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
//
// If `end` is negative, then the offset will start from the end of array.
// If `end` is omitted, then the sequence will have everything from start up
// until the end of the array.
func (a *Array) Range(start int, end ...int) []any {
	a.lazyInit()
	return a.TArray.Range(start, end...)
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
func (a *Array) SubSlice(offset int, length ...int) []any {
	a.lazyInit()
	return a.TArray.SubSlice(offset, length...)
}

// Append is alias of PushRight, please See PushRight.
func (a *Array) Append(value ...any) *Array {
	a.lazyInit()
	a.TArray.Append(value...)
	return a
}

// Len returns the length of array.
func (a *Array) Len() int {
	a.lazyInit()
	return a.TArray.Len()
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *Array) Slice() []any {
	a.lazyInit()
	return a.TArray.Slice()
}

// Interfaces returns current array as []any.
func (a *Array) Interfaces() []any {
	return a.Slice()
}

// Clone returns a new array, which is a copy of current array.
func (a *Array) Clone() (newArray *Array) {
	a.lazyInit()
	return &Array{TArray: a.TArray.Clone()}
}

// Clear deletes all items of current array.
func (a *Array) Clear() *Array {
	a.lazyInit()
	a.TArray.Clear()
	return a
}

// Contains checks whether a value exists in the array.
func (a *Array) Contains(value any) bool {
	a.lazyInit()
	return a.TArray.Contains(value)
}

// Search searches array by `value`, returns the index of `value`,
// or returns -1 if not exists.
func (a *Array) Search(value any) int {
	a.lazyInit()
	return a.TArray.Search(value)
}

// Unique uniques the array, clear repeated items.
// Example: [1,1,2,3,2] -> [1,2,3]
func (a *Array) Unique() *Array {
	a.lazyInit()
	a.TArray.Unique()
	return a
}

// LockFunc locks writing by callback function `f`.
func (a *Array) LockFunc(f func(array []any)) *Array {
	a.lazyInit()
	a.TArray.LockFunc(f)
	return a
}

// RLockFunc locks reading by callback function `f`.
func (a *Array) RLockFunc(f func(array []any)) *Array {
	a.lazyInit()
	a.TArray.RLockFunc(f)
	return a
}

// Merge merges `array` into current array.
// The parameter `array` can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (a *Array) Merge(array any) *Array {
	a.lazyInit()
	return a.Append(gconv.Interfaces(array)...)
}

// Fill fills an array with num entries of the value `value`,
// keys starting at the `startIndex` parameter.
func (a *Array) Fill(startIndex int, num int, value any) error {
	a.lazyInit()
	return a.TArray.Fill(startIndex, num, value)
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by `size`.
// The last chunk may contain less than size elements.
func (a *Array) Chunk(size int) [][]any {
	a.lazyInit()
	return a.TArray.Chunk(size)
}

// Pad pads array to the specified length with `value`.
// If size is positive then the array is padded on the right, or negative on the left.
// If the absolute value of `size` is less than or equal to the length of the array
// then no padding takes place.
func (a *Array) Pad(size int, val any) *Array {
	a.lazyInit()
	a.TArray.Pad(size, val)
	return a
}

// Rand randomly returns one item from array(no deleting).
func (a *Array) Rand() (value any, found bool) {
	a.lazyInit()
	return a.TArray.Rand()
}

// Rands randomly returns `size` items from array(no deleting).
func (a *Array) Rands(size int) []any {
	a.lazyInit()
	return a.TArray.Rands(size)
}

// Shuffle randomly shuffles the array.
func (a *Array) Shuffle() *Array {
	a.lazyInit()
	a.TArray.Shuffle()
	return a
}

// Reverse makes array with elements in reverse order.
func (a *Array) Reverse() *Array {
	a.lazyInit()
	a.TArray.Reverse()
	return a
}

// Join joins array elements with a string `glue`.
func (a *Array) Join(glue string) string {
	a.lazyInit()
	return a.TArray.Join(glue)
}

// CountValues counts the number of occurrences of all values in the array.
func (a *Array) CountValues() map[any]int {
	a.lazyInit()
	return a.TArray.CountValues()
}

// Iterator is alias of IteratorAsc.
func (a *Array) Iterator(f func(k int, v any) bool) {
	a.IteratorAsc(f)
}

// IteratorAsc iterates the array readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *Array) IteratorAsc(f func(k int, v any) bool) {
	a.lazyInit()
	a.TArray.IteratorAsc(f)
}

// IteratorDesc iterates the array readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *Array) IteratorDesc(f func(k int, v any) bool) {
	a.lazyInit()
	a.TArray.IteratorDesc(f)
}

// String returns current array as a string, which implements like json.Marshal does.
func (a *Array) String() string {
	if a == nil {
		return ""
	}
	a.lazyInit()
	return a.TArray.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (a Array) MarshalJSON() ([]byte, error) {
	a.lazyInit()
	return a.TArray.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (a *Array) UnmarshalJSON(b []byte) error {
	a.lazyInit()
	return a.TArray.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for array.
func (a *Array) UnmarshalValue(value any) error {
	a.lazyInit()
	return a.TArray.UnmarshalValue(value)
}

// Filter iterates array and filters elements using custom callback function.
// It removes the element from array if callback function `filter` returns true,
// it or else does nothing and continues iterating.
func (a *Array) Filter(filter func(index int, value any) bool) *Array {
	a.lazyInit()
	a.TArray.Filter(filter)
	return a
}

// FilterNil removes all nil value of the array.
func (a *Array) FilterNil() *Array {
	a.lazyInit()
	a.TArray.FilterNil()
	return a
}

// FilterEmpty removes all empty value of the array.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (a *Array) FilterEmpty() *Array {
	a.lazyInit()
	a.TArray.FilterEmpty()
	return a
}

// Walk applies a user supplied function `f` to every item of array.
func (a *Array) Walk(f func(value any) any) *Array {
	a.lazyInit()
	a.TArray.Walk(f)
	return a
}

// IsEmpty checks whether the array is empty.
func (a *Array) IsEmpty() bool {
	a.lazyInit()
	return a.TArray.IsEmpty()
}

// DeepCopy implements interface for deep copy of current type.
func (a *Array) DeepCopy() any {
	if a == nil {
		return nil
	}
	a.lazyInit()
	return &Array{
		TArray: a.TArray.DeepCopy().(*TArray[any]),
	}
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
	"fmt"
	"sort"
	"sync"

	"github.com/gogf/gf/v2/util/gconv"
)

// SortedArray is a golang sorted array with rich features.
// It is using increasing order in default, which can be changed by
// setting it a custom comparator.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type SortedArray struct {
	*SortedTArray[any]
	once sync.Once
}

// lazyInit lazily initializes the array.
func (a *SortedArray) lazyInit() {
	a.once.Do(func() {
		if a.SortedTArray == nil {
			a.SortedTArray = NewSortedTArraySize[any](0, nil, false)
		}
	})
}

// NewSortedArray creates and returns an empty sorted array.
// The parameter `safe` is used to specify whether using array in concurrent-safety, which is false in default.
// The parameter `comparator` used to compare values to sort in array,
// if it returns value < 0, means `a` < `b`; the `a` will be inserted before `b`;
// if it returns value = 0, means `a` = `b`; the `a` will be replaced by     `b`;
// if it returns value > 0, means `a` > `b`; the `a` will be inserted after  `b`;
func NewSortedArray(comparator func(a, b any) int, safe ...bool) *SortedArray {
	return NewSortedArraySize(0, comparator, safe...)
}

// NewSortedArraySize create and returns an sorted array with given size and cap.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedArraySize(cap int, comparator func(a, b any) int, safe ...bool) *SortedArray {
	return &SortedArray{
		SortedTArray: NewSortedTArraySize(cap, comparator, safe...),
	}
}

// NewSortedArrayRange creates and returns an array by a range from `start` to `end`
// with step value `step`.
func NewSortedArrayRange(start, end, step int, comparator func(a, b any) int, safe ...bool) *SortedArray {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]any, 0)
	index := 0
	for i := start; i <= end; i += step {
		slice = append(slice, i)
		index++
	}
	return NewSortedArrayFrom(slice, comparator, safe...)
}

// NewSortedArrayFrom creates and returns an sorted array with given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedArrayFrom(array []any, comparator func(a, b any) int, safe ...bool) *SortedArray {
	a := NewSortedArraySize(0, comparator, safe...)
	a.array = array
	sort.Slice(a.array, func(i, j int) bool {
		return a.getComparator()(a.array[i], a.array[j]) < 0
	})
	return a
}

// NewSortedArrayFromCopy creates and returns an sorted array from a copy of given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedArrayFromCopy(array []any, comparator func(a, b any) int, safe ...bool) *SortedArray {
	newArray := make([]any, len(array))
	copy(newArray, array)
	return NewSortedArrayFrom(newArray, comparator, safe...)
}

// At returns the value by the specified index.
// If the given `index` is out of range of the array, it returns `nil`.
func (a *SortedArray) At(index int) (value any) {
	a.lazyInit()
	return a.SortedTArray.At(index)
}

// SetArray sets the underlying slice array with the given `array`.
func (a *SortedArray) SetArray(array []any) *SortedArray {
	a.lazyInit()
	a.SortedTArray.SetArray(array)
	return a
}

// SetComparator sets/changes the comparator for sorting.
// It resorts the array as the comparator is changed.
func (a *SortedArray) SetComparator(comparator func(a, b any) int) {
	a.lazyInit()
	a.SortedTArray.SetComparator(comparator)
}

// Sort sorts the array in increasing order.
// The parameter `reverse` controls whether sort
// in increasing order(default) or decreasing order
func (a *SortedArray) Sort() *SortedArray {
	a.lazyInit()
	a.SortedTArray.Sort()
	return a
}

// Add adds one or multiple values to sorted array, the array always keeps sorted.
// It's alias of function Append, see Append.
func (a *SortedArray) Add(values ...any) *SortedArray {
	a.lazyInit()
	a.SortedTArray.Add(values...)
	return a
}

// Append adds one or multiple values to sorted array, the array always keeps sorted.
func (a *SortedArray) Append(values ...any) *SortedArray {
	a.SortedTArray.Append(values...)
	return a
}

// Get returns the value by the specified index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *SortedArray) Get(index int) (value any, found bool) {
	a.lazyInit()
	return a.SortedTArray.Get(index)
}

// Remove removes an item by index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *SortedArray) Remove(index int) (value any, found bool) {
	a.lazyInit()
	return a.SortedTArray.Remove(index)
}

// RemoveValue removes an item by value.
// It returns true if value is found in the array, or else false if not found.
func (a *SortedArray) RemoveValue(value any) bool {
	a.lazyInit()
	return a.SortedTArray.RemoveValue(value)
}

// RemoveValues removes an item by `values`.
func (a *SortedArray) RemoveValues(values ...any) {
	a.lazyInit()
	a.SortedTArray.RemoveValues(values...)
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the `found` is false.
func (a *SortedArray) PopLeft() (value any, found bool) {
	a.lazyInit()
	return a.SortedTArray.PopLeft()
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the `found` is false.
func (a *SortedArray) PopRight() (value any, found bool) {
	a.lazyInit()
	return a.SortedTArray.PopRight()
}

// PopRand randomly pops and return an item out of array.
// Note that if the array is empty, the `found` is false.
func (a *SortedArray) PopRand() (value any, found bool) {
	a.lazyInit()
	return a.SortedTArray.PopRand()
}

// PopRands randomly pops and returns `size` items out of array.
func (a *SortedArray) PopRands(size int) []any {
	a.lazyInit()
	return a.SortedTArray.PopRands(size)
}

// PopLefts pops and returns `size` items from the beginning of array.
func (a *SortedArray) PopLefts(size int) []any {
	a.lazyInit()
	return a.SortedTArray.PopLefts(size)
}

// PopRights pops and returns `size` items from the end of array.
func (a *SortedArray) PopRights(size int) []any {
	a.lazyInit()
	return a.SortedTArray.PopRights(size)
}

// Range picks and returns items by range, like array[start:end].
// Notice, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
//
// If `end` is negative, then the offset will start from the end of array.
// If `end` is omitted, then the sequence will have everything from start up
// until the end of the array.
func (a *SortedArray) Range(start int, end ...int) []any {
	return a.SortedTArray.Range(start, end...)
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
func (a *SortedArray) SubSlice(offset int, length ...int) []any {
	a.lazyInit()
	return a.SortedTArray.SubSlice(offset, length...)
}

// Sum returns the sum of values in an array.
func (a *SortedArray) Sum() (sum int) {
	a.lazyInit()
	return a.SortedTArray.Sum()
}

// Len returns the length of array.
func (a *SortedArray) Len() int {
	a.lazyInit()
	return a.SortedTArray.Len()
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *SortedArray) Slice() []any {
	a.lazyInit()
	return a.SortedTArray.Slice()
}

// Interfaces returns current array as []any.
func (a *SortedArray) Interfaces() []any {
	a.lazyInit()
	return a.SortedTArray.Interfaces()
}

// Contains checks whether a value exists in the array.
func (a *SortedArray) Contains(value any) bool {
	a.lazyInit()
	return a.SortedTArray.Contains(value)
}

// Search searches array by `value`, returns the index of `value`,
// or returns -1 if not exists.
func (a *SortedArray) Search(value any) (index int) {
	a.lazyInit()
	return a.SortedTArray.Search(value)
}

// SetUnique sets unique mark to the array,
// which means it does not contain any repeated items.
// It also does unique check, remove all repeated items.
func (a *SortedArray) SetUnique(unique bool) *SortedArray {
	a.lazyInit()
	a.SortedTArray.SetUnique(unique)
	return a
}

// Unique uniques the array, clear repeated items.
func (a *SortedArray) Unique() *SortedArray {
	a.lazyInit()
	a.SortedTArray.Unique()
	return a
}

// Clone returns a new array, which is a copy of current array.
func (a *SortedArray) Clone() (newArray *SortedArray) {
	a.lazyInit()
	return &SortedArray{
		SortedTArray: a.SortedTArray.Clone(),
	}
}

// Clear deletes all items of current array.
func (a *SortedArray) Clear() *SortedArray {
	a.lazyInit()
	a.SortedTArray.Clear()
	return a
}

// LockFunc locks writing by callback function `f`.
func (a *SortedArray) LockFunc(f func(array []any)) *SortedArray {
	a.lazyInit()
	a.SortedTArray.LockFunc(f)
	return a
}

// RLockFunc locks reading by callback function `f`.
func (a *SortedArray) RLockFunc(f func(array []any)) *SortedArray {
	a.lazyInit()
	a.SortedTArray.RLockFunc(f)
	return a
}

// Merge merges `array` into current array.
// The parameter `array` can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (a *SortedArray) Merge(array any) *SortedArray {
	return a.Add(gconv.Interfaces(array)...)
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by `size`.
// The last chunk may contain less than size elements.
func (a *SortedArray) Chunk(size int) [][]any {
	a.lazyInit()
	return a.SortedTArray.Chunk(size)
}

// Rand randomly returns one item from array(no deleting).
func (a *SortedArray) Rand() (value any, found bool) {
	a.lazyInit()
	return a.SortedTArray.Rand()
}

// Rands randomly returns `size` items from array(no deleting).
func (a *SortedArray) Rands(size int) []any {
	a.lazyInit()
	return a.SortedTArray.Rands(size)
}

// Join joins array elements with a string `glue`.
func (a *SortedArray) Join(glue string) string {
	a.lazyInit()
	return a.SortedTArray.Join(glue)
}

// CountValues counts the number of occurrences of all values in the array.
func (a *SortedArray) CountValues() map[any]int {
	a.lazyInit()
	return a.SortedTArray.CountValues()
}

// Iterator is alias of IteratorAsc.
func (a *SortedArray) Iterator(f func(k int, v any) bool) {
	a.lazyInit()
	a.SortedTArray.Iterator(f)
}

// IteratorAsc iterates the array readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *SortedArray) IteratorAsc(f func(k int, v any) bool) {
	a.lazyInit()
	a.SortedTArray.IteratorAsc(f)
}

// IteratorDesc iterates the array readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *SortedArray) IteratorDesc(f func(k int, v any) bool) {
	a.lazyInit()
	a.SortedTArray.IteratorDesc(f)
}

// String returns current array as a string, which implements like json.Marshal does.
func (a *SortedArray) String() string {
	if a == nil {
		return ""
	}
	a.lazyInit()
	return a.SortedTArray.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (a SortedArray) MarshalJSON() ([]byte, error) {
	a.lazyInit()
	return a.SortedTArray.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
// Note that the comparator is set as string comparator in default.
func (a *SortedArray) UnmarshalJSON(b []byte) error {
	a.lazyInit()
	return a.SortedTArray.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for array.
// Note that the comparator is set as string comparator in default.
func (a *SortedArray) UnmarshalValue(value any) (err error) {
	a.lazyInit()
	return a.SortedTArray.UnmarshalValue(value)
}

// FilterNil removes all nil value of the array.
func (a *SortedArray) FilterNil() *SortedArray {
	a.lazyInit()
	a.SortedTArray.FilterNil()
	return a
}

// Filter iterates array and filters elements using custom callback function.
// It removes the element from array if callback function `filter` returns true,
// it or else does nothing and continues iterating.
func (a *SortedArray) Filter(filter func(index int, value any) bool) *SortedArray {
	a.lazyInit()
	a.SortedTArray.Filter(filter)
	return a
}

// FilterEmpty removes all empty value of the array.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (a *SortedArray) FilterEmpty() *SortedArray {
	a.lazyInit()
	a.SortedTArray.FilterEmpty()
	return a
}

// Walk applies a user supplied function `f` to every item of array.
func (a *SortedArray) Walk(f func(value any) any) *SortedArray {
	a.lazyInit()
	a.SortedTArray.Walk(f)
	return a
}

// IsEmpty checks whether the array is empty.
func (a *SortedArray) IsEmpty() bool {
	a.lazyInit()
	return a.SortedTArray.IsEmpty()
}

// DeepCopy implements interface for deep copy of current type.
func (a *SortedArray) DeepCopy() any {
	a.lazyInit()
	return &SortedArray{
		SortedTArray: a.SortedTArray.DeepCopy().(*SortedTArray[any]),
	}
}

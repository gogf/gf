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

// SortedIntArray is a golang sorted int array with rich features.
// It is using increasing order in default, which can be changed by
// setting it a custom comparator.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type SortedIntArray struct {
	*SortedTArray[int]
}

// lazyInit lazily initializes the array.
func (a *SortedIntArray) lazyInit() {
	if a.SortedTArray == nil {
		a.SortedTArray = NewSortedTArraySize(0, defaultComparatorInt, false)
		a.SetSorter(quickSortInt)
	}
}

// NewSortedIntArray creates and returns an empty sorted array.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArray(safe ...bool) *SortedIntArray {
	return NewSortedIntArraySize(0, safe...)
}

// NewSortedIntArrayComparator creates and returns an empty sorted array with specified comparator.
// The parameter `safe` is used to specify whether using array in concurrent-safety which is false in default.
func NewSortedIntArrayComparator(comparator func(a, b int) int, safe ...bool) *SortedIntArray {
	array := NewSortedIntArray(safe...)
	array.comparator = comparator
	return array
}

// NewSortedIntArraySize create and returns an sorted array with given size and cap.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArraySize(cap int, safe ...bool) *SortedIntArray {
	a := NewSortedTArraySize(cap, defaultComparatorInt, safe...)
	a.SetSorter(quickSortInt)
	return &SortedIntArray{
		SortedTArray: a,
	}
}

// NewSortedIntArrayRange creates and returns an array by a range from `start` to `end`
// with step value `step`.
func NewSortedIntArrayRange(start, end, step int, safe ...bool) *SortedIntArray {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]int, 0)
	index := 0
	for i := start; i <= end; i += step {
		slice = append(slice, i)
		index++
	}
	return NewSortedIntArrayFrom(slice, safe...)
}

// NewSortedIntArrayFrom creates and returns an sorted array with given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArrayFrom(array []int, safe ...bool) *SortedIntArray {
	a := NewSortedIntArraySize(0, safe...)
	a.array = array
	a.sorter(a.array, defaultComparatorInt)
	return a
}

// NewSortedIntArrayFromCopy creates and returns an sorted array from a copy of given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArrayFromCopy(array []int, safe ...bool) *SortedIntArray {
	newArray := make([]int, len(array))
	copy(newArray, array)
	return NewSortedIntArrayFrom(newArray, safe...)
}

// At returns the value by the specified index.
// If the given `index` is out of range of the array, it returns `0`.
func (a *SortedIntArray) At(index int) (value int) {
	a.lazyInit()
	return a.SortedTArray.At(index)
}

// SetArray sets the underlying slice array with the given `array`.
func (a *SortedIntArray) SetArray(array []int) *SortedIntArray {
	a.lazyInit()
	a.SortedTArray.SetArray(array)
	return a
}

// Sort sorts the array in increasing order.
// The parameter `reverse` controls whether sort
// in increasing order(default) or decreasing order.
func (a *SortedIntArray) Sort() *SortedIntArray {
	a.lazyInit()
	a.SortedTArray.Sort()
	return a
}

// Add adds one or multiple values to sorted array, the array always keeps sorted.
// It's alias of function Append, see Append.
func (a *SortedIntArray) Add(values ...int) *SortedIntArray {
	a.lazyInit()
	return a.Append(values...)
}

// Append adds one or multiple values to sorted array, the array always keeps sorted.
func (a *SortedIntArray) Append(values ...int) *SortedIntArray {
	a.lazyInit()
	a.SortedTArray.Append(values...)
	return a
}

// Get returns the value by the specified index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *SortedIntArray) Get(index int) (value int, found bool) {
	a.lazyInit()
	return a.SortedTArray.Get(index)
}

// Remove removes an item by index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *SortedIntArray) Remove(index int) (value int, found bool) {
	a.lazyInit()
	return a.SortedTArray.Remove(index)
}

// RemoveValue removes an item by value.
// It returns true if value is found in the array, or else false if not found.
func (a *SortedIntArray) RemoveValue(value int) bool {
	a.lazyInit()
	return a.SortedTArray.RemoveValue(value)
}

// RemoveValues removes an item by `values`.
func (a *SortedIntArray) RemoveValues(values ...int) {
	a.lazyInit()
	a.SortedTArray.RemoveValues(values...)
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the `found` is false.
func (a *SortedIntArray) PopLeft() (value int, found bool) {
	a.lazyInit()
	return a.SortedTArray.PopLeft()
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the `found` is false.
func (a *SortedIntArray) PopRight() (value int, found bool) {
	a.lazyInit()
	return a.SortedTArray.PopRight()
}

// PopRand randomly pops and return an item out of array.
// Note that if the array is empty, the `found` is false.
func (a *SortedIntArray) PopRand() (value int, found bool) {
	a.lazyInit()
	return a.SortedTArray.PopRand()
}

// PopRands randomly pops and returns `size` items out of array.
// If the given `size` is greater than size of the array, it returns all elements of the array.
// Note that if given `size` <= 0 or the array is empty, it returns nil.
func (a *SortedIntArray) PopRands(size int) []int {
	a.lazyInit()
	return a.SortedTArray.PopRands(size)
}

// PopLefts pops and returns `size` items from the beginning of array.
// If the given `size` is greater than size of the array, it returns all elements of the array.
// Note that if given `size` <= 0 or the array is empty, it returns nil.
func (a *SortedIntArray) PopLefts(size int) []int {
	a.lazyInit()
	return a.SortedTArray.PopLefts(size)
}

// PopRights pops and returns `size` items from the end of array.
// If the given `size` is greater than size of the array, it returns all elements of the array.
// Note that if given `size` <= 0 or the array is empty, it returns nil.
func (a *SortedIntArray) PopRights(size int) []int {
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
func (a *SortedIntArray) Range(start int, end ...int) []int {
	a.lazyInit()
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
func (a *SortedIntArray) SubSlice(offset int, length ...int) []int {
	a.lazyInit()
	return a.SortedTArray.SubSlice(offset, length...)
}

// Len returns the length of array.
func (a *SortedIntArray) Len() int {
	a.lazyInit()
	return a.SortedTArray.Len()
}

// Sum returns the sum of values in an array.
func (a *SortedIntArray) Sum() (sum int) {
	a.lazyInit()
	return a.SortedTArray.Sum()
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *SortedIntArray) Slice() []int {
	a.lazyInit()
	return a.SortedTArray.Slice()
}

// Interfaces returns current array as []any.
func (a *SortedIntArray) Interfaces() []any {
	a.lazyInit()
	return a.SortedTArray.Interfaces()
}

// Contains checks whether a value exists in the array.
func (a *SortedIntArray) Contains(value int) bool {
	a.lazyInit()
	return a.SortedTArray.Contains(value)
}

// Search searches array by `value`, returns the index of `value`,
// or returns -1 if not exists.
func (a *SortedIntArray) Search(value int) (index int) {
	a.lazyInit()
	return a.SortedTArray.Search(value)
}

// SetUnique sets unique mark to the array,
// which means it does not contain any repeated items.
// It also do unique check, remove all repeated items.
func (a *SortedIntArray) SetUnique(unique bool) *SortedIntArray {
	a.lazyInit()
	a.SortedTArray.SetUnique(unique)
	return a
}

// Unique uniques the array, clear repeated items.
func (a *SortedIntArray) Unique() *SortedIntArray {
	a.lazyInit()
	a.SortedTArray.Unique()
	return a
}

// Clone returns a new array, which is a copy of current array.
func (a *SortedIntArray) Clone() (newArray *SortedIntArray) {
	a.lazyInit()
	return &SortedIntArray{
		SortedTArray: a.SortedTArray.Clone(),
	}
}

// Clear deletes all items of current array.
func (a *SortedIntArray) Clear() *SortedIntArray {
	a.lazyInit()
	a.SortedTArray.Clear()
	return a
}

// LockFunc locks writing by callback function `f`.
func (a *SortedIntArray) LockFunc(f func(array []int)) *SortedIntArray {
	a.lazyInit()
	a.SortedTArray.LockFunc(f)
	return a
}

// RLockFunc locks reading by callback function `f`.
func (a *SortedIntArray) RLockFunc(f func(array []int)) *SortedIntArray {
	a.lazyInit()
	a.SortedTArray.RLockFunc(f)
	return a
}

// Merge merges `array` into current array.
// The parameter `array` can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (a *SortedIntArray) Merge(array any) *SortedIntArray {
	a.lazyInit()
	return a.Add(gconv.Ints(array)...)
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by `size`.
// The last chunk may contain less than size elements.
func (a *SortedIntArray) Chunk(size int) [][]int {
	a.lazyInit()
	return a.SortedTArray.Chunk(size)
}

// Rand randomly returns one item from array(no deleting).
func (a *SortedIntArray) Rand() (value int, found bool) {
	a.lazyInit()
	return a.SortedTArray.Rand()
}

// Rands randomly returns `size` items from array(no deleting).
func (a *SortedIntArray) Rands(size int) []int {
	a.lazyInit()
	return a.SortedTArray.Rands(size)
}

// Join joins array elements with a string `glue`.
func (a *SortedIntArray) Join(glue string) string {
	a.lazyInit()
	return a.SortedTArray.Join(glue)
}

// CountValues counts the number of occurrences of all values in the array.
func (a *SortedIntArray) CountValues() map[int]int {
	a.lazyInit()
	return a.SortedTArray.CountValues()
}

// Iterator is alias of IteratorAsc.
func (a *SortedIntArray) Iterator(f func(k int, v int) bool) {
	a.lazyInit()
	a.SortedTArray.Iterator(f)
}

// IteratorAsc iterates the array readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *SortedIntArray) IteratorAsc(f func(k int, v int) bool) {
	a.lazyInit()
	a.SortedTArray.IteratorAsc(f)
}

// IteratorDesc iterates the array readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *SortedIntArray) IteratorDesc(f func(k int, v int) bool) {
	a.lazyInit()
	a.SortedTArray.IteratorDesc(f)
}

// String returns current array as a string, which implements like json.Marshal does.
func (a *SortedIntArray) String() string {
	if a == nil {
		return ""
	}
	a.lazyInit()
	return "[" + a.Join(",") + "]"
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (a *SortedIntArray) MarshalJSON() ([]byte, error) {
	a.lazyInit()
	return a.SortedTArray.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (a *SortedIntArray) UnmarshalJSON(b []byte) error {
	a.lazyInit()
	if a.comparator == nil || a.sorter == nil {
		a.comparator = defaultComparatorInt
		a.sorter = quickSortInt
		a.array = make([]int, 0)
	}

	return a.SortedTArray.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for array.
func (a *SortedIntArray) UnmarshalValue(value any) (err error) {
	a.lazyInit()
	if a.comparator == nil || a.sorter == nil {
		a.comparator = defaultComparatorInt
		a.sorter = quickSortInt
	}

	return a.SortedTArray.UnmarshalValue(value)
}

// Filter iterates array and filters elements using custom callback function.
// It removes the element from array if callback function `filter` returns true,
// it or else does nothing and continues iterating.
func (a *SortedIntArray) Filter(filter func(index int, value int) bool) *SortedIntArray {
	a.lazyInit()
	a.SortedTArray.Filter(filter)
	return a
}

// FilterEmpty removes all zero value of the array.
func (a *SortedIntArray) FilterEmpty() *SortedIntArray {
	a.lazyInit()
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.array) == 0 {
		return a
	}

	if a.array[0] != 0 && a.array[len(a.array)-1] != 0 {
		a.SortedTArray.FilterEmpty()
		return a
	}

	for i := 0; i < len(a.array); {
		if a.array[i] == 0 {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			break
		}
	}
	for i := len(a.array) - 1; i >= 0; {
		if a.array[i] == 0 {
			a.array = append(a.array[:i], a.array[i+1:]...)
			i--
		} else {
			break
		}
	}
	return a
}

// Walk applies a user supplied function `f` to every item of array.
func (a *SortedIntArray) Walk(f func(value int) int) *SortedIntArray {
	a.lazyInit()
	a.SortedTArray.Walk(f)
	return a
}

// IsEmpty checks whether the array is empty.
func (a *SortedIntArray) IsEmpty() bool {
	a.lazyInit()
	return a.SortedTArray.IsEmpty()
}

// DeepCopy implements interface for deep copy of current type.
func (a *SortedIntArray) DeepCopy() any {
	a.lazyInit()
	return &SortedIntArray{
		SortedTArray: a.SortedTArray.DeepCopy().(*SortedTArray[int]),
	}
}

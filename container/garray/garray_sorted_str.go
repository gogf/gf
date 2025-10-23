// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
	"bytes"
	"strings"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// SortedStrArray is a golang sorted string array with rich features.
// It is using increasing order in default, which can be changed by
// setting it a custom comparator.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type SortedStrArray struct {
	*SortedTArray[string]
}

// lazyInit lazily initializes the array.
func (a *SortedStrArray) lazyInit() {
	if a.SortedTArray == nil {
		a.SortedTArray = NewSortedTArraySize(0, defaultComparatorStr, false)
		a.SetSorter(quickSortStr)
	}
}

// NewSortedStrArray creates and returns an empty sorted array.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedStrArray(safe ...bool) *SortedStrArray {
	return NewSortedStrArraySize(0, safe...)
}

// NewSortedStrArrayComparator creates and returns an empty sorted array with specified comparator.
// The parameter `safe` is used to specify whether using array in concurrent-safety which is false in default.
func NewSortedStrArrayComparator(comparator func(a, b string) int, safe ...bool) *SortedStrArray {
	array := NewSortedStrArray(safe...)
	array.comparator = comparator
	return array
}

// NewSortedStrArraySize create and returns an sorted array with given size and cap.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedStrArraySize(cap int, safe ...bool) *SortedStrArray {
	a := NewSortedTArraySize(cap, defaultComparatorStr, safe...)
	a.SetSorter(quickSortStr)
	return &SortedStrArray{
		SortedTArray: a,
	}
}

// NewSortedStrArrayFrom creates and returns an sorted array with given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedStrArrayFrom(array []string, safe ...bool) *SortedStrArray {
	a := NewSortedStrArraySize(0, safe...)
	a.array = array
	quickSortStr(a.array, a.getComparator())
	return a
}

// NewSortedStrArrayFromCopy creates and returns an sorted array from a copy of given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedStrArrayFromCopy(array []string, safe ...bool) *SortedStrArray {
	newArray := make([]string, len(array))
	copy(newArray, array)
	return NewSortedStrArrayFrom(newArray, safe...)
}

// SetArray sets the underlying slice array with the given `array`.
func (a *SortedStrArray) SetArray(array []string) *SortedStrArray {
	a.lazyInit()
	a.SortedTArray.SetArray(array)
	return a
}

// At returns the value by the specified index.
// If the given `index` is out of range of the array, it returns an empty string.
func (a *SortedStrArray) At(index int) (value string) {
	a.lazyInit()
	return a.SortedTArray.At(index)
}

// Sort sorts the array in increasing order.
// The parameter `reverse` controls whether sort
// in increasing order(default) or decreasing order.
func (a *SortedStrArray) Sort() *SortedStrArray {
	a.lazyInit()
	a.SortedTArray.Sort()
	return a
}

// Add adds one or multiple values to sorted array, the array always keeps sorted.
// It's alias of function Append, see Append.
func (a *SortedStrArray) Add(values ...string) *SortedStrArray {
	a.lazyInit()
	a.SortedTArray.Add(values...)
	return a
}

// Append adds one or multiple values to sorted array, the array always keeps sorted.
func (a *SortedStrArray) Append(values ...string) *SortedStrArray {
	a.lazyInit()
	a.SortedTArray.Append(values...)
	return a
}

// Get returns the value by the specified index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *SortedStrArray) Get(index int) (value string, found bool) {
	a.lazyInit()
	return a.SortedTArray.Get(index)
}

// Remove removes an item by index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *SortedStrArray) Remove(index int) (value string, found bool) {
	a.lazyInit()
	return a.SortedTArray.Remove(index)
}

// RemoveValue removes an item by value.
// It returns true if value is found in the array, or else false if not found.
func (a *SortedStrArray) RemoveValue(value string) bool {
	a.lazyInit()
	return a.SortedTArray.RemoveValue(value)
}

// RemoveValues removes an item by `values`.
func (a *SortedStrArray) RemoveValues(values ...string) {
	a.lazyInit()
	a.SortedTArray.RemoveValues(values...)
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the `found` is false.
func (a *SortedStrArray) PopLeft() (value string, found bool) {
	a.lazyInit()
	return a.SortedTArray.PopLeft()
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the `found` is false.
func (a *SortedStrArray) PopRight() (value string, found bool) {
	a.lazyInit()
	return a.SortedTArray.PopRight()
}

// PopRand randomly pops and return an item out of array.
// Note that if the array is empty, the `found` is false.
func (a *SortedStrArray) PopRand() (value string, found bool) {
	a.lazyInit()
	return a.SortedTArray.PopRand()
}

// PopRands randomly pops and returns `size` items out of array.
// If the given `size` is greater than size of the array, it returns all elements of the array.
// Note that if given `size` <= 0 or the array is empty, it returns nil.
func (a *SortedStrArray) PopRands(size int) []string {
	a.lazyInit()
	return a.SortedTArray.PopRands(size)
}

// PopLefts pops and returns `size` items from the beginning of array.
// If the given `size` is greater than size of the array, it returns all elements of the array.
// Note that if given `size` <= 0 or the array is empty, it returns nil.
func (a *SortedStrArray) PopLefts(size int) []string {
	a.lazyInit()
	return a.SortedTArray.PopLefts(size)
}

// PopRights pops and returns `size` items from the end of array.
// If the given `size` is greater than size of the array, it returns all elements of the array.
// Note that if given `size` <= 0 or the array is empty, it returns nil.
func (a *SortedStrArray) PopRights(size int) []string {
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
func (a *SortedStrArray) Range(start int, end ...int) []string {
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
func (a *SortedStrArray) SubSlice(offset int, length ...int) []string {
	a.lazyInit()
	return a.SortedTArray.SubSlice(offset, length...)
}

// Sum returns the sum of values in an array.
func (a *SortedStrArray) Sum() (sum int) {
	a.lazyInit()
	return a.SortedTArray.Sum()
}

// Len returns the length of array.
func (a *SortedStrArray) Len() int {
	a.lazyInit()
	return a.SortedTArray.Len()
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *SortedStrArray) Slice() []string {
	a.lazyInit()
	return a.SortedTArray.Slice()
}

// Interfaces returns current array as []any.
func (a *SortedStrArray) Interfaces() []any {
	a.lazyInit()
	return a.SortedTArray.Interfaces()
}

// Contains checks whether a value exists in the array.
func (a *SortedStrArray) Contains(value string) bool {
	a.lazyInit()
	return a.SortedTArray.Contains(value)
}

// ContainsI checks whether a value exists in the array with case-insensitively.
// Note that it internally iterates the whole array to do the comparison with case-insensitively.
func (a *SortedStrArray) ContainsI(value string) bool {
	a.lazyInit()
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return false
	}
	for _, v := range a.array {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}

// Search searches array by `value`, returns the index of `value`,
// or returns -1 if not exists.
func (a *SortedStrArray) Search(value string) (index int) {
	a.lazyInit()
	return a.SortedTArray.Search(value)
}

// SetUnique sets unique mark to the array,
// which means it does not contain any repeated items.
// It also do unique check, remove all repeated items.
func (a *SortedStrArray) SetUnique(unique bool) *SortedStrArray {
	a.lazyInit()
	a.SortedTArray.SetUnique(unique)
	return a
}

// Unique uniques the array, clear repeated items.
func (a *SortedStrArray) Unique() *SortedStrArray {
	a.lazyInit()
	a.SortedTArray.Unique()
	return a
}

// Clone returns a new array, which is a copy of current array.
func (a *SortedStrArray) Clone() (newArray *SortedStrArray) {
	a.lazyInit()
	return &SortedStrArray{
		SortedTArray: a.SortedTArray.Clone(),
	}
}

// Clear deletes all items of current array.
func (a *SortedStrArray) Clear() *SortedStrArray {
	a.lazyInit()
	a.SortedTArray.Clear()
	return a
}

// LockFunc locks writing by callback function `f`.
func (a *SortedStrArray) LockFunc(f func(array []string)) *SortedStrArray {
	a.lazyInit()
	a.SortedTArray.LockFunc(f)
	return a
}

// RLockFunc locks reading by callback function `f`.
func (a *SortedStrArray) RLockFunc(f func(array []string)) *SortedStrArray {
	a.lazyInit()
	a.SortedTArray.RLockFunc(f)
	return a
}

// Merge merges `array` into current array.
// The parameter `array` can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (a *SortedStrArray) Merge(array any) *SortedStrArray {
	a.lazyInit()
	return a.Add(gconv.Strings(array)...)
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by `size`.
// The last chunk may contain less than size elements.
func (a *SortedStrArray) Chunk(size int) [][]string {
	a.lazyInit()
	return a.SortedTArray.Chunk(size)
}

// Rand randomly returns one item from array(no deleting).
func (a *SortedStrArray) Rand() (value string, found bool) {
	a.lazyInit()
	return a.SortedTArray.Rand()
}

// Rands randomly returns `size` items from array(no deleting).
func (a *SortedStrArray) Rands(size int) []string {
	a.lazyInit()
	return a.SortedTArray.Rands(size)
}

// Join joins array elements with a string `glue`.
func (a *SortedStrArray) Join(glue string) string {
	a.lazyInit()
	return a.SortedTArray.Join(glue)
}

// CountValues counts the number of occurrences of all values in the array.
func (a *SortedStrArray) CountValues() map[string]int {
	a.lazyInit()
	return a.SortedTArray.CountValues()
}

// Iterator is alias of IteratorAsc.
func (a *SortedStrArray) Iterator(f func(k int, v string) bool) {
	a.lazyInit()
	a.SortedTArray.Iterator(f)
}

// IteratorAsc iterates the array readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *SortedStrArray) IteratorAsc(f func(k int, v string) bool) {
	a.lazyInit()
	a.SortedTArray.IteratorAsc(f)
}

// IteratorDesc iterates the array readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *SortedStrArray) IteratorDesc(f func(k int, v string) bool) {
	a.lazyInit()
	a.SortedTArray.IteratorDesc(f)
}

// String returns current array as a string, which implements like json.Marshal does.
func (a *SortedStrArray) String() string {
	if a == nil {
		return ""
	}
	a.lazyInit()
	a.mu.RLock()
	defer a.mu.RUnlock()
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('[')
	for k, v := range a.array {
		buffer.WriteString(`"` + gstr.QuoteMeta(v, `"\`) + `"`)
		if k != len(a.array)-1 {
			buffer.WriteByte(',')
		}
	}
	buffer.WriteByte(']')
	return buffer.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (a SortedStrArray) MarshalJSON() ([]byte, error) {
	a.lazyInit()
	return a.SortedTArray.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (a *SortedStrArray) UnmarshalJSON(b []byte) error {
	a.lazyInit()
	if a.comparator == nil || a.sorter == nil {
		a.comparator = defaultComparatorStr
		a.sorter = quickSortStr
		a.array = make([]string, 0)
	}
	return a.SortedTArray.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for array.
func (a *SortedStrArray) UnmarshalValue(value any) (err error) {
	a.lazyInit()
	if a.comparator == nil || a.sorter == nil {
		a.comparator = defaultComparatorStr
		a.sorter = quickSortStr
	}

	return a.SortedTArray.UnmarshalValue(value)
}

// Filter iterates array and filters elements using custom callback function.
// It removes the element from array if callback function `filter` returns true,
// it or else does nothing and continues iterating.
func (a *SortedStrArray) Filter(filter func(index int, value string) bool) *SortedStrArray {
	a.lazyInit()
	a.SortedTArray.Filter(filter)
	return a
}

// FilterEmpty removes all empty string value of the array.
func (a *SortedStrArray) FilterEmpty() *SortedStrArray {
	a.lazyInit()
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.array) == 0 {
		return a
	}

	if a.array[0] != "" && a.array[len(a.array)-1] != "" {
		a.SortedTArray.FilterEmpty()
		return a
	}

	for i := 0; i < len(a.array); {
		if a.array[i] == "" {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			break
		}
	}
	for i := len(a.array) - 1; i >= 0; {
		if a.array[i] == "" {
			a.array = append(a.array[:i], a.array[i+1:]...)
			i--
		} else {
			break
		}
	}
	return a
}

// Walk applies a user supplied function `f` to every item of array.
func (a *SortedStrArray) Walk(f func(value string) string) *SortedStrArray {
	a.lazyInit()
	a.SortedTArray.Walk(f)
	return a
}

// IsEmpty checks whether the array is empty.
func (a *SortedStrArray) IsEmpty() bool {
	a.lazyInit()
	return a.SortedTArray.IsEmpty()
}

// DeepCopy implements interface for deep copy of current type.
func (a *SortedStrArray) DeepCopy() any {
	a.lazyInit()
	return &SortedStrArray{
		SortedTArray: a.SortedTArray.DeepCopy().(*SortedTArray[string]),
	}
}

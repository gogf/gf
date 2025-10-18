// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
	"bytes"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// StrArray is a golang string array with rich features.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type StrArray struct {
	*TArray[string]
}

// NewStrArray creates and returns an empty array.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewStrArray(safe ...bool) *StrArray {
	return NewStrArraySize(0, 0, safe...)
}

// NewStrArraySize create and returns an array with given size and cap.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewStrArraySize(size int, cap int, safe ...bool) *StrArray {
	return &StrArray{
		TArray: NewTArraySize[string](size, cap, safe...),
	}
}

// NewStrArrayFrom creates and returns an array with given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewStrArrayFrom(array []string, safe ...bool) *StrArray {
	return &StrArray{
		TArray: NewTArrayFrom(array, safe...),
	}
}

// NewStrArrayFromCopy creates and returns an array from a copy of given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewStrArrayFromCopy(array []string, safe ...bool) *StrArray {
	newArray := make([]string, len(array))
	copy(newArray, array)
	return NewStrArrayFrom(newArray, safe...)
}

// lazyInit lazily initializes the array.
func (a *StrArray) lazyInit() {
	if a.TArray == nil {
		a.TArray = NewTArray[string](false)
	}
}

// At returns the value by the specified index.
// If the given `index` is out of range of the array, it returns an empty string.
func (a *StrArray) At(index int) (value string) {
	a.lazyInit()
	return a.TArray.At(index)
}

// Get returns the value by the specified index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *StrArray) Get(index int) (value string, found bool) {
	a.lazyInit()
	return a.TArray.Get(index)
}

// Set sets value to specified index.
func (a *StrArray) Set(index int, value string) error {
	a.lazyInit()
	return a.TArray.Set(index, value)
}

// SetArray sets the underlying slice array with the given `array`.
func (a *StrArray) SetArray(array []string) *StrArray {
	a.lazyInit()
	a.TArray.SetArray(array)
	return a
}

// Replace replaces the array items by given `array` from the beginning of array.
func (a *StrArray) Replace(array []string) *StrArray {
	a.lazyInit()
	a.TArray.Replace(array)
	return a
}

// Sum returns the sum of values in an array.
func (a *StrArray) Sum() (sum int) {
	a.lazyInit()
	return a.TArray.Sum()
}

// Sort sorts the array in increasing order.
// The parameter `reverse` controls whether sort
// in increasing order(default) or decreasing order
func (a *StrArray) Sort(reverse ...bool) *StrArray {
	a.lazyInit()

	a.TArray.mu.Lock()
	defer a.TArray.mu.Unlock()
	if len(reverse) > 0 && reverse[0] {
		sort.Slice(a.array, func(i, j int) bool {
			return strings.Compare(a.array[i], a.array[j]) >= 0
		})
	} else {
		sort.Strings(a.array)
	}
	return a
}

// SortFunc sorts the array by custom function `less`.
func (a *StrArray) SortFunc(less func(v1, v2 string) bool) *StrArray {
	a.lazyInit()
	a.TArray.SortFunc(less)
	return a
}

// InsertBefore inserts the `values` to the front of `index`.
func (a *StrArray) InsertBefore(index int, values ...string) error {
	a.lazyInit()
	return a.TArray.InsertBefore(index, values...)
}

// InsertAfter inserts the `values` to the back of `index`.
func (a *StrArray) InsertAfter(index int, values ...string) error {
	a.lazyInit()
	return a.TArray.InsertAfter(index, values...)
}

// Remove removes an item by index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *StrArray) Remove(index int) (value string, found bool) {
	a.lazyInit()
	return a.TArray.Remove(index)
}

// RemoveValue removes an item by value.
// It returns true if value is found in the array, or else false if not found.
func (a *StrArray) RemoveValue(value string) bool {
	a.lazyInit()
	return a.TArray.RemoveValue(value)
}

// RemoveValues removes multiple items by `values`.
func (a *StrArray) RemoveValues(values ...string) {
	a.lazyInit()
	a.TArray.RemoveValues(values...)
}

// PushLeft pushes one or multiple items to the beginning of array.
func (a *StrArray) PushLeft(value ...string) *StrArray {
	a.lazyInit()
	a.TArray.PushLeft(value...)
	return a
}

// PushRight pushes one or multiple items to the end of array.
// It equals to Append.
func (a *StrArray) PushRight(value ...string) *StrArray {
	a.lazyInit()
	a.TArray.PushRight(value...)
	return a
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the `found` is false.
func (a *StrArray) PopLeft() (value string, found bool) {
	a.lazyInit()
	return a.TArray.PopLeft()
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the `found` is false.
func (a *StrArray) PopRight() (value string, found bool) {
	a.lazyInit()
	return a.TArray.PopRight()
}

// PopRand randomly pops and return an item out of array.
// Note that if the array is empty, the `found` is false.
func (a *StrArray) PopRand() (value string, found bool) {
	a.lazyInit()
	return a.TArray.PopRand()
}

// PopRands randomly pops and returns `size` items out of array.
// If the given `size` is greater than size of the array, it returns all elements of the array.
// Note that if given `size` <= 0 or the array is empty, it returns nil.
func (a *StrArray) PopRands(size int) []string {
	a.lazyInit()
	return a.TArray.PopRands(size)
}

// PopLefts pops and returns `size` items from the beginning of array.
// If the given `size` is greater than size of the array, it returns all elements of the array.
// Note that if given `size` <= 0 or the array is empty, it returns nil.
func (a *StrArray) PopLefts(size int) []string {
	a.lazyInit()
	return a.TArray.PopLefts(size)
}

// PopRights pops and returns `size` items from the end of array.
// If the given `size` is greater than size of the array, it returns all elements of the array.
// Note that if given `size` <= 0 or the array is empty, it returns nil.
func (a *StrArray) PopRights(size int) []string {
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
func (a *StrArray) Range(start int, end ...int) []string {
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
func (a *StrArray) SubSlice(offset int, length ...int) []string {
	a.lazyInit()
	return a.TArray.SubSlice(offset, length...)
}

// Append is alias of PushRight,please See PushRight.
func (a *StrArray) Append(value ...string) *StrArray {
	a.lazyInit()
	a.TArray.Append(value...)
	return a
}

// Len returns the length of array.
func (a *StrArray) Len() int {
	a.lazyInit()
	return a.TArray.Len()
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *StrArray) Slice() []string {
	a.lazyInit()
	return a.TArray.Slice()
}

// Interfaces returns current array as []any.
func (a *StrArray) Interfaces() []any {
	a.lazyInit()
	return a.TArray.Interfaces()
}

// Clone returns a new array, which is a copy of current array.
func (a *StrArray) Clone() (newArray *StrArray) {
	a.lazyInit()
	return &StrArray{
		TArray: a.TArray.Clone(),
	}
}

// Clear deletes all items of current array.
func (a *StrArray) Clear() *StrArray {
	a.lazyInit()
	a.TArray.Clear()
	return a
}

// Contains checks whether a value exists in the array.
func (a *StrArray) Contains(value string) bool {
	a.lazyInit()
	return a.TArray.Contains(value)
}

// ContainsI checks whether a value exists in the array with case-insensitively.
// Note that it internally iterates the whole array to do the comparison with case-insensitively.
func (a *StrArray) ContainsI(value string) bool {
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
func (a *StrArray) Search(value string) int {
	a.lazyInit()
	return a.TArray.Search(value)
}

// Unique uniques the array, clear repeated items.
// Example: [1,1,2,3,2] -> [1,2,3]
func (a *StrArray) Unique() *StrArray {
	a.lazyInit()
	a.TArray.Unique()
	return a
}

// LockFunc locks writing by callback function `f`.
func (a *StrArray) LockFunc(f func(array []string)) *StrArray {
	a.lazyInit()
	a.TArray.LockFunc(f)
	return a
}

// RLockFunc locks reading by callback function `f`.
func (a *StrArray) RLockFunc(f func(array []string)) *StrArray {
	a.lazyInit()
	a.TArray.RLockFunc(f)
	return a
}

// Merge merges `array` into current array.
// The parameter `array` can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (a *StrArray) Merge(array any) *StrArray {
	return a.Append(gconv.Strings(array)...)
}

// Fill fills an array with num entries of the value `value`,
// keys starting at the `startIndex` parameter.
func (a *StrArray) Fill(startIndex int, num int, value string) error {
	a.lazyInit()
	return a.TArray.Fill(startIndex, num, value)
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by `size`.
// The last chunk may contain less than size elements.
func (a *StrArray) Chunk(size int) [][]string {
	a.lazyInit()
	return a.TArray.Chunk(size)
}

// Pad pads array to the specified length with `value`.
// If size is positive then the array is padded on the right, or negative on the left.
// If the absolute value of `size` is less than or equal to the length of the array
// then no padding takes place.
func (a *StrArray) Pad(size int, value string) *StrArray {
	a.lazyInit()
	a.TArray.Pad(size, value)
	return a
}

// Rand randomly returns one item from array(no deleting).
func (a *StrArray) Rand() (value string, found bool) {
	a.lazyInit()
	return a.TArray.Rand()
}

// Rands randomly returns `size` items from array(no deleting).
func (a *StrArray) Rands(size int) []string {
	a.lazyInit()
	return a.TArray.Rands(size)
}

// Shuffle randomly shuffles the array.
func (a *StrArray) Shuffle() *StrArray {
	a.lazyInit()
	a.TArray.Shuffle()
	return a
}

// Reverse makes array with elements in reverse order.
func (a *StrArray) Reverse() *StrArray {
	a.lazyInit()
	a.TArray.Reverse()
	return a
}

// Join joins array elements with a string `glue`.
func (a *StrArray) Join(glue string) string {
	a.lazyInit()
	return a.TArray.Join(glue)
}

// CountValues counts the number of occurrences of all values in the array.
func (a *StrArray) CountValues() map[string]int {
	a.lazyInit()
	return a.TArray.CountValues()
}

// Iterator is alias of IteratorAsc.
func (a *StrArray) Iterator(f func(k int, v string) bool) {
	a.IteratorAsc(f)
}

// IteratorAsc iterates the array readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *StrArray) IteratorAsc(f func(k int, v string) bool) {
	a.lazyInit()
	a.TArray.IteratorAsc(f)
}

// IteratorDesc iterates the array readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *StrArray) IteratorDesc(f func(k int, v string) bool) {
	a.lazyInit()
	a.TArray.IteratorDesc(f)
}

// String returns current array as a string, which implements like json.Marshal does.
func (a *StrArray) String() string {
	if a == nil {
		return ""
	}

	a.lazyInit()

	a.TArray.mu.RLock()
	defer a.TArray.mu.RUnlock()
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('[')
	for k, v := range a.TArray.array {
		buffer.WriteString(`"` + gstr.QuoteMeta(v, `"\`) + `"`)
		if k != len(a.TArray.array)-1 {
			buffer.WriteByte(',')
		}
	}
	buffer.WriteByte(']')
	return buffer.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (a StrArray) MarshalJSON() ([]byte, error) {
	a.lazyInit()
	return a.TArray.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (a *StrArray) UnmarshalJSON(b []byte) error {
	a.lazyInit()
	return a.TArray.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for array.
func (a *StrArray) UnmarshalValue(value any) error {
	a.lazyInit()
	return a.TArray.UnmarshalValue(value)
}

// Filter iterates array and filters elements using custom callback function.
// It removes the element from array if callback function `filter` returns true,
// it or else does nothing and continues iterating.
func (a *StrArray) Filter(filter func(index int, value string) bool) *StrArray {
	a.lazyInit()
	a.TArray.Filter(filter)
	return a
}

// FilterEmpty removes all empty string value of the array.
func (a *StrArray) FilterEmpty() *StrArray {
	a.lazyInit()
	a.TArray.FilterEmpty()
	return a
}

// Walk applies a user supplied function `f` to every item of array.
func (a *StrArray) Walk(f func(value string) string) *StrArray {
	a.lazyInit()
	a.TArray.Walk(f)
	return a
}

// IsEmpty checks whether the array is empty.
func (a *StrArray) IsEmpty() bool {
	a.lazyInit()
	return a.TArray.IsEmpty()
}

// DeepCopy implements interface for deep copy of current type.
func (a *StrArray) DeepCopy() any {
	if a == nil {
		return nil
	}
	a.lazyInit()
	return &StrArray{
		TArray: a.TArray.DeepCopy().(*TArray[string]),
	}
}

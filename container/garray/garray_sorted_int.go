// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/internal/json"
	"math"
	"sort"

	"github.com/gogf/gf/internal/rwmutex"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
)

// SortedIntArray is a golang sorted int array with rich features.
// It is using increasing order in default, which can be changed by
// setting it a custom comparator.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type SortedIntArray struct {
	mu         rwmutex.RWMutex
	array      []int
	unique     bool               // Whether enable unique feature(false)
	comparator func(a, b int) int // Comparison function(it returns -1: a < b; 0: a == b; 1: a > b)
}

// NewSortedIntArray creates and returns an empty sorted array.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArray(safe ...bool) *SortedIntArray {
	return NewSortedIntArraySize(0, safe...)
}

// NewSortedIntArrayComparator creates and returns an empty sorted array with specified comparator.
// The parameter <safe> is used to specify whether using array in concurrent-safety which is false in default.
func NewSortedIntArrayComparator(comparator func(a, b int) int, safe ...bool) *SortedIntArray {
	array := NewSortedIntArray(safe...)
	array.comparator = comparator
	return array
}

// NewSortedIntArraySize create and returns an sorted array with given size and cap.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArraySize(cap int, safe ...bool) *SortedIntArray {
	return &SortedIntArray{
		mu:         rwmutex.Create(safe...),
		array:      make([]int, 0, cap),
		comparator: defaultComparatorInt,
	}
}

// NewSortedIntArrayRange creates and returns a array by a range from <start> to <end>
// with step value <step>.
func NewSortedIntArrayRange(start, end, step int, safe ...bool) *SortedIntArray {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]int, (end-start+1)/step)
	index := 0
	for i := start; i <= end; i += step {
		slice[index] = i
		index++
	}
	return NewSortedIntArrayFrom(slice, safe...)
}

// NewIntArrayFrom creates and returns an sorted array with given slice <array>.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArrayFrom(array []int, safe ...bool) *SortedIntArray {
	a := NewSortedIntArraySize(0, safe...)
	a.array = array
	sort.Ints(a.array)
	return a
}

// NewSortedIntArrayFromCopy creates and returns an sorted array from a copy of given slice <array>.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArrayFromCopy(array []int, safe ...bool) *SortedIntArray {
	newArray := make([]int, len(array))
	copy(newArray, array)
	return NewSortedIntArrayFrom(newArray, safe...)
}

// SetArray sets the underlying slice array with the given <array>.
func (a *SortedIntArray) SetArray(array []int) *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array = array
	quickSortInt(a.array, a.getComparator())
	return a
}

// Sort sorts the array in increasing order.
// The parameter <reverse> controls whether sort
// in increasing order(default) or decreasing order.
func (a *SortedIntArray) Sort() *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	quickSortInt(a.array, a.getComparator())
	return a
}

// Add adds one or multiple values to sorted array, the array always keeps sorted.
// It's alias of function Append, see Append.
func (a *SortedIntArray) Add(values ...int) *SortedIntArray {
	return a.Append(values...)
}

// Append adds one or multiple values to sorted array, the array always keeps sorted.
func (a *SortedIntArray) Append(values ...int) *SortedIntArray {
	if len(values) == 0 {
		return a
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, value := range values {
		index, cmp := a.binSearch(value, false)
		if a.unique && cmp == 0 {
			continue
		}
		if index < 0 {
			a.array = append(a.array, value)
			continue
		}
		if cmp > 0 {
			index++
		}
		rear := append([]int{}, a.array[index:]...)
		a.array = append(a.array[0:index], value)
		a.array = append(a.array, rear...)
	}
	return a
}

// Get returns the value by the specified index.
// If the given <index> is out of range of the array, the <found> is false.
func (a *SortedIntArray) Get(index int) (value int, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if index < 0 || index >= len(a.array) {
		return 0, false
	}
	return a.array[index], true
}

// Remove removes an item by index.
// If the given <index> is out of range of the array, the <found> is false.
func (a *SortedIntArray) Remove(index int) (value int, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doRemoveWithoutLock(index)
}

// doRemoveWithoutLock removes an item by index without lock.
func (a *SortedIntArray) doRemoveWithoutLock(index int) (value int, found bool) {
	if index < 0 || index >= len(a.array) {
		return 0, false
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		value := a.array[0]
		a.array = a.array[1:]
		return value, true
	} else if index == len(a.array)-1 {
		value := a.array[index]
		a.array = a.array[:index]
		return value, true
	}
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
	value = a.array[index]
	a.array = append(a.array[:index], a.array[index+1:]...)
	return value, true
}

// RemoveValue removes an item by value.
// It returns true if value is found in the array, or else false if not found.
func (a *SortedIntArray) RemoveValue(value int) bool {
	if i := a.Search(value); i != -1 {
		_, found := a.Remove(i)
		return found
	}
	return false
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the <found> is false.
func (a *SortedIntArray) PopLeft() (value int, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.array) == 0 {
		return 0, false
	}
	value = a.array[0]
	a.array = a.array[1:]
	return value, true
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the <found> is false.
func (a *SortedIntArray) PopRight() (value int, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	index := len(a.array) - 1
	if index < 0 {
		return 0, false
	}
	value = a.array[index]
	a.array = a.array[:index]
	return value, true
}

// PopRand randomly pops and return an item out of array.
// Note that if the array is empty, the <found> is false.
func (a *SortedIntArray) PopRand() (value int, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doRemoveWithoutLock(grand.Intn(len(a.array)))
}

// PopRands randomly pops and returns <size> items out of array.
// If the given <size> is greater than size of the array, it returns all elements of the array.
// Note that if given <size> <= 0 or the array is empty, it returns nil.
func (a *SortedIntArray) PopRands(size int) []int {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	if size >= len(a.array) {
		size = len(a.array)
	}
	array := make([]int, size)
	for i := 0; i < size; i++ {
		array[i], _ = a.doRemoveWithoutLock(grand.Intn(len(a.array)))
	}
	return array
}

// PopLefts pops and returns <size> items from the beginning of array.
// If the given <size> is greater than size of the array, it returns all elements of the array.
// Note that if given <size> <= 0 or the array is empty, it returns nil.
func (a *SortedIntArray) PopLefts(size int) []int {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	if size >= len(a.array) {
		array := a.array
		a.array = a.array[:0]
		return array
	}
	value := a.array[0:size]
	a.array = a.array[size:]
	return value
}

// PopRights pops and returns <size> items from the end of array.
// If the given <size> is greater than size of the array, it returns all elements of the array.
// Note that if given <size> <= 0 or the array is empty, it returns nil.
func (a *SortedIntArray) PopRights(size int) []int {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	index := len(a.array) - size
	if index <= 0 {
		array := a.array
		a.array = a.array[:0]
		return array
	}
	value := a.array[index:]
	a.array = a.array[:index]
	return value
}

// Range picks and returns items by range, like array[start:end].
// Notice, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
//
// If <end> is negative, then the offset will start from the end of array.
// If <end> is omitted, then the sequence will have everything from start up
// until the end of the array.
func (a *SortedIntArray) Range(start int, end ...int) []int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	offsetEnd := len(a.array)
	if len(end) > 0 && end[0] < offsetEnd {
		offsetEnd = end[0]
	}
	if start > offsetEnd {
		return nil
	}
	if start < 0 {
		start = 0
	}
	array := ([]int)(nil)
	if a.mu.IsSafe() {
		array = make([]int, offsetEnd-start)
		copy(array, a.array[start:offsetEnd])
	} else {
		array = a.array[start:offsetEnd]
	}
	return array
}

// SubSlice returns a slice of elements from the array as specified
// by the <offset> and <size> parameters.
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
	a.mu.RLock()
	defer a.mu.RUnlock()
	size := len(a.array)
	if len(length) > 0 {
		size = length[0]
	}
	if offset > len(a.array) {
		return nil
	}
	if offset < 0 {
		offset = len(a.array) + offset
		if offset < 0 {
			return nil
		}
	}
	if size < 0 {
		offset += size
		size = -size
		if offset < 0 {
			return nil
		}
	}
	end := offset + size
	if end > len(a.array) {
		end = len(a.array)
		size = len(a.array) - offset
	}
	if a.mu.IsSafe() {
		s := make([]int, size)
		copy(s, a.array[offset:])
		return s
	} else {
		return a.array[offset:end]
	}
}

// Len returns the length of array.
func (a *SortedIntArray) Len() int {
	a.mu.RLock()
	length := len(a.array)
	a.mu.RUnlock()
	return length
}

// Sum returns the sum of values in an array.
func (a *SortedIntArray) Sum() (sum int) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		sum += v
	}
	return
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *SortedIntArray) Slice() []int {
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

// Interfaces returns current array as []interface{}.
func (a *SortedIntArray) Interfaces() []interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()
	array := make([]interface{}, len(a.array))
	for k, v := range a.array {
		array[k] = v
	}
	return array
}

// Contains checks whether a value exists in the array.
func (a *SortedIntArray) Contains(value int) bool {
	return a.Search(value) != -1
}

// Search searches array by <value>, returns the index of <value>,
// or returns -1 if not exists.
func (a *SortedIntArray) Search(value int) (index int) {
	if i, r := a.binSearch(value, true); r == 0 {
		return i
	}
	return -1
}

// Binary search.
// It returns the last compared index and the result.
// If <result> equals to 0, it means the value at <index> is equals to <value>.
// If <result> lesser than 0, it means the value at <index> is lesser than <value>.
// If <result> greater than 0, it means the value at <index> is greater than <value>.
func (a *SortedIntArray) binSearch(value int, lock bool) (index int, result int) {
	if lock {
		a.mu.RLock()
		defer a.mu.RUnlock()
	}
	if len(a.array) == 0 {
		return -1, -2
	}
	min := 0
	max := len(a.array) - 1
	mid := 0
	cmp := -2
	for min <= max {
		mid = min + int((max-min)/2)
		cmp = a.getComparator()(value, a.array[mid])
		switch {
		case cmp < 0:
			max = mid - 1
		case cmp > 0:
			min = mid + 1
		default:
			return mid, cmp
		}
	}
	return mid, cmp
}

// SetUnique sets unique mark to the array,
// which means it does not contain any repeated items.
// It also do unique check, remove all repeated items.
func (a *SortedIntArray) SetUnique(unique bool) *SortedIntArray {
	oldUnique := a.unique
	a.unique = unique
	if unique && oldUnique != unique {
		a.Unique()
	}
	return a
}

// Unique uniques the array, clear repeated items.
func (a *SortedIntArray) Unique() *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.array) == 0 {
		return a
	}
	i := 0
	for {
		if i == len(a.array)-1 {
			break
		}
		if a.getComparator()(a.array[i], a.array[i+1]) == 0 {
			a.array = append(a.array[:i+1], a.array[i+1+1:]...)
		} else {
			i++
		}
	}
	return a
}

// Clone returns a new array, which is a copy of current array.
func (a *SortedIntArray) Clone() (newArray *SortedIntArray) {
	a.mu.RLock()
	array := make([]int, len(a.array))
	copy(array, a.array)
	a.mu.RUnlock()
	return NewSortedIntArrayFrom(array, a.mu.IsSafe())
}

// Clear deletes all items of current array.
func (a *SortedIntArray) Clear() *SortedIntArray {
	a.mu.Lock()
	if len(a.array) > 0 {
		a.array = make([]int, 0)
	}
	a.mu.Unlock()
	return a
}

// LockFunc locks writing by callback function <f>.
func (a *SortedIntArray) LockFunc(f func(array []int)) *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	f(a.array)
	return a
}

// RLockFunc locks reading by callback function <f>.
func (a *SortedIntArray) RLockFunc(f func(array []int)) *SortedIntArray {
	a.mu.RLock()
	defer a.mu.RUnlock()
	f(a.array)
	return a
}

// Merge merges <array> into current array.
// The parameter <array> can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (a *SortedIntArray) Merge(array interface{}) *SortedIntArray {
	return a.Add(gconv.Ints(array)...)
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by <size>.
// The last chunk may contain less than size elements.
func (a *SortedIntArray) Chunk(size int) [][]int {
	if size < 1 {
		return nil
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
		n = append(n, a.array[i*size:end])
		i++
	}
	return n
}

// Rand randomly returns one item from array(no deleting).
func (a *SortedIntArray) Rand() (value int, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return 0, false
	}
	return a.array[grand.Intn(len(a.array))], true
}

// Rands randomly returns <size> items from array(no deleting).
func (a *SortedIntArray) Rands(size int) []int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	array := make([]int, size)
	for i := 0; i < size; i++ {
		array[i] = a.array[grand.Intn(len(a.array))]
	}
	return array
}

// Join joins array elements with a string <glue>.
func (a *SortedIntArray) Join(glue string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	for k, v := range a.array {
		buffer.WriteString(gconv.String(v))
		if k != len(a.array)-1 {
			buffer.WriteString(glue)
		}
	}
	return buffer.String()
}

// CountValues counts the number of occurrences of all values in the array.
func (a *SortedIntArray) CountValues() map[int]int {
	m := make(map[int]int)
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		m[v]++
	}
	return m
}

// Iterator is alias of IteratorAsc.
func (a *SortedIntArray) Iterator(f func(k int, v int) bool) {
	a.IteratorAsc(f)
}

// IteratorAsc iterates the array readonly in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (a *SortedIntArray) IteratorAsc(f func(k int, v int) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for k, v := range a.array {
		if !f(k, v) {
			break
		}
	}
}

// IteratorDesc iterates the array readonly in descending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (a *SortedIntArray) IteratorDesc(f func(k int, v int) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for i := len(a.array) - 1; i >= 0; i-- {
		if !f(i, a.array[i]) {
			break
		}
	}
}

// String returns current array as a string, which implements like json.Marshal does.
func (a *SortedIntArray) String() string {
	return "[" + a.Join(",") + "]"
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (a SortedIntArray) MarshalJSON() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return json.Marshal(a.array)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (a *SortedIntArray) UnmarshalJSON(b []byte) error {
	if a.comparator == nil {
		a.array = make([]int, 0)
		a.comparator = defaultComparatorInt
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := json.Unmarshal(b, &a.array); err != nil {
		return err
	}
	if a.array != nil {
		sort.Ints(a.array)
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for array.
func (a *SortedIntArray) UnmarshalValue(value interface{}) (err error) {
	if a.comparator == nil {
		a.comparator = defaultComparatorInt
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	switch value.(type) {
	case string, []byte:
		err = json.Unmarshal(gconv.Bytes(value), &a.array)
	default:
		a.array = gconv.SliceInt(value)
	}
	if a.array != nil {
		sort.Ints(a.array)
	}
	return err
}

// FilterEmpty removes all zero value of the array.
func (a *SortedIntArray) FilterEmpty() *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
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
		} else {
			break
		}
	}
	return a
}

// Walk applies a user supplied function <f> to every item of array.
func (a *SortedIntArray) Walk(f func(value int) int) *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Keep the array always sorted.
	defer quickSortInt(a.array, a.getComparator())

	for i, v := range a.array {
		a.array[i] = f(v)
	}
	return a
}

// IsEmpty checks whether the array is empty.
func (a *SortedIntArray) IsEmpty() bool {
	return a.Len() == 0
}

// getComparator returns the comparator if it's previously set,
// or else it returns a default comparator.
func (a *SortedIntArray) getComparator() func(a, b int) int {
	if a.comparator == nil {
		return defaultComparatorInt
	}
	return a.comparator
}

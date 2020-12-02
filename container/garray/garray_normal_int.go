// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gogf/gf/internal/json"
	"math"
	"sort"

	"github.com/gogf/gf/internal/rwmutex"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
)

// IntArray is a golang int array with rich features.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type IntArray struct {
	mu    rwmutex.RWMutex
	array []int
}

// NewIntArray creates and returns an empty array.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewIntArray(safe ...bool) *IntArray {
	return NewIntArraySize(0, 0, safe...)
}

// NewIntArraySize create and returns an array with given size and cap.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewIntArraySize(size int, cap int, safe ...bool) *IntArray {
	return &IntArray{
		mu:    rwmutex.Create(safe...),
		array: make([]int, size, cap),
	}
}

// NewIntArrayRange creates and returns a array by a range from <start> to <end>
// with step value <step>.
func NewIntArrayRange(start, end, step int, safe ...bool) *IntArray {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]int, (end-start+1)/step)
	index := 0
	for i := start; i <= end; i += step {
		slice[index] = i
		index++
	}
	return NewIntArrayFrom(slice, safe...)
}

// NewIntArrayFrom creates and returns an array with given slice <array>.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewIntArrayFrom(array []int, safe ...bool) *IntArray {
	return &IntArray{
		mu:    rwmutex.Create(safe...),
		array: array,
	}
}

// NewIntArrayFromCopy creates and returns an array from a copy of given slice <array>.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewIntArrayFromCopy(array []int, safe ...bool) *IntArray {
	newArray := make([]int, len(array))
	copy(newArray, array)
	return &IntArray{
		mu:    rwmutex.Create(safe...),
		array: newArray,
	}
}

// Get returns the value by the specified index.
// If the given <index> is out of range of the array, the <found> is false.
func (a *IntArray) Get(index int) (value int, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if index < 0 || index >= len(a.array) {
		return 0, false
	}
	return a.array[index], true
}

// Set sets value to specified index.
func (a *IntArray) Set(index int, value int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", index, len(a.array)))
	}
	a.array[index] = value
	return nil
}

// SetArray sets the underlying slice array with the given <array>.
func (a *IntArray) SetArray(array []int) *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array = array
	return a
}

// Replace replaces the array items by given <array> from the beginning of array.
func (a *IntArray) Replace(array []int) *IntArray {
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
func (a *IntArray) Sum() (sum int) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		sum += v
	}
	return
}

// Sort sorts the array in increasing order.
// The parameter <reverse> controls whether sort in increasing order(default) or decreasing order.
func (a *IntArray) Sort(reverse ...bool) *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(reverse) > 0 && reverse[0] {
		sort.Slice(a.array, func(i, j int) bool {
			if a.array[i] < a.array[j] {
				return false
			}
			return true
		})
	} else {
		sort.Ints(a.array)
	}
	return a
}

// SortFunc sorts the array by custom function <less>.
func (a *IntArray) SortFunc(less func(v1, v2 int) bool) *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	sort.Slice(a.array, func(i, j int) bool {
		return less(a.array[i], a.array[j])
	})
	return a
}

// InsertBefore inserts the <value> to the front of <index>.
func (a *IntArray) InsertBefore(index int, value int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", index, len(a.array)))
	}
	rear := append([]int{}, a.array[index:]...)
	a.array = append(a.array[0:index], value)
	a.array = append(a.array, rear...)
	return nil
}

// InsertAfter inserts the <value> to the back of <index>.
func (a *IntArray) InsertAfter(index int, value int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", index, len(a.array)))
	}
	rear := append([]int{}, a.array[index+1:]...)
	a.array = append(a.array[0:index+1], value)
	a.array = append(a.array, rear...)
	return nil
}

// Remove removes an item by index.
// If the given <index> is out of range of the array, the <found> is false.
func (a *IntArray) Remove(index int) (value int, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doRemoveWithoutLock(index)
}

// doRemoveWithoutLock removes an item by index without lock.
func (a *IntArray) doRemoveWithoutLock(index int) (value int, found bool) {
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
func (a *IntArray) RemoveValue(value int) bool {
	if i := a.Search(value); i != -1 {
		_, found := a.Remove(i)
		return found
	}
	return false
}

// PushLeft pushes one or multiple items to the beginning of array.
func (a *IntArray) PushLeft(value ...int) *IntArray {
	a.mu.Lock()
	a.array = append(value, a.array...)
	a.mu.Unlock()
	return a
}

// PushRight pushes one or multiple items to the end of array.
// It equals to Append.
func (a *IntArray) PushRight(value ...int) *IntArray {
	a.mu.Lock()
	a.array = append(a.array, value...)
	a.mu.Unlock()
	return a
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the <found> is false.
func (a *IntArray) PopLeft() (value int, found bool) {
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
func (a *IntArray) PopRight() (value int, found bool) {
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
func (a *IntArray) PopRand() (value int, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doRemoveWithoutLock(grand.Intn(len(a.array)))
}

// PopRands randomly pops and returns <size> items out of array.
// If the given <size> is greater than size of the array, it returns all elements of the array.
// Note that if given <size> <= 0 or the array is empty, it returns nil.
func (a *IntArray) PopRands(size int) []int {
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
func (a *IntArray) PopLefts(size int) []int {
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
func (a *IntArray) PopRights(size int) []int {
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
func (a *IntArray) Range(start int, end ...int) []int {
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
func (a *IntArray) SubSlice(offset int, length ...int) []int {
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

// See PushRight.
func (a *IntArray) Append(value ...int) *IntArray {
	a.mu.Lock()
	a.array = append(a.array, value...)
	a.mu.Unlock()
	return a
}

// Len returns the length of array.
func (a *IntArray) Len() int {
	a.mu.RLock()
	length := len(a.array)
	a.mu.RUnlock()
	return length
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *IntArray) Slice() []int {
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
func (a *IntArray) Interfaces() []interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()
	array := make([]interface{}, len(a.array))
	for k, v := range a.array {
		array[k] = v
	}
	return array
}

// Clone returns a new array, which is a copy of current array.
func (a *IntArray) Clone() (newArray *IntArray) {
	a.mu.RLock()
	array := make([]int, len(a.array))
	copy(array, a.array)
	a.mu.RUnlock()
	return NewIntArrayFrom(array, a.mu.IsSafe())
}

// Clear deletes all items of current array.
func (a *IntArray) Clear() *IntArray {
	a.mu.Lock()
	if len(a.array) > 0 {
		a.array = make([]int, 0)
	}
	a.mu.Unlock()
	return a
}

// Contains checks whether a value exists in the array.
func (a *IntArray) Contains(value int) bool {
	return a.Search(value) != -1
}

// Search searches array by <value>, returns the index of <value>,
// or returns -1 if not exists.
func (a *IntArray) Search(value int) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return -1
	}
	result := -1
	for index, v := range a.array {
		if v == value {
			result = index
			break
		}
	}
	return result
}

// Unique uniques the array, clear repeated items.
// Example: [1,1,2,3,2] -> [1,2,3]
func (a *IntArray) Unique() *IntArray {
	a.mu.Lock()
	for i := 0; i < len(a.array)-1; i++ {
		for j := i + 1; j < len(a.array); {
			if a.array[i] == a.array[j] {
				a.array = append(a.array[:j], a.array[j+1:]...)
			} else {
				j++
			}
		}
	}
	a.mu.Unlock()
	return a
}

// LockFunc locks writing by callback function <f>.
func (a *IntArray) LockFunc(f func(array []int)) *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	f(a.array)
	return a
}

// RLockFunc locks reading by callback function <f>.
func (a *IntArray) RLockFunc(f func(array []int)) *IntArray {
	a.mu.RLock()
	defer a.mu.RUnlock()
	f(a.array)
	return a
}

// Merge merges <array> into current array.
// The parameter <array> can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (a *IntArray) Merge(array interface{}) *IntArray {
	return a.Append(gconv.Ints(array)...)
}

// Fill fills an array with num entries of the value <value>,
// keys starting at the <startIndex> parameter.
func (a *IntArray) Fill(startIndex int, num int, value int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if startIndex < 0 || startIndex > len(a.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", startIndex, len(a.array)))
	}
	for i := startIndex; i < startIndex+num; i++ {
		if i > len(a.array)-1 {
			a.array = append(a.array, value)
		} else {
			a.array[i] = value
		}
	}
	return nil
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by <size>.
// The last chunk may contain less than size elements.
func (a *IntArray) Chunk(size int) [][]int {
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

// Pad pads array to the specified length with <value>.
// If size is positive then the array is padded on the right, or negative on the left.
// If the absolute value of <size> is less than or equal to the length of the array
// then no padding takes place.
func (a *IntArray) Pad(size int, value int) *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size == 0 || (size > 0 && size < len(a.array)) || (size < 0 && size > -len(a.array)) {
		return a
	}
	n := size
	if size < 0 {
		n = -size
	}
	n -= len(a.array)
	tmp := make([]int, n)
	for i := 0; i < n; i++ {
		tmp[i] = value
	}
	if size > 0 {
		a.array = append(a.array, tmp...)
	} else {
		a.array = append(tmp, a.array...)
	}
	return a
}

// Rand randomly returns one item from array(no deleting).
func (a *IntArray) Rand() (value int, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return 0, false
	}
	return a.array[grand.Intn(len(a.array))], true
}

// Rands randomly returns <size> items from array(no deleting).
func (a *IntArray) Rands(size int) []int {
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

// Shuffle randomly shuffles the array.
func (a *IntArray) Shuffle() *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, v := range grand.Perm(len(a.array)) {
		a.array[i], a.array[v] = a.array[v], a.array[i]
	}
	return a
}

// Reverse makes array with elements in reverse order.
func (a *IntArray) Reverse() *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, j := 0, len(a.array)-1; i < j; i, j = i+1, j-1 {
		a.array[i], a.array[j] = a.array[j], a.array[i]
	}
	return a
}

// Join joins array elements with a string <glue>.
func (a *IntArray) Join(glue string) string {
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
func (a *IntArray) CountValues() map[int]int {
	m := make(map[int]int)
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		m[v]++
	}
	return m
}

// Iterator is alias of IteratorAsc.
func (a *IntArray) Iterator(f func(k int, v int) bool) {
	a.IteratorAsc(f)
}

// IteratorAsc iterates the array readonly in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (a *IntArray) IteratorAsc(f func(k int, v int) bool) {
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
func (a *IntArray) IteratorDesc(f func(k int, v int) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for i := len(a.array) - 1; i >= 0; i-- {
		if !f(i, a.array[i]) {
			break
		}
	}
}

// String returns current array as a string, which implements like json.Marshal does.
func (a *IntArray) String() string {
	return "[" + a.Join(",") + "]"
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (a IntArray) MarshalJSON() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return json.Marshal(a.array)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (a *IntArray) UnmarshalJSON(b []byte) error {
	if a.array == nil {
		a.array = make([]int, 0)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := json.Unmarshal(b, &a.array); err != nil {
		return err
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for array.
func (a *IntArray) UnmarshalValue(value interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	switch value.(type) {
	case string, []byte:
		return json.Unmarshal(gconv.Bytes(value), &a.array)
	default:
		a.array = gconv.SliceInt(value)
	}
	return nil
}

// FilterEmpty removes all zero value of the array.
func (a *IntArray) FilterEmpty() *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := 0; i < len(a.array); {
		if a.array[i] == 0 {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			i++
		}
	}
	return a
}

// Walk applies a user supplied function <f> to every item of array.
func (a *IntArray) Walk(f func(value int) int) *IntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, v := range a.array {
		a.array[i] = f(v)
	}
	return a
}

// IsEmpty checks whether the array is empty.
func (a *IntArray) IsEmpty() bool {
	return a.Len() == 0
}

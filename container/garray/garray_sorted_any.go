// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gutil"
	"math"
	"sort"

	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/internal/rwmutex"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
)

// It's using increasing order in default.
type SortedArray struct {
	mu         *rwmutex.RWMutex
	array      []interface{}
	unique     *gtype.Bool                // Whether enable unique feature(false)
	comparator func(a, b interface{}) int // Comparison function(it returns -1: a < b; 0: a == b; 1: a > b)
}

// NewSortedArray creates and returns an empty sorted array.
// The parameter <safe> is used to specify whether using array in concurrent-safety, which is false in default.
// The parameter <comparator> used to compare values to sort in array,
// if it returns value < 0, means v1 < v2;
// if it returns value = 0, means v1 = v2;
// if it returns value > 0, means v1 > v2;
func NewSortedArray(comparator func(a, b interface{}) int, safe ...bool) *SortedArray {
	return NewSortedArraySize(0, comparator, safe...)
}

// NewSortedArraySize create and returns an sorted array with given size and cap.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedArraySize(cap int, comparator func(a, b interface{}) int, safe ...bool) *SortedArray {
	return &SortedArray{
		mu:         rwmutex.New(safe...),
		unique:     gtype.NewBool(),
		array:      make([]interface{}, 0, cap),
		comparator: comparator,
	}
}

// NewSortedArrayRange creates and returns a array by a range from <start> to <end>
// with step value <step>.
func NewSortedArrayRange(start, end, step int, comparator func(a, b interface{}) int, safe ...bool) *SortedArray {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]interface{}, (end-start+1)/step)
	index := 0
	for i := start; i <= end; i += step {
		slice[index] = i
		index++
	}
	return NewSortedArrayFrom(slice, comparator, safe...)
}

// NewSortedArrayFrom creates and returns an sorted array with given slice <array>.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedArrayFrom(array []interface{}, comparator func(a, b interface{}) int, safe ...bool) *SortedArray {
	a := NewSortedArraySize(0, comparator, safe...)
	a.array = array
	sort.Slice(a.array, func(i, j int) bool {
		return a.comparator(a.array[i], a.array[j]) < 0
	})
	return a
}

// NewSortedArrayFromCopy creates and returns an sorted array from a copy of given slice <array>.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedArrayFromCopy(array []interface{}, comparator func(a, b interface{}) int, safe ...bool) *SortedArray {
	newArray := make([]interface{}, len(array))
	copy(newArray, array)
	return NewSortedArrayFrom(newArray, comparator, safe...)
}

// SetArray sets the underlying slice array with the given <array>.
func (a *SortedArray) SetArray(array []interface{}) *SortedArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array = array
	sort.Slice(a.array, func(i, j int) bool {
		return a.comparator(a.array[i], a.array[j]) < 0
	})
	return a
}

// SetComparator sets/changes the comparator for sorting.
func (a *SortedArray) SetComparator(comparator func(a, b interface{}) int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.comparator = comparator
	// Resort the array if comparator is changed.
	sort.Slice(a.array, func(i, j int) bool {
		return a.comparator(a.array[i], a.array[j]) < 0
	})
}

// Sort sorts the array in increasing order.
// The parameter <reverse> controls whether sort
// in increasing order(default) or decreasing order
func (a *SortedArray) Sort() *SortedArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	sort.Slice(a.array, func(i, j int) bool {
		return a.comparator(a.array[i], a.array[j]) < 0
	})
	return a
}

// Add adds one or multiple values to sorted array, the array always keeps sorted.
func (a *SortedArray) Add(values ...interface{}) *SortedArray {
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
		rear := append([]interface{}{}, a.array[index:]...)
		a.array = append(a.array[0:index], value)
		a.array = append(a.array, rear...)
	}
	return a
}

// Get returns the value of the specified index,
// the caller should notice the boundary of the array.
func (a *SortedArray) Get(index int) interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()
	value := a.array[index]
	return value
}

// Remove removes an item by index.
func (a *SortedArray) Remove(index int) interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return nil
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		value := a.array[0]
		a.array = a.array[1:]
		return value
	} else if index == len(a.array)-1 {
		value := a.array[index]
		a.array = a.array[:index]
		return value
	}
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
	value := a.array[index]
	a.array = append(a.array[:index], a.array[index+1:]...)
	return value
}

// RemoveValue removes an item by value.
// It returns true if value is found in the array, or else false if not found.
func (a *SortedArray) RemoveValue(value interface{}) bool {
	if i := a.Search(value); i != -1 {
		a.Remove(i)
		return true
	}
	return false
}

// PopLeft pops and returns an item from the beginning of array.
func (a *SortedArray) PopLeft() interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()
	value := a.array[0]
	a.array = a.array[1:]
	return value
}

// PopRight pops and returns an item from the end of array.
func (a *SortedArray) PopRight() interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()
	index := len(a.array) - 1
	value := a.array[index]
	a.array = a.array[:index]
	return value
}

// PopRand randomly pops and return an item out of array.
func (a *SortedArray) PopRand() interface{} {
	return a.Remove(grand.Intn(len(a.array)))
}

// PopRands randomly pops and returns <size> items out of array.
func (a *SortedArray) PopRands(size int) []interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size > len(a.array) {
		size = len(a.array)
	}
	array := make([]interface{}, size)
	for i := 0; i < size; i++ {
		index := grand.Intn(len(a.array))
		array[i] = a.array[index]
		a.array = append(a.array[:index], a.array[index+1:]...)
	}
	return array
}

// PopLefts pops and returns <size> items from the beginning of array.
func (a *SortedArray) PopLefts(size int) []interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()
	length := len(a.array)
	if size > length {
		size = length
	}
	value := a.array[0:size]
	a.array = a.array[size:]
	return value
}

// PopRights pops and returns <size> items from the end of array.
func (a *SortedArray) PopRights(size int) []interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()
	index := len(a.array) - size
	if index < 0 {
		index = 0
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
func (a *SortedArray) Range(start int, end ...int) []interface{} {
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
	array := ([]interface{})(nil)
	if a.mu.IsSafe() {
		array = make([]interface{}, offsetEnd-start)
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
func (a *SortedArray) SubSlice(offset int, length ...int) []interface{} {
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
		s := make([]interface{}, size)
		copy(s, a.array[offset:])
		return s
	} else {
		return a.array[offset:end]
	}
}

// Sum returns the sum of values in an array.
func (a *SortedArray) Sum() (sum int) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		sum += gconv.Int(v)
	}
	return
}

// Len returns the length of array.
func (a *SortedArray) Len() int {
	a.mu.RLock()
	length := len(a.array)
	a.mu.RUnlock()
	return length
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *SortedArray) Slice() []interface{} {
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

// Interfaces returns current array as []interface{}.
func (a *SortedArray) Interfaces() []interface{} {
	return a.Slice()
}

// Contains checks whether a value exists in the array.
func (a *SortedArray) Contains(value interface{}) bool {
	return a.Search(value) != -1
}

// Search searches array by <value>, returns the index of <value>,
// or returns -1 if not exists.
func (a *SortedArray) Search(value interface{}) (index int) {
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
func (a *SortedArray) binSearch(value interface{}, lock bool) (index int, result int) {
	if len(a.array) == 0 {
		return -1, -2
	}
	if lock {
		a.mu.RLock()
		defer a.mu.RUnlock()
	}
	min := 0
	max := len(a.array) - 1
	mid := 0
	cmp := -2
	for min <= max {
		mid = int((min + max) / 2)
		cmp = a.comparator(value, a.array[mid])
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
func (a *SortedArray) SetUnique(unique bool) *SortedArray {
	oldUnique := a.unique.Val()
	a.unique.Set(unique)
	if unique && oldUnique != unique {
		a.Unique()
	}
	return a
}

// Unique uniques the array, clear repeated items.
func (a *SortedArray) Unique() *SortedArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	i := 0
	for {
		if i == len(a.array)-1 {
			break
		}
		if a.comparator(a.array[i], a.array[i+1]) == 0 {
			a.array = append(a.array[:i+1], a.array[i+1+1:]...)
		} else {
			i++
		}
	}
	return a
}

// Clone returns a new array, which is a copy of current array.
func (a *SortedArray) Clone() (newArray *SortedArray) {
	a.mu.RLock()
	array := make([]interface{}, len(a.array))
	copy(array, a.array)
	a.mu.RUnlock()
	return NewSortedArrayFrom(array, a.comparator, !a.mu.IsSafe())
}

// Clear deletes all items of current array.
func (a *SortedArray) Clear() *SortedArray {
	a.mu.Lock()
	if len(a.array) > 0 {
		a.array = make([]interface{}, 0)
	}
	a.mu.Unlock()
	return a
}

// LockFunc locks writing by callback function <f>.
func (a *SortedArray) LockFunc(f func(array []interface{})) *SortedArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	f(a.array)
	return a
}

// RLockFunc locks reading by callback function <f>.
func (a *SortedArray) RLockFunc(f func(array []interface{})) *SortedArray {
	a.mu.RLock()
	defer a.mu.RUnlock()
	f(a.array)
	return a
}

// Merge merges <array> into current array.
// The parameter <array> can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (a *SortedArray) Merge(array interface{}) *SortedArray {
	switch v := array.(type) {
	case *Array:
		a.Add(gconv.Interfaces(v.Slice())...)
	case *IntArray:
		a.Add(gconv.Interfaces(v.Slice())...)
	case *StrArray:
		a.Add(gconv.Interfaces(v.Slice())...)
	case *SortedArray:
		a.Add(gconv.Interfaces(v.Slice())...)
	case *SortedIntArray:
		a.Add(gconv.Interfaces(v.Slice())...)
	case *SortedStrArray:
		a.Add(gconv.Interfaces(v.Slice())...)
	default:
		a.Add(gconv.Interfaces(array)...)
	}
	return a
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by <size>.
// The last chunk may contain less than size elements.
func (a *SortedArray) Chunk(size int) [][]interface{} {
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
		n = append(n, a.array[i*size:end])
		i++
	}
	return n
}

// Rand randomly returns one item from array(no deleting).
func (a *SortedArray) Rand() interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.array[grand.Intn(len(a.array))]
}

// Rands randomly returns <size> items from array(no deleting).
func (a *SortedArray) Rands(size int) []interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if size > len(a.array) {
		size = len(a.array)
	}
	n := make([]interface{}, size)
	for i, v := range grand.Perm(len(a.array)) {
		n[i] = a.array[v]
		if i == size-1 {
			break
		}
	}
	return n
}

// Join joins array elements with a string <glue>.
func (a *SortedArray) Join(glue string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
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
func (a *SortedArray) CountValues() map[interface{}]int {
	m := make(map[interface{}]int)
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		m[v]++
	}
	return m
}

// Iterator is alias of IteratorAsc.
func (a *SortedArray) Iterator(f func(k int, v interface{}) bool) {
	a.IteratorAsc(f)
}

// IteratorAsc iterates the array in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (a *SortedArray) IteratorAsc(f func(k int, v interface{}) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for k, v := range a.array {
		if !f(k, v) {
			break
		}
	}
}

// IteratorDesc iterates the array in descending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (a *SortedArray) IteratorDesc(f func(k int, v interface{}) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for i := len(a.array) - 1; i >= 0; i-- {
		if !f(i, a.array[i]) {
			break
		}
	}
}

// String returns current array as a string, which implements like json.Marshal does.
func (a *SortedArray) String() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('[')
	s := ""
	for k, v := range a.array {
		s = gconv.String(v)
		if gstr.IsNumeric(s) {
			buffer.WriteString(s)
		} else {
			buffer.WriteString(`"` + gstr.QuoteMeta(s, `"\`) + `"`)
		}
		if k != len(a.array)-1 {
			buffer.WriteByte(',')
		}
	}
	buffer.WriteByte(']')
	return buffer.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (a *SortedArray) MarshalJSON() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return json.Marshal(a.array)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (a *SortedArray) UnmarshalJSON(b []byte) error {
	if a.mu == nil {
		a.mu = rwmutex.New()
		a.array = make([]interface{}, 0)
		a.unique = gtype.NewBool()
		// Note that the comparator is string comparator in default.
		a.comparator = gutil.ComparatorString
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := json.Unmarshal(b, &a.array); err != nil {
		return err
	}
	if a.comparator != nil && a.array != nil {
		sort.Slice(a.array, func(i, j int) bool {
			return a.comparator(a.array[i], a.array[j]) < 0
		})
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for array.
func (a *SortedArray) UnmarshalValue(value interface{}) (err error) {
	if a.mu == nil {
		a.mu = rwmutex.New()
		a.unique = gtype.NewBool()
		// Note that the comparator is string comparator in default.
		a.comparator = gutil.ComparatorString
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	switch value.(type) {
	case string, []byte:
		err = json.Unmarshal(gconv.Bytes(value), &a.array)
	default:
		a.array = gconv.SliceAny(value)
	}
	if a.comparator != nil && a.array != nil {
		sort.Slice(a.array, func(i, j int) bool {
			return a.comparator(a.array[i], a.array[j]) < 0
		})
	}
	return err
}

// FilterNil removes all nil value of the array.
func (a *SortedArray) FilterNil() *SortedArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := 0; i < len(a.array); {
		if empty.IsNil(a.array[i]) {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			break
		}
	}
	for i := len(a.array) - 1; i >= 0; {
		if empty.IsNil(a.array[i]) {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			break
		}
	}
	return a
}

// FilterEmpty removes all empty value of the array.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (a *SortedArray) FilterEmpty() *SortedArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := 0; i < len(a.array); {
		if empty.IsEmpty(a.array[i]) {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			break
		}
	}
	for i := len(a.array) - 1; i >= 0; {
		if empty.IsEmpty(a.array[i]) {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			break
		}
	}
	return a
}

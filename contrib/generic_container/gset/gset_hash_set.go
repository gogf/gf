// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gset provides kinds of concurrent-safe/unsafe sets.
package gset

import (
	"bytes"
	"strings"

	"github.com/gogf/gf/contrib/generic_container/v2/conv"
	"github.com/gogf/gf/contrib/generic_container/v2/internal/json"
	"github.com/gogf/gf/contrib/generic_container/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type HashSet[T comparable] struct {
	mu   rwmutex.RWMutex
	data map[T]struct{}
}

// New create and returns a new set, which contains un-repeated items.
// The parameter `safe` is used to specify whether using set in concurrent-safety,
// which is false in default.
func New[T comparable](safe ...bool) *HashSet[T] {
	return NewSet[T](safe...)
}

// NewSet create and returns a new set, which contains un-repeated items.
// Also see New.
func NewSet[T comparable](safe ...bool) *HashSet[T] {
	return &HashSet[T]{
		data: make(map[T]struct{}),
		mu:   rwmutex.Create(safe...),
	}
}

// NewFrom returns a new set from `items`.
// Parameter `items` can be either a variable of any type, or a slice.
func NewFrom[T comparable](items []T, safe ...bool) *HashSet[T] {
	m := make(map[T]struct{})
	for _, v := range items {
		m[v] = struct{}{}
	}
	return &HashSet[T]{
		data: m,
		mu:   rwmutex.Create(safe...),
	}
}

// Iterator iterates the set readonly with given callback function `f`,
// if `f` returns true then continue iterating; or false to stop.
func (set *HashSet[T]) Iterator(f func(v T) bool) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k := range set.data {
		if !f(k) {
			break
		}
	}
}

// Add adds one or multiple items to the set.
func (set *HashSet[T]) Add(items ...T) {
	set.mu.Lock()
	if set.data == nil {
		set.data = make(map[T]struct{})
	}
	for _, v := range items {
		set.data[v] = struct{}{}
	}
	set.mu.Unlock()
}

// AddIfNotExist checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set,
// or else it does nothing and returns false.
//
// Note that, if `item` is nil, it does nothing and returns false.
func (set *HashSet[T]) AddIfNotExist(item T) bool {
	if g.IsNil(item) {
		return false
	}
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[T]struct{})
		}
		if _, ok := set.data[item]; !ok {
			set.data[item] = struct{}{}
			return true
		}
	}
	return false
}

// AddIfNotExistFunc checks whether item exists in the set,
// it adds the item to set and returns true if it does not exist in the set and
// function `f` returns true, or else it does nothing and returns false.
//
// Note that, if `item` is nil, it does nothing and returns false. The function `f`
// is executed without writing lock.
func (set *HashSet[T]) AddIfNotExistFunc(item T, f func() bool) bool {
	if g.IsNil(item) {
		return false
	}
	if !set.Contains(item) {
		if f() {
			set.mu.Lock()
			defer set.mu.Unlock()
			if set.data == nil {
				set.data = make(map[T]struct{})
			}
			if _, ok := set.data[item]; !ok {
				set.data[item] = struct{}{}
				return true
			}
		}
	}
	return false
}

// AddIfNotExistFuncLock checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and
// function `f` returns true, or else it does nothing and returns false.
//
// Note that, if `item` is nil, it does nothing and returns false. The function `f`
// is executed within writing lock.
func (set *HashSet[T]) AddIfNotExistFuncLock(item T, f func() bool) bool {
	if g.IsNil(item) {
		return false
	}
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[T]struct{})
		}
		if f() {
			if _, ok := set.data[item]; !ok {
				set.data[item] = struct{}{}
				return true
			}
		}
	}
	return false
}

// Contains checks whether the set contains `item`.
func (set *HashSet[T]) Contains(item T) bool {
	var ok bool
	set.mu.RLock()
	if set.data != nil {
		_, ok = set.data[item]
	}
	set.mu.RUnlock()
	return ok
}

// ContainsI checks whether a value exists in the set with case-insensitively.
// Note that it internally iterates the whole set to do the comparison with case-insensitively.
func (set *HashSet[T]) ContainsI(item T) bool {
	if g.IsNil(item) {
		return false
	}
	itemStr, ok := any(item).(string)
	if ok {
		set.mu.RLock()
		defer set.mu.RUnlock()
		for k := range set.data {
			if strings.EqualFold(any(k).(string), itemStr) {
				return true
			}
		}
	}
	return set.Contains(item)
}

// Remove deletes `item` from set.
func (set *HashSet[T]) Remove(item T) {
	set.mu.Lock()
	if set.data != nil {
		delete(set.data, item)
	}
	set.mu.Unlock()
}

// Size returns the size of the set.
func (set *HashSet[T]) Size() int {
	set.mu.RLock()
	l := len(set.data)
	set.mu.RUnlock()
	return l
}

// Clear deletes all items of the set.
func (set *HashSet[T]) Clear() {
	set.mu.Lock()
	set.data = make(map[T]struct{})
	set.mu.Unlock()
}

// Slice returns the an of items of the set as slice.
func (set *HashSet[T]) Slice() []T {
	set.mu.RLock()
	var (
		i   = 0
		ret = make([]T, len(set.data))
	)
	for item := range set.data {
		ret[i] = item
		i++
	}
	set.mu.RUnlock()
	return ret
}

// Join joins items with a string `glue`.
func (set *HashSet[T]) Join(glue string) string {
	set.mu.RLock()
	defer set.mu.RUnlock()
	if len(set.data) == 0 {
		return ""
	}
	var (
		l      = len(set.data)
		i      = 0
		buffer = bytes.NewBuffer(nil)
	)
	for k := range set.data {
		buffer.WriteString(gconv.String(k))
		if i != l-1 {
			buffer.WriteString(glue)
		}
		i++
	}
	return buffer.String()
}

// String returns items as a string, which implements like json.Marshal does.
func (set *HashSet[T]) String() string {
	if set == nil {
		return ""
	}
	set.mu.RLock()
	defer set.mu.RUnlock()
	var (
		s      string
		l      = len(set.data)
		i      = 0
		buffer = bytes.NewBuffer(nil)
	)
	buffer.WriteByte('[')
	for k := range set.data {
		s = gconv.String(k)
		if gstr.IsNumeric(s) {
			buffer.WriteString(s)
		} else {
			buffer.WriteString(`"` + gstr.QuoteMeta(s, `"\`) + `"`)
		}
		if i != l-1 {
			buffer.WriteByte(',')
		}
		i++
	}
	buffer.WriteByte(']')
	return buffer.String()
}

// LockFunc locks writing with callback function `f`.
func (set *HashSet[T]) LockFunc(f func(m map[T]struct{})) {
	set.mu.Lock()
	defer set.mu.Unlock()
	f(set.data)
}

// RLockFunc locks reading with callback function `f`.
func (set *HashSet[T]) RLockFunc(f func(m map[T]struct{})) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	f(set.data)
}

// Equal checks whether the two sets equal.
func (set *HashSet[T]) Equal(other *HashSet[T]) bool {
	if set == other {
		return true
	}
	set.mu.RLock()
	defer set.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	if len(set.data) != len(other.data) {
		return false
	}
	for key := range set.data {
		if _, ok := other.data[key]; !ok {
			return false
		}
	}
	return true
}

// IsSubsetOf checks whether the current set is a sub-set of `other`.
func (set *HashSet[T]) IsSubsetOf(other *HashSet[T]) bool {
	if set == other {
		return true
	}
	set.mu.RLock()
	defer set.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	for key := range set.data {
		if _, ok := other.data[key]; !ok {
			return false
		}
	}
	return true
}

// Union returns a new set which is the union of `set` and `others`.
// Which means, all the items in `newSet` are in `set` or in `others`.
func (set *HashSet[T]) Union(others ...*HashSet[T]) (newSet *HashSet[T]) {
	newSet = NewSet[T]()
	set.mu.RLock()
	defer set.mu.RUnlock()
	for _, other := range others {
		if set != other {
			other.mu.RLock()
		}
		for k, v := range set.data {
			newSet.data[k] = v
		}
		if set != other {
			for k, v := range other.data {
				newSet.data[k] = v
			}
		}
		if set != other {
			other.mu.RUnlock()
		}
	}

	return
}

// Diff returns a new set which is the difference set from `set` to `others`.
// Which means, all the items in `newSet` are in `set` but not in `others`.
func (set *HashSet[T]) Diff(others ...*HashSet[T]) (newSet *HashSet[T]) {
	newSet = NewSet[T]()
	set.mu.RLock()
	defer set.mu.RUnlock()
	for _, other := range others {
		if set == other {
			continue
		}
		other.mu.RLock()
		for k, v := range set.data {
			if _, ok := other.data[k]; !ok {
				newSet.data[k] = v
			}
		}
		other.mu.RUnlock()
	}
	return
}

// Intersect returns a new set which is the intersection from `set` to `others`.
// Which means, all the items in `newSet` are in `set` and also in `others`.
func (set *HashSet[T]) Intersect(others ...*HashSet[T]) (newSet *HashSet[T]) {
	newSet = NewSet[T]()
	set.mu.RLock()
	defer set.mu.RUnlock()
	for _, other := range others {
		if set != other {
			other.mu.RLock()
		}
		for k, v := range set.data {
			if _, ok := other.data[k]; ok {
				newSet.data[k] = v
			}
		}
		if set != other {
			other.mu.RUnlock()
		}
	}
	return
}

// Complement returns a new set which is the complement from `set` to `full`.
// Which means, all the items in `newSet` are in `full` and not in `set`.
//
// It returns the difference between `full` and `set`
// if the given set `full` is not the full set of `set`.
func (set *HashSet[T]) Complement(full *HashSet[T]) (newSet *HashSet[T]) {
	newSet = NewSet[T]()
	set.mu.RLock()
	defer set.mu.RUnlock()
	if set != full {
		full.mu.RLock()
		defer full.mu.RUnlock()
	}
	for k, v := range full.data {
		if _, ok := set.data[k]; !ok {
			newSet.data[k] = v
		}
	}
	return
}

// Merge adds items from `others` sets into `set`.
func (set *HashSet[T]) Merge(others ...*HashSet[T]) *HashSet[T] {
	set.mu.Lock()
	defer set.mu.Unlock()
	for _, other := range others {
		if set != other {
			other.mu.RLock()
		}
		for k, v := range other.data {
			set.data[k] = v
		}
		if set != other {
			other.mu.RUnlock()
		}
	}
	return set
}

// Sum sums items.
// Note: The items should be converted to int type,
// or you'd get a result that you unexpected.
func (set *HashSet[T]) Sum() (sum int) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k := range set.data {
		sum += gconv.Int(k)
	}
	return
}

// Pop randomly pops an item from set.
func (set *HashSet[T]) Pop() (value T) {
	set.mu.Lock()
	defer set.mu.Unlock()
	for k := range set.data {
		delete(set.data, k)
		return k
	}
	return
}

// Pops randomly pops `size` items from set.
// It returns all items if size == -1.
func (set *HashSet[T]) Pops(size int) []T {
	set.mu.Lock()
	defer set.mu.Unlock()
	if size > len(set.data) || size == -1 {
		size = len(set.data)
	}
	if size <= 0 {
		return nil
	}
	index := 0
	array := make([]T, size)
	for k := range set.data {
		delete(set.data, k)
		array[index] = k
		index++
		if index == size {
			break
		}
	}
	return array
}

// Walk applies a user supplied function `f` to every item of set.
func (set *HashSet[T]) Walk(f func(item T) T) *HashSet[T] {
	set.mu.Lock()
	defer set.mu.Unlock()
	m := make(map[T]struct{}, len(set.data))
	for k, v := range set.data {
		m[f(k)] = v
	}
	set.data = m
	return set
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (set HashSet[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Slice())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (set *HashSet[T]) UnmarshalJSON(b []byte) error {
	set.mu.Lock()
	defer set.mu.Unlock()
	if set.data == nil {
		set.data = make(map[T]struct{})
	}
	var array []T
	if err := json.UnmarshalUseNumber(b, &array); err != nil {
		return err
	}
	for _, v := range array {
		set.data[v] = struct{}{}
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for set.
func (set *HashSet[T]) UnmarshalValue(value interface{}) (err error) {
	set.mu.Lock()
	defer set.mu.Unlock()
	if set.data == nil {
		set.data = make(map[T]struct{})
	}
	var array []T
	switch value.(type) {
	case string, []byte:
		err = json.UnmarshalUseNumber(gconv.Bytes(value), &array)
	default:
		array = conv.SliceAny[T](value)
	}
	for _, v := range array {
		set.data[v] = struct{}{}
	}
	return
}

// DeepCopy implements interface for deep copy of current type.
func (set *HashSet[T]) DeepCopy() interface{} {
	if set == nil {
		return nil
	}
	set.mu.RLock()
	defer set.mu.RUnlock()
	data := make([]T, 0)
	for k := range set.data {
		data = append(data, k)
	}
	return NewFrom[T](data, set.mu.IsSafe())
}

// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gset provides kinds of concurrent-safe/unsafe sets.
package gset

import (
	"bytes"
	"encoding/json"
	"github.com/gogf/gf/internal/rwmutex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

type Set struct {
	mu   rwmutex.RWMutex
	data map[interface{}]struct{}
}

// New create and returns a new set, which contains un-repeated items.
// The parameter <safe> is used to specify whether using set in concurrent-safety,
// which is false in default.
func New(safe ...bool) *Set {
	return NewSet(safe...)
}

// See New.
func NewSet(safe ...bool) *Set {
	return &Set{
		data: make(map[interface{}]struct{}),
		mu:   rwmutex.Create(safe...),
	}
}

// NewFrom returns a new set from <items>.
// Parameter <items> can be either a variable of any type, or a slice.
func NewFrom(items interface{}, safe ...bool) *Set {
	m := make(map[interface{}]struct{})
	for _, v := range gconv.Interfaces(items) {
		m[v] = struct{}{}
	}
	return &Set{
		data: m,
		mu:   rwmutex.Create(safe...),
	}
}

// Iterator iterates the set with given callback function <f>,
// if <f> returns true then continue iterating; or false to stop.
func (set *Set) Iterator(f func(v interface{}) bool) *Set {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k, _ := range set.data {
		if !f(k) {
			break
		}
	}
	return set
}

// Add adds one or multiple items to the set.
func (set *Set) Add(item ...interface{}) *Set {
	set.mu.Lock()
	if set.data == nil {
		set.data = make(map[interface{}]struct{})
	}
	for _, v := range item {
		set.data[v] = struct{}{}
	}
	set.mu.Unlock()
	return set
}

// AddIfNotExistFunc adds the returned value of callback function <f> to the set
// if <item> does not exit in the set.
func (set *Set) AddIfNotExistFunc(item interface{}, f func() interface{}) *Set {
	if !set.Contains(item) {
		set.doAddWithLockCheck(item, f())
	}
	return set
}

// AddIfNotExistFuncLock adds the returned value of callback function <f> to the set
// if <item> does not exit in the set.
//
// Note that the callback function <f> is executed in the mutex.Lock of the set.
func (set *Set) AddIfNotExistFuncLock(item interface{}, f func() interface{}) *Set {
	if !set.Contains(item) {
		set.doAddWithLockCheck(item, f)
	}
	return set
}

// doAddWithLockCheck checks whether item exists with mutex.Lock,
// if not exists, it adds item to the set or else just returns the existing value.
//
// If <value> is type of <func() interface {}>,
// it will be executed with mutex.Lock of the set,
// and its return value will be added to the set.
//
// It returns item successfully added..
func (set *Set) doAddWithLockCheck(item interface{}, value interface{}) interface{} {
	set.mu.Lock()
	defer set.mu.Unlock()
	if set.data == nil {
		set.data = make(map[interface{}]struct{})
	}
	if _, ok := set.data[item]; !ok && value != nil {
		if f, ok := value.(func() interface{}); ok {
			item = f()
		} else {
			item = value
		}
	}
	if item != nil {
		set.data[item] = struct{}{}
	}
	return item
}

// Contains checks whether the set contains <item>.
func (set *Set) Contains(item interface{}) bool {
	var ok bool
	set.mu.RLock()
	if set.data != nil {
		_, ok = set.data[item]
	}
	set.mu.RUnlock()
	return ok
}

// Remove deletes <item> from set.
func (set *Set) Remove(item interface{}) *Set {
	set.mu.Lock()
	if set.data != nil {
		delete(set.data, item)
	}
	set.mu.Unlock()
	return set
}

// Size returns the size of the set.
func (set *Set) Size() int {
	set.mu.RLock()
	l := len(set.data)
	set.mu.RUnlock()
	return l
}

// Clear deletes all items of the set.
func (set *Set) Clear() *Set {
	set.mu.Lock()
	set.data = make(map[interface{}]struct{})
	set.mu.Unlock()
	return set
}

// Slice returns the a of items of the set as slice.
func (set *Set) Slice() []interface{} {
	set.mu.RLock()
	var (
		i   = 0
		ret = make([]interface{}, len(set.data))
	)
	for item := range set.data {
		ret[i] = item
		i++
	}
	set.mu.RUnlock()
	return ret
}

// Join joins items with a string <glue>.
func (set *Set) Join(glue string) string {
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
	for k, _ := range set.data {
		buffer.WriteString(gconv.String(k))
		if i != l-1 {
			buffer.WriteString(glue)
		}
		i++
	}
	return buffer.String()
}

// String returns items as a string, which implements like json.Marshal does.
func (set *Set) String() string {
	set.mu.RLock()
	defer set.mu.RUnlock()
	var (
		s      = ""
		l      = len(set.data)
		i      = 0
		buffer = bytes.NewBuffer(nil)
	)
	buffer.WriteByte('[')
	for k, _ := range set.data {
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

// LockFunc locks writing with callback function <f>.
func (set *Set) LockFunc(f func(m map[interface{}]struct{})) {
	set.mu.Lock()
	defer set.mu.Unlock()
	f(set.data)
}

// RLockFunc locks reading with callback function <f>.
func (set *Set) RLockFunc(f func(m map[interface{}]struct{})) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	f(set.data)
}

// Equal checks whether the two sets equal.
func (set *Set) Equal(other *Set) bool {
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

// IsSubsetOf checks whether the current set is a sub-set of <other>.
func (set *Set) IsSubsetOf(other *Set) bool {
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

// Union returns a new set which is the union of <set> and <others>.
// Which means, all the items in <newSet> are in <set> or in <others>.
func (set *Set) Union(others ...*Set) (newSet *Set) {
	newSet = NewSet()
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

// Diff returns a new set which is the difference set from <set> to <others>.
// Which means, all the items in <newSet> are in <set> but not in <others>.
func (set *Set) Diff(others ...*Set) (newSet *Set) {
	newSet = NewSet()
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

// Intersect returns a new set which is the intersection from <set> to <others>.
// Which means, all the items in <newSet> are in <set> and also in <others>.
func (set *Set) Intersect(others ...*Set) (newSet *Set) {
	newSet = NewSet()
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

// Complement returns a new set which is the complement from <set> to <full>.
// Which means, all the items in <newSet> are in <full> and not in <set>.
//
// It returns the difference between <full> and <set>
// if the given set <full> is not the full set of <set>.
func (set *Set) Complement(full *Set) (newSet *Set) {
	newSet = NewSet()
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

// Merge adds items from <others> sets into <set>.
func (set *Set) Merge(others ...*Set) *Set {
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
func (set *Set) Sum() (sum int) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k, _ := range set.data {
		sum += gconv.Int(k)
	}
	return
}

// Pops randomly pops an item from set.
func (set *Set) Pop() interface{} {
	set.mu.Lock()
	defer set.mu.Unlock()
	for k, _ := range set.data {
		delete(set.data, k)
		return k
	}
	return nil
}

// Pops randomly pops <size> items from set.
// It returns all items if size == -1.
func (set *Set) Pops(size int) []interface{} {
	set.mu.Lock()
	defer set.mu.Unlock()
	if size > len(set.data) || size == -1 {
		size = len(set.data)
	}
	if size <= 0 {
		return nil
	}
	index := 0
	array := make([]interface{}, size)
	for k, _ := range set.data {
		delete(set.data, k)
		array[index] = k
		index++
		if index == size {
			break
		}
	}
	return array
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (set *Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Slice())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (set *Set) UnmarshalJSON(b []byte) error {
	set.mu.Lock()
	defer set.mu.Unlock()
	if set.data == nil {
		set.data = make(map[interface{}]struct{})
	}
	var array []interface{}
	if err := json.Unmarshal(b, &array); err != nil {
		return err
	}
	for _, v := range array {
		set.data[v] = struct{}{}
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for set.
func (set *Set) UnmarshalValue(value interface{}) (err error) {
	set.mu.Lock()
	defer set.mu.Unlock()
	if set.data == nil {
		set.data = make(map[interface{}]struct{})
	}
	var array []interface{}
	switch value.(type) {
	case string, []byte:
		err = json.Unmarshal(gconv.Bytes(value), &array)
	default:
		array = gconv.SliceAny(value)
	}
	for _, v := range array {
		set.data[v] = struct{}{}
	}
	return
}

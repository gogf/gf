// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gset

import (
	"bytes"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/internal/rwmutex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type StrSet struct {
	mu   rwmutex.RWMutex
	data map[string]struct{}
}

// New create and returns a new set, which contains un-repeated items.
// The parameter <safe> is used to specify whether using set in concurrent-safety,
// which is false in default.
func NewStrSet(safe ...bool) *StrSet {
	return &StrSet{
		mu:   rwmutex.Create(safe...),
		data: make(map[string]struct{}),
	}
}

// NewStrSetFrom returns a new set from <items>.
func NewStrSetFrom(items []string, safe ...bool) *StrSet {
	m := make(map[string]struct{})
	for _, v := range items {
		m[v] = struct{}{}
	}
	return &StrSet{
		mu:   rwmutex.Create(safe...),
		data: m,
	}
}

// Iterator iterates the set readonly with given callback function <f>,
// if <f> returns true then continue iterating; or false to stop.
func (set *StrSet) Iterator(f func(v string) bool) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k, _ := range set.data {
		if !f(k) {
			break
		}
	}
}

// Add adds one or multiple items to the set.
func (set *StrSet) Add(item ...string) {
	set.mu.Lock()
	if set.data == nil {
		set.data = make(map[string]struct{})
	}
	for _, v := range item {
		set.data[v] = struct{}{}
	}
	set.mu.Unlock()
}

// AddIfNotExist checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set,
// or else it does nothing and returns false.
func (set *StrSet) AddIfNotExist(item string) bool {
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[string]struct{})
		}
		if _, ok := set.data[item]; !ok {
			set.data[item] = struct{}{}
			return true
		}
	}
	return false
}

// AddIfNotExistFunc checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and
// function <f> returns true, or else it does nothing and returns false.
//
// Note that, the function <f> is executed without writing lock.
func (set *StrSet) AddIfNotExistFunc(item string, f func() bool) bool {
	if !set.Contains(item) {
		if f() {
			set.mu.Lock()
			defer set.mu.Unlock()
			if set.data == nil {
				set.data = make(map[string]struct{})
			}
			if _, ok := set.data[item]; !ok {
				set.data[item] = struct{}{}
				return true
			}
		}
	}
	return false
}

// AddIfNotExistFunc checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and
// function <f> returns true, or else it does nothing and returns false.
//
// Note that, the function <f> is executed without writing lock.
func (set *StrSet) AddIfNotExistFuncLock(item string, f func() bool) bool {
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[string]struct{})
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

// Contains checks whether the set contains <item>.
func (set *StrSet) Contains(item string) bool {
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
func (set *StrSet) ContainsI(item string) bool {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k, _ := range set.data {
		if strings.EqualFold(k, item) {
			return true
		}
	}
	return false
}

// Remove deletes <item> from set.
func (set *StrSet) Remove(item string) {
	set.mu.Lock()
	if set.data != nil {
		delete(set.data, item)
	}
	set.mu.Unlock()
}

// Size returns the size of the set.
func (set *StrSet) Size() int {
	set.mu.RLock()
	l := len(set.data)
	set.mu.RUnlock()
	return l
}

// Clear deletes all items of the set.
func (set *StrSet) Clear() {
	set.mu.Lock()
	set.data = make(map[string]struct{})
	set.mu.Unlock()
}

// Slice returns the a of items of the set as slice.
func (set *StrSet) Slice() []string {
	set.mu.RLock()
	var (
		i   = 0
		ret = make([]string, len(set.data))
	)
	for item := range set.data {
		ret[i] = item
		i++
	}

	set.mu.RUnlock()
	return ret
}

// Join joins items with a string <glue>.
func (set *StrSet) Join(glue string) string {
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
		buffer.WriteString(k)
		if i != l-1 {
			buffer.WriteString(glue)
		}
		i++
	}
	return buffer.String()
}

// String returns items as a string, which implements like json.Marshal does.
func (set *StrSet) String() string {
	set.mu.RLock()
	defer set.mu.RUnlock()
	var (
		l      = len(set.data)
		i      = 0
		buffer = bytes.NewBuffer(nil)
	)
	for k, _ := range set.data {
		buffer.WriteString(`"` + gstr.QuoteMeta(k, `"\`) + `"`)
		if i != l-1 {
			buffer.WriteByte(',')
		}
		i++
	}
	return buffer.String()
}

// LockFunc locks writing with callback function <f>.
func (set *StrSet) LockFunc(f func(m map[string]struct{})) {
	set.mu.Lock()
	defer set.mu.Unlock()
	f(set.data)
}

// RLockFunc locks reading with callback function <f>.
func (set *StrSet) RLockFunc(f func(m map[string]struct{})) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	f(set.data)
}

// Equal checks whether the two sets equal.
func (set *StrSet) Equal(other *StrSet) bool {
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
func (set *StrSet) IsSubsetOf(other *StrSet) bool {
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

// Union returns a new set which is the union of <set> and <other>.
// Which means, all the items in <newSet> are in <set> or in <other>.
func (set *StrSet) Union(others ...*StrSet) (newSet *StrSet) {
	newSet = NewStrSet()
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

// Diff returns a new set which is the difference set from <set> to <other>.
// Which means, all the items in <newSet> are in <set> but not in <other>.
func (set *StrSet) Diff(others ...*StrSet) (newSet *StrSet) {
	newSet = NewStrSet()
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

// Intersect returns a new set which is the intersection from <set> to <other>.
// Which means, all the items in <newSet> are in <set> and also in <other>.
func (set *StrSet) Intersect(others ...*StrSet) (newSet *StrSet) {
	newSet = NewStrSet()
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
func (set *StrSet) Complement(full *StrSet) (newSet *StrSet) {
	newSet = NewStrSet()
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
func (set *StrSet) Merge(others ...*StrSet) *StrSet {
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
func (set *StrSet) Sum() (sum int) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k, _ := range set.data {
		sum += gconv.Int(k)
	}
	return
}

// Pops randomly pops an item from set.
func (set *StrSet) Pop() string {
	set.mu.Lock()
	defer set.mu.Unlock()
	for k, _ := range set.data {
		delete(set.data, k)
		return k
	}
	return ""
}

// Pops randomly pops <size> items from set.
// It returns all items if size == -1.
func (set *StrSet) Pops(size int) []string {
	set.mu.Lock()
	defer set.mu.Unlock()
	if size > len(set.data) || size == -1 {
		size = len(set.data)
	}
	if size <= 0 {
		return nil
	}
	index := 0
	array := make([]string, size)
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

// Walk applies a user supplied function <f> to every item of set.
func (set *StrSet) Walk(f func(item string) string) *StrSet {
	set.mu.Lock()
	defer set.mu.Unlock()
	m := make(map[string]struct{}, len(set.data))
	for k, v := range set.data {
		m[f(k)] = v
	}
	set.data = m
	return set
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (set *StrSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Slice())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (set *StrSet) UnmarshalJSON(b []byte) error {
	set.mu.Lock()
	defer set.mu.Unlock()
	if set.data == nil {
		set.data = make(map[string]struct{})
	}
	var array []string
	if err := json.Unmarshal(b, &array); err != nil {
		return err
	}
	for _, v := range array {
		set.data[v] = struct{}{}
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for set.
func (set *StrSet) UnmarshalValue(value interface{}) (err error) {
	set.mu.Lock()
	defer set.mu.Unlock()
	if set.data == nil {
		set.data = make(map[string]struct{})
	}
	var array []string
	switch value.(type) {
	case string, []byte:
		err = json.Unmarshal(gconv.Bytes(value), &array)
	default:
		array = gconv.SliceStr(value)
	}
	for _, v := range array {
		set.data[v] = struct{}{}
	}
	return
}

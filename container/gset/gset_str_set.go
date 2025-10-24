// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gset

import (
	"strings"
	"sync"
)

// StrSet is consisted of string items.
type StrSet struct {
	*TSet[string]
	once sync.Once
}

// NewStrSet create and returns a new set, which contains un-repeated items.
// The parameter `safe` is used to specify whether using set in concurrent-safety,
// which is false in default.
func NewStrSet(safe ...bool) *StrSet {
	return &StrSet{
		TSet: NewTSet[string](safe...),
	}
}

// NewStrSetFrom returns a new set from `items`.
func NewStrSetFrom(items []string, safe ...bool) *StrSet {
	return &StrSet{
		TSet: NewTSetFrom(items, safe...),
	}
}

// lazyInit lazily initializes the set.
func (a *StrSet) lazyInit() {
	a.once.Do(func() {
		if a.TSet == nil {
			a.TSet = NewTSet[string]()
		}
	})
}

// Iterator iterates the set readonly with given callback function `f`,
// if `f` returns true then continue iterating; or false to stop.
func (set *StrSet) Iterator(f func(v string) bool) {
	set.lazyInit()
	set.TSet.Iterator(f)
}

// Add adds one or multiple items to the set.
func (set *StrSet) Add(item ...string) {
	set.lazyInit()
	set.TSet.Add(item...)
}

// AddIfNotExist checks whether item exists in the set,
// it adds the item to set and returns true if it does not exist in the set,
// or else it does nothing and returns false.
func (set *StrSet) AddIfNotExist(item string) bool {
	set.lazyInit()
	return set.TSet.AddIfNotExist(item)
}

// AddIfNotExistFunc checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and
// function `f` returns true, or else it does nothing and returns false.
//
// Note that, the function `f` is executed without writing lock.
func (set *StrSet) AddIfNotExistFunc(item string, f func() bool) bool {
	set.lazyInit()
	return set.TSet.AddIfNotExistFunc(item, f)
}

// AddIfNotExistFuncLock checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and
// function `f` returns true, or else it does nothing and returns false.
//
// Note that, the function `f` is executed without writing lock.
func (set *StrSet) AddIfNotExistFuncLock(item string, f func() bool) bool {
	set.lazyInit()
	return set.TSet.AddIfNotExistFuncLock(item, f)
}

// Contains checks whether the set contains `item`.
func (set *StrSet) Contains(item string) bool {
	set.lazyInit()
	return set.TSet.Contains(item)
}

// ContainsI checks whether a value exists in the set with case-insensitively.
// Note that it internally iterates the whole set to do the comparison with case-insensitively.
func (set *StrSet) ContainsI(item string) bool {
	set.lazyInit()
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k := range set.data {
		if strings.EqualFold(k, item) {
			return true
		}
	}
	return false
}

// Remove deletes `item` from set.
func (set *StrSet) Remove(item string) {
	set.lazyInit()
	set.TSet.Remove(item)
}

// Size returns the size of the set.
func (set *StrSet) Size() int {
	set.lazyInit()
	return set.TSet.Size()
}

// Clear deletes all items of the set.
func (set *StrSet) Clear() {
	set.lazyInit()
	set.TSet.Clear()
}

// Slice returns the an of items of the set as slice.
func (set *StrSet) Slice() []string {
	set.lazyInit()
	return set.TSet.Slice()
}

// Join joins items with a string `glue`.
func (set *StrSet) Join(glue string) string {
	set.lazyInit()
	return set.TSet.Join(glue)
}

// String returns items as a string, which implements like json.Marshal does.
func (set *StrSet) String() string {
	if set == nil {
		return ""
	}
	set.lazyInit()
	return set.TSet.String()
}

// LockFunc locks writing with callback function `f`.
func (set *StrSet) LockFunc(f func(m map[string]struct{})) {
	set.lazyInit()
	set.TSet.LockFunc(f)
}

// RLockFunc locks reading with callback function `f`.
func (set *StrSet) RLockFunc(f func(m map[string]struct{})) {
	set.lazyInit()
	set.TSet.RLockFunc(f)
}

// Equal checks whether the two sets equal.
func (set *StrSet) Equal(other *StrSet) bool {
	set.lazyInit()
	other.lazyInit()
	return set.TSet.Equal(other.TSet)
}

// IsSubsetOf checks whether the current set is a sub-set of `other`.
func (set *StrSet) IsSubsetOf(other *StrSet) bool {
	if set == other {
		return true
	}

	set.lazyInit()
	other.lazyInit()

	return set.TSet.IsSubsetOf(other.TSet)
}

// Union returns a new set which is the union of `set` and `other`.
// Which means, all the items in `newSet` are in `set` or in `other`.
func (set *StrSet) Union(others ...*StrSet) (newSet *StrSet) {
	set.lazyInit()
	return &StrSet{
		TSet: set.TSet.Union(set.toTSetSlice(others)...),
	}
}

// Diff returns a new set which is the difference set from `set` to `other`.
// Which means, all the items in `newSet` are in `set` but not in `other`.
func (set *StrSet) Diff(others ...*StrSet) (newSet *StrSet) {
	set.lazyInit()
	return &StrSet{
		TSet: set.TSet.Diff(set.toTSetSlice(others)...),
	}
}

// Intersect returns a new set which is the intersection from `set` to `other`.
// Which means, all the items in `newSet` are in `set` and also in `other`.
func (set *StrSet) Intersect(others ...*StrSet) (newSet *StrSet) {
	set.lazyInit()
	return &StrSet{
		TSet: set.TSet.Intersect(set.toTSetSlice(others)...),
	}
}

// Complement returns a new set which is the complement from `set` to `full`.
// Which means, all the items in `newSet` are in `full` and not in `set`.
//
// It returns the difference between `full` and `set`
// if the given set `full` is not the full set of `set`.
func (set *StrSet) Complement(full *StrSet) (newSet *StrSet) {
	set.lazyInit()
	return &StrSet{
		TSet: set.TSet.Complement(full.TSet),
	}
}

// Merge adds items from `others` sets into `set`.
func (set *StrSet) Merge(others ...*StrSet) *StrSet {
	set.lazyInit()
	set.TSet.Merge(set.toTSetSlice(others)...)
	return set
}

// Sum sums items.
// Note: The items should be converted to int type,
// or you'd get a result that you unexpected.
func (set *StrSet) Sum() (sum int) {
	set.lazyInit()
	return set.TSet.Sum()
}

// Pop randomly pops an item from set.
func (set *StrSet) Pop() string {
	set.lazyInit()
	return set.TSet.Pop()
}

// Pops randomly pops `size` items from set.
// It returns all items if size == -1.
func (set *StrSet) Pops(size int) []string {
	set.lazyInit()
	return set.TSet.Pops(size)
}

// Walk applies a user supplied function `f` to every item of set.
func (set *StrSet) Walk(f func(item string) string) *StrSet {
	set.lazyInit()
	set.TSet.Walk(f)
	return set
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (set StrSet) MarshalJSON() ([]byte, error) {
	set.lazyInit()
	return set.TSet.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (set *StrSet) UnmarshalJSON(b []byte) error {
	set.lazyInit()
	return set.TSet.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for set.
func (set *StrSet) UnmarshalValue(value any) (err error) {
	set.lazyInit()
	return set.TSet.UnmarshalValue(value)
}

// DeepCopy implements interface for deep copy of current type.
func (set *StrSet) DeepCopy() any {
	if set == nil {
		return nil
	}
	set.lazyInit()
	return &StrSet{
		TSet: set.TSet.DeepCopy().(*TSet[string]),
	}
}

// toTSetSlice converts []*StrSet to []*TSet[string]
func (set *StrSet) toTSetSlice(sets []*StrSet) (tSets []*TSet[string]) {
	tSets = make([]*TSet[string], len(sets))
	for i, v := range sets {
		if v == nil {
			continue
		}
		v.lazyInit()
		tSets[i] = v.TSet
	}
	return
}

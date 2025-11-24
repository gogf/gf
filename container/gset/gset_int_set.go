// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gset

import (
	"sync"
)

// IntSet is consisted of int items.
type IntSet struct {
	*TSet[int]
	once sync.Once
}

// NewIntSet create and returns a new set, which contains un-repeated items.
// The parameter `safe` is used to specify whether using set in concurrent-safety,
// which is false in default.
func NewIntSet(safe ...bool) *IntSet {
	return &IntSet{
		TSet: NewTSet[int](safe...),
	}
}

// NewIntSetFrom returns a new set from `items`.
func NewIntSetFrom(items []int, safe ...bool) *IntSet {
	return &IntSet{
		TSet: NewTSetFrom(items, safe...),
	}
}

// lazyInit lazily initializes the set.
func (a *IntSet) lazyInit() {
	a.once.Do(func() {
		if a.TSet == nil {
			a.TSet = NewTSet[int]()
		}
	})
}

// Iterator iterates the set readonly with given callback function `f`,
// if `f` returns true then continue iterating; or false to stop.
func (set *IntSet) Iterator(f func(v int) bool) {
	set.lazyInit()
	set.TSet.Iterator(f)
}

// Add adds one or multiple items to the set.
func (set *IntSet) Add(item ...int) {
	set.lazyInit()
	set.TSet.Add(item...)
}

// AddIfNotExist checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set,
// or else it does nothing and returns false.
//
// Note that, if `item` is nil, it does nothing and returns false.
func (set *IntSet) AddIfNotExist(item int) bool {
	set.lazyInit()
	return set.TSet.AddIfNotExist(item)
}

// AddIfNotExistFunc checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and
// function `f` returns true, or else it does nothing and returns false.
//
// Note that, the function `f` is executed without writing lock.
func (set *IntSet) AddIfNotExistFunc(item int, f func() bool) bool {
	set.lazyInit()
	return set.TSet.AddIfNotExistFunc(item, f)
}

// AddIfNotExistFuncLock checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and
// function `f` returns true, or else it does nothing and returns false.
//
// Note that, the function `f` is executed without writing lock.
func (set *IntSet) AddIfNotExistFuncLock(item int, f func() bool) bool {
	set.lazyInit()
	return set.TSet.AddIfNotExistFuncLock(item, f)
}

// Contains checks whether the set contains `item`.
func (set *IntSet) Contains(item int) bool {
	set.lazyInit()
	return set.TSet.Contains(item)
}

// Remove deletes `item` from set.
func (set *IntSet) Remove(item int) {
	set.lazyInit()
	set.TSet.Remove(item)
}

// Size returns the size of the set.
func (set *IntSet) Size() int {
	set.lazyInit()
	return set.TSet.Size()
}

// Clear deletes all items of the set.
func (set *IntSet) Clear() {
	set.lazyInit()
	set.TSet.Clear()
}

// Slice returns the an of items of the set as slice.
func (set *IntSet) Slice() []int {
	set.lazyInit()
	return set.TSet.Slice()
}

// Join joins items with a string `glue`.
func (set *IntSet) Join(glue string) string {
	set.lazyInit()
	return set.TSet.Join(glue)
}

// String returns items as a string, which implements like json.Marshal does.
func (set *IntSet) String() string {
	if set == nil {
		return ""
	}
	set.lazyInit()
	return set.TSet.String()
}

// LockFunc locks writing with callback function `f`.
func (set *IntSet) LockFunc(f func(m map[int]struct{})) {
	set.lazyInit()
	set.TSet.LockFunc(f)
}

// RLockFunc locks reading with callback function `f`.
func (set *IntSet) RLockFunc(f func(m map[int]struct{})) {
	set.lazyInit()
	set.TSet.RLockFunc(f)
}

// Equal checks whether the two sets equal.
func (set *IntSet) Equal(other *IntSet) bool {
	set.lazyInit()
	other.lazyInit()
	return set.TSet.Equal(other.TSet)
}

// IsSubsetOf checks whether the current set is a sub-set of `other`.
func (set *IntSet) IsSubsetOf(other *IntSet) bool {
	if set == other {
		return true
	}

	set.lazyInit()
	other.lazyInit()

	return set.TSet.IsSubsetOf(other.TSet)
}

// Union returns a new set which is the union of `set` and `other`.
// Which means, all the items in `newSet` are in `set` or in `other`.
func (set *IntSet) Union(others ...*IntSet) (newSet *IntSet) {
	set.lazyInit()
	return &IntSet{
		TSet: set.TSet.Union(set.toTSetSlice(others)...),
	}
}

// Diff returns a new set which is the difference set from `set` to `other`.
// Which means, all the items in `newSet` are in `set` but not in `other`.
func (set *IntSet) Diff(others ...*IntSet) (newSet *IntSet) {
	set.lazyInit()
	return &IntSet{
		TSet: set.TSet.Diff(set.toTSetSlice(others)...),
	}
}

// Intersect returns a new set which is the intersection from `set` to `other`.
// Which means, all the items in `newSet` are in `set` and also in `other`.
func (set *IntSet) Intersect(others ...*IntSet) (newSet *IntSet) {
	set.lazyInit()
	return &IntSet{
		TSet: set.TSet.Intersect(set.toTSetSlice(others)...),
	}
}

// Complement returns a new set which is the complement from `set` to `full`.
// Which means, all the items in `newSet` are in `full` and not in `set`.
//
// It returns the difference between `full` and `set`
// if the given set `full` is not the full set of `set`.
func (set *IntSet) Complement(full *IntSet) (newSet *IntSet) {
	set.lazyInit()
	if full == nil {
		return &IntSet{
			TSet: NewTSet[int](),
		}
	}
	full.lazyInit()
	return &IntSet{
		TSet: set.TSet.Complement(full.TSet),
	}
}

// Merge adds items from `others` sets into `set`.
func (set *IntSet) Merge(others ...*IntSet) *IntSet {
	set.lazyInit()
	set.TSet.Merge(set.toTSetSlice(others)...)
	return set
}

// Sum sums items.
// Note: The items should be converted to int type,
// or you'd get a result that you unexpected.
func (set *IntSet) Sum() (sum int) {
	set.lazyInit()
	return set.TSet.Sum()
}

// Pop randomly pops an item from set.
func (set *IntSet) Pop() int {
	set.lazyInit()
	return set.TSet.Pop()
}

// Pops randomly pops `size` items from set.
// It returns all items if size == -1.
func (set *IntSet) Pops(size int) []int {
	set.lazyInit()
	return set.TSet.Pops(size)
}

// Walk applies a user supplied function `f` to every item of set.
func (set *IntSet) Walk(f func(item int) int) *IntSet {
	set.lazyInit()
	set.TSet.Walk(f)
	return set
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (set IntSet) MarshalJSON() ([]byte, error) {
	set.lazyInit()
	return set.TSet.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (set *IntSet) UnmarshalJSON(b []byte) error {
	set.lazyInit()
	return set.TSet.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for set.
func (set *IntSet) UnmarshalValue(value any) (err error) {
	set.lazyInit()
	return set.TSet.UnmarshalValue(value)
}

// DeepCopy implements interface for deep copy of current type.
func (set *IntSet) DeepCopy() any {
	if set == nil {
		return nil
	}
	set.lazyInit()
	return &IntSet{
		TSet: set.TSet.DeepCopy().(*TSet[int]),
	}
}

// toTSetSlice converts []*IntSet to []*TSet[int]
func (set *IntSet) toTSetSlice(sets []*IntSet) (tSets []*TSet[int]) {
	tSets = make([]*TSet[int], len(sets))
	for i, v := range sets {
		if v == nil {
			continue
		}
		v.lazyInit()
		tSets[i] = v.TSet
	}
	return
}

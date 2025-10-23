// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gset provides kinds of concurrent-safe/unsafe sets.
package gset

import (
	"github.com/gogf/gf/v2/util/gconv"
)

// Set is consisted of any items.
type Set struct {
	*TSet[any]
}

// New create and returns a new set, which contains un-repeated items.
// The parameter `safe` is used to specify whether using set in concurrent-safety,
// which is false in default.
func New(safe ...bool) *Set {
	return NewSet(safe...)
}

// NewSet create and returns a new set, which contains un-repeated items.
// Also see New.
func NewSet(safe ...bool) *Set {
	return &Set{
		TSet: NewTSet[any](safe...),
	}
}

// NewFrom returns a new set from `items`.
// Parameter `items` can be either a variable of any type, or a slice.
func NewFrom(items any, safe ...bool) *Set {
	return &Set{
		TSet: NewTSetFrom[any](gconv.Interfaces(items), safe...),
	}
}

// lazyInit lazily initializes the set.
func (a *Set) lazyInit() {
	if a.TSet == nil {
		a.TSet = NewTSet[any]()
	}
}

// Iterator iterates the set readonly with given callback function `f`,
// if `f` returns true then continue iterating; or false to stop.
func (set *Set) Iterator(f func(v any) bool) {
	set.lazyInit()
	set.TSet.Iterator(f)
}

// Add adds one or multiple items to the set.
func (set *Set) Add(items ...any) {
	set.lazyInit()
	set.TSet.Add(items...)
}

// AddIfNotExist checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set,
// or else it does nothing and returns false.
//
// Note that, if `item` is nil, it does nothing and returns false.
func (set *Set) AddIfNotExist(item any) bool {
	set.lazyInit()
	return set.TSet.AddIfNotExist(item)
}

// AddIfNotExistFunc checks whether item exists in the set,
// it adds the item to set and returns true if it does not exist in the set and
// function `f` returns true, or else it does nothing and returns false.
//
// Note that, if `item` is nil, it does nothing and returns false. The function `f`
// is executed without writing lock.
func (set *Set) AddIfNotExistFunc(item any, f func() bool) bool {
	set.lazyInit()
	return set.TSet.AddIfNotExistFunc(item, f)
}

// AddIfNotExistFuncLock checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and
// function `f` returns true, or else it does nothing and returns false.
//
// Note that, if `item` is nil, it does nothing and returns false. The function `f`
// is executed within writing lock.
func (set *Set) AddIfNotExistFuncLock(item any, f func() bool) bool {
	set.lazyInit()
	return set.TSet.AddIfNotExistFuncLock(item, f)
}

// Contains checks whether the set contains `item`.
func (set *Set) Contains(item any) bool {
	set.lazyInit()
	return set.TSet.Contains(item)
}

// Remove deletes `item` from set.
func (set *Set) Remove(item any) {
	set.lazyInit()
	set.TSet.Remove(item)
}

// Size returns the size of the set.
func (set *Set) Size() int {
	set.lazyInit()
	return set.TSet.Size()
}

// Clear deletes all items of the set.
func (set *Set) Clear() {
	set.lazyInit()
	set.TSet.Clear()
}

// Slice returns all items of the set as slice.
func (set *Set) Slice() []any {
	set.lazyInit()
	return set.TSet.Slice()
}

// Join joins items with a string `glue`.
func (set *Set) Join(glue string) string {
	set.lazyInit()
	return set.TSet.Join(glue)
}

// String returns items as a string, which implements like json.Marshal does.
func (set *Set) String() string {
	if set == nil {
		return ""
	}
	set.lazyInit()
	return set.TSet.String()
}

// LockFunc locks writing with callback function `f`.
func (set *Set) LockFunc(f func(m map[any]struct{})) {
	set.lazyInit()
	set.TSet.LockFunc(f)
}

// RLockFunc locks reading with callback function `f`.
func (set *Set) RLockFunc(f func(m map[any]struct{})) {
	set.lazyInit()
	set.TSet.RLockFunc(f)
}

// Equal checks whether the two sets equal.
func (set *Set) Equal(other *Set) bool {
	set.lazyInit()
	other.lazyInit()
	return set.TSet.Equal(other.TSet)
}

// IsSubsetOf checks whether the current set is a sub-set of `other`.
func (set *Set) IsSubsetOf(other *Set) bool {
	if set == other {
		return true
	}

	set.lazyInit()
	other.lazyInit()

	return set.TSet.IsSubsetOf(other.TSet)
}

// Union returns a new set which is the union of `set` and `others`.
// Which means, all the items in `newSet` are in `set` or in `others`.
func (set *Set) Union(others ...*Set) (newSet *Set) {
	set.lazyInit()
	tOthers := make([]*TSet[any], len(others))
	for _, o := range others {
		tOthers = append(tOthers, o.TSet)
	}

	return &Set{
		TSet: set.TSet.Union(tOthers...),
	}
}

// Diff returns a new set which is the difference set from `set` to `others`.
// Which means, all the items in `newSet` are in `set` but not in `others`.
func (set *Set) Diff(others ...*Set) (newSet *Set) {
	set.lazyInit()
	tOthers := make([]*TSet[any], len(others))
	for _, o := range others {
		tOthers = append(tOthers, o.TSet)
	}
	return &Set{
		TSet: set.TSet.Diff(tOthers...),
	}
}

// Intersect returns a new set which is the intersection from `set` to `others`.
// Which means, all the items in `newSet` are in `set` and also in `others`.
func (set *Set) Intersect(others ...*Set) (newSet *Set) {
	set.lazyInit()
	tOthers := make([]*TSet[any], len(others))
	for _, o := range others {
		tOthers = append(tOthers, o.TSet)
	}
	return &Set{
		TSet: set.TSet.Intersect(tOthers...),
	}
}

// Complement returns a new set which is the complement from `set` to `full`.
// Which means, all the items in `newSet` are in `full` and not in `set`.
//
// It returns the difference between `full` and `set`
// if the given set `full` is not the full set of `set`.
func (set *Set) Complement(full *Set) (newSet *Set) {
	set.lazyInit()
	return &Set{
		TSet: set.TSet.Complement(full.TSet),
	}
}

// Merge adds items from `others` sets into `set`.
func (set *Set) Merge(others ...*Set) *Set {
	set.lazyInit()
	tOthers := make([]*TSet[any], len(others))
	for _, o := range others {
		tOthers = append(tOthers, o.TSet)
	}
	set.TSet.Merge(tOthers...)
	return set
}

// Sum sums items.
// Note: The items should be converted to int type,
// or you'd get a result that you unexpected.
func (set *Set) Sum() (sum int) {
	set.lazyInit()
	return set.TSet.Sum()
}

// Pop randomly pops an item from set.
func (set *Set) Pop() any {
	set.lazyInit()
	return set.TSet.Pop()
}

// Pops randomly pops `size` items from set.
// It returns all items if size == -1.
func (set *Set) Pops(size int) []any {
	set.lazyInit()
	return set.TSet.Pops(size)
}

// Walk applies a user supplied function `f` to every item of set.
func (set *Set) Walk(f func(item any) any) *Set {
	set.lazyInit()
	set.TSet.Walk(f)
	return set
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (set Set) MarshalJSON() ([]byte, error) {
	set.lazyInit()
	return set.TSet.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (set *Set) UnmarshalJSON(b []byte) error {
	set.lazyInit()
	return set.TSet.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for set.
func (set *Set) UnmarshalValue(value any) (err error) {
	set.lazyInit()
	return set.TSet.UnmarshalValue(value)
}

// DeepCopy implements interface for deep copy of current type.
func (set *Set) DeepCopy() any {
	if set == nil {
		return nil
	}
	set.lazyInit()
	return &Set{
		TSet: set.TSet.DeepCopy().(*TSet[any]),
	}
}

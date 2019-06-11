// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gset

import (
    "github.com/gogf/gf/g/internal/rwmutex"
    "github.com/gogf/gf/g/util/gconv"
    "strings"
)

type IntSet struct {
	mu *rwmutex.RWMutex
	m  map[int]struct{}
}

// New create and returns a new set, which contains un-repeated items.
// The param <unsafe> used to specify whether using set in un-concurrent-safety,
// which is false in default.
func NewIntSet(unsafe...bool) *IntSet {
	return &IntSet{
		m  : make(map[int]struct{}),
		mu : rwmutex.New(unsafe...),
	}
}

// NewIntSetFrom returns a new set from <items>.
func NewIntSetFrom(items []int, unsafe...bool) *IntSet {
	m := make(map[int]struct{})
	for _, v := range items {
		m[v] = struct{}{}
	}
	return &IntSet{
		m  : m,
		mu : rwmutex.New(unsafe...),
	}
}

// Iterator iterates the set with given callback function <f>,
// if <f> returns true then continue iterating; or false to stop.
func (set *IntSet) Iterator(f func (v int) bool) *IntSet {
    set.mu.RLock()
    defer set.mu.RUnlock()
	for k, _ := range set.m {
		if !f(k) {
			break
		}
	}
    return set
}

// Add adds one or multiple items to the set.
func (set *IntSet) Add(item...int) *IntSet {
	set.mu.Lock()
    for _, v := range item {
        set.m[v] = struct{}{}
    }
	set.mu.Unlock()
	return set
}

// Contains checks whether the set contains <item>.
func (set *IntSet) Contains(item int) bool {
	set.mu.RLock()
	_, exists := set.m[item]
	set.mu.RUnlock()
	return exists
}

// Remove deletes <item> from set.
func (set *IntSet) Remove(item int) *IntSet {
	set.mu.Lock()
	delete(set.m, item)
	set.mu.Unlock()
	return set
}

// Size returns the size of the set.
func (set *IntSet) Size() int {
	set.mu.RLock()
	l := len(set.m)
	set.mu.RUnlock()
	return l
}

// Clear deletes all items of the set.
func (set *IntSet) Clear() *IntSet {
	set.mu.Lock()
	set.m = make(map[int]struct{})
	set.mu.Unlock()
    return set
}

// Slice returns the a of items of the set as slice.
func (set *IntSet) Slice() []int {
	set.mu.RLock()
	ret := make([]int, len(set.m))
	i := 0
	for k, _ := range set.m {
		ret[i] = k
		i++
	}
	set.mu.RUnlock()
	return ret
}

// Join joins items with a string <glue>.
func (set *IntSet) Join(glue string) string {
    return strings.Join(gconv.Strings(set.Slice()), ",")
}

// String returns items as a string, which are joined by char ','.
func (set *IntSet) String() string {
    return set.Join(",")
}

// LockFunc locks writing with callback function <f>.
func (set *IntSet) LockFunc(f func(m map[int]struct{})) {
	set.mu.Lock()
	defer set.mu.Unlock()
	f(set.m)
}

// RLockFunc locks reading with callback function <f>.
func (set *IntSet) RLockFunc(f func(m map[int]struct{})) {
    set.mu.RLock()
    defer set.mu.RUnlock()
    f(set.m)
}

// Equal checks whether the two sets equal.
func (set *IntSet) Equal(other *IntSet) bool {
    if set == other {
        return true
    }
    set.mu.RLock()
    defer set.mu.RUnlock()
    other.mu.RLock()
    defer other.mu.RUnlock()
	if len(set.m) != len(other.m) {
		return false
	}
	for key := range set.m {
		if _, ok := other.m[key]; !ok {
			return false
		}
	}
	return true
}

// IsSubsetOf checks whether the current set is a sub-set of <other>.
func (set *IntSet) IsSubsetOf(other *IntSet) bool {
    if set == other {
        return true
    }
    set.mu.RLock()
    defer set.mu.RUnlock()
    other.mu.RLock()
    defer other.mu.RUnlock()
    for key := range set.m {
        if _, ok := other.m[key]; !ok {
            return false
        }
    }
    return true
}

// Union returns a new set which is the union of <set> and <other>.
// Which means, all the items in <newSet> are in <set> or in <other>.
func (set *IntSet) Union(others ... *IntSet) (newSet *IntSet) {
    newSet = NewIntSet(true)
    set.mu.RLock()
    defer set.mu.RUnlock()
    for _, other := range others {
        if set != other {
            other.mu.RLock()
        }
        for k, v := range set.m {
            newSet.m[k] = v
        }
        if set != other {
            for k, v := range other.m {
                newSet.m[k] = v
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
func (set *IntSet) Diff(others...*IntSet) (newSet *IntSet) {
    newSet = NewIntSet(true)
    set.mu.RLock()
    defer set.mu.RUnlock()
    for _, other := range others {
        if set == other {
            continue
        }
        other.mu.RLock()
        for k, v := range set.m {
            if _, ok := other.m[k]; !ok {
                newSet.m[k] = v
            }
        }
        other.mu.RUnlock()
    }
    return
}

// Intersect returns a new set which is the intersection from <set> to <other>.
// Which means, all the items in <newSet> are in <set> and also in <other>.
func (set *IntSet) Intersect(others...*IntSet) (newSet *IntSet) {
    newSet = NewIntSet(true)
    set.mu.RLock()
    defer set.mu.RUnlock()
    for _, other := range others {
        if set != other {
            other.mu.RLock()
        }
        for k, v := range set.m {
            if _, ok := other.m[k]; ok {
                newSet.m[k] = v
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
func (set *IntSet) Complement(full *IntSet) (newSet *IntSet) {
    newSet = NewIntSet(true)
    set.mu.RLock()
    defer set.mu.RUnlock()
    if set != full {
        full.mu.RLock()
        defer full.mu.RUnlock()
    }
    for k, v := range full.m {
        if _, ok := set.m[k]; !ok {
            newSet.m[k] = v
        }
    }
    return
}

// Merge adds items from <others> sets into <set>.
func (set *IntSet) Merge(others ... *IntSet) *IntSet {
	set.mu.Lock()
	defer set.mu.Unlock()
	for _, other := range others {
		if set != other {
			other.mu.RLock()
		}
		for k, v := range other.m {
			set.m[k] = v
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
func (set *IntSet) Sum() (sum int) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k, _ := range set.m {
		sum += k
	}
	return
}

// Pops randomly pops an item from set.
func (set *IntSet) Pop(size int) int {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k, _ := range set.m {
		return k
	}
	return 0
}

// Pops randomly pops <size> items from set.
func (set *IntSet) Pops(size int) []int {
	set.mu.RLock()
	defer set.mu.RUnlock()
	if size > len(set.m) {
		size = len(set.m)
	}
	index := 0
	array := make([]int, size)
	for k, _ := range set.m {
		array[index] = k
		index++
		if index == size {
			break
		}
	}
	return array
}
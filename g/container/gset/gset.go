// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gset provides kinds of concurrent-safe/unsafe sets.
package gset

import (
    "github.com/gogf/gf/g/internal/rwmutex"
    "github.com/gogf/gf/g/util/gconv"
    "strings"
)

type Set struct {
    mu *rwmutex.RWMutex
    m  map[interface{}]struct{}
}

// New create and returns a new set, which contains un-repeated items.
// The param <unsafe> used to specify whether using set in un-concurrent-safety,
// which is false in default.
func New(unsafe...bool) *Set {
    return NewSet(unsafe...)
}

// See New.
func NewSet(unsafe...bool) *Set {
    return &Set{
        m  : make(map[interface{}]struct{}),
        mu : rwmutex.New(unsafe...),
    }
}

// NewFrom returns a new set from <items>.
// Parameter <items> can be either a variable of any type, or a slice.
func NewFrom(items interface{}, unsafe...bool) *Set {
	m := make(map[interface{}]struct{})
	for _, v := range gconv.Interfaces(items) {
		m[v] = struct{}{}
	}
	return &Set{
		m  : m,
		mu : rwmutex.New(unsafe...),
	}
}

// Iterator iterates the set with given callback function <f>,
// if <f> returns true then continue iterating; or false to stop.
func (set *Set) Iterator(f func (v interface{}) bool) *Set {
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
func (set *Set) Add(item...interface{}) *Set {
    set.mu.Lock()
    for _, v := range item {
        set.m[v] = struct{}{}
    }
    set.mu.Unlock()
    return set
}

// Contains checks whether the set contains <item>.
func (set *Set) Contains(item interface{}) bool {
    set.mu.RLock()
    _, exists := set.m[item]
    set.mu.RUnlock()
    return exists
}

// Remove deletes <item> from set.
func (set *Set) Remove(item interface{}) *Set {
    set.mu.Lock()
    delete(set.m, item)
    set.mu.Unlock()
    return set
}

// Size returns the size of the set.
func (set *Set) Size() int {
    set.mu.RLock()
    l := len(set.m)
    set.mu.RUnlock()
    return l
}

// Clear deletes all items of the set.
func (set *Set) Clear() *Set {
    set.mu.Lock()
    set.m = make(map[interface{}]struct{})
    set.mu.Unlock()
    return set
}

// Slice returns the a of items of the set as slice.
func (set *Set) Slice() []interface{} {
    set.mu.RLock()
    i   := 0
    ret := make([]interface{}, len(set.m))
    for item := range set.m {
        ret[i] = item
        i++
    }
    set.mu.RUnlock()
    return ret
}

// Join joins items with a string <glue>.
func (set *Set) Join(glue string) string {
    return strings.Join(gconv.Strings(set.Slice()), ",")
}

// String returns items as a string, which are joined by char ','.
func (set *Set) String() string {
    return set.Join(",")
}

// LockFunc locks writing with callback function <f>.
func (set *Set) LockFunc(f func(m map[interface{}]struct{})) {
    set.mu.Lock()
    defer set.mu.Unlock()
    f(set.m)
}

// RLockFunc locks reading with callback function <f>.
func (set *Set) RLockFunc(f func(m map[interface{}]struct{})) {
    set.mu.RLock()
    defer set.mu.RUnlock()
    f(set.m)
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
func (set *Set) IsSubsetOf(other *Set) bool {
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

// Union returns a new set which is the union of <set> and <others>.
// Which means, all the items in <newSet> are in <set> or in <others>.
func (set *Set) Union(others ... *Set) (newSet *Set) {
    newSet = NewSet(true)
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

// Diff returns a new set which is the difference set from <set> to <others>.
// Which means, all the items in <newSet> are in <set> but not in <others>.
func (set *Set) Diff(others...*Set) (newSet *Set) {
    newSet = NewSet(true)
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

// Intersect returns a new set which is the intersection from <set> to <others>.
// Which means, all the items in <newSet> are in <set> and also in <others>.
func (set *Set) Intersect(others...*Set) (newSet *Set) {
    newSet = NewSet(true)
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
func (set *Set) Complement(full *Set) (newSet *Set) {
    newSet = NewSet(true)
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
func (set *Set) Merge(others ... *Set) *Set {
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
func (set *Set) Sum() (sum int) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k, _ := range set.m {
		sum += gconv.Int(k)
	}
	return
}

// Pops randomly pops an item from set.
func (set *Set) Pop(size int) interface{} {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k, _ := range set.m {
		return k
	}
	return nil
}

// Pops randomly pops <size> items from set.
func (set *Set) Pops(size int) []interface{} {
	set.mu.RLock()
	defer set.mu.RUnlock()
	if size > len(set.m) {
		size = len(set.m)
	}
	index := 0
	array := make([]interface{}, size)
	for k, _ := range set.m {
		array[index] = k
		index++
		if index == size {
			break
		}
	}
	return array
}
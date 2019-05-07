// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.
//

package gmap

import (
    "github.com/gogf/gf/g/internal/rwmutex"
)

// Deprecated, use Map instead.
type StringStringMap struct {
	mu *rwmutex.RWMutex
	m  map[string]string
}

// Deprecated, use New instead.
// NewStringStringMap returns an empty StringStringMap object.
// The param <unsafe> used to specify whether using map with un-concurrent-safety,
// which is false in default, means concurrent-safe.
func NewStringStringMap(unsafe...bool) *StringStringMap {
	return &StringStringMap{
		m  : make(map[string]string),
        mu : rwmutex.New(unsafe...),
	}
}

// Deprecated, use NewFrom instead.
// NewStringStringMapFrom returns an StringStringMap object from given map <m>.
// Notice that, the param map is a type of pointer,
// there might be some concurrent-safe issues when changing the map outside.
func NewStringStringMapFrom(m map[string]string, unsafe...bool) *StringStringMap {
    return &StringStringMap{
        m  : m,
        mu : rwmutex.New(unsafe...),
    }
}

// Deprecated, use NewFromArray instead.
// NewStringStringMapFromArray returns an StringStringMap object from given array.
// The param <keys> given as the keys of the map,
// and <values> as its corresponding values.
//
// If length of <keys> is greater than that of <values>,
// the corresponding overflow map values will be the default value of its type.
func NewStringStringMapFromArray(keys []string, values []string, unsafe...bool) *StringStringMap {
    m := make(map[string]string)
    l := len(values)
    for i, k := range keys {
        if i < l {
            m[k] = values[i]
        } else {
            m[k] = ""
        }
    }
    return &StringStringMap{
        m  : m,
        mu : rwmutex.New(unsafe...),
    }
}

// Iterator iterates the hash map with custom callback function <f>.
// If f returns true, then continue iterating; or false to stop.
func (gm *StringStringMap) Iterator(f func (k string, v string) bool) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	for k, v := range gm.m {
		if !f(k, v) {
			break
		}
	}
}

// Clone returns a new hash map with copy of current map data.
func (gm *StringStringMap) Clone() *StringStringMap {
    return NewStringStringMapFrom(gm.Map(), !gm.mu.IsSafe())
}

// Map returns a copy of the data of the hash map.
func (gm *StringStringMap) Map() map[string]string {
    m := make(map[string]string)
    gm.mu.RLock()
    for k, v := range gm.m {
        m[k] = v
    }
    gm.mu.RUnlock()
    return m
}

// Set sets key-value to the hash map.
func (gm *StringStringMap) Set(key string, val string) {
	gm.mu.Lock()
	gm.m[key] = val
	gm.mu.Unlock()
}

// BatchSet batch sets key-values to the hash map.
func (gm *StringStringMap) BatchSet(m map[string]string) {
    gm.mu.Lock()
    for k, v := range m {
        gm.m[k] = v
    }
    gm.mu.Unlock()
}

// Get returns the value by given <key>.
func (gm *StringStringMap) Get(key string) string {
	gm.mu.RLock()
	val, _ := gm.m[key]
	gm.mu.RUnlock()
	return val
}

// doSetWithLockCheck checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given <key>,
// or else just return the existing value.
//
// It returns value with given <key>.
func (gm *StringStringMap) doSetWithLockCheck(key string, value string) string {
	gm.mu.Lock()
	if v, ok := gm.m[key]; ok {
		gm.mu.Unlock()
		return v
	}
	gm.m[key] = value
	gm.mu.Unlock()
	return value
}

// GetOrSet returns the value by key,
// or set value with given <value> if not exist and returns this value.
func (gm *StringStringMap) GetOrSet(key string, value string) string {
	gm.mu.RLock()
	v, ok := gm.m[key]
	gm.mu.RUnlock()
	if !ok {
		return gm.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with return value of callback function <f> if not exist
// and returns this value.
func (gm *StringStringMap) GetOrSetFunc(key string, f func() string) string {
	gm.mu.RLock()
	v, ok := gm.m[key]
	gm.mu.RUnlock()
	if !ok {
		return gm.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock returns the value by key,
// or sets value with return value of callback function <f> if not exist
// and returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function <f>
// with mutex.Lock of the hash map.
func (gm *StringStringMap) GetOrSetFuncLock(key string, f func() string) string {
	gm.mu.RLock()
	val, ok := gm.m[key]
	gm.mu.RUnlock()
	if !ok {
		gm.mu.Lock()
		defer gm.mu.Unlock()
		if v, ok := gm.m[key]; ok {
			return v
		}
		val       = f()
		gm.m[key] = val
		return val
	} else {
		return val
	}
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (gm *StringStringMap) SetIfNotExist(key string, value string) bool {
	if !gm.Contains(key) {
		gm.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (gm *StringStringMap) SetIfNotExistFunc(key string, f func() string) bool {
	if !gm.Contains(key) {
		gm.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function <f>, then return true.
// It returns false if <key> exists, and <value> would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function <f> with mutex.Lock of the hash map.
func (gm *StringStringMap) SetIfNotExistFuncLock(key string, f func() string) bool {
	if !gm.Contains(key) {
		gm.mu.Lock()
		defer gm.mu.Unlock()
		if _, ok := gm.m[key]; !ok {
			gm.m[key] = f()
		}
		return true
	}
	return false
}

// BatchRemove batch deletes values of the map by keys.
func (gm *StringStringMap) BatchRemove(keys []string) {
    gm.mu.Lock()
    for _, key := range keys {
        delete(gm.m, key)
    }
    gm.mu.Unlock()
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (gm *StringStringMap) Remove(key string) string {
	gm.mu.Lock()
	val, exists := gm.m[key]
	if exists {
		delete(gm.m, key)
	}
	gm.mu.Unlock()
	return val
}

// Keys returns all keys of the map as a slice.
func (gm *StringStringMap) Keys() []string {
	gm.mu.RLock()
	keys := make([]string, 0)
	for key, _ := range gm.m {
		keys = append(keys, key)
	}
    gm.mu.RUnlock()
	return keys
}

// Values returns all values of the map as a slice.
func (gm *StringStringMap) Values() []string {
	gm.mu.RLock()
	vals := make([]string, 0)
	for _, val := range gm.m {
		vals = append(vals, val)
	}
	gm.mu.RUnlock()
	return vals
}

// Contains checks whether a key exists.
// It returns true if the <key> exists, or else false.
func (gm *StringStringMap) Contains(key string) bool {
	gm.mu.RLock()
	_, exists := gm.m[key]
	gm.mu.RUnlock()
	return exists
}

// Size returns the size of the map.
func (gm *StringStringMap) Size() int {
	gm.mu.RLock()
	length := len(gm.m)
	gm.mu.RUnlock()
	return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (gm *StringStringMap) IsEmpty() bool {
	gm.mu.RLock()
	empty := len(gm.m) == 0
	gm.mu.RUnlock()
	return empty
}

// Clear deletes all data of the map, it will remake a new underlying map data map.
func (gm *StringStringMap) Clear() {
    gm.mu.Lock()
    gm.m = make(map[string]string)
    gm.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> and mutex.Lock.
func (gm *StringStringMap) LockFunc(f func(m map[string]string)) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	f(gm.m)
}

// RLockFunc locks reading with given callback function <f> and mutex.RLock.
func (gm *StringStringMap) RLockFunc(f func(m map[string]string)) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	f(gm.m)
}

// Flip exchanges key-value of the map, it will change key-value to value-key.
func (gm *StringStringMap) Flip() {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	n := make(map[string]string, len(gm.m))
	for k, v := range gm.m {
		n[v] = k
	}
	gm.m = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <gm>.
func (gm *StringStringMap) Merge(other *StringStringMap) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	if other != gm {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	for k, v := range other.m {
		gm.m[k] = v
	}
}

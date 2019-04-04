// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmap provides concurrent-safe/unsafe maps.
package gmap

import "github.com/gogf/gf/g/internal/rwmutex"

type Map struct {
    mu *rwmutex.RWMutex
    m  map[interface{}]interface{}
}

// New returns an empty hash map.
// The param <unsafe> used to specify whether using map with un-concurrent-safety,
// which is false in default, means concurrent-safe.
func New(unsafe ...bool) *Map {
    return NewMap(unsafe...)
}

// Alias of New. See New.
func NewMap(unsafe ...bool) *Map {
    return &Map{
        m  : make(map[interface{}]interface{}),
        mu : rwmutex.New(unsafe...),
    }
}

// NewFrom returns a hash map from given map <m>.
// Notice that, the param map is a type of pointer,
// there might be some concurrent-safe issues when changing the map outside.
func NewFrom(m map[interface{}]interface{}, unsafe...bool) *Map {
    return &Map{
        m  : m,
        mu : rwmutex.New(unsafe...),
    }
}

// NewFromArray returns a hash map from given array.
// The param <keys> given as the keys of the map,
// and <values> as its corresponding values.
//
// If length of <keys> is greater than that of <values>,
// the corresponding overflow map values will be the default value of its type.
func NewFromArray(keys []interface{}, values []interface{}, unsafe...bool) *Map {
    m := make(map[interface{}]interface{})
    l := len(values)
    for i, k := range keys {
        if i < l {
            m[k] = values[i]
        } else {
            m[k] = interface{}(nil)
        }
    }
    return &Map{
        m  : m,
        mu : rwmutex.New(unsafe...),
    }
}

// Iterator iterates the hash map with custom callback function <f>.
// If f returns true, then continue iterating; or false to stop.
func (gm *Map) Iterator(f func (k interface{}, v interface{}) bool) {
    gm.mu.RLock()
    defer gm.mu.RUnlock()
    for k, v := range gm.m {
        if !f(k, v) {
            break
        }
    }
}

// Clone returns a new hash map with copy of current map data.
func (gm *Map) Clone(unsafe ...bool) *Map {
    return NewFrom(gm.Map(), unsafe ...)
}

// Map returns a copy of the data of the hash map.
func (gm *Map) Map() map[interface{}]interface{} {
    m := make(map[interface{}]interface{})
    gm.mu.RLock()
    for k, v := range gm.m {
        m[k] = v
    }
    gm.mu.RUnlock()
    return m
}

// Set sets key-value to the hash map.
func (gm *Map) Set(key interface{}, val interface{}) {
    gm.mu.Lock()
    gm.m[key] = val
    gm.mu.Unlock()
}

// BatchSet batch sets key-values to the hash map.
func (gm *Map) BatchSet(m map[interface{}]interface{}) {
    gm.mu.Lock()
    for k, v := range m {
        gm.m[k] = v
    }
    gm.mu.Unlock()
}

// Get returns the value by given <key>.
func (gm *Map) Get(key interface{}) interface{} {
    gm.mu.RLock()
    val, _ := gm.m[key]
    gm.mu.RUnlock()
    return val
}

// doSetWithLockCheck checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given <key>,
// or else just return the existing value.
//
// When setting value, if <value> is type of <func() interface {}>,
// it will be executed with mutex.Lock of the hash map,
// and its return value will be set to the map with <key>.
//
// It returns value with given <key>.
func (gm *Map) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
    gm.mu.Lock()
    defer gm.mu.Unlock()
    if v, ok := gm.m[key]; ok {
        return v
    }
    if f, ok := value.(func() interface {}); ok {
        value = f()
    }
    gm.m[key] = value
    return value
}

// GetOrSet returns the value by key,
// or set value with given <value> if not exist and returns this value.
func (gm *Map) GetOrSet(key interface{}, value interface{}) interface{} {
    if v := gm.Get(key); v == nil {
        return gm.doSetWithLockCheck(key, value)
    } else {
        return v
    }
}

// GetOrSetFunc returns the value by key,
// or sets value with return value of callback function <f> if not exist
// and returns this value.
func (gm *Map) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
    if v := gm.Get(key); v == nil {
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
func (gm *Map) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
    if v := gm.Get(key); v == nil {
        return gm.doSetWithLockCheck(key, f)
    } else {
        return v
    }
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (gm *Map) SetIfNotExist(key interface{}, value interface{}) bool {
    if !gm.Contains(key) {
        gm.doSetWithLockCheck(key, value)
        return true
    }
    return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (gm *Map) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
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
func (gm *Map) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !gm.Contains(key) {
		gm.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

// BatchRemove batch deletes values of the map by keys.
func (gm *Map) BatchRemove(keys []interface{}) {
    gm.mu.Lock()
    for _, key := range keys {
        delete(gm.m, key)
    }
    gm.mu.Unlock()
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (gm *Map) Remove(key interface{}) interface{} {
    gm.mu.Lock()
    val, exists := gm.m[key]
    if exists {
        delete(gm.m, key)
    }
    gm.mu.Unlock()
    return val
}

// Keys returns all keys of the map as a slice.
func (gm *Map) Keys() []interface{} {
    gm.mu.RLock()
    keys := make([]interface{}, 0)
    for key, _ := range gm.m {
        keys = append(keys, key)
    }
    gm.mu.RUnlock()
    return keys
}

// Values returns all values of the map as a slice.
func (gm *Map) Values() []interface{} {
    gm.mu.RLock()
    vals := make([]interface{}, 0)
    for _, val := range gm.m {
        vals = append(vals, val)
    }
    gm.mu.RUnlock()
    return vals
}

// Contains checks whether a key exists.
// It returns true if the <key> exists, or else false.
func (gm *Map) Contains(key interface{}) bool {
    gm.mu.RLock()
    _, exists := gm.m[key]
    gm.mu.RUnlock()
    return exists
}

// Size returns the size of the map.
func (gm *Map) Size() int {
    gm.mu.RLock()
    length := len(gm.m)
    gm.mu.RUnlock()
    return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (gm *Map) IsEmpty() bool {
    gm.mu.RLock()
    empty := len(gm.m) == 0
    gm.mu.RUnlock()
    return empty
}

// Clear deletes all data of the map, it will remake a new underlying map data map.
func (gm *Map) Clear() {
    gm.mu.Lock()
    gm.m = make(map[interface{}]interface{})
    gm.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> and mutex.Lock.
func (gm *Map) LockFunc(f func(m map[interface{}]interface{})) {
    gm.mu.Lock()
    defer gm.mu.Unlock()
    f(gm.m)
}

// RLockFunc locks reading with given callback function <f> and mutex.RLock.
func (gm *Map) RLockFunc(f func(m map[interface{}]interface{})) {
    gm.mu.RLock()
    defer gm.mu.RUnlock()
    f(gm.m)
}

// Flip exchanges key-value of the map, it will change key-value to value-key.
func (gm *Map) Flip() {
    gm.mu.Lock()
    defer gm.mu.Unlock()
    n := make(map[interface{}]interface{}, len(gm.m))
    for i, v := range gm.m {
        n[v] = i
    }
    gm.m = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <gm>.
func (gm *Map) Merge(other *Map) {
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
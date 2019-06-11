// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import (
	"github.com/gogf/gf/g/container/gvar"
	"github.com/gogf/gf/g/internal/rwmutex"
)

type AnyAnyMap struct {
    mu   *rwmutex.RWMutex
    data map[interface{}]interface{}
}

// NewAnyAnyMap returns an empty hash map.
// The param <unsafe> used to specify whether using map in un-concurrent-safety,
// which is false in default, means concurrent-safe.
func NewAnyAnyMap(unsafe ...bool) *AnyAnyMap {
	return &AnyAnyMap{
		mu   : rwmutex.New(unsafe...),
		data : make(map[interface{}]interface{}),
	}
}

// NewAnyAnyMapFrom returns a hash map from given map <data>.
// Note that, the param <data> map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewAnyAnyMapFrom(data map[interface{}]interface{}, unsafe...bool) *AnyAnyMap {
    return &AnyAnyMap{
        mu   : rwmutex.New(unsafe...),
	    data : data,
    }
}

// Iterator iterates the hash map with custom callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (m *AnyAnyMap) Iterator(f func (k interface{}, v interface{}) bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    for k, v := range m.data {
        if !f(k, v) {
            break
        }
    }
}

// Clone returns a new hash map with copy of current map data.
func (m *AnyAnyMap) Clone(unsafe ...bool) *AnyAnyMap {
    return NewFrom(m.Map(), unsafe ...)
}

// Map returns a copy of the data of the hash map.
func (m *AnyAnyMap) Map() map[interface{}]interface{} {
    m.mu.RLock()
	data := make(map[interface{}]interface{}, len(m.data))
    for k, v := range m.data {
	    data[k] = v
    }
    m.mu.RUnlock()
    return data
}

// Set sets key-value to the hash map.
func (m *AnyAnyMap) Set(key interface{}, val interface{}) {
    m.mu.Lock()
    m.data[key] = val
    m.mu.Unlock()
}

// Sets batch sets key-values to the hash map.
func (m *AnyAnyMap) Sets(data map[interface{}]interface{}) {
    m.mu.Lock()
    for k, v := range data {
        m.data[k] = v
    }
    m.mu.Unlock()
}

// Search searches the map with given <key>.
// Second return parameter <found> is true if key was found, otherwise false.
func (m *AnyAnyMap) Search(key interface{}) (value interface{}, found bool) {
	m.mu.RLock()
	value, found = m.data[key]
	m.mu.RUnlock()
	return
}

// Get returns the value by given <key>.
func (m *AnyAnyMap) Get(key interface{}) interface{} {
    m.mu.RLock()
    val, _ := m.data[key]
    m.mu.RUnlock()
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
func (m *AnyAnyMap) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
    m.mu.Lock()
    defer m.mu.Unlock()
    if v, ok := m.data[key]; ok {
        return v
    }
    if f, ok := value.(func() interface {}); ok {
        value = f()
    }
    m.data[key] = value
    return value
}

// GetOrSet returns the value by key,
// or set value with given <value> if not exist and returns this value.
func (m *AnyAnyMap) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
        return m.doSetWithLockCheck(key, value)
    } else {
        return v
    }
}

// GetOrSetFunc returns the value by key,
// or sets value with return value of callback function <f> if not exist
// and returns this value.
func (m *AnyAnyMap) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
        return m.doSetWithLockCheck(key, f())
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
func (m *AnyAnyMap) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
        return m.doSetWithLockCheck(key, f)
    } else {
        return v
    }
}

// GetVar returns a gvar.Var with the value by given <key>.
// The returned gvar.Var is un-concurrent safe.
func (m *AnyAnyMap) GetVar(key interface{}) *gvar.Var {
	return gvar.New(m.Get(key), true)
}

// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
// The returned gvar.Var is un-concurrent safe.
func (m *AnyAnyMap) GetVarOrSet(key interface{}, value interface{}) *gvar.Var {
	return gvar.New(m.GetOrSet(key, value), true)
}

// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
// The returned gvar.Var is un-concurrent safe.
func (m *AnyAnyMap) GetVarOrSetFunc(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(m.GetOrSetFunc(key, f), true)
}

// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
// The returned gvar.Var is un-concurrent safe.
func (m *AnyAnyMap) GetVarOrSetFuncLock(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(m.GetOrSetFuncLock(key, f), true)
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *AnyAnyMap) SetIfNotExist(key interface{}, value interface{}) bool {
    if !m.Contains(key) {
        m.doSetWithLockCheck(key, value)
        return true
    }
    return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *AnyAnyMap) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function <f>, then return true.
// It returns false if <key> exists, and <value> would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function <f> with mutex.Lock of the hash map.
func (m *AnyAnyMap) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (m *AnyAnyMap) Remove(key interface{}) interface{} {
    m.mu.Lock()
    val, exists := m.data[key]
    if exists {
        delete(m.data, key)
    }
    m.mu.Unlock()
    return val
}

// Removes batch deletes values of the map by keys.
func (m *AnyAnyMap) Removes(keys []interface{}) {
	m.mu.Lock()
	for _, key := range keys {
		delete(m.data, key)
	}
	m.mu.Unlock()
}

// Keys returns all keys of the map as a slice.
func (m *AnyAnyMap) Keys() []interface{} {
    m.mu.RLock()
    keys  := make([]interface{}, len(m.data))
    index := 0
    for key := range m.data {
        keys[index] = key
        index++
    }
    m.mu.RUnlock()
    return keys
}

// Values returns all values of the map as a slice.
func (m *AnyAnyMap) Values() []interface{} {
    m.mu.RLock()
    values := make([]interface{}, len(m.data))
	index  := 0
    for _, value := range m.data {
        values[index] = value
        index++
    }
    m.mu.RUnlock()
    return values
}

// Contains checks whether a key exists.
// It returns true if the <key> exists, or else false.
func (m *AnyAnyMap) Contains(key interface{}) bool {
    m.mu.RLock()
    _, exists := m.data[key]
    m.mu.RUnlock()
    return exists
}

// Size returns the size of the map.
func (m *AnyAnyMap) Size() int {
    m.mu.RLock()
    length := len(m.data)
    m.mu.RUnlock()
    return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *AnyAnyMap) IsEmpty() bool {
    m.mu.RLock()
    empty := len(m.data) == 0
    m.mu.RUnlock()
    return empty
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *AnyAnyMap) Clear() {
    m.mu.Lock()
    m.data = make(map[interface{}]interface{})
    m.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (m *AnyAnyMap) LockFunc(f func(m map[interface{}]interface{})) {
    m.mu.Lock()
    defer m.mu.Unlock()
    f(m.data)
}

// RLockFunc locks reading with given callback function <f> within RWMutex.RLock.
func (m *AnyAnyMap) RLockFunc(f func(m map[interface{}]interface{})) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    f(m.data)
}

// Flip exchanges key-value of the map to value-key.
func (m *AnyAnyMap) Flip() {
    m.mu.Lock()
    defer m.mu.Unlock()
    n := make(map[interface{}]interface{}, len(m.data))
    for k, v := range m.data {
        n[v] = k
    }
    m.data = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <m>.
func (m *AnyAnyMap) Merge(other *AnyAnyMap) {
    m.mu.Lock()
    defer m.mu.Unlock()
    if other != m {
	    other.mu.RLock()
        defer other.mu.RUnlock()
    }
    for k, v := range other.data {
        m.data[k] = v
    }
}
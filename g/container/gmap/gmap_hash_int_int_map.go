// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import (
    "github.com/gogf/gf/g/internal/rwmutex"
)

type IntIntMap struct {
	mu   *rwmutex.RWMutex
	data map[int]int
}

// NewIntIntMap returns an empty IntIntMap object.
// The param <unsafe> used to specify whether using map in un-concurrent-safety,
// which is false in default, means concurrent-safe.
func NewIntIntMap(unsafe...bool) *IntIntMap {
	return &IntIntMap{
        mu   : rwmutex.New(unsafe...),
		data : make(map[int]int),
    }
}

// NewIntIntMapFrom returns a hash map from given map <data>.
// Note that, the param <data> map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewIntIntMapFrom(data map[int]int, unsafe...bool) *IntIntMap {
    return &IntIntMap{
        mu   : rwmutex.New(unsafe...),
	    data : data,
    }
}

// Iterator iterates the hash map with custom callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (m *IntIntMap) Iterator(f func (k int, v int) bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    for k, v := range m.data {
        if !f(k, v) {
            break
        }
    }
}

// Clone returns a new hash map with copy of current map data.
func (m *IntIntMap) Clone() *IntIntMap {
    return NewIntIntMapFrom(m.Map(), !m.mu.IsSafe())
}

// Map returns a copy of the data of the hash map.
func (m *IntIntMap) Map() map[int]int {
	m.mu.RLock()
	data := make(map[int]int, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
    m.mu.RUnlock()
	return data
}

// Set sets key-value to the hash map.
func (m *IntIntMap) Set(key int, val int) {
	m.mu.Lock()
	m.data[key] = val
	m.mu.Unlock()
}

// Sets batch sets key-values to the hash map.
func (m *IntIntMap) Sets(data map[int]int) {
	m.mu.Lock()
	for k, v := range data {
		m.data[k] = v
	}
	m.mu.Unlock()
}

// Search searches the map with given <key>.
// Second return parameter <found> is true if key was found, otherwise false.
func (m *IntIntMap) Search(key int) (value int, found bool) {
	m.mu.RLock()
	value, found = m.data[key]
	m.mu.RUnlock()
	return
}

// Get returns the value by given <key>.
func (m *IntIntMap) Get(key int) (int) {
	m.mu.RLock()
	val, _ := m.data[key]
	m.mu.RUnlock()
	return val
}

// doSetWithLockCheck checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given <key>,
// or else just return the existing value.
//
// It returns value with given <key>.
func (m *IntIntMap) doSetWithLockCheck(key int, value int) int {
    m.mu.Lock()
    if v, ok := m.data[key]; ok {
        m.mu.Unlock()
        return v
    }
    m.data[key] = value
    m.mu.Unlock()
    return value
}

// GetOrSet returns the value by key,
// or set value with given <value> if not exist and returns this value.
func (m *IntIntMap) GetOrSet(key int, value int) int {
	if v, ok := m.Search(key); !ok {
        return m.doSetWithLockCheck(key, value)
    } else {
        return v
    }
}

// GetOrSetFunc returns the value by key,
// or sets value with return value of callback function <f> if not exist and returns this value.
func (m *IntIntMap) GetOrSetFunc(key int, f func() int) int {
	if v, ok := m.Search(key); !ok {
        return m.doSetWithLockCheck(key, f())
    } else {
        return v
    }
}

// GetOrSetFuncLock returns the value by key,
// or sets value with return value of callback function <f> if not exist and returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function <f>
// with mutex.Lock of the hash map.
func (m *IntIntMap) GetOrSetFuncLock(key int, f func() int) int {
	if v, ok := m.Search(key); !ok {
        m.mu.Lock()
        defer m.mu.Unlock()
        if v, ok = m.data[key]; ok {
            return v
        }
        v = f()
        m.data[key] = v
        return v
    } else {
        return v
    }
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *IntIntMap) SetIfNotExist(key int, value int) bool {
    if !m.Contains(key) {
        m.doSetWithLockCheck(key, value)
        return true
    }
    return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *IntIntMap) SetIfNotExistFunc(key int, f func() int) bool {
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
func (m *IntIntMap) SetIfNotExistFuncLock(key int, f func() int) bool {
	if !m.Contains(key) {
		m.mu.Lock()
		defer m.mu.Unlock()
		if _, ok := m.data[key]; !ok {
			m.data[key] = f()
		}
		return true
	}
	return false
}

// Removes batch deletes values of the map by keys.
func (m *IntIntMap) Removes(keys []int) {
    m.mu.Lock()
    for _, key := range keys {
        delete(m.data, key)
    }
    m.mu.Unlock()
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (m *IntIntMap) Remove(key int) int {
    m.mu.Lock()
    val, exists := m.data[key]
    if exists {
        delete(m.data, key)
    }
    m.mu.Unlock()
    return val
}

// Keys returns all keys of the map as a slice.
func (m *IntIntMap) Keys() []int {
    m.mu.RLock()
    keys  := make([]int, len(m.data))
    index := 0
    for key := range m.data {
        keys[index] = key
        index++
    }
    m.mu.RUnlock()
    return keys
}

// Values returns all values of the map as a slice.
func (m *IntIntMap) Values() []int {
    m.mu.RLock()
    values := make([]int, len(m.data))
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
func (m *IntIntMap) Contains(key int) bool {
    m.mu.RLock()
    _, exists := m.data[key]
    m.mu.RUnlock()
    return exists
}

// Size returns the size of the map.
func (m *IntIntMap) Size() int {
    m.mu.RLock()
    length := len(m.data)
    m.mu.RUnlock()
    return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *IntIntMap) IsEmpty() bool {
    m.mu.RLock()
    empty := len(m.data) == 0
    m.mu.RUnlock()
    return empty
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *IntIntMap) Clear() {
    m.mu.Lock()
    m.data = make(map[int]int)
    m.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (m *IntIntMap) LockFunc(f func(m map[int]int)) {
    m.mu.Lock()
    defer m.mu.Unlock()
    f(m.data)
}

// RLockFunc locks reading with given callback function <f> within RWMutex.RLock.
func (m *IntIntMap) RLockFunc(f func(m map[int]int)) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    f(m.data)
}

// Flip exchanges key-value of the map to value-key.
func (m *IntIntMap) Flip() {
    m.mu.Lock()
    defer m.mu.Unlock()
    n := make(map[int]int, len(m.data))
    for k, v := range m.data {
        n[v] = k
    }
    m.data = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <m>.
func (m *IntIntMap) Merge(other *IntIntMap) {
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

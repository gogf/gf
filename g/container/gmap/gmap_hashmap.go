// Copyright 2017-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with m file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import "github.com/gogf/gf/g/internal/rwmutex"

type HashMap struct {
    mu   *rwmutex.RWMutex
    data map[interface{}]interface{}
}

// New returns an empty hash map.
// The param <unsafe> used to specify whether using map with un-concurrent-safety,
// which is false in default, means concurrent-safe.
func New(unsafe ...bool) *HashMap {
	return &HashMap{
		mu   : rwmutex.New(unsafe...),
		data : make(map[interface{}]interface{}),
	}
}

// NewFrom returns a hash map from given map <m>.
// Notice that, the param map is a type of pointer,
// there might be some concurrent-safe issues when changing the map outside.
func NewFrom(data map[interface{}]interface{}, unsafe...bool) *HashMap {
    return &HashMap {
	    mu   : rwmutex.New(unsafe...),
        data : data,
    }
}

// NewFromArray returns a hash map from given array.
// The param <keys> given as the keys of the map,
// and <values> as its corresponding values.
//
// If length of <keys> is greater than that of <values>,
// the corresponding overflow map values will be the default value of its type.
func NewFromArray(keys []interface{}, values []interface{}, unsafe...bool) *HashMap {
    m := make(map[interface{}]interface{})
    l := len(values)
    for i, k := range keys {
        if i < l {
            m[k] = values[i]
        } else {
            m[k] = interface{}(nil)
        }
    }
    return &HashMap{
        mu   : rwmutex.New(unsafe...),
	    data : m,
    }
}

// Iterator iterates the hash map with custom callback function <f>.
// If f returns true, then continue iterating; or false to stop.
func (m *HashMap) Iterator(f func (k interface{}, v interface{}) bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    for k, v := range m.data {
        if !f(k, v) {
            break
        }
    }
}

// Clone returns a new hash map with copy of current map data.
func (m *HashMap) Clone(unsafe ...bool) *HashMap {
    return NewFrom(m.Map(), unsafe ...)
}

// Map returns a copy of the data of the hash map.
func (m *HashMap) Map() map[interface{}]interface{} {
    data := make(map[interface{}]interface{})
    m.mu.RLock()
    for k, v := range m.data {
	    data[k] = v
    }
    m.mu.RUnlock()
    return data
}

// Set sets key-value to the hash map.
func (m *HashMap) Set(key interface{}, val interface{}) {
    m.mu.Lock()
    m.data[key] = val
    m.mu.Unlock()
}

// BatchSet batch sets key-values to the hash map.
func (m *HashMap) BatchSet(data map[interface{}]interface{}) {
    m.mu.Lock()
    for k, v := range data {
        m.data[k] = v
    }
    m.mu.Unlock()
}

// Get returns the value by given <key>.
func (m *HashMap) Get(key interface{}) interface{} {
    m.mu.RLock()
    val, _ := m.data[key]
    m.mu.RUnlock()
    return val
}

// doSetWithLockCheck checks whether value of the key exists within mutex.Lock,
// if the value not exists, it sets value to the map with given <key>,
// or else just returns the existing value.
//
// When setting value, if <value> is type of <func() interface {}>,
// it will be executed within mutex.Lock of the hash map,
// and its return value will be set to the map with <key>.
//
// It returns value of given <key>.
func (m *HashMap) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
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
func (m *HashMap) GetOrSet(key interface{}, value interface{}) interface{} {
    if v := m.Get(key); v == nil {
        return m.doSetWithLockCheck(key, value)
    } else {
        return v
    }
}

// GetOrSetFunc returns the value by key,
// or sets value with return value of callback function <f> if not exist
// and returns this value.
func (m *HashMap) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
    if v := m.Get(key); v == nil {
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
func (m *HashMap) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
    if v := m.Get(key); v == nil {
        return m.doSetWithLockCheck(key, f)
    } else {
        return v
    }
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *HashMap) SetIfNotExist(key interface{}, value interface{}) bool {
    if !m.Contains(key) {
        m.doSetWithLockCheck(key, value)
        return true
    }
    return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *HashMap) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
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
func (m *HashMap) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

// BatchRemove batch deletes values of the map by keys.
func (m *HashMap) BatchRemove(keys []interface{}) {
    m.mu.Lock()
    for _, key := range keys {
        delete(m.data, key)
    }
    m.mu.Unlock()
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (m *HashMap) Remove(key interface{}) interface{} {
    m.mu.Lock()
    val, exists := m.data[key]
    if exists {
        delete(m.data, key)
    }
    m.mu.Unlock()
    return val
}

// Keys returns all keys of the map as a slice.
func (m *HashMap) Keys() []interface{} {
    m.mu.RLock()
    keys := make([]interface{}, 0)
    for key, _ := range m.data {
        keys = append(keys, key)
    }
    m.mu.RUnlock()
    return keys
}

// Values returns all values of the map as a slice.
func (m *HashMap) Values() []interface{} {
    m.mu.RLock()
    vals := make([]interface{}, 0)
    for _, val := range m.data {
        vals = append(vals, val)
    }
    m.mu.RUnlock()
    return vals
}

// Contains checks whether a key exists.
// It returns true if the <key> exists, or else false.
func (m *HashMap) Contains(key interface{}) bool {
    m.mu.RLock()
    _, exists := m.data[key]
    m.mu.RUnlock()
    return exists
}

// Size returns the size of the map.
func (m *HashMap) Size() int {
    m.mu.RLock()
    length := len(m.data)
    m.mu.RUnlock()
    return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *HashMap) IsEmpty() bool {
    m.mu.RLock()
    empty := len(m.data) == 0
    m.mu.RUnlock()
    return empty
}

// Clear deletes all data of the map, it will remake a new underlying map data map.
func (m *HashMap) Clear() {
    m.mu.Lock()
    m.data = make(map[interface{}]interface{})
    m.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> and mutex.Lock.
func (m *HashMap) LockFunc(f func(m map[interface{}]interface{})) {
    m.mu.Lock()
    defer m.mu.Unlock()
    f(m.data)
}

// RLockFunc locks reading with given callback function <f> and mutex.RLock.
func (m *HashMap) RLockFunc(f func(m map[interface{}]interface{})) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    f(m.data)
}

// Flip exchanges key-value of the map, it will change key-value to value-key.
func (m *HashMap) Flip() {
    m.mu.Lock()
    defer m.mu.Unlock()
    n := make(map[interface{}]interface{}, len(m.data))
    for i, v := range m.data {
        n[v] = i
    }
    m.data = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <gm>.
func (m *HashMap) Merge(other *HashMap) {
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
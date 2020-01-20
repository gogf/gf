// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.
//

package gmap

import (
	"encoding/json"
	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/internal/empty"

	"github.com/gogf/gf/internal/rwmutex"
)

type StrStrMap struct {
	mu   *rwmutex.RWMutex
	data map[string]string
}

// NewStrStrMap returns an empty StrStrMap object.
// The parameter <safe> is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewStrStrMap(safe ...bool) *StrStrMap {
	return &StrStrMap{
		data: make(map[string]string),
		mu:   rwmutex.New(safe...),
	}
}

// NewStrStrMapFrom creates and returns a hash map from given map <data>.
// Note that, the param <data> map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewStrStrMapFrom(data map[string]string, safe ...bool) *StrStrMap {
	return &StrStrMap{
		mu:   rwmutex.New(safe...),
		data: data,
	}
}

// Iterator iterates the hash map with custom callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (m *StrStrMap) Iterator(f func(k string, v string) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Clone returns a new hash map with copy of current map data.
func (m *StrStrMap) Clone() *StrStrMap {
	return NewStrStrMapFrom(m.MapCopy(), !m.mu.IsSafe())
}

// Map returns the underlying data map.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (m *StrStrMap) Map() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.mu.IsSafe() {
		return m.data
	}
	data := make(map[string]string, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// MapStrAny returns a copy of the underlying data of the map as map[string]interface{}.
func (m *StrStrMap) MapStrAny() map[string]interface{} {
	m.mu.RLock()
	data := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	m.mu.RUnlock()
	return data
}

// MapCopy returns a copy of the underlying data of the hash map.
func (m *StrStrMap) MapCopy() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[string]string, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// FilterEmpty deletes all key-value pair of which the value is empty.
func (m *StrStrMap) FilterEmpty() {
	m.mu.Lock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
	m.mu.Unlock()
}

// Set sets key-value to the hash map.
func (m *StrStrMap) Set(key string, val string) {
	m.mu.Lock()
	m.data[key] = val
	m.mu.Unlock()
}

// Sets batch sets key-values to the hash map.
func (m *StrStrMap) Sets(data map[string]string) {
	m.mu.Lock()
	for k, v := range data {
		m.data[k] = v
	}
	m.mu.Unlock()
}

// Search searches the map with given <key>.
// Second return parameter <found> is true if key was found, otherwise false.
func (m *StrStrMap) Search(key string) (value string, found bool) {
	m.mu.RLock()
	value, found = m.data[key]
	m.mu.RUnlock()
	return
}

// Get returns the value by given <key>.
func (m *StrStrMap) Get(key string) string {
	m.mu.RLock()
	val, _ := m.data[key]
	m.mu.RUnlock()
	return val
}

// Pop retrieves and deletes an item from the map.
func (m *StrStrMap) Pop() (key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, value = range m.data {
		delete(m.data, key)
		return
	}
	return
}

// Pops retrieves and deletes <size> items from the map.
// It returns all items if size == -1.
func (m *StrStrMap) Pops(size int) map[string]string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if size > len(m.data) || size == -1 {
		size = len(m.data)
	}
	if size == 0 {
		return nil
	}
	index := 0
	newMap := make(map[string]string, size)
	for k, v := range m.data {
		delete(m.data, k)
		newMap[k] = v
		index++
		if index == size {
			break
		}
	}
	return newMap
}

// doSetWithLockCheck checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given <key>,
// or else just return the existing value.
//
// It returns value with given <key>.
func (m *StrStrMap) doSetWithLockCheck(key string, value string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.data[key]; ok {
		return v
	}
	m.data[key] = value
	return value
}

// GetOrSet returns the value by key,
// or sets value with given <value> if it does not exist and then returns this value.
func (m *StrStrMap) GetOrSet(key string, value string) string {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
func (m *StrStrMap) GetOrSetFunc(key string, f func() string) string {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function <f>
// with mutex.Lock of the hash map.
func (m *StrStrMap) GetOrSetFuncLock(key string, f func() string) string {
	if v, ok := m.Search(key); !ok {
		m.mu.Lock()
		defer m.mu.Unlock()
		if v, ok = m.data[key]; ok {
			return v
		}
		v = f()
		if v != "" {
			m.data[key] = v
		}
		return v
	} else {
		return v
	}
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *StrStrMap) SetIfNotExist(key string, value string) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *StrStrMap) SetIfNotExistFunc(key string, f func() string) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function <f> with mutex.Lock of the hash map.
func (m *StrStrMap) SetIfNotExistFuncLock(key string, f func() string) bool {
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
func (m *StrStrMap) Removes(keys []string) {
	m.mu.Lock()
	for _, key := range keys {
		delete(m.data, key)
	}
	m.mu.Unlock()
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (m *StrStrMap) Remove(key string) string {
	m.mu.Lock()
	val, exists := m.data[key]
	if exists {
		delete(m.data, key)
	}
	m.mu.Unlock()
	return val
}

// Keys returns all keys of the map as a slice.
func (m *StrStrMap) Keys() []string {
	m.mu.RLock()
	keys := make([]string, len(m.data))
	index := 0
	for key := range m.data {
		keys[index] = key
		index++
	}
	m.mu.RUnlock()
	return keys
}

// Values returns all values of the map as a slice.
func (m *StrStrMap) Values() []string {
	m.mu.RLock()
	values := make([]string, len(m.data))
	index := 0
	for _, value := range m.data {
		values[index] = value
		index++
	}
	m.mu.RUnlock()
	return values
}

// Contains checks whether a key exists.
// It returns true if the <key> exists, or else false.
func (m *StrStrMap) Contains(key string) bool {
	m.mu.RLock()
	_, exists := m.data[key]
	m.mu.RUnlock()
	return exists
}

// Size returns the size of the map.
func (m *StrStrMap) Size() int {
	m.mu.RLock()
	length := len(m.data)
	m.mu.RUnlock()
	return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *StrStrMap) IsEmpty() bool {
	return m.Size() == 0
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *StrStrMap) Clear() {
	m.mu.Lock()
	m.data = make(map[string]string)
	m.mu.Unlock()
}

// Replace the data of the map with given <data>.
func (m *StrStrMap) Replace(data map[string]string) {
	m.mu.Lock()
	m.data = data
	m.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (m *StrStrMap) LockFunc(f func(m map[string]string)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

// RLockFunc locks reading with given callback function <f> within RWMutex.RLock.
func (m *StrStrMap) RLockFunc(f func(m map[string]string)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f(m.data)
}

// Flip exchanges key-value of the map to value-key.
func (m *StrStrMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[string]string, len(m.data))
	for k, v := range m.data {
		n[v] = k
	}
	m.data = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <m>.
func (m *StrStrMap) Merge(other *StrStrMap) {
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

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m *StrStrMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (m *StrStrMap) UnmarshalJSON(b []byte) error {
	if m.mu == nil {
		m.mu = rwmutex.New()
		m.data = make(map[string]string)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := json.Unmarshal(b, &m.data); err != nil {
		return err
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (m *StrStrMap) UnmarshalValue(value interface{}) (err error) {
	if m.mu == nil {
		m.mu = rwmutex.New()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = gconv.MapStrStr(value)
	return
}

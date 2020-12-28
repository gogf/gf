// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.
//

package gmap

import (
	"github.com/gogf/gf/internal/json"

	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/internal/rwmutex"
	"github.com/gogf/gf/util/gconv"
)

type StrIntMap struct {
	mu   rwmutex.RWMutex
	data map[string]int
}

// NewStrIntMap returns an empty StrIntMap object.
// The parameter <safe> is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewStrIntMap(safe ...bool) *StrIntMap {
	return &StrIntMap{
		mu:   rwmutex.Create(safe...),
		data: make(map[string]int),
	}
}

// NewStrIntMapFrom creates and returns a hash map from given map <data>.
// Note that, the param <data> map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewStrIntMapFrom(data map[string]int, safe ...bool) *StrIntMap {
	return &StrIntMap{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

// Iterator iterates the hash map readonly with custom callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (m *StrIntMap) Iterator(f func(k string, v int) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Clone returns a new hash map with copy of current map data.
func (m *StrIntMap) Clone() *StrIntMap {
	return NewStrIntMapFrom(m.MapCopy(), m.mu.IsSafe())
}

// Map returns the underlying data map.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (m *StrIntMap) Map() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.mu.IsSafe() {
		return m.data
	}
	data := make(map[string]int, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// MapStrAny returns a copy of the underlying data of the map as map[string]interface{}.
func (m *StrIntMap) MapStrAny() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// MapCopy returns a copy of the underlying data of the hash map.
func (m *StrIntMap) MapCopy() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[string]int, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// FilterEmpty deletes all key-value pair of which the value is empty.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (m *StrIntMap) FilterEmpty() {
	m.mu.Lock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
	m.mu.Unlock()
}

// Set sets key-value to the hash map.
func (m *StrIntMap) Set(key string, val int) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[string]int)
	}
	m.data[key] = val
	m.mu.Unlock()
}

// Sets batch sets key-values to the hash map.
func (m *StrIntMap) Sets(data map[string]int) {
	m.mu.Lock()
	if m.data == nil {
		m.data = data
	} else {
		for k, v := range data {
			m.data[k] = v
		}
	}
	m.mu.Unlock()
}

// Search searches the map with given <key>.
// Second return parameter <found> is true if key was found, otherwise false.
func (m *StrIntMap) Search(key string) (value int, found bool) {
	m.mu.RLock()
	if m.data != nil {
		value, found = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Get returns the value by given <key>.
func (m *StrIntMap) Get(key string) (value int) {
	m.mu.RLock()
	if m.data != nil {
		value, _ = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Pop retrieves and deletes an item from the map.
func (m *StrIntMap) Pop() (key string, value int) {
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
func (m *StrIntMap) Pops(size int) map[string]int {
	m.mu.Lock()
	defer m.mu.Unlock()
	if size > len(m.data) || size == -1 {
		size = len(m.data)
	}
	if size == 0 {
		return nil
	}
	var (
		index  = 0
		newMap = make(map[string]int, size)
	)
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
func (m *StrIntMap) doSetWithLockCheck(key string, value int) int {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[string]int)
	}
	if v, ok := m.data[key]; ok {
		m.mu.Unlock()
		return v
	}
	m.data[key] = value
	m.mu.Unlock()
	return value
}

// GetOrSet returns the value by key,
// or sets value with given <value> if it does not exist and then returns this value.
func (m *StrIntMap) GetOrSet(key string, value int) int {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
func (m *StrIntMap) GetOrSetFunc(key string, f func() int) int {
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
func (m *StrIntMap) GetOrSetFuncLock(key string, f func() int) int {
	if v, ok := m.Search(key); !ok {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.data == nil {
			m.data = make(map[string]int)
		}
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

// SetIfNotExist sets <value> to the map if the <key> does not exist, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *StrIntMap) SetIfNotExist(key string, value int) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *StrIntMap) SetIfNotExistFunc(key string, f func() int) bool {
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
func (m *StrIntMap) SetIfNotExistFuncLock(key string, f func() int) bool {
	if !m.Contains(key) {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.data == nil {
			m.data = make(map[string]int)
		}
		if _, ok := m.data[key]; !ok {
			m.data[key] = f()
		}
		return true
	}
	return false
}

// Removes batch deletes values of the map by keys.
func (m *StrIntMap) Removes(keys []string) {
	m.mu.Lock()
	if m.data != nil {
		for _, key := range keys {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (m *StrIntMap) Remove(key string) (value int) {
	m.mu.Lock()
	if m.data != nil {
		var ok bool
		if value, ok = m.data[key]; ok {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
	return
}

// Keys returns all keys of the map as a slice.
func (m *StrIntMap) Keys() []string {
	m.mu.RLock()
	var (
		keys  = make([]string, len(m.data))
		index = 0
	)
	for key := range m.data {
		keys[index] = key
		index++
	}
	m.mu.RUnlock()
	return keys
}

// Values returns all values of the map as a slice.
func (m *StrIntMap) Values() []int {
	m.mu.RLock()
	var (
		values = make([]int, len(m.data))
		index  = 0
	)
	for _, value := range m.data {
		values[index] = value
		index++
	}
	m.mu.RUnlock()
	return values
}

// Contains checks whether a key exists.
// It returns true if the <key> exists, or else false.
func (m *StrIntMap) Contains(key string) bool {
	var ok bool
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return ok
}

// Size returns the size of the map.
func (m *StrIntMap) Size() int {
	m.mu.RLock()
	length := len(m.data)
	m.mu.RUnlock()
	return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *StrIntMap) IsEmpty() bool {
	return m.Size() == 0
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *StrIntMap) Clear() {
	m.mu.Lock()
	m.data = make(map[string]int)
	m.mu.Unlock()
}

// Replace the data of the map with given <data>.
func (m *StrIntMap) Replace(data map[string]int) {
	m.mu.Lock()
	m.data = data
	m.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (m *StrIntMap) LockFunc(f func(m map[string]int)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

// RLockFunc locks reading with given callback function <f> within RWMutex.RLock.
func (m *StrIntMap) RLockFunc(f func(m map[string]int)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f(m.data)
}

// Flip exchanges key-value of the map to value-key.
func (m *StrIntMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[string]int, len(m.data))
	for k, v := range m.data {
		n[gconv.String(v)] = gconv.Int(k)
	}
	m.data = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <m>.
func (m *StrIntMap) Merge(other *StrIntMap) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = other.MapCopy()
		return
	}
	if other != m {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	for k, v := range other.data {
		m.data[k] = v
	}
}

// String returns the map as a string.
func (m *StrIntMap) String() string {
	b, _ := m.MarshalJSON()
	return gconv.UnsafeBytesToStr(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m *StrIntMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (m *StrIntMap) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[string]int)
	}
	if err := json.Unmarshal(b, &m.data); err != nil {
		return err
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (m *StrIntMap) UnmarshalValue(value interface{}) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[string]int)
	}
	switch value.(type) {
	case string, []byte:
		return json.Unmarshal(gconv.Bytes(value), &m.data)
	default:
		for k, v := range gconv.Map(value) {
			m.data[k] = gconv.Int(v)
		}
	}
	return
}

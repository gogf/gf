// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import (
	"reflect"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/deepcopy"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/util/gconv"
)

// KVMap wraps map type `map[K]V` and provides more map features.
type KVMap[K comparable, V any] struct {
	mu   rwmutex.RWMutex
	data map[K]V

	doSetWithLockCheckFn func(key K, value V) (val V)
}

// NewKVMap creates and returns an empty hash map.
// The parameter `safe` is used to specify whether to use the map in concurrent-safety mode,
// which is false by default.
func NewKVMap[K comparable, V any](safe ...bool) *KVMap[K, V] {
	m := &KVMap[K, V]{
		mu:   rwmutex.Create(safe...),
		data: make(map[K]V),
	}
	m.doSetWithLockCheckFn = m.doSetWithLockCheck
	return m
}

// NewKVMapFrom creates and returns a hash map from given map `data`.
// Note that, the param `data` map will be set as the underlying data map (no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewKVMapFrom[K comparable, V any](data map[K]V, safe ...bool) *KVMap[K, V] {
	m := &KVMap[K, V]{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
	m.doSetWithLockCheckFn = m.doSetWithLockCheck
	return m
}

// Iterator iterates the hash map readonly with custom callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *KVMap[K, V]) Iterator(f func(k K, v V) bool) {
	for k, v := range m.Map() {
		if !f(k, v) {
			break
		}
	}
}

// Clone returns a new hash map with copy of current map data.
func (m *KVMap[K, V]) Clone(safe ...bool) *KVMap[K, V] {
	if len(safe) == 0 {
		return NewKVMapFrom(m.MapCopy(), m.mu.IsSafe())
	}
	return NewKVMapFrom(m.MapCopy(), safe...)
}

// Map returns the underlying data map.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (m *KVMap[K, V]) Map() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.mu.IsSafe() {
		return m.data
	}
	data := make(map[K]V, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// MapCopy returns a shallow copy of the underlying data of the hash map.
func (m *KVMap[K, V]) MapCopy() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[K]V, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// MapStrAny returns a copy of the underlying data of the map as map[string]any.
func (m *KVMap[K, V]) MapStrAny() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[string]any, len(m.data))
	for k, v := range m.data {
		data[gconv.String(k)] = v
	}
	return data
}

// FilterEmpty deletes all key-value pair of which the value is empty.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (m *KVMap[K, V]) FilterEmpty() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
}

// FilterNil deletes all key-value pair of which the value is nil.
func (m *KVMap[K, V]) FilterNil() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if empty.IsNil(v) {
			delete(m.data, k)
		}
	}
}

// Set sets key-value to the hash map.
func (m *KVMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[K]V)
	}
	m.data[key] = value
	m.mu.Unlock()
}

// Sets batch sets key-values to the hash map.
func (m *KVMap[K, V]) Sets(data map[K]V) {
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

// Search searches the map with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (m *KVMap[K, V]) Search(key K) (value V, found bool) {
	m.mu.RLock()
	if m.data != nil {
		value, found = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Get returns the value by given `key`.
func (m *KVMap[K, V]) Get(key K) (value V) {
	m.mu.RLock()
	if m.data != nil {
		value = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Pop retrieves and deletes an item from the map.
func (m *KVMap[K, V]) Pop() (key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, value = range m.data {
		delete(m.data, key)
		return
	}
	return
}

// Pops retrieves and deletes `size` items from the map.
// It returns all items if size == -1.
func (m *KVMap[K, V]) Pops(size int) map[K]V {
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
		newMap = make(map[K]V, size)
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
// if not exists, set value to the map with given `key`,
// or else just return the existing value.
//
// When setting value, if `value` is type of `func() interface {}`,
// it will be executed with mutex.Lock of the hash map,
// and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (m *KVMap[K, V]) doSetWithLockCheck(key K, value V) (val V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data == nil {
		m.data = make(map[K]V)
	}
	if v, ok := m.data[key]; ok {
		return v
	}

	if any(value) != nil {
		m.data[key] = value
	}
	return value
}

// GetOrSet returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (m *KVMap[K, V]) GetOrSet(key K, value V) V {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheckFn(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
func (m *KVMap[K, V]) GetOrSetFunc(key K, f func() V) V {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheckFn(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`
// with mutex.Lock of the hash map.
func (m *KVMap[K, V]) GetOrSetFuncLock(key K, f func() V) V {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheckFn(key, f())
	} else {
		return v
	}
}

// GetVar returns a Var with the value by given `key`.
// The returned Var is un-concurrent safe.
func (m *KVMap[K, V]) GetVar(key K) *gvar.Var {
	return gvar.New(m.Get(key))
}

// GetVarOrSet returns a Var with result from GetOrSet.
// The returned Var is un-concurrent safe.
func (m *KVMap[K, V]) GetVarOrSet(key K, value V) *gvar.Var {
	return gvar.New(m.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a Var with result from GetOrSetFunc.
// The returned Var is un-concurrent safe.
func (m *KVMap[K, V]) GetVarOrSetFunc(key K, f func() V) *gvar.Var {
	return gvar.New(m.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a Var with result from GetOrSetFuncLock.
// The returned Var is un-concurrent safe.
func (m *KVMap[K, V]) GetVarOrSetFuncLock(key K, f func() V) *gvar.Var {
	return gvar.New(m.GetOrSetFuncLock(key, f))
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *KVMap[K, V]) SetIfNotExist(key K, value V) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheckFn(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *KVMap[K, V]) SetIfNotExistFunc(key K, f func() V) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheckFn(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function `f` with mutex.Lock of the hash map.
func (m *KVMap[K, V]) SetIfNotExistFuncLock(key K, f func() V) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheckFn(key, f())
		return true
	}
	return false
}

// Remove deletes value from map by given `key`, and return this deleted value.
func (m *KVMap[K, V]) Remove(key K) (value V) {
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

// Removes batch deletes values of the map by keys.
func (m *KVMap[K, V]) Removes(keys []K) {
	m.mu.Lock()
	if m.data != nil {
		for _, key := range keys {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
}

// Keys returns all keys of the map as a slice.
func (m *KVMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var (
		keys  = make([]K, len(m.data))
		index = 0
	)
	for key := range m.data {
		keys[index] = key
		index++
	}
	return keys
}

// Values returns all values of the map as a slice.
func (m *KVMap[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var (
		values = make([]V, len(m.data))
		index  = 0
	)
	for _, value := range m.data {
		values[index] = value
		index++
	}
	return values
}

// Contains checks whether a key exists.
// It returns true if the `key` exists, or else false.
func (m *KVMap[K, V]) Contains(key K) bool {
	var ok bool
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return ok
}

// Size returns the size of the map.
func (m *KVMap[K, V]) Size() int {
	m.mu.RLock()
	length := len(m.data)
	m.mu.RUnlock()
	return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *KVMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *KVMap[K, V]) Clear() {
	m.mu.Lock()
	m.data = make(map[K]V)
	m.mu.Unlock()
}

// Replace the data of the map with given `data`.
func (m *KVMap[K, V]) Replace(data map[K]V) {
	m.mu.Lock()
	m.data = data
	m.mu.Unlock()
}

// LockFunc locks writing with given callback function `f` within RWMutex.Lock.
func (m *KVMap[K, V]) LockFunc(f func(m map[K]V)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

// RLockFunc locks reading with given callback function `f` within RWMutex.RLock.
func (m *KVMap[K, V]) RLockFunc(f func(m map[K]V)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f(m.data)
}

// Flip exchanges key-value of the map to value-key.
func (m *KVMap[K, V]) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[K]V, len(m.data))
	for k, v := range m.data {
		var (
			k0 K
			v0 V
		)
		if err := gconv.Scan(v, &k0); err != nil {
			continue
		}
		if err := gconv.Scan(k, &v0); err != nil {
			continue
		}
		n[k0] = v0
	}
	m.data = n
}

// Merge merges two hash maps.
// The `other` map will be merged into the map `m`.
func (m *KVMap[K, V]) Merge(other *KVMap[K, V]) {
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
func (m *KVMap[K, V]) String() string {
	if m == nil {
		return ""
	}
	b, _ := m.MarshalJSON()
	return string(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m KVMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(gconv.Map(m.Map()))
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (m *KVMap[K, V]) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]V)
	}
	var data map[string]V
	if err := json.UnmarshalUseNumber(b, &data); err != nil {
		return err
	}
	if err := gconv.Scan(data, &m.data); err != nil {
		return err
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (m *KVMap[K, V]) UnmarshalValue(value any) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]V)
	}
	data := gconv.Map(value)
	if err := gconv.Scan(data, &m.data); err != nil {
		return err
	}
	return
}

// DeepCopy implements interface for deep copy of current type.
func (m *KVMap[K, V]) DeepCopy() any {
	if m == nil {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[K]V, len(m.data))
	for k, v := range m.data {
		data[k] = deepcopy.Copy(v).(V)
	}
	return NewKVMapFrom(data, m.mu.IsSafe())
}

// IsSubOf checks whether the current map is a sub-map of `other`.
func (m *KVMap[K, V]) IsSubOf(other *KVMap[K, V]) bool {
	if m == other {
		return true
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	for key, value := range m.data {
		otherValue, ok := other.data[key]
		if !ok {
			return false
		}

		if !reflect.DeepEqual(otherValue, value) {
			return false
		}
	}
	return true
}

// Diff compares current map `m` with map `other` and returns their different keys.
// The returned `addedKeys` are the keys that are in map `m` but not in map `other`.
// The returned `removedKeys` are the keys that are in map `other` but not in map `m`.
// The returned `updatedKeys` are the keys that are both in map `m` and `other` but their values and not equal (`!=`).
func (m *KVMap[K, V]) Diff(other *KVMap[K, V]) (addedKeys, removedKeys, updatedKeys []K) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()

	for key := range m.data {
		if _, ok := other.data[key]; !ok {
			removedKeys = append(removedKeys, key)
		} else if !reflect.DeepEqual(m.data[key], other.data[key]) {
			updatedKeys = append(updatedKeys, key)
		}
	}
	for key := range other.data {
		if _, ok := m.data[key]; !ok {
			addedKeys = append(addedKeys, key)
		}
	}
	return
}

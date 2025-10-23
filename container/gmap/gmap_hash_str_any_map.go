// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.
//

package gmap

import (
	"sync"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/util/gconv"
)

// StrAnyMap implements map[string]any with RWMutex that has switch.
type StrAnyMap struct {
	*KVMap[string, any]
	once sync.Once
}

// NewStrAnyMap returns an empty StrAnyMap object.
// The parameter `safe` is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewStrAnyMap(safe ...bool) *StrAnyMap {
	m := &StrAnyMap{
		KVMap: NewKVMap[string, any](safe...),
	}
	m.doSetWithLockCheckFn = m.doSetWithLockCheck
	return m
}

// NewStrAnyMapFrom creates and returns a hash map from given map `data`.
// Note that, the param `data` map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewStrAnyMapFrom(data map[string]any, safe ...bool) *StrAnyMap {
	m := &StrAnyMap{
		KVMap: NewKVMapFrom(data, safe...),
	}
	m.doSetWithLockCheckFn = m.doSetWithLockCheck
	return m
}

// lazyInit lazily initializes the map.
func (m *StrAnyMap) lazyInit() {
	m.once.Do(func() {
		if m.KVMap == nil {
			m.KVMap = NewKVMap[string, any](false)
			m.doSetWithLockCheckFn = m.doSetWithLockCheck
		}
	})
}

// Iterator iterates the hash map readonly with custom callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *StrAnyMap) Iterator(f func(k string, v any) bool) {
	m.lazyInit()
	m.KVMap.Iterator(f)
}

// Clone returns a new hash map with copy of current map data.
func (m *StrAnyMap) Clone(safe ...bool) *StrAnyMap {
	m.lazyInit()
	return NewStrAnyMapFrom(m.MapCopy(), safe...)
}

// Map returns the underlying data map.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (m *StrAnyMap) Map() map[string]any {
	m.lazyInit()
	return m.KVMap.Map()
}

// MapStrAny returns a copy of the underlying data of the map as map[string]any.
func (m *StrAnyMap) MapStrAny() map[string]any {
	return m.Map()
}

// MapCopy returns a copy of the underlying data of the hash map.
func (m *StrAnyMap) MapCopy() map[string]any {
	m.lazyInit()
	return m.KVMap.MapCopy()
}

// FilterEmpty deletes all key-value pair of which the value is empty.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (m *StrAnyMap) FilterEmpty() {
	m.lazyInit()
	m.KVMap.FilterEmpty()
}

// FilterNil deletes all key-value pair of which the value is nil.
func (m *StrAnyMap) FilterNil() {
	m.lazyInit()
	m.KVMap.FilterNil()
}

// Set sets key-value to the hash map.
func (m *StrAnyMap) Set(key string, val any) {
	m.lazyInit()
	m.KVMap.Set(key, val)
}

// Sets batch sets key-values to the hash map.
func (m *StrAnyMap) Sets(data map[string]any) {
	m.lazyInit()
	m.KVMap.Sets(data)
}

// Search searches the map with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (m *StrAnyMap) Search(key string) (value any, found bool) {
	m.lazyInit()
	return m.KVMap.Search(key)
}

// Get returns the value by given `key`.
func (m *StrAnyMap) Get(key string) (value any) {
	m.lazyInit()
	return m.KVMap.Get(key)
}

// Pop retrieves and deletes an item from the map.
func (m *StrAnyMap) Pop() (key string, value any) {
	m.lazyInit()
	return m.KVMap.Pop()
}

// Pops retrieves and deletes `size` items from the map.
// It returns all items if size == -1.
func (m *StrAnyMap) Pops(size int) map[string]any {
	m.lazyInit()
	return m.KVMap.Pops(size)
}

// GetOrSet returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (m *StrAnyMap) GetOrSet(key string, value any) any {
	m.lazyInit()
	return m.KVMap.GetOrSet(key, value)
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
func (m *StrAnyMap) GetOrSetFunc(key string, f func() any) any {
	m.lazyInit()
	return m.KVMap.GetOrSetFunc(key, f)
}

// GetOrSetFuncLock returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`
// with mutex.Lock of the hash map.
func (m *StrAnyMap) GetOrSetFuncLock(key string, f func() any) any {
	m.lazyInit()
	return m.KVMap.GetOrSetFuncLock(key, f)
}

// GetVar returns a Var with the value by given `key`.
// The returned Var is un-concurrent safe.
func (m *StrAnyMap) GetVar(key string) *gvar.Var {
	m.lazyInit()
	return m.KVMap.GetVar(key)
}

// GetVarOrSet returns a Var with result from GetVarOrSet.
// The returned Var is un-concurrent safe.
func (m *StrAnyMap) GetVarOrSet(key string, value any) *gvar.Var {
	m.lazyInit()
	return m.KVMap.GetVarOrSet(key, value)
}

// GetVarOrSetFunc returns a Var with result from GetOrSetFunc.
// The returned Var is un-concurrent safe.
func (m *StrAnyMap) GetVarOrSetFunc(key string, f func() any) *gvar.Var {
	m.lazyInit()
	return m.KVMap.GetVarOrSetFunc(key, f)
}

// GetVarOrSetFuncLock returns a Var with result from GetOrSetFuncLock.
// The returned Var is un-concurrent safe.
func (m *StrAnyMap) GetVarOrSetFuncLock(key string, f func() any) *gvar.Var {
	m.lazyInit()
	return m.KVMap.GetVarOrSetFuncLock(key, f)
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *StrAnyMap) SetIfNotExist(key string, value any) bool {
	m.lazyInit()
	return m.KVMap.SetIfNotExist(key, value)
}

// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *StrAnyMap) SetIfNotExistFunc(key string, f func() any) bool {
	m.lazyInit()
	return m.KVMap.SetIfNotExistFunc(key, f)
}

// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function `f` with mutex.Lock of the hash map.
func (m *StrAnyMap) SetIfNotExistFuncLock(key string, f func() any) bool {
	m.lazyInit()
	return m.KVMap.SetIfNotExistFuncLock(key, f)
}

// Removes batch deletes values of the map by keys.
func (m *StrAnyMap) Removes(keys []string) {
	m.lazyInit()
	m.KVMap.Removes(keys)
}

// Remove deletes value from map by given `key`, and return this deleted value.
func (m *StrAnyMap) Remove(key string) (value any) {
	m.lazyInit()
	return m.KVMap.Remove(key)
}

// Keys returns all keys of the map as a slice.
func (m *StrAnyMap) Keys() []string {
	m.lazyInit()
	return m.KVMap.Keys()
}

// Values returns all values of the map as a slice.
func (m *StrAnyMap) Values() []any {
	m.lazyInit()
	return m.KVMap.Values()
}

// Contains checks whether a key exists.
// It returns true if the `key` exists, or else false.
func (m *StrAnyMap) Contains(key string) bool {
	m.lazyInit()
	return m.KVMap.Contains(key)
}

// Size returns the size of the map.
func (m *StrAnyMap) Size() int {
	m.lazyInit()
	return m.KVMap.Size()
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *StrAnyMap) IsEmpty() bool {
	m.lazyInit()
	return m.KVMap.IsEmpty()
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *StrAnyMap) Clear() {
	m.lazyInit()
	m.KVMap.Clear()
}

// Replace the data of the map with given `data`.
func (m *StrAnyMap) Replace(data map[string]any) {
	m.lazyInit()
	m.KVMap.Replace(data)
}

// LockFunc locks writing with given callback function `f` within RWMutex.Lock.
func (m *StrAnyMap) LockFunc(f func(m map[string]any)) {
	m.lazyInit()
	m.KVMap.LockFunc(f)
}

// RLockFunc locks reading with given callback function `f` within RWMutex.RLock.
func (m *StrAnyMap) RLockFunc(f func(m map[string]any)) {
	m.lazyInit()
	m.KVMap.RLockFunc(f)
}

// Flip exchanges key-value of the map to value-key.
func (m *StrAnyMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[string]any, len(m.data))
	for k, v := range m.data {
		n[gconv.String(v)] = k
	}
	m.data = n
}

// Merge merges two hash maps.
// The `other` map will be merged into the map `m`.
func (m *StrAnyMap) Merge(other *StrAnyMap) {
	m.lazyInit()
	m.KVMap.Merge(other.KVMap)
}

// String returns the map as a string.
func (m *StrAnyMap) String() string {
	if m == nil {
		return ""
	}
	m.lazyInit()
	return m.KVMap.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m StrAnyMap) MarshalJSON() ([]byte, error) {
	m.lazyInit()
	return m.KVMap.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (m *StrAnyMap) UnmarshalJSON(b []byte) error {
	m.lazyInit()
	return m.KVMap.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (m *StrAnyMap) UnmarshalValue(value any) (err error) {
	m.lazyInit()
	return m.KVMap.UnmarshalValue(value)
}

// DeepCopy implements interface for deep copy of current type.
func (m *StrAnyMap) DeepCopy() any {
	m.lazyInit()
	return &StrAnyMap{
		KVMap: m.KVMap.DeepCopy().(*KVMap[string, any]),
	}
}

// IsSubOf checks whether the current map is a sub-map of `other`.
func (m *StrAnyMap) IsSubOf(other *StrAnyMap) bool {
	m.lazyInit()
	return m.KVMap.IsSubOf(other.KVMap)
}

// Diff compares current map `m` with map `other` and returns their different keys.
// The returned `addedKeys` are the keys that are in map `m` but not in map `other`.
// The returned `removedKeys` are the keys that are in map `other` but not in map `m`.
// The returned `updatedKeys` are the keys that are both in map `m` and `other` but their values and not equal (`!=`).
func (m *StrAnyMap) Diff(other *StrAnyMap) (addedKeys, removedKeys, updatedKeys []string) {
	m.lazyInit()
	return m.KVMap.Diff(other.KVMap)
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
func (m *StrAnyMap) doSetWithLockCheck(key string, value any) any {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[string]any)
	}
	if v, ok := m.data[key]; ok {
		return v
	}
	if f, ok := value.(func() any); ok {
		value = f()
	}
	if value != nil {
		m.data[key] = value
	}
	return value
}

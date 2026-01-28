// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import (
	"sync"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/util/gconv"
)

// ListMap is a map that preserves insertion-order.
//
// It is backed by a hash table to store values and doubly-linked list to store ordering.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array
type ListMap struct {
	*ListKVMap[any, any]
	once sync.Once
}

type gListMapNode = gListKVMapNode[any, any]

// NewListMap returns an empty link map.
// ListMap is backed by a hash table to store values and doubly-linked list to store ordering.
// The parameter `safe` is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewListMap(safe ...bool) *ListMap {
	return &ListMap{
		ListKVMap: NewListKVMap[any, any](safe...),
	}
}

// NewListMapFrom returns a link map from given map `data`.
// Note that, the param `data` map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewListMapFrom(data map[any]any, safe ...bool) *ListMap {
	m := NewListMap(safe...)
	m.Sets(data)
	return m
}

// lazyInit lazily initializes the list map.
func (m *ListMap) lazyInit() {
	m.once.Do(func() {
		if m.ListKVMap == nil {
			m.ListKVMap = NewListKVMap[any, any](false)
		}
	})
}

// Iterator is alias of IteratorAsc.
func (m *ListMap) Iterator(f func(key, value any) bool) {
	m.IteratorAsc(f)
}

// IteratorAsc iterates the map readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *ListMap) IteratorAsc(f func(key any, value any) bool) {
	m.lazyInit()
	m.ListKVMap.IteratorAsc(f)
}

// IteratorDesc iterates the map readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *ListMap) IteratorDesc(f func(key any, value any) bool) {
	m.lazyInit()
	m.ListKVMap.IteratorDesc(f)
}

// Clone returns a new link map with copy of current map data.
func (m *ListMap) Clone(safe ...bool) *ListMap {
	return NewListMapFrom(m.Map(), safe...)
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *ListMap) Clear() {
	m.lazyInit()
	m.ListKVMap.Clear()
}

// Replace the data of the map with given `data`.
func (m *ListMap) Replace(data map[any]any) {
	m.lazyInit()
	m.ListKVMap.Replace(data)
}

// Map returns a copy of the underlying data of the map.
func (m *ListMap) Map() map[any]any {
	m.lazyInit()
	return m.ListKVMap.Map()
}

// MapStrAny returns a copy of the underlying data of the map as map[string]any.
func (m *ListMap) MapStrAny() map[string]any {
	m.lazyInit()
	return m.ListKVMap.MapStrAny()
}

// FilterEmpty deletes all key-value pair of which the value is empty.
func (m *ListMap) FilterEmpty() {
	m.lazyInit()
	m.ListKVMap.FilterEmpty()
}

// Set sets key-value to the map.
func (m *ListMap) Set(key any, value any) {
	m.lazyInit()
	m.ListKVMap.Set(key, value)
}

// Sets batch sets key-values to the map.
func (m *ListMap) Sets(data map[any]any) {
	m.lazyInit()
	m.ListKVMap.Sets(data)
}

// Search searches the map with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (m *ListMap) Search(key any) (value any, found bool) {
	m.lazyInit()
	return m.ListKVMap.Search(key)
}

// Get returns the value by given `key`.
func (m *ListMap) Get(key any) (value any) {
	m.lazyInit()
	return m.ListKVMap.Get(key)
}

// Pop retrieves and deletes an item from the map.
func (m *ListMap) Pop() (key, value any) {
	m.lazyInit()
	return m.ListKVMap.Pop()
}

// Pops retrieves and deletes `size` items from the map.
// It returns all items if size == -1.
func (m *ListMap) Pops(size int) map[any]any {
	m.lazyInit()
	return m.ListKVMap.Pops(size)
}

// GetOrSet returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (m *ListMap) GetOrSet(key any, value any) any {
	m.lazyInit()
	return m.ListKVMap.GetOrSet(key, value)
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
func (m *ListMap) GetOrSetFunc(key any, f func() any) any {
	m.lazyInit()
	return m.ListKVMap.GetOrSetFunc(key, f)
}

// GetOrSetFuncLock returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`
// with mutex.Lock of the map.
func (m *ListMap) GetOrSetFuncLock(key any, f func() any) any {
	m.lazyInit()
	return m.ListKVMap.GetOrSetFuncLock(key, f)
}

// GetVar returns a Var with the value by given `key`.
// The returned Var is un-concurrent safe.
func (m *ListMap) GetVar(key any) *gvar.Var {
	m.lazyInit()
	return m.ListKVMap.GetVar(key)
}

// GetVarOrSet returns a Var with result from GetVarOrSet.
// The returned Var is un-concurrent safe.
func (m *ListMap) GetVarOrSet(key any, value any) *gvar.Var {
	m.lazyInit()
	return m.ListKVMap.GetVarOrSet(key, value)
}

// GetVarOrSetFunc returns a Var with result from GetOrSetFunc.
// The returned Var is un-concurrent safe.
func (m *ListMap) GetVarOrSetFunc(key any, f func() any) *gvar.Var {
	m.lazyInit()
	return m.ListKVMap.GetVarOrSetFunc(key, f)
}

// GetVarOrSetFuncLock returns a Var with result from GetOrSetFuncLock.
// The returned Var is un-concurrent safe.
func (m *ListMap) GetVarOrSetFuncLock(key any, f func() any) *gvar.Var {
	m.lazyInit()
	return m.ListKVMap.GetVarOrSetFuncLock(key, f)
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *ListMap) SetIfNotExist(key any, value any) bool {
	m.lazyInit()
	return m.ListKVMap.SetIfNotExist(key, value)
}

// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *ListMap) SetIfNotExistFunc(key any, f func() any) bool {
	m.lazyInit()
	return m.ListKVMap.SetIfNotExistFunc(key, f)
}

// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function `f` with mutex.Lock of the map.
func (m *ListMap) SetIfNotExistFuncLock(key any, f func() any) bool {
	m.lazyInit()
	return m.ListKVMap.SetIfNotExistFuncLock(key, f)
}

// Remove deletes value from map by given `key`, and return this deleted value.
func (m *ListMap) Remove(key any) (value any) {
	m.lazyInit()
	return m.ListKVMap.Remove(key)
}

// Removes batch deletes values of the map by keys.
func (m *ListMap) Removes(keys []any) {
	m.lazyInit()
	m.ListKVMap.Removes(keys)
}

// Keys returns all keys of the map as a slice in ascending order.
func (m *ListMap) Keys() []any {
	m.lazyInit()
	return m.ListKVMap.Keys()
}

// Values returns all values of the map as a slice.
func (m *ListMap) Values() []any {
	m.lazyInit()
	return m.ListKVMap.Values()
}

// Contains checks whether a key exists.
// It returns true if the `key` exists, or else false.
func (m *ListMap) Contains(key any) (ok bool) {
	m.lazyInit()
	return m.ListKVMap.Contains(key)
}

// Size returns the size of the map.
func (m *ListMap) Size() (size int) {
	m.lazyInit()
	return m.ListKVMap.Size()
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *ListMap) IsEmpty() bool {
	m.lazyInit()
	return m.ListKVMap.IsEmpty()
}

// Flip exchanges key-value of the map to value-key.
func (m *ListMap) Flip() {
	data := m.Map()
	m.Clear()
	for key, value := range data {
		m.Set(value, key)
	}
}

// Merge merges two link maps.
// The `other` map will be merged into the map `m`.
func (m *ListMap) Merge(other *ListMap) {
	m.lazyInit()
	other.lazyInit()
	m.ListKVMap.Merge(other.ListKVMap)
}

// String returns the map as a string.
func (m *ListMap) String() string {
	m.lazyInit()
	return m.ListKVMap.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m ListMap) MarshalJSON() (jsonBytes []byte, err error) {
	return m.ListKVMap.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (m *ListMap) UnmarshalJSON(b []byte) error {
	m.lazyInit()
	return m.ListKVMap.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (m *ListMap) UnmarshalValue(value any) (err error) {
	m.lazyInit()

	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range gconv.Map(value) {
		if e, ok := m.data[k]; !ok {
			m.data[k] = m.list.PushBack(&gListMapNode{k, v})
		} else {
			e.Value = &gListMapNode{k, v}
		}
	}
	return
}

// DeepCopy implements interface for deep copy of current type.
func (m *ListMap) DeepCopy() any {
	if m == nil {
		return nil
	}
	m.lazyInit()
	return &ListMap{
		ListKVMap: m.ListKVMap.DeepCopy().(*ListKVMap[any, any]),
	}
}

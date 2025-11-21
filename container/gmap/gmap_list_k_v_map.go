// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import (
	"bytes"
	"fmt"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/deepcopy"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/util/gconv"
)

// ListKVMap is a map that preserves insertion-order.
//
// It is backed by a hash table to store values and doubly-linked list to store ordering.
//
// Thread-safety is optional and controlled by the `safe` parameter during initialization.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array
type ListKVMap[K comparable, V any] struct {
	mu   rwmutex.RWMutex
	data map[K]*glist.TElement[*gListKVMapNode[K, V]]
	list *glist.TList[*gListKVMapNode[K, V]]
}

type gListKVMapNode[K comparable, V any] struct {
	key   K
	value V
}

// NewListKVMap returns an empty link map.
// ListKVMap is backed by a hash table to store values and doubly-linked list to store ordering.
// The parameter `safe` is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewListKVMap[K comparable, V any](safe ...bool) *ListKVMap[K, V] {
	return &ListKVMap[K, V]{
		mu:   rwmutex.Create(safe...),
		data: make(map[K]*glist.TElement[*gListKVMapNode[K, V]]),
		list: glist.NewT[*gListKVMapNode[K, V]](),
	}
}

// NewListKVMapFrom returns a link map from given map `data`.
// Note that, the param `data` map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewListKVMapFrom[K comparable, V any](data map[K]V, safe ...bool) *ListKVMap[K, V] {
	m := NewListKVMap[K, V](safe...)
	m.Sets(data)
	return m
}

// Iterator is alias of IteratorAsc.
func (m *ListKVMap[K, V]) Iterator(f func(key K, value V) bool) {
	m.IteratorAsc(f)
}

// IteratorAsc iterates the map readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *ListKVMap[K, V]) IteratorAsc(f func(key K, value V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.list != nil {
		m.list.IteratorAsc(func(e *glist.TElement[*gListKVMapNode[K, V]]) bool {
			return f(e.Value.key, e.Value.value)
		})
	}
}

// IteratorDesc iterates the map readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *ListKVMap[K, V]) IteratorDesc(f func(key K, value V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.list != nil {
		m.list.IteratorDesc(func(e *glist.TElement[*gListKVMapNode[K, V]]) bool {
			return f(e.Value.key, e.Value.value)
		})
	}
}

// Clone returns a new link map with copy of current map data.
func (m *ListKVMap[K, V]) Clone(safe ...bool) *ListKVMap[K, V] {
	return NewListKVMapFrom(m.Map(), safe...)
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *ListKVMap[K, V]) Clear() {
	m.mu.Lock()
	m.data = make(map[K]*glist.TElement[*gListKVMapNode[K, V]])
	m.list = glist.NewT[*gListKVMapNode[K, V]]()
	m.mu.Unlock()
}

// Replace the data of the map with given `data`.
func (m *ListKVMap[K, V]) Replace(data map[K]V) {
	m.mu.Lock()
	m.data = make(map[K]*glist.TElement[*gListKVMapNode[K, V]])
	m.list = glist.NewT[*gListKVMapNode[K, V]]()
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListKVMapNode[K, V]{key, value})
		} else {
			e.Value = &gListKVMapNode[K, V]{key, value}
		}
	}
	m.mu.Unlock()
}

// Map returns a copy of the underlying data of the map.
func (m *ListKVMap[K, V]) Map() map[K]V {
	m.mu.RLock()
	var data map[K]V
	if m.list != nil {
		data = make(map[K]V, len(m.data))
		m.list.IteratorAsc(func(e *glist.TElement[*gListKVMapNode[K, V]]) bool {
			data[e.Value.key] = e.Value.value
			return true
		})
	}
	m.mu.RUnlock()
	return data
}

// MapStrAny returns a copy of the underlying data of the map as map[string]any.
func (m *ListKVMap[K, V]) MapStrAny() map[string]any {
	m.mu.RLock()
	var data map[string]any
	if m.list != nil {
		data = make(map[string]any, len(m.data))
		m.list.IteratorAsc(func(e *glist.TElement[*gListKVMapNode[K, V]]) bool {
			data[gconv.String(e.Value.key)] = e.Value.value
			return true
		})
	}
	m.mu.RUnlock()
	return data
}

// FilterEmpty deletes all key-value pair of which the value is empty.
func (m *ListKVMap[K, V]) FilterEmpty() {
	m.mu.Lock()
	if m.list != nil {
		var keys = make([]K, 0, m.list.Size())
		m.list.IteratorAsc(func(e *glist.TElement[*gListKVMapNode[K, V]]) bool {
			if empty.IsEmpty(e.Value.value) {
				keys = append(keys, e.Value.key)
			}
			return true
		})

		if len(keys) > 0 {
			for _, key := range keys {
				if e, ok := m.data[key]; ok {
					delete(m.data, key)
					m.list.Remove(e)
				}
			}
		}
	}
	m.mu.Unlock()
}

// Set sets key-value to the map.
func (m *ListKVMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[K]*glist.TElement[*gListKVMapNode[K, V]])
		m.list = glist.NewT[*gListKVMapNode[K, V]]()
	}
	if e, ok := m.data[key]; !ok {
		m.data[key] = m.list.PushBack(&gListKVMapNode[K, V]{key, value})
	} else {
		e.Value = &gListKVMapNode[K, V]{key, value}
	}
	m.mu.Unlock()
}

// Sets batch sets key-values to the map.
func (m *ListKVMap[K, V]) Sets(data map[K]V) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[K]*glist.TElement[*gListKVMapNode[K, V]])
		m.list = glist.NewT[*gListKVMapNode[K, V]]()
	}
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListKVMapNode[K, V]{key, value})
		} else {
			e.Value = &gListKVMapNode[K, V]{key, value}
		}
	}
	m.mu.Unlock()
}

// Search searches the map with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (m *ListKVMap[K, V]) Search(key K) (value V, found bool) {
	m.mu.RLock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.value
			found = ok
		}
	}
	m.mu.RUnlock()
	return
}

// Get returns the value by given `key`.
func (m *ListKVMap[K, V]) Get(key K) (value V) {
	m.mu.RLock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.value
		}
	}
	m.mu.RUnlock()
	return
}

// Pop retrieves and deletes an item from the map.
func (m *ListKVMap[K, V]) Pop() (key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, e := range m.data {
		value = e.Value.value
		delete(m.data, k)
		m.list.Remove(e)
		return k, value
	}
	return
}

// Pops retrieves and deletes `size` items from the map.
// It returns all items if size == -1.
func (m *ListKVMap[K, V]) Pops(size int) map[K]V {
	m.mu.Lock()
	defer m.mu.Unlock()
	if size > len(m.data) || size == -1 {
		size = len(m.data)
	}
	if size == 0 {
		return nil
	}
	index := 0
	newMap := make(map[K]V, size)
	for k, e := range m.data {
		value := e.Value.value
		delete(m.data, k)
		m.list.Remove(e)
		newMap[k] = value
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
// it will be executed with mutex.Lock of the map,
// and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (m *ListKVMap[K, V]) doSetWithLockCheck(key K, value V) V {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.doSetWithLockCheckWithoutLock(key, value)
}

func (m *ListKVMap[K, V]) doSetWithLockCheckWithoutLock(key K, value V) V {
	if m.data == nil {
		m.data = make(map[K]*glist.TElement[*gListKVMapNode[K, V]])
		m.list = glist.NewT[*gListKVMapNode[K, V]]()
	}
	if e, ok := m.data[key]; ok {
		return e.Value.value
	}
	if any(value) != nil {
		m.data[key] = m.list.PushBack(&gListKVMapNode[K, V]{key, value})
	}
	return value
}

// GetOrSet returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (m *ListKVMap[K, V]) GetOrSet(key K, value V) V {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
func (m *ListKVMap[K, V]) GetOrSetFunc(key K, f func() V) V {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`
// with mutex.Lock of the map.
func (m *ListKVMap[K, V]) GetOrSetFuncLock(key K, f func() V) V {
	if v, ok := m.Search(key); !ok {
		m.mu.Lock()
		defer m.mu.Unlock()

		return m.doSetWithLockCheckWithoutLock(key, f())
	} else {
		return v
	}
}

// GetVar returns a Var with the value by given `key`.
// The returned Var is un-concurrent safe.
func (m *ListKVMap[K, V]) GetVar(key K) *gvar.Var {
	return gvar.New(m.Get(key))
}

// GetVarOrSet returns a Var with result from GetVarOrSet.
// The returned Var is un-concurrent safe.
func (m *ListKVMap[K, V]) GetVarOrSet(key K, value V) *gvar.Var {
	return gvar.New(m.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a Var with result from GetOrSetFunc.
// The returned Var is un-concurrent safe.
func (m *ListKVMap[K, V]) GetVarOrSetFunc(key K, f func() V) *gvar.Var {
	return gvar.New(m.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a Var with result from GetOrSetFuncLock.
// The returned Var is un-concurrent safe.
func (m *ListKVMap[K, V]) GetVarOrSetFuncLock(key K, f func() V) *gvar.Var {
	return gvar.New(m.GetOrSetFuncLock(key, f))
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *ListKVMap[K, V]) SetIfNotExist(key K, value V) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *ListKVMap[K, V]) SetIfNotExistFunc(key K, f func() V) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function `f` with mutex.Lock of the map.
func (m *ListKVMap[K, V]) SetIfNotExistFuncLock(key K, f func() V) bool {
	if !m.Contains(key) {
		m.mu.Lock()
		defer m.mu.Unlock()

		m.doSetWithLockCheckWithoutLock(key, f())
		return true
	}
	return false
}

// Remove deletes value from map by given `key`, and return this deleted value.
func (m *ListKVMap[K, V]) Remove(key K) (value V) {
	m.mu.Lock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.value
			delete(m.data, key)
			m.list.Remove(e)
		}
	}
	m.mu.Unlock()
	return
}

// Removes batch deletes values of the map by keys.
func (m *ListKVMap[K, V]) Removes(keys []K) {
	m.mu.Lock()
	if m.data != nil {
		for _, key := range keys {
			if e, ok := m.data[key]; ok {
				delete(m.data, key)
				m.list.Remove(e)
			}
		}
	}
	m.mu.Unlock()
}

// Keys returns all keys of the map as a slice in ascending order.
func (m *ListKVMap[K, V]) Keys() []K {
	m.mu.RLock()
	var (
		keys  = make([]K, m.list.Len())
		index = 0
	)
	if m.list != nil {
		m.list.IteratorAsc(func(e *glist.TElement[*gListKVMapNode[K, V]]) bool {
			keys[index] = e.Value.key
			index++
			return true
		})
	}
	m.mu.RUnlock()
	return keys
}

// Values returns all values of the map as a slice.
func (m *ListKVMap[K, V]) Values() []V {
	m.mu.RLock()
	var (
		values = make([]V, m.list.Len())
		index  = 0
	)
	if m.list != nil {
		m.list.IteratorAsc(func(e *glist.TElement[*gListKVMapNode[K, V]]) bool {
			values[index] = e.Value.value
			index++
			return true
		})
	}
	m.mu.RUnlock()
	return values
}

// Contains checks whether a key exists.
// It returns true if the `key` exists, or else false.
func (m *ListKVMap[K, V]) Contains(key K) (ok bool) {
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Size returns the size of the map.
func (m *ListKVMap[K, V]) Size() (size int) {
	m.mu.RLock()
	size = len(m.data)
	m.mu.RUnlock()
	return
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *ListKVMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

// Flip exchanges key-value of the map to value-key.
func (m *ListKVMap[K, V]) Flip() error {
	data := m.Map()
	m.Clear()
	for key, value := range data {
		var (
			newKey   K
			newValue V
		)

		if err := gconv.Scan(value, &newKey); err != nil {
			return err
		}

		if err := gconv.Scan(key, &newValue); err != nil {
			return err
		}
		m.Set(newKey, newValue)
	}

	return nil
}

// Merge merges two link maps.
// The `other` map will be merged into the map `m`.
func (m *ListKVMap[K, V]) Merge(other *ListKVMap[K, V]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*glist.TElement[*gListKVMapNode[K, V]])
		m.list = glist.NewT[*gListKVMapNode[K, V]]()
	}
	if other != m {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	var node *gListKVMapNode[K, V]
	other.list.IteratorAsc(func(e *glist.TElement[*gListKVMapNode[K, V]]) bool {
		node = e.Value
		if e, ok := m.data[node.key]; !ok {
			m.data[node.key] = m.list.PushBack(&gListKVMapNode[K, V]{node.key, node.value})
		} else {
			e.Value = &gListKVMapNode[K, V]{node.key, node.value}
		}
		return true
	})
}

// String returns the map as a string.
func (m *ListKVMap[K, V]) String() string {
	if m == nil {
		return ""
	}
	b, _ := m.MarshalJSON()
	return string(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m ListKVMap[K, V]) MarshalJSON() (jsonBytes []byte, err error) {
	if m.data == nil {
		return []byte("null"), nil
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('{')
	m.Iterator(func(key K, value V) bool {
		valueBytes, valueJSONErr := json.Marshal(value)
		if valueJSONErr != nil {
			err = valueJSONErr
			return false
		}
		if buffer.Len() > 1 {
			buffer.WriteByte(',')
		}
		fmt.Fprintf(buffer, `"%v":%s`, key, valueBytes)
		return true
	})
	buffer.WriteByte('}')
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (m *ListKVMap[K, V]) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*glist.TElement[*gListKVMapNode[K, V]])
		m.list = glist.NewT[*gListKVMapNode[K, V]]()
	}
	var data map[string]V
	if err := json.UnmarshalUseNumber(b, &data); err != nil {
		return err
	}
	var kvData map[K]V
	if err := gconv.Scan(data, &kvData); err != nil {
		return err
	}
	for key, value := range kvData {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListKVMapNode[K, V]{key, value})
		} else {
			e.Value = &gListKVMapNode[K, V]{key, value}
		}
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (m *ListKVMap[K, V]) UnmarshalValue(value any) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*glist.TElement[*gListKVMapNode[K, V]])
		m.list = glist.NewT[*gListKVMapNode[K, V]]()
	}
	var dataMap map[K]V
	if err = gconv.Scan(value, &dataMap); err != nil {
		return
	}
	for k, v := range dataMap {
		if e, ok := m.data[k]; !ok {
			m.data[k] = m.list.PushBack(&gListKVMapNode[K, V]{k, v})
		} else {
			e.Value = &gListKVMapNode[K, V]{k, v}
		}
	}
	return
}

// DeepCopy implements interface for deep copy of current type.
func (m *ListKVMap[K, V]) DeepCopy() any {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[any]any, len(m.data))
	if m.list != nil {
		m.list.IteratorAsc(func(e *glist.TElement[*gListKVMapNode[K, V]]) bool {
			data[e.Value.key] = deepcopy.Copy(e.Value.value)
			return true
		})
	}
	return NewListKVMapFrom(data, m.mu.IsSafe())
}

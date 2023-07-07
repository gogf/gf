// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import (
	"bytes"
	json2 "encoding/json"
	"fmt"

	"github.com/gogf/gf/contrib/generic_container/v2/conv"
	"github.com/gogf/gf/contrib/generic_container/v2/glist"
	"github.com/gogf/gf/contrib/generic_container/v2/internal/deepcopy"
	"github.com/gogf/gf/contrib/generic_container/v2/internal/empty"
	"github.com/gogf/gf/contrib/generic_container/v2/internal/json"
	"github.com/gogf/gf/contrib/generic_container/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/util/gconv"
)

// ListMap is a map that preserves insertion-order.
//
// It is backed by a hash table to store values and doubly-linked list to store ordering.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array
type ListMap[K comparable, V comparable] struct {
	mu   rwmutex.RWMutex
	data map[K]*glist.Element[*gListMapNode[K, V]]
	list *glist.List[*gListMapNode[K, V]]
}

type gListMapNode[K comparable, V comparable] struct {
	key   K
	value V
}

// NewListMap returns an empty link map.
// ListMap is backed by a hash table to store values and doubly-linked list to store ordering.
// The parameter `safe` is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewListMap[K comparable, V comparable](safe ...bool) *ListMap[K, V] {
	return &ListMap[K, V]{
		mu:   rwmutex.Create(safe...),
		data: make(map[K]*glist.Element[*gListMapNode[K, V]]),
		list: glist.New[*gListMapNode[K, V]](),
	}
}

// NewListMapFrom returns a link map from given map `data`.
// Note that, the param `data` map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewListMapFrom[K comparable, V comparable](data map[K]V, safe ...bool) *ListMap[K, V] {
	m := NewListMap[K, V](safe...)
	m.Sets(data)
	return m
}

// Iterator is alias of IteratorAsc.
func (m *ListMap[K, V]) Iterator(f func(key K, value V) bool) {
	m.IteratorAsc(f)
}

// IteratorAsc iterates the map readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *ListMap[K, V]) IteratorAsc(f func(key K, value V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.list != nil {
		var node *gListMapNode[K, V]
		m.list.IteratorAsc(func(e *glist.Element[*gListMapNode[K, V]]) bool {
			node = e.Value
			return f(node.key, node.value)
		})
	}
}

// IteratorDesc iterates the map readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *ListMap[K, V]) IteratorDesc(f func(key K, value interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.list != nil {
		var node *gListMapNode[K, V]
		m.list.IteratorDesc(func(e *glist.Element[*gListMapNode[K, V]]) bool {
			node = e.Value
			return f(node.key, node.value)
		})
	}
}

// Clone returns a new link map with copy of current map data.
func (m *ListMap[K, V]) Clone(safe ...bool) Map[K, V] {
	return NewListMapFrom[K, V](m.Map(), safe...)
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *ListMap[K, V]) Clear() {
	m.mu.Lock()
	m.data = make(map[K]*glist.Element[*gListMapNode[K, V]])
	m.list = glist.New[*gListMapNode[K, V]]()
	m.mu.Unlock()
}

// Replace the data of the map with given `data`.
func (m *ListMap[K, V]) Replace(data map[K]V) {
	m.mu.Lock()
	m.data = make(map[K]*glist.Element[*gListMapNode[K, V]])
	m.list = glist.New[*gListMapNode[K, V]]()
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
		} else {
			e.Value = &gListMapNode[K, V]{key, value}
		}
	}
	m.mu.Unlock()
}

// Map returns a copy of the underlying data of the map.
func (m *ListMap[K, V]) Map() map[K]V {
	m.mu.RLock()
	var node *gListMapNode[K, V]
	var data map[K]V
	if m.list != nil {
		data = make(map[K]V, len(m.data))
		m.list.IteratorAsc(func(e *glist.Element[*gListMapNode[K, V]]) bool {
			node = e.Value
			data[node.key] = node.value
			return true
		})
	}
	m.mu.RUnlock()
	return data
}

// MapStrAny returns a copy of the underlying data of the map as map[string]V.
func (m *ListMap[K, V]) MapStrAny() map[string]V {
	m.mu.RLock()
	var node *gListMapNode[K, V]
	var data map[string]V
	if m.list != nil {
		data = make(map[string]V, len(m.data))
		m.list.IteratorAsc(func(e *glist.Element[*gListMapNode[K, V]]) bool {
			node = e.Value
			data[gconv.String(node.key)] = node.value
			return true
		})
	}
	m.mu.RUnlock()
	return data
}

// FilterEmpty deletes all key-value pair of which the value is empty.
func (m *ListMap[K, V]) FilterEmpty() {
	m.mu.Lock()
	if m.list != nil {
		var (
			keys = make([]K, 0)
			node *gListMapNode[K, V]
		)
		m.list.IteratorAsc(func(e *glist.Element[*gListMapNode[K, V]]) bool {
			node = e.Value
			if empty.IsEmpty(node.value) {
				keys = append(keys, node.key)
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
func (m *ListMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[K]*glist.Element[*gListMapNode[K, V]])
		m.list = glist.New[*gListMapNode[K, V]]()
	}
	if e, ok := m.data[key]; !ok {
		m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
	} else {
		e.Value = &gListMapNode[K, V]{key, value}
	}
	m.mu.Unlock()
}

// Sets batch sets key-values to the map.
func (m *ListMap[K, V]) Sets(data map[K]V) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[K]*glist.Element[*gListMapNode[K, V]])
		m.list = glist.New[*gListMapNode[K, V]]()
	}
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
		} else {
			e.Value = &gListMapNode[K, V]{key, value}
		}
	}
	m.mu.Unlock()
}

// Search searches the map with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (m *ListMap[K, V]) Search(key K) (value V, found bool) {
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
func (m *ListMap[K, V]) Get(key K) (value V) {
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
func (m *ListMap[K, V]) Pop() (key K, value V) {
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
func (m *ListMap[K, V]) Pops(size int) map[K]V {
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
func (m *ListMap[K, V]) doSetWithLockCheck(key K, value V) V {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*glist.Element[*gListMapNode[K, V]])
		m.list = glist.New[*gListMapNode[K, V]]()
	}
	if e, ok := m.data[key]; ok {
		return e.Value.value
	}
	if f, ok := any(value).(func() V); ok {
		value = f()
	}
	if any(value) != nil {
		m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
	}
	return value
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
func (m *ListMap[K, V]) doSetWithLockCheckFunc(key K, f func() V) V {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*glist.Element[*gListMapNode[K, V]])
		m.list = glist.New[*gListMapNode[K, V]]()
	}
	if e, ok := m.data[key]; ok {
		return e.Value.value
	}
	var value V
	value = f()
	if any(value) != nil {
		m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
	}
	return value
}

// GetOrSet returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (m *ListMap[K, V]) GetOrSet(key K, value V) V {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
func (m *ListMap[K, V]) GetOrSetFunc(key K, f func() V) V {
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
func (m *ListMap[K, V]) GetOrSetFuncLock(key K, f func() V) V {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheckFunc(key, f)
	} else {
		return v
	}
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *ListMap[K, V]) SetIfNotExist(key K, value V) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *ListMap[K, V]) SetIfNotExistFunc(key K, f func() V) bool {
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
func (m *ListMap[K, V]) SetIfNotExistFuncLock(key K, f func() V) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheckFunc(key, f)
		return true
	}
	return false
}

// Remove deletes value from map by given `key`, and return this deleted value.
func (m *ListMap[K, V]) Remove(key K) (value V) {
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
func (m *ListMap[K, V]) Removes(keys []K) {
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
func (m *ListMap[K, V]) Keys() []K {
	m.mu.RLock()
	var (
		keys  = make([]K, m.list.Len())
		index = 0
	)
	if m.list != nil {
		m.list.IteratorAsc(func(e *glist.Element[*gListMapNode[K, V]]) bool {
			keys[index] = e.Value.key
			index++
			return true
		})
	}
	m.mu.RUnlock()
	return keys
}

// Values returns all values of the map as a slice.
func (m *ListMap[K, V]) Values() []V {
	m.mu.RLock()
	var (
		values = make([]V, m.list.Len())
		index  = 0
	)
	if m.list != nil {
		m.list.IteratorAsc(func(e *glist.Element[*gListMapNode[K, V]]) bool {
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
func (m *ListMap[K, V]) Contains(key K) (ok bool) {
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Size returns the size of the map.
func (m *ListMap[K, V]) Size() (size int) {
	m.mu.RLock()
	size = len(m.data)
	m.mu.RUnlock()
	return
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *ListMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

// Flip exchanges key-value of current map to value-key and return the new map, without modifying current map.
func (m *ListMap[K, V]) Flip() *ListMap[V, K] {
	data := m.Map()
	result := NewListMap[V, K](m.mu.IsSafe())
	for key, value := range data {
		result.Set(value, key)
	}
	return result
}

// Merge merges two link maps.
// The `other` map will be merged into the map `m`.
func (m *ListMap[K, V]) Merge(other *ListMap[K, V]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*glist.Element[*gListMapNode[K, V]])
		m.list = glist.New[*gListMapNode[K, V]]()
	}
	if other != m {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	var node *gListMapNode[K, V]
	other.list.IteratorAsc(func(e *glist.Element[*gListMapNode[K, V]]) bool {
		node = e.Value
		if e, ok := m.data[node.key]; !ok {
			m.data[node.key] = m.list.PushBack(&gListMapNode[K, V]{node.key, node.value})
		} else {
			e.Value = &gListMapNode[K, V]{node.key, node.value}
		}
		return true
	})
}

// String returns the map as a string.
func (m *ListMap[K, V]) String() string {
	if m == nil {
		return ""
	}
	b, _ := m.MarshalJSON()
	return string(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m ListMap[K, V]) MarshalJSON() (jsonBytes []byte, err error) {
	if m.data == nil {
		return []byte("null"), nil
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('{')
	m.Iterator(func(key K, value V) bool {
		valueBytes, valueJsonErr := json.Marshal(value)
		if valueJsonErr != nil {
			err = valueJsonErr
			return false
		}
		if buffer.Len() > 1 {
			buffer.WriteByte(',')
		}
		buffer.WriteString(fmt.Sprintf(`"%v":%s`, key, valueBytes))
		return true
	})
	buffer.WriteByte('}')
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (m *ListMap[K, V]) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*glist.Element[*gListMapNode[K, V]])
		m.list = glist.New[*gListMapNode[K, V]]()
	}
	var data map[K]V
	if err := json.UnmarshalUseNumber(b, &data); err != nil {
		return err
	}
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
		} else {
			e.Value = &gListMapNode[K, V]{key, value}
		}
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (m *ListMap[K, V]) UnmarshalValue(value interface{}) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*glist.Element[*gListMapNode[K, V]])
		m.list = glist.New[*gListMapNode[K, V]]()
	}
	for k, v := range gconv.Map(value) {
		kt := conv.Convert[K](k)
		var vt V
		switch v.(type) {
		case string, []byte, json2.Number:
			var ok bool
			if vt, ok = v.(V); !ok {
				if err = json.UnmarshalUseNumber(gconv.Bytes(v), &vt); err != nil {
					return err
				}
			}
		default:
			vt, _ = v.(V)
		}
		if e, ok := m.data[kt]; !ok {
			m.data[kt] = m.list.PushBack(&gListMapNode[K, V]{kt, vt})
		} else {
			e.Value = &gListMapNode[K, V]{kt, vt}
		}
	}
	return
}

// DeepCopy implements interface for deep copy of current type.
func (m *ListMap[K, V]) DeepCopy() interface{} {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[K]V, len(m.data))
	if m.list != nil {
		var node *gListMapNode[K, V]
		m.list.IteratorAsc(func(e *glist.Element[*gListMapNode[K, V]]) bool {
			node = e.Value
			data[node.key] = deepcopy.Copy(node.value).(V)
			return true
		})
	}
	return NewListMapFrom(data, m.mu.IsSafe())
}

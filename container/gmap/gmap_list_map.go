// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import (
	"github.com/gogf/gf/internal/json"

	"github.com/gogf/gf/internal/empty"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/rwmutex"
)

type ListMap struct {
	mu   rwmutex.RWMutex
	data map[interface{}]*glist.Element
	list *glist.List
}

type gListMapNode struct {
	key   interface{}
	value interface{}
}

// NewListMap returns an empty link map.
// ListMap is backed by a hash table to store values and doubly-linked list to store ordering.
// The parameter <safe> is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewListMap(safe ...bool) *ListMap {
	return &ListMap{
		mu:   rwmutex.Create(safe...),
		data: make(map[interface{}]*glist.Element),
		list: glist.New(),
	}
}

// NewListMapFrom returns a link map from given map <data>.
// Note that, the param <data> map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewListMapFrom(data map[interface{}]interface{}, safe ...bool) *ListMap {
	m := NewListMap(safe...)
	m.Sets(data)
	return m
}

// Iterator is alias of IteratorAsc.
func (m *ListMap) Iterator(f func(key, value interface{}) bool) {
	m.IteratorAsc(f)
}

// IteratorAsc iterates the map readonly in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (m *ListMap) IteratorAsc(f func(key interface{}, value interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.list != nil {
		node := (*gListMapNode)(nil)
		m.list.IteratorAsc(func(e *glist.Element) bool {
			node = e.Value.(*gListMapNode)
			return f(node.key, node.value)
		})
	}
}

// IteratorDesc iterates the map readonly in descending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (m *ListMap) IteratorDesc(f func(key interface{}, value interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.list != nil {
		node := (*gListMapNode)(nil)
		m.list.IteratorDesc(func(e *glist.Element) bool {
			node = e.Value.(*gListMapNode)
			return f(node.key, node.value)
		})
	}
}

// Clone returns a new link map with copy of current map data.
func (m *ListMap) Clone(safe ...bool) *ListMap {
	return NewListMapFrom(m.Map(), safe...)
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *ListMap) Clear() {
	m.mu.Lock()
	m.data = make(map[interface{}]*glist.Element)
	m.list = glist.New()
	m.mu.Unlock()
}

// Replace the data of the map with given <data>.
func (m *ListMap) Replace(data map[interface{}]interface{}) {
	m.mu.Lock()
	m.data = make(map[interface{}]*glist.Element)
	m.list = glist.New()
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListMapNode{key, value})
		} else {
			e.Value = &gListMapNode{key, value}
		}
	}
	m.mu.Unlock()
}

// Map returns a copy of the underlying data of the map.
func (m *ListMap) Map() map[interface{}]interface{} {
	m.mu.RLock()
	var node *gListMapNode
	var data map[interface{}]interface{}
	if m.list != nil {
		data = make(map[interface{}]interface{}, len(m.data))
		m.list.IteratorAsc(func(e *glist.Element) bool {
			node = e.Value.(*gListMapNode)
			data[node.key] = node.value
			return true
		})
	}
	m.mu.RUnlock()
	return data
}

// MapStrAny returns a copy of the underlying data of the map as map[string]interface{}.
func (m *ListMap) MapStrAny() map[string]interface{} {
	m.mu.RLock()
	var node *gListMapNode
	var data map[string]interface{}
	if m.list != nil {
		data = make(map[string]interface{}, len(m.data))
		m.list.IteratorAsc(func(e *glist.Element) bool {
			node = e.Value.(*gListMapNode)
			data[gconv.String(node.key)] = node.value
			return true
		})
	}
	m.mu.RUnlock()
	return data
}

// FilterEmpty deletes all key-value pair of which the value is empty.
func (m *ListMap) FilterEmpty() {
	m.mu.Lock()
	if m.list != nil {
		keys := make([]interface{}, 0)
		node := (*gListMapNode)(nil)
		m.list.IteratorAsc(func(e *glist.Element) bool {
			node = e.Value.(*gListMapNode)
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
func (m *ListMap) Set(key interface{}, value interface{}) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[interface{}]*glist.Element)
		m.list = glist.New()
	}
	if e, ok := m.data[key]; !ok {
		m.data[key] = m.list.PushBack(&gListMapNode{key, value})
	} else {
		e.Value = &gListMapNode{key, value}
	}
	m.mu.Unlock()
}

// Sets batch sets key-values to the map.
func (m *ListMap) Sets(data map[interface{}]interface{}) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[interface{}]*glist.Element)
		m.list = glist.New()
	}
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListMapNode{key, value})
		} else {
			e.Value = &gListMapNode{key, value}
		}
	}
	m.mu.Unlock()
}

// Search searches the map with given <key>.
// Second return parameter <found> is true if key was found, otherwise false.
func (m *ListMap) Search(key interface{}) (value interface{}, found bool) {
	m.mu.RLock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.(*gListMapNode).value
			found = ok
		}
	}
	m.mu.RUnlock()
	return
}

// Get returns the value by given <key>.
func (m *ListMap) Get(key interface{}) (value interface{}) {
	m.mu.RLock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.(*gListMapNode).value
		}
	}
	m.mu.RUnlock()
	return
}

// Pop retrieves and deletes an item from the map.
func (m *ListMap) Pop() (key, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, e := range m.data {
		value = e.Value.(*gListMapNode).value
		delete(m.data, k)
		m.list.Remove(e)
		return k, value
	}
	return
}

// Pops retrieves and deletes <size> items from the map.
// It returns all items if size == -1.
func (m *ListMap) Pops(size int) map[interface{}]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if size > len(m.data) || size == -1 {
		size = len(m.data)
	}
	if size == 0 {
		return nil
	}
	index := 0
	newMap := make(map[interface{}]interface{}, size)
	for k, e := range m.data {
		value := e.Value.(*gListMapNode).value
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
// if not exists, set value to the map with given <key>,
// or else just return the existing value.
//
// When setting value, if <value> is type of <func() interface {}>,
// it will be executed with mutex.Lock of the map,
// and its return value will be set to the map with <key>.
//
// It returns value with given <key>.
func (m *ListMap) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[interface{}]*glist.Element)
		m.list = glist.New()
	}
	if e, ok := m.data[key]; ok {
		return e.Value.(*gListMapNode).value
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value != nil {
		m.data[key] = m.list.PushBack(&gListMapNode{key, value})
	}
	return value
}

// GetOrSet returns the value by key,
// or sets value with given <value> if it does not exist and then returns this value.
func (m *ListMap) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
func (m *ListMap) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
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
// with mutex.Lock of the map.
func (m *ListMap) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

// GetVar returns a Var with the value by given <key>.
// The returned Var is un-concurrent safe.
func (m *ListMap) GetVar(key interface{}) *gvar.Var {
	return gvar.New(m.Get(key))
}

// GetVarOrSet returns a Var with result from GetVarOrSet.
// The returned Var is un-concurrent safe.
func (m *ListMap) GetVarOrSet(key interface{}, value interface{}) *gvar.Var {
	return gvar.New(m.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a Var with result from GetOrSetFunc.
// The returned Var is un-concurrent safe.
func (m *ListMap) GetVarOrSetFunc(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(m.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a Var with result from GetOrSetFuncLock.
// The returned Var is un-concurrent safe.
func (m *ListMap) GetVarOrSetFuncLock(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(m.GetOrSetFuncLock(key, f))
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *ListMap) SetIfNotExist(key interface{}, value interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *ListMap) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
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
// it executes function <f> with mutex.Lock of the map.
func (m *ListMap) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (m *ListMap) Remove(key interface{}) (value interface{}) {
	m.mu.Lock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.(*gListMapNode).value
			delete(m.data, key)
			m.list.Remove(e)
		}
	}
	m.mu.Unlock()
	return
}

// Removes batch deletes values of the map by keys.
func (m *ListMap) Removes(keys []interface{}) {
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
func (m *ListMap) Keys() []interface{} {
	m.mu.RLock()
	var (
		keys  = make([]interface{}, m.list.Len())
		index = 0
	)
	if m.list != nil {
		m.list.IteratorAsc(func(e *glist.Element) bool {
			keys[index] = e.Value.(*gListMapNode).key
			index++
			return true
		})
	}
	m.mu.RUnlock()
	return keys
}

// Values returns all values of the map as a slice.
func (m *ListMap) Values() []interface{} {
	m.mu.RLock()
	var (
		values = make([]interface{}, m.list.Len())
		index  = 0
	)
	if m.list != nil {
		m.list.IteratorAsc(func(e *glist.Element) bool {
			values[index] = e.Value.(*gListMapNode).value
			index++
			return true
		})
	}
	m.mu.RUnlock()
	return values
}

// Contains checks whether a key exists.
// It returns true if the <key> exists, or else false.
func (m *ListMap) Contains(key interface{}) (ok bool) {
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Size returns the size of the map.
func (m *ListMap) Size() (size int) {
	m.mu.RLock()
	size = len(m.data)
	m.mu.RUnlock()
	return
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *ListMap) IsEmpty() bool {
	return m.Size() == 0
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
// The <other> map will be merged into the map <m>.
func (m *ListMap) Merge(other *ListMap) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[interface{}]*glist.Element)
		m.list = glist.New()
	}
	if other != m {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	node := (*gListMapNode)(nil)
	other.list.IteratorAsc(func(e *glist.Element) bool {
		node = e.Value.(*gListMapNode)
		if e, ok := m.data[node.key]; !ok {
			m.data[node.key] = m.list.PushBack(&gListMapNode{node.key, node.value})
		} else {
			e.Value = &gListMapNode{node.key, node.value}
		}
		return true
	})
}

// String returns the map as a string.
func (m *ListMap) String() string {
	b, _ := m.MarshalJSON()
	return gconv.UnsafeBytesToStr(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m *ListMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(gconv.Map(m.Map()))
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (m *ListMap) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[interface{}]*glist.Element)
		m.list = glist.New()
	}
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListMapNode{key, value})
		} else {
			e.Value = &gListMapNode{key, value}
		}
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (m *ListMap) UnmarshalValue(value interface{}) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[interface{}]*glist.Element)
		m.list = glist.New()
	}
	for k, v := range gconv.Map(value) {
		if e, ok := m.data[k]; !ok {
			m.data[k] = m.list.PushBack(&gListMapNode{k, v})
		} else {
			e.Value = &gListMapNode{k, v}
		}
	}
	return
}

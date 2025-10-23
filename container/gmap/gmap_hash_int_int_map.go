// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

// IntIntMap implements map[int]int with RWMutex that has switch.
type IntIntMap struct {
	*KVMap[int, int]
}

// NewIntIntMap returns an empty IntIntMap object.
// The parameter `safe` is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewIntIntMap(safe ...bool) *IntIntMap {
	return &IntIntMap{
		KVMap: NewKVMap[int, int](safe...),
	}
}

// NewIntIntMapFrom creates and returns a hash map from given map `data`.
// Note that, the param `data` map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewIntIntMapFrom(data map[int]int, safe ...bool) *IntIntMap {
	return &IntIntMap{
		KVMap: NewKVMapFrom(data, safe...),
	}
}

// lazyInit lazily initializes the map.
func (m *IntIntMap) lazyInit() {
	if m.KVMap == nil {
		m.KVMap = NewKVMap[int, int](false)
	}
}

// Iterator iterates the hash map readonly with custom callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *IntIntMap) Iterator(f func(k int, v int) bool) {
	m.lazyInit()
	m.KVMap.Iterator(f)
}

// Clone returns a new hash map with copy of current map data.
func (m *IntIntMap) Clone(safe ...bool) *IntIntMap {
	m.lazyInit()
	return &IntIntMap{KVMap: m.KVMap.Clone(safe...)}
}

// Map returns the underlying data map.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (m *IntIntMap) Map() map[int]int {
	m.lazyInit()
	return m.KVMap.Map()
}

// MapStrAny returns a copy of the underlying data of the map as map[string]any.
func (m *IntIntMap) MapStrAny() map[string]any {
	m.lazyInit()
	return m.KVMap.MapStrAny()
}

// MapCopy returns a copy of the underlying data of the hash map.
func (m *IntIntMap) MapCopy() map[int]int {
	m.lazyInit()
	return m.KVMap.MapCopy()
}

// FilterEmpty deletes all key-value pair of which the value is empty.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (m *IntIntMap) FilterEmpty() {
	m.lazyInit()
	m.KVMap.FilterEmpty()
}

// Set sets key-value to the hash map.
func (m *IntIntMap) Set(key int, val int) {
	m.lazyInit()
	m.KVMap.Set(key, val)
}

// Sets batch sets key-values to the hash map.
func (m *IntIntMap) Sets(data map[int]int) {
	m.lazyInit()
	m.KVMap.Sets(data)
}

// Search searches the map with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (m *IntIntMap) Search(key int) (value int, found bool) {
	m.lazyInit()
	return m.KVMap.Search(key)
}

// Get returns the value by given `key`.
func (m *IntIntMap) Get(key int) (value int) {
	m.lazyInit()
	return m.KVMap.Get(key)
}

// Pop retrieves and deletes an item from the map.
func (m *IntIntMap) Pop() (key, value int) {
	m.lazyInit()
	return m.KVMap.Pop()
}

// Pops retrieves and deletes `size` items from the map.
// It returns all items if size == -1.
func (m *IntIntMap) Pops(size int) map[int]int {
	m.lazyInit()
	return m.KVMap.Pops(size)
}

// GetOrSet returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (m *IntIntMap) GetOrSet(key int, value int) int {
	m.lazyInit()
	return m.KVMap.GetOrSet(key, value)
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist and returns this value.
func (m *IntIntMap) GetOrSetFunc(key int, f func() int) int {
	m.lazyInit()
	return m.KVMap.GetOrSetFunc(key, f)
}

// GetOrSetFuncLock returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist and returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`
// with mutex.Lock of the hash map.
func (m *IntIntMap) GetOrSetFuncLock(key int, f func() int) int {
	m.lazyInit()
	return m.KVMap.GetOrSetFuncLock(key, f)
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *IntIntMap) SetIfNotExist(key int, value int) bool {
	m.lazyInit()
	return m.KVMap.SetIfNotExist(key, value)
}

// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *IntIntMap) SetIfNotExistFunc(key int, f func() int) bool {
	m.lazyInit()
	return m.KVMap.SetIfNotExistFunc(key, f)
}

// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function `f` with mutex.Lock of the hash map.
func (m *IntIntMap) SetIfNotExistFuncLock(key int, f func() int) bool {
	m.lazyInit()
	return m.KVMap.SetIfNotExistFuncLock(key, f)
}

// Removes batch deletes values of the map by keys.
func (m *IntIntMap) Removes(keys []int) {
	m.lazyInit()
	m.KVMap.Removes(keys)
}

// Remove deletes value from map by given `key`, and return this deleted value.
func (m *IntIntMap) Remove(key int) (value int) {
	m.lazyInit()
	return m.KVMap.Remove(key)
}

// Keys returns all keys of the map as a slice.
func (m *IntIntMap) Keys() []int {
	m.lazyInit()
	return m.KVMap.Keys()
}

// Values returns all values of the map as a slice.
func (m *IntIntMap) Values() []int {
	m.lazyInit()
	return m.KVMap.Values()
}

// Contains checks whether a key exists.
// It returns true if the `key` exists, or else false.
func (m *IntIntMap) Contains(key int) bool {
	m.lazyInit()
	return m.KVMap.Contains(key)
}

// Size returns the size of the map.
func (m *IntIntMap) Size() int {
	m.lazyInit()
	return m.KVMap.Size()
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *IntIntMap) IsEmpty() bool {
	m.lazyInit()
	return m.KVMap.IsEmpty()
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *IntIntMap) Clear() {
	m.lazyInit()
	m.KVMap.Clear()
}

// Replace the data of the map with given `data`.
func (m *IntIntMap) Replace(data map[int]int) {
	m.lazyInit()
	m.KVMap.Replace(data)
}

// LockFunc locks writing with given callback function `f` within RWMutex.Lock.
func (m *IntIntMap) LockFunc(f func(m map[int]int)) {
	m.lazyInit()
	m.KVMap.LockFunc(f)
}

// RLockFunc locks reading with given callback function `f` within RWMutex.RLock.
func (m *IntIntMap) RLockFunc(f func(m map[int]int)) {
	m.lazyInit()
	m.KVMap.RLockFunc(f)
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
// The `other` map will be merged into the map `m`.
func (m *IntIntMap) Merge(other *IntIntMap) {
	m.lazyInit()
	m.KVMap.Merge(other.KVMap)
}

// String returns the map as a string.
func (m *IntIntMap) String() string {
	if m == nil {
		return ""
	}
	m.lazyInit()
	return m.KVMap.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m IntIntMap) MarshalJSON() ([]byte, error) {
	m.lazyInit()
	return m.KVMap.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (m *IntIntMap) UnmarshalJSON(b []byte) error {
	m.lazyInit()
	return m.KVMap.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (m *IntIntMap) UnmarshalValue(value any) (err error) {
	m.lazyInit()
	return m.KVMap.UnmarshalValue(value)
}

// DeepCopy implements interface for deep copy of current type.
func (m *IntIntMap) DeepCopy() any {
	m.lazyInit()
	return &IntIntMap{
		KVMap: m.KVMap.DeepCopy().(*KVMap[int, int]),
	}
}

// IsSubOf checks whether the current map is a sub-map of `other`.
func (m *IntIntMap) IsSubOf(other *IntIntMap) bool {
	m.lazyInit()
	return m.KVMap.IsSubOf(other.KVMap)
}

// Diff compares current map `m` with map `other` and returns their different keys.
// The returned `addedKeys` are the keys that are in map `m` but not in map `other`.
// The returned `removedKeys` are the keys that are in map `other` but not in map `m`.
// The returned `updatedKeys` are the keys that are both in map `m` and `other` but their values and not equal (`!=`).
func (m *IntIntMap) Diff(other *IntIntMap) (addedKeys, removedKeys, updatedKeys []int) {
	m.lazyInit()
	return m.KVMap.Diff(other.KVMap)
}

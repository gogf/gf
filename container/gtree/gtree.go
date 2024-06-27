// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gtree provides concurrent-safe/unsafe tree containers.
//
// Some implements are from: https://github.com/emirpasic/gods
package gtree

import "github.com/gogf/gf/v2/container/gvar"

// iTree defines the interface for basic operations of a tree.
type iTree interface {
	// Set inserts node into the tree.
	Set(key interface{}, value interface{})

	// Sets batch sets key-values to the tree.
	Sets(data map[interface{}]interface{})

	// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
	// It returns false if `key` exists, and `value` would be ignored.
	SetIfNotExist(key interface{}, value interface{}) bool

	// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
	// It returns false if `key` exists, and `value` would be ignored.
	SetIfNotExistFunc(key interface{}, f func() interface{}) bool

	// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
	// It returns false if `key` exists, and `value` would be ignored.
	// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
	// it executes function `f` with mutex.Lock of the hash map.
	SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool

	// Get searches the node in the tree by `key` and returns its value or nil if key is not found in tree.
	Get(key interface{}) (value interface{})

	// GetOrSet returns the value by key,
	// or sets value with given `value` if it does not exist and then returns this value.
	GetOrSet(key interface{}, value interface{}) interface{}

	// GetOrSetFunc returns the value by key,
	// or sets value with returned value of callback function `f` if it does not exist
	// and then returns this value.
	GetOrSetFunc(key interface{}, f func() interface{}) interface{}

	// GetOrSetFuncLock returns the value by key,
	// or sets value with returned value of callback function `f` if it does not exist
	// and then returns this value.
	// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`
	// with mutex.Lock of the hash map.
	GetOrSetFuncLock(key interface{}, f func() interface{}) interface{}

	// GetVar returns a gvar.Var with the value by given `key`.
	// The returned gvar.Var is un-concurrent safe.
	GetVar(key interface{}) *gvar.Var

	// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
	// The returned gvar.Var is un-concurrent safe.
	GetVarOrSet(key interface{}, value interface{}) *gvar.Var

	// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
	// The returned gvar.Var is un-concurrent safe.
	GetVarOrSetFunc(key interface{}, f func() interface{}) *gvar.Var

	// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
	// The returned gvar.Var is un-concurrent safe.
	GetVarOrSetFuncLock(key interface{}, f func() interface{}) *gvar.Var

	// Search searches the tree with given `key`.
	// Second return parameter `found` is true if key was found, otherwise false.
	Search(key interface{}) (value interface{}, found bool)

	// Contains checks whether `key` exists in the tree.
	Contains(key interface{}) bool

	// Size returns number of nodes in the tree.
	Size() int

	// IsEmpty returns true if tree does not contain any nodes.
	IsEmpty() bool

	// Remove removes the node from the tree by key.
	// Key should adhere to the comparator's type assertion, otherwise method panics.
	Remove(key interface{}) (value interface{})

	// Removes batch deletes values of the tree by `keys`.
	Removes(keys []interface{})

	// Clear removes all nodes from the tree.
	Clear()

	// Keys returns all keys in asc order.
	Keys() []interface{}

	// Values returns all values in asc order based on the key.
	Values() []interface{}

	// Replace the data of the tree with given `data`.
	Replace(data map[interface{}]interface{})

	// Print prints the tree to stdout.
	Print()

	// String returns a string representation of container
	String() string

	// MarshalJSON implements the interface MarshalJSON for json.Marshal.
	MarshalJSON() (jsonBytes []byte, err error)

	Map() map[interface{}]interface{}
	MapStrAny() map[string]interface{}

	// Iterator is alias of IteratorAsc.
	Iterator(f func(key, value interface{}) bool)

	// IteratorFrom is alias of IteratorAscFrom.
	IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool)

	// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	IteratorAsc(f func(key, value interface{}) bool)

	// IteratorAscFrom iterates the tree readonly in ascending order with given callback function `f`.
	// The parameter `key` specifies the start entry for iterating. The `match` specifies whether
	// starting iterating if the `key` is fully matched, or else using index searching iterating.
	// If `f` returns true, then it continues iterating; or false to stop.
	IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool)

	// IteratorDesc iterates the tree readonly in descending order with given callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	IteratorDesc(f func(key, value interface{}) bool)

	// IteratorDescFrom iterates the tree readonly in descending order with given callback function `f`.
	// The parameter `key` specifies the start entry for iterating. The `match` specifies whether
	// starting iterating if the `key` is fully matched, or else using index searching iterating.
	// If `f` returns true, then it continues iterating; or false to stop.
	IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool)
}

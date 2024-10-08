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
	// Set sets key-value pair into the tree.
	Set(key interface{}, value interface{})

	// Sets batch sets key-values to the tree.
	Sets(data map[interface{}]interface{})

	// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
	// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
	SetIfNotExist(key interface{}, value interface{}) bool

	// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
	// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
	SetIfNotExistFunc(key interface{}, f func() interface{}) bool

	// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
	// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
	//
	// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
	// it executes function `f` within mutex.Lock of the hash map.
	SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool

	// Get searches the `key` in the tree and returns its associated `value` or nil if key is not found in tree.
	//
	// Note that, the `nil` value from Get function cannot be used to determine key existence, please use Contains
	// function to do so.
	Get(key interface{}) (value interface{})

	// GetOrSet returns its `value` of `key`, or sets value with given `value` if it does not exist and then returns
	// this value.
	GetOrSet(key interface{}, value interface{}) interface{}

	// GetOrSetFunc returns its `value` of `key`, or sets value with returned value of callback function `f` if it does
	// not exist and then returns this value.
	GetOrSetFunc(key interface{}, f func() interface{}) interface{}

	// GetOrSetFuncLock returns its `value` of `key`, or sets value with returned value of callback function `f` if it
	// does not exist and then returns this value.
	//
	// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f` within mutex.Lock of the
	// hash map.
	GetOrSetFuncLock(key interface{}, f func() interface{}) interface{}

	// GetVar returns a gvar.Var with the value by given `key`.
	// Note that, the returned gvar.Var is un-concurrent safe.
	//
	// Also see function Get.
	GetVar(key interface{}) *gvar.Var

	// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
	// Note that, the returned gvar.Var is un-concurrent safe.
	//
	// Also see function GetOrSet.
	GetVarOrSet(key interface{}, value interface{}) *gvar.Var

	// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
	// Note that, the returned gvar.Var is un-concurrent safe.
	//
	// Also see function GetOrSetFunc.
	GetVarOrSetFunc(key interface{}, f func() interface{}) *gvar.Var

	// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
	// Note that, the returned gvar.Var is un-concurrent safe.
	//
	// Also see function GetOrSetFuncLock.
	GetVarOrSetFuncLock(key interface{}, f func() interface{}) *gvar.Var

	// Search searches the tree with given `key`.
	// Second return parameter `found` is true if key was found, otherwise false.
	Search(key interface{}) (value interface{}, found bool)

	// Contains checks and returns whether given `key` exists in the tree.
	Contains(key interface{}) bool

	// Size returns number of nodes in the tree.
	Size() int

	// IsEmpty returns true if tree does not contain any nodes.
	IsEmpty() bool

	// Remove removes the node from the tree by `key`, and returns its associated value of `key`.
	// The given `key` should adhere to the comparator's type assertion, otherwise method panics.
	Remove(key interface{}) (value interface{})

	// Removes batch deletes key-value pairs from the tree by `keys`.
	Removes(keys []interface{})

	// Clear removes all nodes from the tree.
	Clear()

	// Keys returns all keys from the tree in order by its comparator.
	Keys() []interface{}

	// Values returns all values from the true in order by its comparator based on the key.
	Values() []interface{}

	// Replace clears the data of the tree and sets the nodes by given `data`.
	Replace(data map[interface{}]interface{})

	// Print prints the tree to stdout.
	Print()

	// String returns a string representation of container
	String() string

	// MarshalJSON implements the interface MarshalJSON for json.Marshal.
	MarshalJSON() (jsonBytes []byte, err error)

	// Map returns all key-value pairs as map.
	Map() map[interface{}]interface{}

	// MapStrAny returns all key-value items as map[string]any.
	MapStrAny() map[string]interface{}

	// Iterator is alias of IteratorAsc.
	//
	// Also see IteratorAsc.
	Iterator(f func(key, value interface{}) bool)

	// IteratorFrom is alias of IteratorAscFrom.
	//
	// Also see IteratorAscFrom.
	IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool)

	// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
	// If callback function `f` returns true, then it continues iterating; or false to stop.
	IteratorAsc(f func(key, value interface{}) bool)

	// IteratorAscFrom iterates the tree readonly in ascending order with given callback function `f`.
	//
	// The parameter `key` specifies the start entry for iterating.
	// The parameter `match` specifies whether starting iterating only if the `key` is fully matched, or else using
	// index searching iterating.
	// If callback function `f` returns true, then it continues iterating; or false to stop.
	IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool)

	// IteratorDesc iterates the tree readonly in descending order with given callback function `f`.
	//
	// If callback function `f` returns true, then it continues iterating; or false to stop.
	IteratorDesc(f func(key, value interface{}) bool)

	// IteratorDescFrom iterates the tree readonly in descending order with given callback function `f`.
	//
	// The parameter `key` specifies the start entry for iterating.
	// The parameter `match` specifies whether starting iterating only if the `key` is fully matched, or else using
	// index searching iterating.
	// If callback function `f` returns true, then it continues iterating; or false to stop.
	IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool)
}

// iteratorFromGetIndex returns the index of the key in the keys slice.
//
// The parameter `match` specifies whether starting iterating only if the `key` is fully matched,
// or else using index searching iterating.
// If `isIterator` is true, iterator is available; or else not.
func iteratorFromGetIndex(key any, keys []any, match bool) (index int, canIterator bool) {
	if match {
		for i, k := range keys {
			if k == key {
				canIterator = true
				index = i
			}
		}
	} else {
		if i, ok := key.(int); ok {
			canIterator = true
			index = i
		}
	}
	return
}

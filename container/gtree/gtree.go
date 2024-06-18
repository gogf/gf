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

// gtree defines the interface for basic operations of a tree.
type gtree interface {
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

	GetVar(key interface{}) *gvar.Var
	GetVarOrSet(key interface{}, value interface{}) *gvar.Var
	GetVarOrSetFunc(key interface{}, f func() interface{}) *gvar.Var
	GetVarOrSetFuncLock(key interface{}, f func() interface{}) *gvar.Var

	// Search searches the tree with given `key`.
	// Second return parameter `found` is true if key was found, otherwise false.
	Search(key interface{}) (value interface{}, found bool)

	// Contains checks whether `key` exists in the tree.
	Contains(key interface{}) bool

	Size() int
	IsEmpty() bool
	Remove(key interface{}) (value interface{})
	Removes(keys []interface{})
	Clear()

	Keys() []interface{}
	Values() []interface{}
	Replace(data map[interface{}]interface{})

	Print()
	String() string
	MarshalJSON() (jsonBytes []byte, err error)

	Map() map[interface{}]interface{}
	MapStrAny() map[string]interface{}

	// Iterator is alias of IteratorAsc.
	Iterator(f func(key, value interface{}) bool)
	// IteratorFrom is alias of IteratorAscFrom.
	IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool)
	IteratorAsc(f func(key, value interface{}) bool)
	IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool)
	IteratorDesc(f func(key, value interface{}) bool)
	IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool)
}

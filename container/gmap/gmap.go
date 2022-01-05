// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmap provides most commonly used map container which also support concurrent-safe/unsafe switch feature.
package gmap

type (
	Map     = AnyAnyMap // Map is alias of AnyAnyMap.
	HashMap = AnyAnyMap // HashMap is alias of AnyAnyMap.
)

// New creates and returns an empty hash map.
// The parameter `safe` is used to specify whether using map in concurrent-safety,
// which is false in default.
func New(safe ...bool) *Map {
	return NewAnyAnyMap(safe...)
}

// NewFrom creates and returns a hash map from given map `data`.
// Note that, the param `data` map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewFrom(data map[interface{}]interface{}, safe ...bool) *Map {
	return NewAnyAnyMapFrom(data, safe...)
}

// NewHashMap creates and returns an empty hash map.
// The parameter `safe` is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewHashMap(safe ...bool) *Map {
	return NewAnyAnyMap(safe...)
}

// NewHashMapFrom creates and returns a hash map from given map `data`.
// Note that, the param `data` map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewHashMapFrom(data map[interface{}]interface{}, safe ...bool) *Map {
	return NewAnyAnyMapFrom(data, safe...)
}

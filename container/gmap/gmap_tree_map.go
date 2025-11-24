// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

// TreeMap based on red-black tree, alias of TreeKVMap[any,any].
type TreeMap = TreeKVMap[any, any]

// NewTreeMap instantiates a tree map with the custom comparator.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewTreeMap(comparator func(v1, v2 any) int, safe ...bool) *TreeMap {
	return NewTreeKVMap[any, any](comparator, safe...)
}

// NewTreeMapFrom instantiates a tree map with the custom comparator and `data` map.
// Note that, the param `data` map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewTreeMapFrom(comparator func(v1, v2 any) int, data map[any]any, safe ...bool) *TreeMap {
	return NewTreeKVMapFrom(comparator, data, safe...)
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import (
	"github.com/gogf/gf/v2/container/gtree"
)

// TreeKVMap based on red-black tree, alias of RedBlackKVTree.
type TreeKVMap[K comparable, V any] = gtree.RedBlackKVTree[K, V]

// NewTreeKVMap instantiates a tree map with the custom comparator.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewTreeKVMap[K comparable, V any](comparator func(v1, v2 K) int, safe ...bool) *TreeKVMap[K, V] {
	return gtree.NewRedBlackKVTree[K, V](comparator, safe...)
}

// NewTreeKVMapFrom instantiates a tree map with the custom comparator and `data` map.
// Note that, the param `data` map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewTreeKVMapFrom[K comparable, V any](comparator func(v1, v2 K) int, data map[K]V, safe ...bool) *TreeKVMap[K, V] {
	return gtree.NewRedBlackKVTreeFrom(comparator, data, safe...)
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree

import (
	"sync"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/util/gutil"
)

var _ iTree = (*RedBlackTree)(nil)

// RedBlackTree holds elements of the red-black tree.
type RedBlackTree struct {
	*RedBlackKVTree[any, any]
	once sync.Once
}

// RedBlackTreeNode is a single element within the tree.
type RedBlackTreeNode = RedBlackKVTreeNode[any, any]

// NewRedBlackTree instantiates a red-black tree with the custom key comparator.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewRedBlackTree(comparator func(v1, v2 any) int, safe ...bool) *RedBlackTree {
	return &RedBlackTree{
		RedBlackKVTree: NewRedBlackKVTree[any, any](comparator, safe...),
	}
}

// NewRedBlackTreeFrom instantiates a red-black tree with the custom key comparator and `data` map.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewRedBlackTreeFrom(comparator func(v1, v2 any) int, data map[any]any, safe ...bool) *RedBlackTree {
	return &RedBlackTree{
		RedBlackKVTree: NewRedBlackKVTreeFrom(comparator, data, safe...),
	}
}

// lazyInit lazily initializes the tree.
func (tree *RedBlackTree) lazyInit() {
	tree.once.Do(func() {
		if tree.RedBlackKVTree == nil {
			tree.RedBlackKVTree = NewRedBlackKVTree[any, any](gutil.ComparatorTStr, false)
		}
	})
}

// SetComparator sets/changes the comparator for sorting.
func (tree *RedBlackTree) SetComparator(comparator func(a, b any) int) {
	tree.lazyInit()
	tree.RedBlackKVTree.SetComparator(comparator)
}

// Clone clones and returns a new tree from current tree.
func (tree *RedBlackTree) Clone() *RedBlackTree {
	if tree == nil {
		return nil
	}
	return &RedBlackTree{
		RedBlackKVTree: tree.RedBlackKVTree.Clone(),
	}
}

// Set sets key-value pair into the tree.
func (tree *RedBlackTree) Set(key any, value any) {
	tree.lazyInit()
	tree.RedBlackKVTree.Set(key, value)
}

// Sets batch sets key-values to the tree.
func (tree *RedBlackTree) Sets(data map[any]any) {
	tree.lazyInit()
	tree.RedBlackKVTree.Sets(data)
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
func (tree *RedBlackTree) SetIfNotExist(key any, value any) bool {
	tree.lazyInit()
	return tree.RedBlackKVTree.SetIfNotExist(key, value)
}

// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
func (tree *RedBlackTree) SetIfNotExistFunc(key any, f func() any) bool {
	tree.lazyInit()
	return tree.RedBlackKVTree.SetIfNotExistFunc(key, f)
}

// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function `f` within mutex lock.
func (tree *RedBlackTree) SetIfNotExistFuncLock(key any, f func() any) bool {
	tree.lazyInit()
	return tree.RedBlackKVTree.SetIfNotExistFuncLock(key, f)
}

// Get searches the `key` in the tree and returns its associated `value` or nil if key is not found in tree.
//
// Note that, the `nil` value from Get function cannot be used to determine key existence, please use Contains function
// to do so.
func (tree *RedBlackTree) Get(key any) (value any) {
	tree.lazyInit()
	return tree.RedBlackKVTree.Get(key)
}

// GetOrSet returns its `value` of `key`, or sets value with given `value` if it does not exist and then returns
// this value.
func (tree *RedBlackTree) GetOrSet(key any, value any) any {
	tree.lazyInit()
	return tree.RedBlackKVTree.GetOrSet(key, value)
}

// GetOrSetFunc returns its `value` of `key`, or sets value with returned value of callback function `f` if it does not
// exist and then returns this value.
func (tree *RedBlackTree) GetOrSetFunc(key any, f func() any) any {
	tree.lazyInit()
	return tree.RedBlackKVTree.GetOrSetFunc(key, f)
}

// GetOrSetFuncLock returns its `value` of `key`, or sets value with returned value of callback function `f` if it does
// not exist and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`within mutex lock.
func (tree *RedBlackTree) GetOrSetFuncLock(key any, f func() any) any {
	tree.lazyInit()
	return tree.RedBlackKVTree.GetOrSetFuncLock(key, f)
}

// GetVar returns a gvar.Var with the value by given `key`.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function Get.
func (tree *RedBlackTree) GetVar(key any) *gvar.Var {
	tree.lazyInit()
	return tree.RedBlackKVTree.GetVar(key)
}

// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSet.
func (tree *RedBlackTree) GetVarOrSet(key any, value any) *gvar.Var {
	tree.lazyInit()
	return tree.RedBlackKVTree.GetVarOrSet(key, value)
}

// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFunc.
func (tree *RedBlackTree) GetVarOrSetFunc(key any, f func() any) *gvar.Var {
	tree.lazyInit()
	return tree.RedBlackKVTree.GetVarOrSetFunc(key, f)
}

// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFuncLock.
func (tree *RedBlackTree) GetVarOrSetFuncLock(key any, f func() any) *gvar.Var {
	tree.lazyInit()
	return tree.RedBlackKVTree.GetVarOrSetFuncLock(key, f)
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *RedBlackTree) Search(key any) (value any, found bool) {
	tree.lazyInit()
	return tree.RedBlackKVTree.Search(key)
}

// Contains checks and returns whether given `key` exists in the tree.
func (tree *RedBlackTree) Contains(key any) bool {
	tree.lazyInit()
	return tree.RedBlackKVTree.Contains(key)
}

// Size returns number of nodes in the tree.
func (tree *RedBlackTree) Size() int {
	tree.lazyInit()
	return tree.RedBlackKVTree.Size()
}

// IsEmpty returns true if tree does not contain any nodes.
func (tree *RedBlackTree) IsEmpty() bool {
	tree.lazyInit()
	return tree.RedBlackKVTree.IsEmpty()
}

// Remove removes the node from the tree by `key`, and returns its associated value of `key`.
// The given `key` should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RedBlackTree) Remove(key any) (value any) {
	tree.lazyInit()
	return tree.RedBlackKVTree.Remove(key)
}

// Removes batch deletes key-value pairs from the tree by `keys`.
func (tree *RedBlackTree) Removes(keys []any) {
	tree.lazyInit()
	tree.RedBlackKVTree.Removes(keys)
}

// Clear removes all nodes from the tree.
func (tree *RedBlackTree) Clear() {
	tree.lazyInit()
	tree.RedBlackKVTree.Clear()
}

// Keys returns all keys from the tree in order by its comparator.
func (tree *RedBlackTree) Keys() []any {
	tree.lazyInit()
	return tree.RedBlackKVTree.Keys()
}

// Values returns all values from the true in order by its comparator based on the key.
func (tree *RedBlackTree) Values() []any {
	tree.lazyInit()
	return tree.RedBlackKVTree.Values()
}

// Replace clears the data of the tree and sets the nodes by given `data`.
func (tree *RedBlackTree) Replace(data map[any]any) {
	tree.lazyInit()
	tree.RedBlackKVTree.Replace(data)
}

// Print prints the tree to stdout.
func (tree *RedBlackTree) Print() {
	tree.lazyInit()
	tree.RedBlackKVTree.Print()
}

// String returns a string representation of container
func (tree *RedBlackTree) String() string {
	tree.lazyInit()
	return tree.RedBlackKVTree.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree RedBlackTree) MarshalJSON() (jsonBytes []byte, err error) {
	tree.lazyInit()
	return tree.RedBlackKVTree.MarshalJSON()
}

// Map returns all key-value pairs as map.
func (tree *RedBlackTree) Map() map[any]any {
	tree.lazyInit()
	return tree.RedBlackKVTree.Map()
}

// MapStrAny returns all key-value items as map[string]any.
func (tree *RedBlackTree) MapStrAny() map[string]any {
	tree.lazyInit()
	return tree.RedBlackKVTree.MapStrAny()
}

// Iterator is alias of IteratorAsc.
//
// Also see IteratorAsc.
func (tree *RedBlackTree) Iterator(f func(key, value any) bool) {
	tree.lazyInit()
	tree.RedBlackKVTree.Iterator(f)
}

// IteratorFrom is alias of IteratorAscFrom.
//
// Also see IteratorAscFrom.
func (tree *RedBlackTree) IteratorFrom(key any, match bool, f func(key, value any) bool) {
	tree.lazyInit()
	tree.RedBlackKVTree.IteratorFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorAsc(f func(key, value any) bool) {
	tree.lazyInit()
	tree.RedBlackKVTree.IteratorAsc(f)
}

// IteratorAscFrom iterates the tree readonly in ascending order with given callback function `f`.
//
// The parameter `key` specifies the start entry for iterating.
// The parameter `match` specifies whether starting iterating only if the `key` is fully matched, or else using index
// searching iterating.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorAscFrom(key any, match bool, f func(key, value any) bool) {
	tree.lazyInit()
	tree.RedBlackKVTree.IteratorAscFrom(key, match, f)
}

// IteratorDesc iterates the tree readonly in descending order with given callback function `f`.
//
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorDesc(f func(key, value any) bool) {
	tree.lazyInit()
	tree.RedBlackKVTree.IteratorDesc(f)
}

// IteratorDescFrom iterates the tree readonly in descending order with given callback function `f`.
//
// The parameter `key` specifies the start entry for iterating.
// The parameter `match` specifies whether starting iterating only if the `key` is fully matched, or else using index
// searching iterating.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorDescFrom(key any, match bool, f func(key, value any) bool) {
	tree.lazyInit()
	tree.RedBlackKVTree.IteratorDescFrom(key, match, f)
}

// Left returns the minimum element corresponding to the comparator of the tree or nil if the tree is empty.
func (tree *RedBlackTree) Left() *RedBlackTreeNode {
	tree.lazyInit()
	return tree.RedBlackKVTree.Left()
}

// Right returns the maximum element corresponding to the comparator of the tree or nil if the tree is empty.
func (tree *RedBlackTree) Right() *RedBlackTreeNode {
	tree.lazyInit()
	return tree.RedBlackKVTree.Right()
}

// Floor Finds floor node of the input key, returns the floor node or nil if no floor node is found.
// The second returned parameter `found` is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree is larger than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RedBlackTree) Floor(key any) (floor *RedBlackTreeNode, found bool) {
	tree.lazyInit()
	return tree.RedBlackKVTree.Floor(key)
}

// Ceiling finds ceiling node of the input key, returns the ceiling node or nil if no ceiling node is found.
// The second return parameter `found` is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree is smaller than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RedBlackTree) Ceiling(key any) (ceiling *RedBlackTreeNode, found bool) {
	tree.lazyInit()
	return tree.RedBlackKVTree.Ceiling(key)
}

// Flip exchanges key-value of the tree to value-key.
// Note that you should guarantee the value is the same type as key,
// or else the comparator would panic.
//
// If the type of value is different with key, you pass the new `comparator`.
func (tree *RedBlackTree) Flip(comparator ...func(v1, v2 any) int) {
	tree.lazyInit()
	var t = new(RedBlackTree)
	if len(comparator) > 0 {
		t = NewRedBlackTree(comparator[0], tree.mu.IsSafe())
	} else {
		t = NewRedBlackTree(tree.comparator, tree.mu.IsSafe())
	}
	tree.IteratorAsc(func(key, value any) bool {
		t.doSet(value, key)
		return true
	})
	tree.Clear()
	tree.Sets(t.Map())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (tree *RedBlackTree) UnmarshalJSON(b []byte) error {
	tree.lazyInit()
	return tree.RedBlackKVTree.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (tree *RedBlackTree) UnmarshalValue(value any) (err error) {
	tree.lazyInit()
	return tree.RedBlackKVTree.UnmarshalValue(value)
}

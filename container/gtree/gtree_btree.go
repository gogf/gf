// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree

import (
	"sync"

	"github.com/gogf/gf/v2/container/gvar"
)

var _ iTree = (*BTree)(nil)

// BTree holds elements of the B-tree.
type BTree struct {
	*BKVTree[any, any]
	once sync.Once
}

// BTreeEntry represents the key-value pair contained within nodes.
type BTreeEntry = BKVTreeEntry[any, any]

// NewBTree instantiates a B-tree with `m` (maximum number of children) and a custom key comparator.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
// Note that the `m` must be greater or equal than 3, or else it panics.
func NewBTree(m int, comparator func(v1, v2 any) int, safe ...bool) *BTree {
	return &BTree{
		BKVTree: NewBKVTree[any, any](m, comparator, safe...),
	}
}

// NewBTreeFrom instantiates a B-tree with `m` (maximum number of children), a custom key comparator and data map.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewBTreeFrom(m int, comparator func(v1, v2 any) int, data map[any]any, safe ...bool) *BTree {
	return &BTree{
		BKVTree: NewBKVTreeFrom(m, comparator, data, safe...),
	}
}

// Clone clones and returns a new tree from current tree.
func (tree *BTree) Clone() *BTree {
	if tree == nil {
		return nil
	}
	return &BTree{
		BKVTree: tree.BKVTree.Clone(),
	}
}

// Set sets key-value pair into the tree.
func (tree *BTree) Set(key any, value any) {
	tree.BKVTree.Set(key, value)
}

// Sets batch sets key-values to the tree.
func (tree *BTree) Sets(data map[any]any) {
	tree.BKVTree.Sets(data)
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
func (tree *BTree) SetIfNotExist(key any, value any) bool {
	return tree.BKVTree.SetIfNotExist(key, value)
}

// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
func (tree *BTree) SetIfNotExistFunc(key any, f func() any) bool {
	return tree.BKVTree.SetIfNotExistFunc(key, f)
}

// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function `f` within mutex lock.
func (tree *BTree) SetIfNotExistFuncLock(key any, f func() any) bool {
	return tree.BKVTree.SetIfNotExistFuncLock(key, f)
}

// Get searches the `key` in the tree and returns its associated `value` or nil if key is not found in tree.
//
// Note that, the `nil` value from Get function cannot be used to determine key existence, please use Contains function
// to do so.
func (tree *BTree) Get(key any) (value any) {
	return tree.BKVTree.Get(key)
}

// GetOrSet returns its `value` of `key`, or sets value with given `value` if it does not exist and then returns
// this value.
func (tree *BTree) GetOrSet(key any, value any) any {
	return tree.BKVTree.GetOrSet(key, value)
}

// GetOrSetFunc returns its `value` of `key`, or sets value with returned value of callback function `f` if it does not
// exist and then returns this value.
func (tree *BTree) GetOrSetFunc(key any, f func() any) any {
	return tree.BKVTree.GetOrSetFunc(key, f)
}

// GetOrSetFuncLock returns its `value` of `key`, or sets value with returned value of callback function `f` if it does
// not exist and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f` within mutex lock.
func (tree *BTree) GetOrSetFuncLock(key any, f func() any) any {
	return tree.BKVTree.GetOrSetFuncLock(key, f)
}

// GetVar returns a gvar.Var with the value by given `key`.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function Get.
func (tree *BTree) GetVar(key any) *gvar.Var {
	return tree.BKVTree.GetVar(key)
}

// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSet.
func (tree *BTree) GetVarOrSet(key any, value any) *gvar.Var {
	return tree.BKVTree.GetVarOrSet(key, value)
}

// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFunc.
func (tree *BTree) GetVarOrSetFunc(key any, f func() any) *gvar.Var {
	return tree.BKVTree.GetVarOrSetFunc(key, f)
}

// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFuncLock.
func (tree *BTree) GetVarOrSetFuncLock(key any, f func() any) *gvar.Var {
	return tree.BKVTree.GetVarOrSetFuncLock(key, f)
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *BTree) Search(key any) (value any, found bool) {
	return tree.BKVTree.Search(key)
}

// Contains checks and returns whether given `key` exists in the tree.
func (tree *BTree) Contains(key any) bool {
	return tree.BKVTree.Contains(key)
}

// Size returns number of nodes in the tree.
func (tree *BTree) Size() int {
	return tree.BKVTree.Size()
}

// IsEmpty returns true if tree does not contain any nodes
func (tree *BTree) IsEmpty() bool {
	return tree.BKVTree.IsEmpty()
}

// Remove removes the node from the tree by `key`, and returns its associated value of `key`.
// The given `key` should adhere to the comparator's type assertion, otherwise method panics.
func (tree *BTree) Remove(key any) (value any) {
	return tree.BKVTree.Remove(key)
}

// Removes batch deletes key-value pairs from the tree by `keys`.
func (tree *BTree) Removes(keys []any) {
	tree.BKVTree.Removes(keys)
}

// Clear removes all nodes from the tree.
func (tree *BTree) Clear() {
	tree.BKVTree.Clear()
}

// Keys returns all keys from the tree in order by its comparator.
func (tree *BTree) Keys() []any {
	return tree.BKVTree.Keys()
}

// Values returns all values from the true in order by its comparator based on the key.
func (tree *BTree) Values() []any {
	return tree.BKVTree.Values()
}

// Replace clears the data of the tree and sets the nodes by given `data`.
func (tree *BTree) Replace(data map[any]any) {
	tree.BKVTree.Replace(data)
}

// Map returns all key-value pairs as map.
func (tree *BTree) Map() map[any]any {
	return tree.BKVTree.Map()
}

// MapStrAny returns all key-value items as map[string]any.
func (tree *BTree) MapStrAny() map[string]any {
	return tree.BKVTree.MapStrAny()
}

// Print prints the tree to stdout.
func (tree *BTree) Print() {
	tree.BKVTree.Print()
}

// String returns a string representation of container (for debugging purposes)
func (tree *BTree) String() string {
	return tree.BKVTree.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree BTree) MarshalJSON() (jsonBytes []byte, err error) {
	return tree.BKVTree.MarshalJSON()
}

// Iterator is alias of IteratorAsc.
//
// Also see IteratorAsc.
func (tree *BTree) Iterator(f func(key, value any) bool) {
	tree.BKVTree.Iterator(f)
}

// IteratorFrom is alias of IteratorAscFrom.
//
// Also see IteratorAscFrom.
func (tree *BTree) IteratorFrom(key any, match bool, f func(key, value any) bool) {
	tree.BKVTree.IteratorFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorAsc(f func(key, value any) bool) {
	tree.BKVTree.IteratorAsc(f)
}

// IteratorAscFrom iterates the tree readonly in ascending order with given callback function `f`.
//
// The parameter `key` specifies the start entry for iterating.
// The parameter `match` specifies whether starting iterating only if the `key` is fully matched, or else using index
// searching iterating.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorAscFrom(key any, match bool, f func(key, value any) bool) {
	tree.BKVTree.IteratorAscFrom(key, match, f)
}

// IteratorDesc iterates the tree readonly in descending order with given callback function `f`.
//
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorDesc(f func(key, value any) bool) {
	tree.BKVTree.IteratorDesc(f)
}

// IteratorDescFrom iterates the tree readonly in descending order with given callback function `f`.
//
// The parameter `key` specifies the start entry for iterating.
// The parameter `match` specifies whether starting iterating only if the `key` is fully matched, or else using index
// searching iterating.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorDescFrom(key any, match bool, f func(key, value any) bool) {
	tree.BKVTree.IteratorDescFrom(key, match, f)
}

// Height returns the height of the tree.
func (tree *BTree) Height() int {
	return tree.BKVTree.Height()
}

// Left returns the minimum element corresponding to the comparator of the tree or nil if the tree is empty.
func (tree *BTree) Left() *BTreeEntry {
	return tree.BKVTree.Left()
}

// Right returns the maximum element corresponding to the comparator of the tree or nil if the tree is empty.
func (tree *BTree) Right() *BTreeEntry {
	return tree.BKVTree.Right()
}

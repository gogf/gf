// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree

import (
	"fmt"

	"github.com/emirpasic/gods/trees/btree"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

var _ iTree = (*BTree)(nil)

// BTree holds elements of the B-tree.
type BTree struct {
	mu         rwmutex.RWMutex
	comparator func(v1, v2 any) int
	m          int // order (maximum number of children)
	tree       *btree.Tree
}

// BTreeEntry represents the key-value pair contained within nodes.
type BTreeEntry struct {
	Key   any
	Value any
}

// NewBTree instantiates a B-tree with `m` (maximum number of children) and a custom key comparator.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
// Note that the `m` must be greater or equal than 3, or else it panics.
func NewBTree(m int, comparator func(v1, v2 any) int, safe ...bool) *BTree {
	return &BTree{
		mu:         rwmutex.Create(safe...),
		m:          m,
		comparator: comparator,
		tree:       btree.NewWith(m, comparator),
	}
}

// NewBTreeFrom instantiates a B-tree with `m` (maximum number of children), a custom key comparator and data map.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewBTreeFrom(m int, comparator func(v1, v2 any) int, data map[any]any, safe ...bool) *BTree {
	tree := NewBTree(m, comparator, safe...)
	for k, v := range data {
		tree.doSet(k, v)
	}
	return tree
}

// Clone clones and returns a new tree from current tree.
func (tree *BTree) Clone() *BTree {
	newTree := NewBTree(tree.m, tree.comparator, tree.mu.IsSafe())
	newTree.Sets(tree.Map())
	return newTree
}

// Set sets key-value pair into the tree.
func (tree *BTree) Set(key any, value any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.doSet(key, value)
}

// Sets batch sets key-values to the tree.
func (tree *BTree) Sets(data map[any]any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
func (tree *BTree) SetIfNotExist(key any, value any) bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, ok := tree.doGet(key); !ok {
		tree.doSet(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
func (tree *BTree) SetIfNotExistFunc(key any, f func() any) bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, ok := tree.doGet(key); !ok {
		tree.doSet(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function `f` within mutex lock.
func (tree *BTree) SetIfNotExistFuncLock(key any, f func() any) bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, ok := tree.doGet(key); !ok {
		tree.doSet(key, f)
		return true
	}
	return false
}

// Get searches the `key` in the tree and returns its associated `value` or nil if key is not found in tree.
//
// Note that, the `nil` value from Get function cannot be used to determine key existence, please use Contains function
// to do so.
func (tree *BTree) Get(key any) (value any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	value, _ = tree.doGet(key)
	return
}

// GetOrSet returns its `value` of `key`, or sets value with given `value` if it does not exist and then returns
// this value.
func (tree *BTree) GetOrSet(key any, value any) any {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if v, ok := tree.doGet(key); !ok {
		return tree.doSet(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns its `value` of `key`, or sets value with returned value of callback function `f` if it does not
// exist and then returns this value.
func (tree *BTree) GetOrSetFunc(key any, f func() any) any {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if v, ok := tree.doGet(key); !ok {
		return tree.doSet(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock returns its `value` of `key`, or sets value with returned value of callback function `f` if it does
// not exist and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f` within mutex lock.
func (tree *BTree) GetOrSetFuncLock(key any, f func() any) any {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if v, ok := tree.doGet(key); !ok {
		return tree.doSet(key, f)
	} else {
		return v
	}
}

// GetVar returns a gvar.Var with the value by given `key`.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function Get.
func (tree *BTree) GetVar(key any) *gvar.Var {
	return gvar.New(tree.Get(key))
}

// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSet.
func (tree *BTree) GetVarOrSet(key any, value any) *gvar.Var {
	return gvar.New(tree.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFunc.
func (tree *BTree) GetVarOrSetFunc(key any, f func() any) *gvar.Var {
	return gvar.New(tree.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFuncLock.
func (tree *BTree) GetVarOrSetFuncLock(key any, f func() any) *gvar.Var {
	return gvar.New(tree.GetOrSetFuncLock(key, f))
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *BTree) Search(key any) (value any, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Get(key)
}

// Contains checks and returns whether given `key` exists in the tree.
func (tree *BTree) Contains(key any) bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	_, ok := tree.doGet(key)
	return ok
}

// Size returns number of nodes in the tree.
func (tree *BTree) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Size()
}

// IsEmpty returns true if tree does not contain any nodes
func (tree *BTree) IsEmpty() bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Size() == 0
}

// Remove removes the node from the tree by `key`, and returns its associated value of `key`.
// The given `key` should adhere to the comparator's type assertion, otherwise method panics.
func (tree *BTree) Remove(key any) (value any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

// Removes batch deletes key-value pairs from the tree by `keys`.
func (tree *BTree) Removes(keys []any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.doRemove(key)
	}
}

// Clear removes all nodes from the tree.
func (tree *BTree) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.tree.Clear()
}

// Keys returns all keys from the tree in order by its comparator.
func (tree *BTree) Keys() []any {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Keys()
}

// Values returns all values from the true in order by its comparator based on the key.
func (tree *BTree) Values() []any {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Values()
}

// Replace clears the data of the tree and sets the nodes by given `data`.
func (tree *BTree) Replace(data map[any]any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.tree.Clear()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// Map returns all key-value pairs as map.
func (tree *BTree) Map() map[any]any {
	m := make(map[any]any, tree.Size())
	tree.IteratorAsc(func(key, value any) bool {
		m[key] = value
		return true
	})
	return m
}

// MapStrAny returns all key-value items as map[string]any.
func (tree *BTree) MapStrAny() map[string]any {
	m := make(map[string]any, tree.Size())
	tree.IteratorAsc(func(key, value any) bool {
		m[gconv.String(key)] = value
		return true
	})
	return m
}

// Print prints the tree to stdout.
func (tree *BTree) Print() {
	fmt.Println(tree.String())
}

// String returns a string representation of container (for debugging purposes)
func (tree *BTree) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return gstr.Replace(tree.tree.String(), "BTree\n", "")
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree *BTree) MarshalJSON() (jsonBytes []byte, err error) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.MarshalJSON()
}

// Iterator is alias of IteratorAsc.
//
// Also see IteratorAsc.
func (tree *BTree) Iterator(f func(key, value any) bool) {
	tree.IteratorAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
//
// Also see IteratorAscFrom.
func (tree *BTree) IteratorFrom(key any, match bool, f func(key, value any) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorAsc(f func(key, value any) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var (
		ok bool
		it = tree.tree.Iterator()
	)
	for it.Begin(); it.Next(); {
		index, value := it.Key(), it.Value()
		if ok = f(index, value); !ok {
			break
		}
	}
}

// IteratorAscFrom iterates the tree readonly in ascending order with given callback function `f`.
//
// The parameter `key` specifies the start entry for iterating.
// The parameter `match` specifies whether starting iterating only if the `key` is fully matched, or else using index
// searching iterating.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorAscFrom(key any, match bool, f func(key, value any) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var keys = tree.tree.Keys()
	index, canIterator := iteratorFromGetIndex(key, keys, match)
	if !canIterator {
		return
	}
	for ; index < len(keys); index++ {
		f(keys[index], tree.Get(keys[index]))
	}
}

// IteratorDesc iterates the tree readonly in descending order with given callback function `f`.
//
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorDesc(f func(key, value any) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var (
		ok bool
		it = tree.tree.Iterator()
	)
	for it.End(); it.Prev(); {
		index, value := it.Key(), it.Value()
		if ok = f(index, value); !ok {
			break
		}
	}
}

// IteratorDescFrom iterates the tree readonly in descending order with given callback function `f`.
//
// The parameter `key` specifies the start entry for iterating.
// The parameter `match` specifies whether starting iterating only if the `key` is fully matched, or else using index
// searching iterating.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorDescFrom(key any, match bool, f func(key, value any) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var keys = tree.tree.Keys()
	index, canIterator := iteratorFromGetIndex(key, keys, match)
	if !canIterator {
		return
	}
	for ; index >= 0; index-- {
		f(keys[index], tree.Get(keys[index]))
	}
}

// Height returns the height of the tree.
func (tree *BTree) Height() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Height()
}

// Left returns the minimum element corresponding to the comparator of the tree or nil if the tree is empty.
func (tree *BTree) Left() *BTreeEntry {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.tree.Left()
	if node == nil || node.Entries == nil || len(node.Entries) == 0 {
		return nil
	}
	return &BTreeEntry{
		Key:   node.Entries[0].Key,
		Value: node.Entries[0].Value,
	}
}

// Right returns the maximum element corresponding to the comparator of the tree or nil if the tree is empty.
func (tree *BTree) Right() *BTreeEntry {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.tree.Right()
	if node == nil || node.Entries == nil || len(node.Entries) == 0 {
		return nil
	}
	return &BTreeEntry{
		Key:   node.Entries[len(node.Entries)-1].Key,
		Value: node.Entries[len(node.Entries)-1].Value,
	}
}

// doSet inserts key-value pair node into the tree without lock.
// If `key` already exists, then its value is updated with the new value.
// If `value` is type of <func() any>, it will be executed and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (tree *BTree) doSet(key any, value any) any {
	if f, ok := value.(func() any); ok {
		value = f()
	}
	if value == nil {
		return value
	}
	tree.tree.Put(key, value)
	return value
}

// doGet get the value from the tree by key without lock.
func (tree *BTree) doGet(key any) (value any, ok bool) {
	return tree.tree.Get(key)
}

// doRemove removes key from tree and returns its associated value without lock.
// Note that, the given `key` should adhere to the comparator's type assertion, otherwise method panics.
func (tree *BTree) doRemove(key any) (value any) {
	value, _ = tree.tree.Get(key)
	tree.tree.Remove(key)
	return
}

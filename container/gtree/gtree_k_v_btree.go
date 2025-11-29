// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree

import (
	"fmt"

	"github.com/emirpasic/gods/v2/trees/btree"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// BKVTree holds elements of the B-tree.
type BKVTree[K comparable, V any] struct {
	mu         rwmutex.RWMutex
	comparator func(v1, v2 K) int
	m          int // order (maximum number of children)
	tree       *btree.Tree[K, V]
}

// BKVTreeEntry represents the key-value pair contained within nodes.
type BKVTreeEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// NewBKVTree instantiates a B-tree with `m` (maximum number of children) and a custom key comparator.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
// Note that the `m` must be greater or equal than 3, or else it panics.
func NewBKVTree[K comparable, V any](m int, comparator func(v1, v2 K) int, safe ...bool) *BKVTree[K, V] {
	return &BKVTree[K, V]{
		mu:         rwmutex.Create(safe...),
		m:          m,
		comparator: comparator,
		tree:       btree.NewWith[K, V](m, comparator),
	}
}

// NewBKVTreeFrom instantiates a B-tree with `m` (maximum number of children), a custom key comparator and data map.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewBKVTreeFrom[K comparable, V any](m int, comparator func(v1, v2 K) int, data map[K]V, safe ...bool) *BKVTree[K, V] {
	tree := NewBKVTree[K, V](m, comparator, safe...)
	for k, v := range data {
		tree.doSet(k, v)
	}
	return tree
}

// Clone clones and returns a new tree from current tree.
func (tree *BKVTree[K, V]) Clone() *BKVTree[K, V] {
	if tree == nil {
		return nil
	}
	newTree := NewBKVTree[K, V](tree.m, tree.comparator, tree.mu.IsSafe())
	newTree.Sets(tree.Map())
	return newTree
}

// Set sets key-value pair into the tree.
func (tree *BKVTree[K, V]) Set(key K, value V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.doSet(key, value)
}

// Sets batch sets key-values to the tree.
func (tree *BKVTree[K, V]) Sets(data map[K]V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
func (tree *BKVTree[K, V]) SetIfNotExist(key K, value V) bool {
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
func (tree *BKVTree[K, V]) SetIfNotExistFunc(key K, f func() V) bool {
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
func (tree *BKVTree[K, V]) SetIfNotExistFuncLock(key K, f func() V) bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, ok := tree.doGet(key); !ok {
		tree.doSet(key, f())
		return true
	}
	return false
}

// Get searches the `key` in the tree and returns its associated `value` or nil if key is not found in tree.
//
// Note that, the `nil` value from Get function cannot be used to determine key existence, please use Contains function
// to do so.
func (tree *BKVTree[K, V]) Get(key K) (value V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	value, _ = tree.doGet(key)
	return
}

// GetOrSet returns its `value` of `key`, or sets value with given `value` if it does not exist and then returns
// this value.
func (tree *BKVTree[K, V]) GetOrSet(key K, value V) V {
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
func (tree *BKVTree[K, V]) GetOrSetFunc(key K, f func() V) V {
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
func (tree *BKVTree[K, V]) GetOrSetFuncLock(key K, f func() V) V {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if v, ok := tree.doGet(key); !ok {
		return tree.doSet(key, f())
	} else {
		return v
	}
}

// GetVar returns a gvar.Var with the value by given `key`.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function Get.
func (tree *BKVTree[K, V]) GetVar(key K) *gvar.Var {
	return gvar.New(tree.Get(key))
}

// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSet.
func (tree *BKVTree[K, V]) GetVarOrSet(key K, value V) *gvar.Var {
	return gvar.New(tree.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFunc.
func (tree *BKVTree[K, V]) GetVarOrSetFunc(key K, f func() V) *gvar.Var {
	return gvar.New(tree.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFuncLock.
func (tree *BKVTree[K, V]) GetVarOrSetFuncLock(key K, f func() V) *gvar.Var {
	return gvar.New(tree.GetOrSetFuncLock(key, f))
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *BKVTree[K, V]) Search(key K) (value V, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Get(key)
}

// Contains checks and returns whether given `key` exists in the tree.
func (tree *BKVTree[K, V]) Contains(key K) bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	_, ok := tree.doGet(key)
	return ok
}

// Size returns number of nodes in the tree.
func (tree *BKVTree[K, V]) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Size()
}

// IsEmpty returns true if tree does not contain any nodes
func (tree *BKVTree[K, V]) IsEmpty() bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Size() == 0
}

// Remove removes the node from the tree by `key`, and returns its associated value of `key`.
// The given `key` should adhere to the comparator's type assertion, otherwise method panics.
func (tree *BKVTree[K, V]) Remove(key K) (value V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

// Removes batch deletes key-value pairs from the tree by `keys`.
func (tree *BKVTree[K, V]) Removes(keys []K) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.doRemove(key)
	}
}

// Clear removes all nodes from the tree.
func (tree *BKVTree[K, V]) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.tree.Clear()
}

// Keys returns all keys from the tree in order by its comparator.
func (tree *BKVTree[K, V]) Keys() []K {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Keys()
}

// Values returns all values from the true in order by its comparator based on the key.
func (tree *BKVTree[K, V]) Values() []V {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Values()
}

// Replace clears the data of the tree and sets the nodes by given `data`.
func (tree *BKVTree[K, V]) Replace(data map[K]V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.tree.Clear()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// Map returns all key-value pairs as map.
func (tree *BKVTree[K, V]) Map() map[K]V {
	m := make(map[K]V, tree.Size())
	tree.IteratorAsc(func(key K, value V) bool {
		m[key] = value
		return true
	})
	return m
}

// MapStrAny returns all key-value items as map[string]any.
func (tree *BKVTree[K, V]) MapStrAny() map[string]any {
	m := make(map[string]any, tree.Size())
	tree.IteratorAsc(func(key K, value V) bool {
		m[gconv.String(key)] = value
		return true
	})
	return m
}

// Print prints the tree to stdout.
func (tree *BKVTree[K, V]) Print() {
	fmt.Println(tree.String())
}

// String returns a string representation of container (for debugging purposes)
func (tree *BKVTree[K, V]) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return gstr.Replace(tree.tree.String(), "BTree\n", "")
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree *BKVTree[K, V]) MarshalJSON() (jsonBytes []byte, err error) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	elements := make(map[string]V)
	it := tree.tree.Iterator()
	for it.Next() {
		elements[gconv.String(it.Key())] = it.Value()
	}
	return json.Marshal(&elements)
}

// Iterator is alias of IteratorAsc.
//
// Also see IteratorAsc.
func (tree *BKVTree[K, V]) Iterator(f func(key K, value V) bool) {
	tree.IteratorAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
//
// Also see IteratorAscFrom.
func (tree *BKVTree[K, V]) IteratorFrom(key K, match bool, f func(key K, value V) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *BKVTree[K, V]) IteratorAsc(f func(key K, value V) bool) {
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
func (tree *BKVTree[K, V]) IteratorAscFrom(key K, match bool, f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var keys = tree.tree.Keys()
	index, canIterator := iteratorFromGetIndexT(key, keys, match)
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
func (tree *BKVTree[K, V]) IteratorDesc(f func(key K, value V) bool) {
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
func (tree *BKVTree[K, V]) IteratorDescFrom(key K, match bool, f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var keys = tree.tree.Keys()
	index, canIterator := iteratorFromGetIndexT(key, keys, match)
	if !canIterator {
		return
	}
	for ; index >= 0; index-- {
		f(keys[index], tree.Get(keys[index]))
	}
}

// Height returns the height of the tree.
func (tree *BKVTree[K, V]) Height() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Height()
}

// Left returns the minimum element corresponding to the comparator of the tree or nil if the tree is empty.
func (tree *BKVTree[K, V]) Left() *BKVTreeEntry[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.tree.Left()
	if node == nil || node.Entries == nil || len(node.Entries) == 0 {
		return nil
	}
	return &BKVTreeEntry[K, V]{
		Key:   node.Entries[0].Key,
		Value: node.Entries[0].Value,
	}
}

// Right returns the maximum element corresponding to the comparator of the tree or nil if the tree is empty.
func (tree *BKVTree[K, V]) Right() *BKVTreeEntry[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.tree.Right()
	if node == nil || node.Entries == nil || len(node.Entries) == 0 {
		return nil
	}
	return &BKVTreeEntry[K, V]{
		Key:   node.Entries[len(node.Entries)-1].Key,
		Value: node.Entries[len(node.Entries)-1].Value,
	}
}

// doSet inserts key-value pair node into the tree without lock.
// If `key` already exists, then its value is updated with the new value.
// If `value` is type of <func() any>, it will be executed and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (tree *BKVTree[K, V]) doSet(key K, value V) V {
	if any(value) == nil {
		return value
	}
	tree.tree.Put(key, value)
	return value
}

// doGet get the value from the tree by key without lock.
func (tree *BKVTree[K, V]) doGet(key K) (value V, ok bool) {
	return tree.tree.Get(key)
}

// doRemove removes key from tree and returns its associated value without lock.
// Note that, the given `key` should adhere to the comparator's type assertion, otherwise method panics.
func (tree *BKVTree[K, V]) doRemove(key K) (value V) {
	value, _ = tree.tree.Get(key)
	tree.tree.Remove(key)
	return
}

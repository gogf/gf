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
	comparator func(v1, v2 interface{}) int
	m          int // order (maximum number of children)
	tree       *btree.Tree
}

// BTreeEntry represents the key-value pair contained within nodes.
type BTreeEntry struct {
	Key   interface{}
	Value interface{}
}

// NewBTree instantiates a B-tree with `m` (maximum number of children) and a custom key comparator.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
// Note that the `m` must be greater or equal than 3, or else it panics.
func NewBTree(m int, comparator func(v1, v2 interface{}) int, safe ...bool) *BTree {
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
func NewBTreeFrom(m int, comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}, safe ...bool) *BTree {
	tree := NewBTree(m, comparator, safe...)
	for k, v := range data {
		tree.doSet(k, v)
	}
	return tree
}

// Clone returns a new tree with a copy of current tree.
func (tree *BTree) Clone() *BTree {
	newTree := NewBTree(tree.m, tree.comparator, tree.mu.IsSafe())
	newTree.Sets(tree.Map())
	return newTree
}

// Set inserts key-value item into the tree.
func (tree *BTree) Set(key interface{}, value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.doSet(key, value)
}

// Sets batch sets key-values to the tree.
func (tree *BTree) Sets(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (tree *BTree) SetIfNotExist(key interface{}, value interface{}) bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, ok := tree.doGet(key); !ok {
		tree.doSet(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (tree *BTree) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, ok := tree.doGet(key); !ok {
		tree.doSet(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function `f` with mutex.Lock of the hash map.
func (tree *BTree) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, ok := tree.doGet(key); !ok {
		tree.doSet(key, f)
		return true
	}
	return false
}

// Get searches the node in the tree by `key` and returns its value or nil if key is not found in tree.
func (tree *BTree) Get(key interface{}) (value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	value, _ = tree.doGet(key)
	return
}

// GetOrSet returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (tree *BTree) GetOrSet(key interface{}, value interface{}) interface{} {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if v, ok := tree.doGet(key); !ok {
		return tree.doSet(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
func (tree *BTree) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if v, ok := tree.doGet(key); !ok {
		return tree.doSet(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`
// with mutex.Lock of the hash map.
func (tree *BTree) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if v, ok := tree.doGet(key); !ok {
		return tree.doSet(key, f)
	} else {
		return v
	}
}

// GetVar returns a gvar.Var with the value by given `key`.
// The returned gvar.Var is un-concurrent safe.
func (tree *BTree) GetVar(key interface{}) *gvar.Var {
	return gvar.New(tree.Get(key))
}

// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
// The returned gvar.Var is un-concurrent safe.
func (tree *BTree) GetVarOrSet(key interface{}, value interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
// The returned gvar.Var is un-concurrent safe.
func (tree *BTree) GetVarOrSetFunc(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
// The returned gvar.Var is un-concurrent safe.
func (tree *BTree) GetVarOrSetFuncLock(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSetFuncLock(key, f))
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *BTree) Search(key interface{}) (value interface{}, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Get(key)
}

// Contains checks whether `key` exists in the tree.
func (tree *BTree) Contains(key interface{}) bool {
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

// Remove removes the node from the tree by `key`.
func (tree *BTree) Remove(key interface{}) (value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

// Removes batch deletes values of the tree by `keys`.
func (tree *BTree) Removes(keys []interface{}) {
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

// Keys returns all keys in asc order.
func (tree *BTree) Keys() []interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Keys()
}

// Values returns all values in asc order based on the key.
func (tree *BTree) Values() []interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Values()
}

// Replace the data of the tree with given `data`.
func (tree *BTree) Replace(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.tree.Clear()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// Map returns all key-value items as map.
func (tree *BTree) Map() map[interface{}]interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	m := make(map[interface{}]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

// MapStrAny returns all key-value items as map[string]interface{}.
func (tree *BTree) MapStrAny() map[string]interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	m := make(map[string]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
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
func (tree *BTree) Iterator(f func(key, value interface{}) bool) {
	tree.IteratorAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
func (tree *BTree) IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorAsc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	it := tree.tree.Iterator()
	for it.Begin(); it.Next(); {
		index, value := it.Key(), it.Value()
		if ok := f(index, value); !ok {
			break
		}
	}
}

// IteratorAscFrom iterates the tree readonly in ascending order with given callback function `f`.
// The parameter `key` specifies the start entry for iterating. The `match` specifies whether
// starting iterating if the `key` is fully matched, or else using index searching iterating.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var keys = tree.tree.Keys()
	index, isIterator := tree.iteratorFromGetIndex(key, keys, match)
	if !isIterator {
		return
	}
	for ; index < len(keys); index++ {
		f(keys[index], tree.Get(keys[index]))
	}
}

// IteratorDesc iterates the tree readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorDesc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	it := tree.tree.Iterator()
	for it.End(); it.Prev(); {
		index, value := it.Key(), it.Value()
		if ok := f(index, value); !ok {
			break
		}
	}
}

// IteratorDescFrom iterates the tree readonly in descending order with given callback function `f`.
// The parameter `key` specifies the start entry for iterating. The `match` specifies whether
// starting iterating if the `key` is fully matched, or else using index searching iterating.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var keys = tree.tree.Keys()
	index, isIterator := tree.iteratorFromGetIndex(key, keys, match)
	if !isIterator {
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

// Left returns the left-most (min) entry or nil if tree is empty.
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

// Right returns the right-most (max) entry or nil if tree is empty.
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

// doSet inserts key-value pair node into the tree.
// If key already exists, then its value is updated with the new value.
// If `value` is type of <func() interface {}>,
// it will be executed and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (tree *BTree) doSet(key interface{}, value interface{}) interface{} {
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value == nil {
		return value
	}
	tree.tree.Put(key, value)
	return value
}

// doGet get the value from the tree by key.
func (tree *BTree) doGet(key interface{}) (value interface{}, ok bool) {
	return tree.tree.Get(key)
}

// doRemove removes the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *BTree) doRemove(key interface{}) (value interface{}) {
	value, _ = tree.tree.Get(key)
	tree.tree.Remove(key)
	return
}

// iteratorFromGetIndex returns the index of the key in the keys slice.
// The parameter `match` specifies whether starting iterating if the `key` is fully matched,
// or else using index searching iterating.
// If `isIterator` is true, iterator is available; or else not.
func (tree *BTree) iteratorFromGetIndex(key interface{}, keys []interface{}, match bool) (index int, isIterator bool) {
	if match {
		for i, k := range keys {
			if k == key {
				isIterator = true
				index = i
			}
		}
	} else {
		if i, ok := key.(int); ok {
			isIterator = true
			index = i
		}
	}
	return
}

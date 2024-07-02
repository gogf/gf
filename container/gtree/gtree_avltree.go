// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree

import (
	"fmt"

	"github.com/emirpasic/gods/trees/avltree"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

var _ iTree = (*AVLTree)(nil)

// AVLTree holds elements of the AVL tree.
type AVLTree struct {
	mu         rwmutex.RWMutex
	root       *AVLTreeNode
	comparator func(v1, v2 interface{}) int
	tree       *avltree.Tree
}

// AVLTreeNode is a single element within the tree.
type AVLTreeNode struct {
	Key   interface{}
	Value interface{}
}

// NewAVLTree instantiates an AVL tree with the custom key comparator.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewAVLTree(comparator func(v1, v2 interface{}) int, safe ...bool) *AVLTree {
	return &AVLTree{
		mu:         rwmutex.Create(safe...),
		comparator: comparator,
		tree:       avltree.NewWith(comparator),
	}
}

// NewAVLTreeFrom instantiates an AVL tree with the custom key comparator and data map.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewAVLTreeFrom(comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}, safe ...bool) *AVLTree {
	tree := NewAVLTree(comparator, safe...)
	for k, v := range data {
		tree.doSet(k, v)
	}
	return tree
}

// Clone returns a new tree with a copy of current tree.
func (tree *AVLTree) Clone() *AVLTree {
	newTree := NewAVLTree(tree.comparator, tree.mu.IsSafe())
	newTree.Sets(tree.Map())
	return newTree
}

// Set inserts node into the tree.
func (tree *AVLTree) Set(key interface{}, value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.doSet(key, value)
}

// Sets batch sets key-values to the tree.
func (tree *AVLTree) Sets(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for key, value := range data {
		tree.doSet(key, value)
	}
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (tree *AVLTree) SetIfNotExist(key interface{}, value interface{}) bool {
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
func (tree *AVLTree) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
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
func (tree *AVLTree) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, ok := tree.doGet(key); !ok {
		tree.doSet(key, f)
		return true
	}
	return false
}

// Get searches the node in the tree by `key` and returns its value or nil if key is not found in tree.
func (tree *AVLTree) Get(key interface{}) (value interface{}) {
	value, _ = tree.Search(key)
	return
}

// GetOrSet returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (tree *AVLTree) GetOrSet(key interface{}, value interface{}) interface{} {
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
func (tree *AVLTree) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
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
func (tree *AVLTree) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
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
func (tree *AVLTree) GetVar(key interface{}) *gvar.Var {
	return gvar.New(tree.Get(key))
}

// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
// The returned gvar.Var is un-concurrent safe.
func (tree *AVLTree) GetVarOrSet(key interface{}, value interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
// The returned gvar.Var is un-concurrent safe.
func (tree *AVLTree) GetVarOrSetFunc(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
// The returned gvar.Var is un-concurrent safe.
func (tree *AVLTree) GetVarOrSetFuncLock(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSetFuncLock(key, f))
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *AVLTree) Search(key interface{}) (value interface{}, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if node, found := tree.doGet(key); found {
		return node, true
	}
	return nil, false
}

// Contains checks whether `key` exists in the tree.
func (tree *AVLTree) Contains(key interface{}) bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	_, ok := tree.doGet(key)
	return ok
}

// Size returns number of nodes in the tree.
func (tree *AVLTree) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Size()
}

// IsEmpty returns true if tree does not contain any nodes.
func (tree *AVLTree) IsEmpty() bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Size() == 0
}

// Remove removes the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *AVLTree) Remove(key interface{}) (value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

// Removes batch deletes values of the tree by `keys`.
func (tree *AVLTree) Removes(keys []interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.doRemove(key)
	}
}

// Clear removes all nodes from the tree.
func (tree *AVLTree) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.tree.Clear()
}

// Keys returns all keys in asc order.
func (tree *AVLTree) Keys() []interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Keys()
}

// Values returns all values in asc order based on the key.
func (tree *AVLTree) Values() []interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Values()
}

// Replace the data of the tree with given `data`.
func (tree *AVLTree) Replace(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.tree.Clear()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// Print prints the tree to stdout.
func (tree *AVLTree) Print() {
	fmt.Println(tree.String())
}

// String returns a string representation of container
func (tree *AVLTree) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return gstr.Replace(tree.tree.String(), "AVLTree\n", "")
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree *AVLTree) MarshalJSON() (jsonBytes []byte, err error) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.MarshalJSON()
}

// Map returns all key-value items as map.
func (tree *AVLTree) Map() map[interface{}]interface{} {
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
func (tree *AVLTree) MapStrAny() map[string]interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	m := make(map[string]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[gconv.String(key)] = value
		return true
	})
	return m
}

// Iterator is alias of IteratorAsc.
func (tree *AVLTree) Iterator(f func(key, value interface{}) bool) {
	tree.IteratorAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
func (tree *AVLTree) IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *AVLTree) IteratorAsc(f func(key, value interface{}) bool) {
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
func (tree *AVLTree) IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
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
func (tree *AVLTree) IteratorDesc(f func(key, value interface{}) bool) {
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
func (tree *AVLTree) IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
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

// Left returns the minimum element of the AVL tree
// or nil if the tree is empty.
func (tree *AVLTree) Left() *AVLTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.tree.Left()
	if node == nil {
		return nil
	}
	return &AVLTreeNode{
		Key:   node.Key,
		Value: node.Value,
	}
}

// Right returns the maximum element of the AVL tree
// or nil if the tree is empty.
func (tree *AVLTree) Right() *AVLTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.tree.Right()
	if node == nil {
		return nil
	}
	return &AVLTreeNode{
		Key:   node.Key,
		Value: node.Value,
	}
}

// Floor Finds floor node of the input key, return the floor node or nil if no floor node is found.
// Second return parameter is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree is larger than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *AVLTree) Floor(key interface{}) (floor *AVLTreeNode, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.tree.Floor(key)
	if !found {
		return nil, false
	}
	return &AVLTreeNode{
		Key:   node.Key,
		Value: node.Value,
	}, true
}

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling node is found.
// Second return parameter is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree is smaller than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *AVLTree) Ceiling(key interface{}) (ceiling *AVLTreeNode, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.tree.Ceiling(key)
	if !found {
		return nil, false
	}
	return &AVLTreeNode{
		Key:   node.Key,
		Value: node.Value,
	}, true
}

// Flip exchanges key-value of the tree to value-key.
// Note that you should guarantee the value is the same type as key,
// or else the comparator would panic.
//
// If the type of value is different with key, you pass the new `comparator`.
func (tree *AVLTree) Flip(comparator ...func(v1, v2 interface{}) int) {
	var t = new(AVLTree)
	if len(comparator) > 0 {
		t = NewAVLTree(comparator[0], tree.mu.IsSafe())
	} else {
		t = NewAVLTree(tree.comparator, tree.mu.IsSafe())
	}
	tree.IteratorAsc(func(key, value interface{}) bool {
		t.doSet(value, key)
		return true
	})
	tree.Clear()
	tree.Sets(t.Map())
}

// doSet sets key-value pair to the tree.
func (tree *AVLTree) doSet(key, value interface{}) interface{} {
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value == nil {
		return value
	}
	tree.tree.Put(key, value)
	return value
}

// doGet retrieves and returns the value of given key from tree.
func (tree *AVLTree) doGet(key interface{}) (value interface{}, found bool) {
	return tree.tree.Get(key)
}

// doRemove removes key from tree.
func (tree *AVLTree) doRemove(key interface{}) (value interface{}) {
	value, _ = tree.tree.Get(key)
	tree.tree.Remove(key)
	return
}

// iteratorFromGetIndex returns the index of the key in the keys slice.
// The parameter `match` specifies whether starting iterating if the `key` is fully matched,
// or else using index searching iterating.
// If `isIterator` is true, iterator is available; or else not.
func (tree *AVLTree) iteratorFromGetIndex(key interface{}, keys []interface{}, match bool) (index int, isIterator bool) {
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

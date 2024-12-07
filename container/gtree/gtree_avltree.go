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
	comparator func(v1, v2 any) int
	tree       *avltree.Tree
}

// AVLTreeNode is a single element within the tree.
type AVLTreeNode struct {
	Key   any
	Value any
}

// NewAVLTree instantiates an AVL tree with the custom key comparator.
//
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewAVLTree(comparator func(v1, v2 any) int, safe ...bool) *AVLTree {
	return &AVLTree{
		mu:         rwmutex.Create(safe...),
		comparator: comparator,
		tree:       avltree.NewWith(comparator),
	}
}

// NewAVLTreeFrom instantiates an AVL tree with the custom key comparator and data map.
//
// The parameter `safe` is used to specify whether using tree in concurrent-safety, which is false in default.
func NewAVLTreeFrom(comparator func(v1, v2 any) int, data map[any]any, safe ...bool) *AVLTree {
	tree := NewAVLTree(comparator, safe...)
	for k, v := range data {
		tree.doSet(k, v)
	}
	return tree
}

// Clone clones and returns a new tree from current tree.
func (tree *AVLTree) Clone() *AVLTree {
	newTree := NewAVLTree(tree.comparator, tree.mu.IsSafe())
	newTree.Sets(tree.Map())
	return newTree
}

// Set sets key-value pair into the tree.
func (tree *AVLTree) Set(key any, value any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.doSet(key, value)
}

// Sets batch sets key-values to the tree.
func (tree *AVLTree) Sets(data map[any]any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for key, value := range data {
		tree.doSet(key, value)
	}
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
func (tree *AVLTree) SetIfNotExist(key any, value any) bool {
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
func (tree *AVLTree) SetIfNotExistFunc(key any, f func() any) bool {
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
func (tree *AVLTree) SetIfNotExistFuncLock(key any, f func() any) bool {
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
func (tree *AVLTree) Get(key any) (value any) {
	value, _ = tree.Search(key)
	return
}

// GetOrSet returns its `value` of `key`, or sets value with given `value` if it does not exist and then returns
// this value.
func (tree *AVLTree) GetOrSet(key any, value any) any {
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
func (tree *AVLTree) GetOrSetFunc(key any, f func() any) any {
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
func (tree *AVLTree) GetOrSetFuncLock(key any, f func() any) any {
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
func (tree *AVLTree) GetVar(key any) *gvar.Var {
	return gvar.New(tree.Get(key))
}

// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSet.
func (tree *AVLTree) GetVarOrSet(key any, value any) *gvar.Var {
	return gvar.New(tree.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFunc.
func (tree *AVLTree) GetVarOrSetFunc(key any, f func() any) *gvar.Var {
	return gvar.New(tree.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFuncLock.
func (tree *AVLTree) GetVarOrSetFuncLock(key any, f func() any) *gvar.Var {
	return gvar.New(tree.GetOrSetFuncLock(key, f))
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *AVLTree) Search(key any) (value any, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if node, found := tree.doGet(key); found {
		return node, true
	}
	return nil, false
}

// Contains checks and returns whether given `key` exists in the tree.
func (tree *AVLTree) Contains(key any) bool {
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

// IsEmpty returns true if the tree does not contain any nodes.
func (tree *AVLTree) IsEmpty() bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Size() == 0
}

// Remove removes the node from the tree by `key`, and returns its associated value of `key`.
// The given `key` should adhere to the comparator's type assertion, otherwise method panics.
func (tree *AVLTree) Remove(key any) (value any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

// Removes batch deletes key-value pairs from the tree by `keys`.
func (tree *AVLTree) Removes(keys []any) {
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

// Keys returns all keys from the tree in order by its comparator.
func (tree *AVLTree) Keys() []any {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Keys()
}

// Values returns all values from the true in order by its comparator based on the key.
func (tree *AVLTree) Values() []any {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Values()
}

// Replace clears the data of the tree and sets the nodes by given `data`.
func (tree *AVLTree) Replace(data map[any]any) {
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

// String returns a string representation of container.
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

// Map returns all key-value pairs as map.
func (tree *AVLTree) Map() map[any]any {
	m := make(map[any]any, tree.Size())
	tree.IteratorAsc(func(key, value any) bool {
		m[key] = value
		return true
	})
	return m
}

// MapStrAny returns all key-value items as map[string]any.
func (tree *AVLTree) MapStrAny() map[string]any {
	m := make(map[string]any, tree.Size())
	tree.IteratorAsc(func(key, value any) bool {
		m[gconv.String(key)] = value
		return true
	})
	return m
}

// Iterator is alias of IteratorAsc.
//
// Also see IteratorAsc.
func (tree *AVLTree) Iterator(f func(key, value any) bool) {
	tree.IteratorAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
//
// Also see IteratorAscFrom.
func (tree *AVLTree) IteratorFrom(key any, match bool, f func(key, value any) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *AVLTree) IteratorAsc(f func(key, value any) bool) {
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
func (tree *AVLTree) IteratorAscFrom(key any, match bool, f func(key, value any) bool) {
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
func (tree *AVLTree) IteratorDesc(f func(key, value any) bool) {
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
func (tree *AVLTree) IteratorDescFrom(key any, match bool, f func(key, value any) bool) {
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

// Left returns the minimum element corresponding to the comparator of the tree or nil if the tree is empty.
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

// Right returns the maximum element corresponding to the comparator of the tree or nil if the tree is empty.
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

// Floor Finds floor node of the input key, returns the floor node or nil if no floor node is found.
// The second returned parameter `found` is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree is larger than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *AVLTree) Floor(key any) (floor *AVLTreeNode, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, ok := tree.tree.Floor(key)
	if !ok {
		return nil, false
	}
	return &AVLTreeNode{
		Key:   node.Key,
		Value: node.Value,
	}, true
}

// Ceiling finds ceiling node of the input key, returns the ceiling node or nil if no ceiling node is found.
// The second return parameter `found` is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree is smaller than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *AVLTree) Ceiling(key any) (ceiling *AVLTreeNode, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, ok := tree.tree.Ceiling(key)
	if !ok {
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
func (tree *AVLTree) Flip(comparator ...func(v1, v2 any) int) {
	var t = new(AVLTree)
	if len(comparator) > 0 {
		t = NewAVLTree(comparator[0], tree.mu.IsSafe())
	} else {
		t = NewAVLTree(tree.comparator, tree.mu.IsSafe())
	}
	tree.IteratorAsc(func(key, value any) bool {
		t.doSet(value, key)
		return true
	})
	tree.Clear()
	tree.Sets(t.Map())
}

// doSet inserts key-value pair node into the tree without lock.
// If `key` already exists, then its value is updated with the new value.
// If `value` is type of <func() any>, it will be executed and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (tree *AVLTree) doSet(key, value any) any {
	if f, ok := value.(func() any); ok {
		value = f()
	}
	if value == nil {
		return value
	}
	tree.tree.Put(key, value)
	return value
}

// doGet retrieves and returns the value of given key from tree without lock.
func (tree *AVLTree) doGet(key any) (value any, found bool) {
	return tree.tree.Get(key)
}

// doRemove removes key from tree and returns its associated value without lock.
// Note that, the given `key` should adhere to the comparator's type assertion, otherwise method panics.
func (tree *AVLTree) doRemove(key any) (value any) {
	value, _ = tree.tree.Get(key)
	tree.tree.Remove(key)
	return
}

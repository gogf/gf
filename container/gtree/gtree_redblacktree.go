// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree

import (
	"fmt"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

var _ iTree = (*RedBlackTree)(nil)

// RedBlackTree holds elements of the red-black tree.
type RedBlackTree struct {
	mu         rwmutex.RWMutex
	comparator func(v1, v2 interface{}) int
	tree       *redblacktree.Tree
}

// RedBlackTreeNode is a single element within the tree.
type RedBlackTreeNode struct {
	Key   interface{}
	Value interface{}
}

// NewRedBlackTree instantiates a red-black tree with the custom key comparator.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewRedBlackTree(comparator func(v1, v2 interface{}) int, safe ...bool) *RedBlackTree {
	return &RedBlackTree{
		mu:         rwmutex.Create(safe...),
		comparator: comparator,
		tree:       redblacktree.NewWith(comparator),
	}
}

// NewRedBlackTreeFrom instantiates a red-black tree with the custom key comparator and `data` map.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewRedBlackTreeFrom(comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}, safe ...bool) *RedBlackTree {
	tree := NewRedBlackTree(comparator, safe...)
	for k, v := range data {
		tree.doSet(k, v)
	}
	return tree
}

// SetComparator sets/changes the comparator for sorting.
func (tree *RedBlackTree) SetComparator(comparator func(a, b interface{}) int) {
	tree.comparator = comparator
	if tree.tree == nil {
		tree.tree = redblacktree.NewWith(comparator)
	}
	size := tree.tree.Size()
	if size > 0 {
		m := tree.Map()
		tree.Sets(m)
	}
}

// Clone returns a new tree with a copy of current tree.
func (tree *RedBlackTree) Clone() *RedBlackTree {
	newTree := NewRedBlackTree(tree.comparator, tree.mu.IsSafe())
	newTree.Sets(tree.Map())
	return newTree
}

// Set inserts node into the tree.
func (tree *RedBlackTree) Set(key interface{}, value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.doSet(key, value)
}

// Sets batch sets key-values to the tree.
func (tree *RedBlackTree) Sets(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for key, value := range data {
		tree.doSet(key, value)
	}
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (tree *RedBlackTree) SetIfNotExist(key interface{}, value interface{}) bool {
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
func (tree *RedBlackTree) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
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
func (tree *RedBlackTree) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, ok := tree.doGet(key); !ok {
		tree.doSet(key, f)
		return true
	}
	return false
}

// Get searches the node in the tree by `key` and returns its value or nil if key is not found in tree.
func (tree *RedBlackTree) Get(key interface{}) (value interface{}) {
	value, _ = tree.Search(key)
	return
}

// GetOrSet returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (tree *RedBlackTree) GetOrSet(key interface{}, value interface{}) interface{} {
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
func (tree *RedBlackTree) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
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
func (tree *RedBlackTree) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
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
func (tree *RedBlackTree) GetVar(key interface{}) *gvar.Var {
	return gvar.New(tree.Get(key))
}

// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
// The returned gvar.Var is un-concurrent safe.
func (tree *RedBlackTree) GetVarOrSet(key interface{}, value interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
// The returned gvar.Var is un-concurrent safe.
func (tree *RedBlackTree) GetVarOrSetFunc(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
// The returned gvar.Var is un-concurrent safe.
func (tree *RedBlackTree) GetVarOrSetFuncLock(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSetFuncLock(key, f))
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *RedBlackTree) Search(key interface{}) (value interface{}, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if node, found := tree.doGet(key); found {
		return node, true
	}
	return nil, false
}

// Contains checks whether `key` exists in the tree.
func (tree *RedBlackTree) Contains(key interface{}) bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	_, ok := tree.doGet(key)
	return ok
}

// Size returns number of nodes in the tree.
func (tree *RedBlackTree) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Size()
}

// IsEmpty returns true if tree does not contain any nodes.
func (tree *RedBlackTree) IsEmpty() bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Size() == 0
}

// Remove removes the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RedBlackTree) Remove(key interface{}) (value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

// Removes batch deletes values of the tree by `keys`.
func (tree *RedBlackTree) Removes(keys []interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.doRemove(key)
	}
}

// Clear removes all nodes from the tree.
func (tree *RedBlackTree) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.tree.Clear()
}

// Keys returns all keys in asc order.
func (tree *RedBlackTree) Keys() []interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Keys()
}

// Values returns all values in asc order based on the key.
func (tree *RedBlackTree) Values() []interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Values()
}

// Replace the data of the tree with given `data`.
func (tree *RedBlackTree) Replace(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.tree.Clear()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// Print prints the tree to stdout.
func (tree *RedBlackTree) Print() {
	fmt.Println(tree.String())
}

// String returns a string representation of container
func (tree *RedBlackTree) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return gstr.Replace(tree.tree.String(), "RedBlackTree\n", "")
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree *RedBlackTree) MarshalJSON() (jsonBytes []byte, err error) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.MarshalJSON()
}

// Map returns all key-value items as map.
func (tree *RedBlackTree) Map() map[interface{}]interface{} {
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
func (tree *RedBlackTree) MapStrAny() map[string]interface{} {
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
func (tree *RedBlackTree) Iterator(f func(key, value interface{}) bool) {
	tree.IteratorAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
func (tree *RedBlackTree) IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorAsc(f func(key, value interface{}) bool) {
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
func (tree *RedBlackTree) IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
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
func (tree *RedBlackTree) IteratorDesc(f func(key, value interface{}) bool) {
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
func (tree *RedBlackTree) IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
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
func (tree *RedBlackTree) Left() *RedBlackTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.tree.Left()
	if node == nil {
		return nil
	}
	return &RedBlackTreeNode{
		Key:   node.Key,
		Value: node.Value,
	}
}

// Right returns the maximum element of the AVL tree
// or nil if the tree is empty.
func (tree *RedBlackTree) Right() *RedBlackTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.tree.Right()
	if node == nil {
		return nil
	}
	return &RedBlackTreeNode{
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
func (tree *RedBlackTree) Floor(key interface{}) (floor *RedBlackTreeNode, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.tree.Floor(key)
	if !found {
		return nil, false
	}
	return &RedBlackTreeNode{
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
func (tree *RedBlackTree) Ceiling(key interface{}) (ceiling *RedBlackTreeNode, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.tree.Ceiling(key)
	if !found {
		return nil, false
	}
	return &RedBlackTreeNode{
		Key:   node.Key,
		Value: node.Value,
	}, true
}

// Flip exchanges key-value of the tree to value-key.
// Note that you should guarantee the value is the same type as key,
// or else the comparator would panic.
//
// If the type of value is different with key, you pass the new `comparator`.
func (tree *RedBlackTree) Flip(comparator ...func(v1, v2 interface{}) int) {
	var t = new(RedBlackTree)
	if len(comparator) > 0 {
		t = NewRedBlackTree(comparator[0], tree.mu.IsSafe())
	} else {
		t = NewRedBlackTree(tree.comparator, tree.mu.IsSafe())
	}
	tree.IteratorAsc(func(key, value interface{}) bool {
		t.doSet(value, key)
		return true
	})
	tree.Clear()
	tree.Sets(t.Map())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (tree *RedBlackTree) UnmarshalJSON(b []byte) error {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.comparator == nil {
		tree.comparator = gutil.ComparatorString
		tree.tree = redblacktree.NewWith(tree.comparator)
	}
	var data map[string]interface{}
	if err := json.UnmarshalUseNumber(b, &data); err != nil {
		return err
	}
	for k, v := range data {
		tree.doSet(k, v)
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (tree *RedBlackTree) UnmarshalValue(value interface{}) (err error) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.comparator == nil {
		tree.comparator = gutil.ComparatorString
		tree.tree = redblacktree.NewWith(tree.comparator)
	}
	for k, v := range gconv.Map(value) {
		tree.doSet(k, v)
	}
	return
}

// doSet sets key-value pair to the tree.
func (tree *RedBlackTree) doSet(key, value interface{}) interface{} {
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
func (tree *RedBlackTree) doGet(key interface{}) (value interface{}, found bool) {
	return tree.tree.Get(key)
}

// doRemove removes key from tree.
func (tree *RedBlackTree) doRemove(key interface{}) (value interface{}) {
	value, _ = tree.tree.Get(key)
	tree.tree.Remove(key)
	return
}

// iteratorFromGetIndex returns the index of the key in the keys slice.
// The parameter `match` specifies whether starting iterating if the `key` is fully matched,
// or else using index searching iterating.
// If `isIterator` is true, iterator is available; or else not.
func (tree *RedBlackTree) iteratorFromGetIndex(key interface{}, keys []interface{}, match bool) (index int, isIterator bool) {
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

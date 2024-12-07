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
	comparator func(v1, v2 any) int
	tree       *redblacktree.Tree
}

// RedBlackTreeNode is a single element within the tree.
type RedBlackTreeNode struct {
	Key   any
	Value any
}

// NewRedBlackTree instantiates a red-black tree with the custom key comparator.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewRedBlackTree(comparator func(v1, v2 any) int, safe ...bool) *RedBlackTree {
	return &RedBlackTree{
		mu:         rwmutex.Create(safe...),
		comparator: comparator,
		tree:       redblacktree.NewWith(comparator),
	}
}

// NewRedBlackTreeFrom instantiates a red-black tree with the custom key comparator and `data` map.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewRedBlackTreeFrom(comparator func(v1, v2 any) int, data map[any]any, safe ...bool) *RedBlackTree {
	tree := NewRedBlackTree(comparator, safe...)
	for k, v := range data {
		tree.doSet(k, v)
	}
	return tree
}

// SetComparator sets/changes the comparator for sorting.
func (tree *RedBlackTree) SetComparator(comparator func(a, b any) int) {
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

// Clone clones and returns a new tree from current tree.
func (tree *RedBlackTree) Clone() *RedBlackTree {
	newTree := NewRedBlackTree(tree.comparator, tree.mu.IsSafe())
	newTree.Sets(tree.Map())
	return newTree
}

// Set sets key-value pair into the tree.
func (tree *RedBlackTree) Set(key any, value any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.doSet(key, value)
}

// Sets batch sets key-values to the tree.
func (tree *RedBlackTree) Sets(data map[any]any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for key, value := range data {
		tree.doSet(key, value)
	}
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and such setting key-value pair operation would be ignored.
func (tree *RedBlackTree) SetIfNotExist(key any, value any) bool {
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
func (tree *RedBlackTree) SetIfNotExistFunc(key any, f func() any) bool {
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
func (tree *RedBlackTree) SetIfNotExistFuncLock(key any, f func() any) bool {
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
func (tree *RedBlackTree) Get(key any) (value any) {
	value, _ = tree.Search(key)
	return
}

// GetOrSet returns its `value` of `key`, or sets value with given `value` if it does not exist and then returns
// this value.
func (tree *RedBlackTree) GetOrSet(key any, value any) any {
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
func (tree *RedBlackTree) GetOrSetFunc(key any, f func() any) any {
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
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`within mutex lock.
func (tree *RedBlackTree) GetOrSetFuncLock(key any, f func() any) any {
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
func (tree *RedBlackTree) GetVar(key any) *gvar.Var {
	return gvar.New(tree.Get(key))
}

// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSet.
func (tree *RedBlackTree) GetVarOrSet(key any, value any) *gvar.Var {
	return gvar.New(tree.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFunc.
func (tree *RedBlackTree) GetVarOrSetFunc(key any, f func() any) *gvar.Var {
	return gvar.New(tree.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
// Note that, the returned gvar.Var is un-concurrent safe.
//
// Also see function GetOrSetFuncLock.
func (tree *RedBlackTree) GetVarOrSetFuncLock(key any, f func() any) *gvar.Var {
	return gvar.New(tree.GetOrSetFuncLock(key, f))
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *RedBlackTree) Search(key any) (value any, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if node, found := tree.doGet(key); found {
		return node, true
	}
	return nil, false
}

// Contains checks and returns whether given `key` exists in the tree.
func (tree *RedBlackTree) Contains(key any) bool {
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

// Remove removes the node from the tree by `key`, and returns its associated value of `key`.
// The given `key` should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RedBlackTree) Remove(key any) (value any) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

// Removes batch deletes key-value pairs from the tree by `keys`.
func (tree *RedBlackTree) Removes(keys []any) {
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

// Keys returns all keys from the tree in order by its comparator.
func (tree *RedBlackTree) Keys() []any {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Keys()
}

// Values returns all values from the true in order by its comparator based on the key.
func (tree *RedBlackTree) Values() []any {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.tree.Values()
}

// Replace clears the data of the tree and sets the nodes by given `data`.
func (tree *RedBlackTree) Replace(data map[any]any) {
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

// Map returns all key-value pairs as map.
func (tree *RedBlackTree) Map() map[any]any {
	m := make(map[any]any, tree.Size())
	tree.IteratorAsc(func(key, value any) bool {
		m[key] = value
		return true
	})
	return m
}

// MapStrAny returns all key-value items as map[string]any.
func (tree *RedBlackTree) MapStrAny() map[string]any {
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
func (tree *RedBlackTree) Iterator(f func(key, value any) bool) {
	tree.IteratorAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
//
// Also see IteratorAscFrom.
func (tree *RedBlackTree) IteratorFrom(key any, match bool, f func(key, value any) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
// If callback function `f` returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorAsc(f func(key, value any) bool) {
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
func (tree *RedBlackTree) IteratorAscFrom(key any, match bool, f func(key, value any) bool) {
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
func (tree *RedBlackTree) IteratorDesc(f func(key, value any) bool) {
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
func (tree *RedBlackTree) IteratorDescFrom(key any, match bool, f func(key, value any) bool) {
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

// Right returns the maximum element corresponding to the comparator of the tree or nil if the tree is empty.
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

// Floor Finds floor node of the input key, returns the floor node or nil if no floor node is found.
// The second returned parameter `found` is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree is larger than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RedBlackTree) Floor(key any) (floor *RedBlackTreeNode, found bool) {
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

// Ceiling finds ceiling node of the input key, returns the ceiling node or nil if no ceiling node is found.
// The second return parameter `found` is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree is smaller than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RedBlackTree) Ceiling(key any) (ceiling *RedBlackTreeNode, found bool) {
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
func (tree *RedBlackTree) Flip(comparator ...func(v1, v2 any) int) {
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
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.comparator == nil {
		tree.comparator = gutil.ComparatorString
		tree.tree = redblacktree.NewWith(tree.comparator)
	}
	var data map[string]any
	if err := json.UnmarshalUseNumber(b, &data); err != nil {
		return err
	}
	for k, v := range data {
		tree.doSet(k, v)
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (tree *RedBlackTree) UnmarshalValue(value any) (err error) {
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

// doSet inserts key-value pair node into the tree without lock.
// If `key` already exists, then its value is updated with the new value.
// If `value` is type of <func() any>, it will be executed and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (tree *RedBlackTree) doSet(key, value any) any {
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
func (tree *RedBlackTree) doGet(key any) (value any, found bool) {
	return tree.tree.Get(key)
}

// doRemove removes key from tree and returns its associated value without lock.
// Note that, the given `key` should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RedBlackTree) doRemove(key any) (value any) {
	value, _ = tree.tree.Get(key)
	tree.tree.Remove(key)
	return
}

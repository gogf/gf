// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree

import (
	"bytes"
	json2 "encoding/json"
	"fmt"

	"github.com/gogf/gf/contrib/generic_container/v2/comparator"
	"github.com/gogf/gf/contrib/generic_container/v2/conv"
	"github.com/gogf/gf/contrib/generic_container/v2/gmap"
	"github.com/gogf/gf/contrib/generic_container/v2/internal/json"
	"github.com/gogf/gf/contrib/generic_container/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/util/gconv"
)

type color bool

const (
	black, red color = true, false
)

// RedBlackTree holds elements of the red-black tree.
type RedBlackTree[K comparable, V comparable] struct {
	mu         rwmutex.RWMutex
	root       *RedBlackTreeNode[K, V]
	size       int
	comparator func(v1, v2 K) int
}

// RedBlackTreeNode is a single element within the tree.
type RedBlackTreeNode[K comparable, V comparable] struct {
	Key    K
	Value  V
	color  color
	left   *RedBlackTreeNode[K, V]
	right  *RedBlackTreeNode[K, V]
	parent *RedBlackTreeNode[K, V]
}

// NewRedBlackTree instantiates a red-black tree with the custom key comparator.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewRedBlackTree[K comparable, V comparable](comparator func(v1, v2 K) int, safe ...bool) *RedBlackTree[K, V] {
	return &RedBlackTree[K, V]{
		mu:         rwmutex.Create(safe...),
		comparator: comparator,
	}
}

// NewRedBlackTreeFrom instantiates a red-black tree with the custom key comparator and `data` map.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewRedBlackTreeFrom[K comparable, V comparable](comparator func(v1, v2 K) int, data map[K]V, safe ...bool) *RedBlackTree[K, V] {
	tree := NewRedBlackTree[K, V](comparator, safe...)
	for k, v := range data {
		tree.doSet(k, v)
	}
	return tree
}

// SetComparator sets/changes the comparator for sorting.
func (tree *RedBlackTree[K, V]) SetComparator(comparator func(a, b K) int) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.comparator = comparator
	if tree.size > 0 {
		data := make(map[K]V, tree.size)
		tree.doIteratorAsc(tree.leftNode(), func(key K, value V) bool {
			data[key] = value
			return true
		})
		// Resort the tree if comparator is changed.
		tree.root = nil
		tree.size = 0
		for k, v := range data {
			tree.doSet(k, v)
		}
	}
}

// Clone returns a new tree with a copy of current tree.
func (tree *RedBlackTree[K, V]) Clone(safe ...bool) gmap.Map[K, V] {
	newTree := NewRedBlackTree[K, V](tree.comparator, safe...)
	newTree.Sets(tree.Map())
	return newTree
}

// Set inserts key-value item into the tree.
func (tree *RedBlackTree[K, V]) Set(key K, value V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.doSet(key, value)
}

// Sets batch sets key-values to the tree.
func (tree *RedBlackTree[K, V]) Sets(data map[K]V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// doSet inserts key-value item into the tree without mutex.
func (tree *RedBlackTree[K, V]) doSet(key K, value V) {
	insertedNode := (*RedBlackTreeNode[K, V])(nil)
	if tree.root == nil {
		// Assert key is of comparator's type for initial tree
		tree.getComparator()(key, key)
		tree.root = &RedBlackTreeNode[K, V]{Key: key, Value: value, color: red}
		insertedNode = tree.root
	} else {
		node := tree.root
		loop := true
		for loop {
			compare := tree.getComparator()(key, node.Key)
			switch {
			case compare == 0:
				// node.Key   = key
				node.Value = value
				return
			case compare < 0:
				if node.left == nil {
					node.left = &RedBlackTreeNode[K, V]{Key: key, Value: value, color: red}
					insertedNode = node.left
					loop = false
				} else {
					node = node.left
				}
			case compare > 0:
				if node.right == nil {
					node.right = &RedBlackTreeNode[K, V]{Key: key, Value: value, color: red}
					insertedNode = node.right
					loop = false
				} else {
					node = node.right
				}
			}
		}
		insertedNode.parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++
}

// Get searches the node in the tree by `key` and returns its value or nil if key is not found in tree.
func (tree *RedBlackTree[K, V]) Get(key K) (value V) {
	value, _ = tree.Search(key)
	return
}

// doSetWithLockCheck checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given `key`,
// or else just return the existing value.
//
// When setting value, if `value` is type of <func() interface {}>,
// it will be executed with mutex.Lock of the hash map,
// and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (tree *RedBlackTree[K, V]) doSetWithLockCheck(key K, value V) V {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if node, found := tree.doSearch(key); found {
		return node.Value
	}
	if any(value) != nil {
		tree.doSet(key, value)
	}
	return value
}

// doSetWithLockCheckFunc checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given `key`,
// or else just return the existing value.
//
// When setting value, if `value` is type of <func() interface {}>,
// it will be executed with mutex.Lock of the hash map,
// and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (tree *RedBlackTree[K, V]) doSetWithLockCheckFunc(key K, f func() V) V {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if node, found := tree.doSearch(key); found {
		return node.Value
	}
	value := f()
	if any(value) != nil {
		tree.doSet(key, value)
	}
	return value
}

// GetOrSet returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (tree *RedBlackTree[K, V]) GetOrSet(key K, value V) V {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
func (tree *RedBlackTree[K, V]) GetOrSetFunc(key K, f func() V) V {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f())
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
func (tree *RedBlackTree[K, V]) GetOrSetFuncLock(key K, f func() V) V {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheckFunc(key, f)
	} else {
		return v
	}
}

// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (tree *RedBlackTree[K, V]) SetIfNotExist(key K, value V) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (tree *RedBlackTree[K, V]) SetIfNotExistFunc(key K, f func() V) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function `f` with mutex.Lock of the hash map.
func (tree *RedBlackTree[K, V]) SetIfNotExistFuncLock(key K, f func() V) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheckFunc(key, f)
		return true
	}
	return false
}

// Contains checks whether `key` exists in the tree.
func (tree *RedBlackTree[K, V]) Contains(key K) bool {
	_, ok := tree.Search(key)
	return ok
}

// doRemove removes the node from the tree by `key` without mutex.
func (tree *RedBlackTree[K, V]) doRemove(key K) (value V) {
	child := (*RedBlackTreeNode[K, V])(nil)
	node, found := tree.doSearch(key)
	if !found {
		return
	}
	value = node.Value
	if node.left != nil && node.right != nil {
		p := node.left.maximumNode()
		node.Key = p.Key
		node.Value = p.Value
		node = p
	}
	if node.left == nil || node.right == nil {
		if node.right == nil {
			child = node.left
		} else {
			child = node.right
		}
		if node.color == black {
			node.color = tree.nodeColor(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.parent == nil && child != nil {
			child.color = black
		}
	}
	tree.size--
	return
}

// Remove removes the node from the tree by `key`.
func (tree *RedBlackTree[K, V]) Remove(key K) (value V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

// Removes batch deletes values of the tree by `keys`.
func (tree *RedBlackTree[K, V]) Removes(keys []K) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.doRemove(key)
	}
}

// IsEmpty returns true if tree does not contain any nodes.
func (tree *RedBlackTree[K, V]) IsEmpty() bool {
	return tree.Size() == 0
}

// Size returns number of nodes in the tree.
func (tree *RedBlackTree[K, V]) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.size
}

// Keys returns all keys in asc order.
func (tree *RedBlackTree[K, V]) Keys() []K {
	var (
		keys  = make([]K, tree.Size())
		index = 0
	)
	tree.IteratorAsc(func(key K, value V) bool {
		keys[index] = key
		index++
		return true
	})
	return keys
}

// Values returns all values in asc order based on the key.
func (tree *RedBlackTree[K, V]) Values() []V {
	var (
		values = make([]V, tree.Size())
		index  = 0
	)
	tree.IteratorAsc(func(key K, value V) bool {
		values[index] = value
		index++
		return true
	})
	return values
}

// Map returns all key-value items as map.
func (tree *RedBlackTree[K, V]) Map() map[K]V {
	m := make(map[K]V, tree.Size())
	tree.IteratorAsc(func(key K, value V) bool {
		m[key] = value
		return true
	})
	return m
}

// MapStrAny returns all key-value items as map[string]V.
func (tree *RedBlackTree[K, V]) MapStrAny() map[string]V {
	m := make(map[string]V, tree.Size())
	tree.IteratorAsc(func(key K, value V) bool {
		m[gconv.String(key)] = value
		return true
	})
	return m
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *RedBlackTree[K, V]) Left() *RedBlackTreeNode[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.leftNode()
	if tree.mu.IsSafe() {
		return &RedBlackTreeNode[K, V]{
			Key:   node.Key,
			Value: node.Value,
		}
	}
	return node
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *RedBlackTree[K, V]) Right() *RedBlackTreeNode[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.rightNode()
	if tree.mu.IsSafe() {
		return &RedBlackTreeNode[K, V]{
			Key:   node.Key,
			Value: node.Value,
		}
	}
	return node
}

// leftNode returns the left-most (min) node or nil if tree is empty.
func (tree *RedBlackTree[K, V]) leftNode() *RedBlackTreeNode[K, V] {
	p := (*RedBlackTreeNode[K, V])(nil)
	n := tree.root
	for n != nil {
		p = n
		n = n.left
	}
	return p
}

// rightNode returns the right-most (max) node or nil if tree is empty.
func (tree *RedBlackTree[K, V]) rightNode() *RedBlackTreeNode[K, V] {
	p := (*RedBlackTreeNode[K, V])(nil)
	n := tree.root
	for n != nil {
		p = n
		n = n.right
	}
	return p
}

// Floor Finds floor node of the input key, return the floor node or nil if no floor node is found.
// Second return parameter is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that its key is smaller than or equal to the given `key`.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
func (tree *RedBlackTree[K, V]) Floor(key K) (floor *RedBlackTreeNode[K, V], found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	n := tree.root
	for n != nil {
		compare := tree.getComparator()(key, n.Key)
		switch {
		case compare == 0:
			return n, true
		case compare < 0:
			n = n.left
		case compare > 0:
			floor, found = n, true
			n = n.right
		}
	}
	if found {
		return
	}
	return nil, false
}

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling node is found.
// Second return parameter is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that its key is larger than or equal to the given `key`.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
func (tree *RedBlackTree[K, V]) Ceiling(key K) (ceiling *RedBlackTreeNode[K, V], found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	n := tree.root
	for n != nil {
		compare := tree.getComparator()(key, n.Key)
		switch {
		case compare == 0:
			return n, true
		case compare > 0:
			n = n.right
		case compare < 0:
			ceiling, found = n, true
			n = n.left
		}
	}
	if found {
		return
	}
	return nil, false
}

// Iterator is alias of IteratorAsc.
func (tree *RedBlackTree[K, V]) Iterator(f func(key K, value V) bool) {
	tree.IteratorAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
func (tree *RedBlackTree[K, V]) IteratorFrom(key K, match bool, f func(key K, value V) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree[K, V]) IteratorAsc(f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorAsc(tree.leftNode(), f)
}

// IteratorAscFrom iterates the tree readonly in ascending order with given callback function `f`.
// The parameter `key` specifies the start entry for iterating. The `match` specifies whether
// starting iterating if the `key` is fully matched, or else using index searching iterating.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree[K, V]) IteratorAscFrom(key K, match bool, f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.doSearch(key)
	if match {
		if found {
			tree.doIteratorAsc(node, f)
		}
	} else {
		tree.doIteratorAsc(node, f)
	}
}

func (tree *RedBlackTree[K, V]) doIteratorAsc(node *RedBlackTreeNode[K, V], f func(key K, value V) bool) {
loop:
	if node == nil {
		return
	}
	if !f(node.Key, node.Value) {
		return
	}
	if node.right != nil {
		node = node.right
		for node.left != nil {
			node = node.left
		}
		goto loop
	}
	if node.parent != nil {
		old := node
		for node.parent != nil {
			node = node.parent
			if tree.getComparator()(old.Key, node.Key) <= 0 {
				goto loop
			}
		}
	}
}

// IteratorDesc iterates the tree readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree[K, V]) IteratorDesc(f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorDesc(tree.rightNode(), f)
}

// IteratorDescFrom iterates the tree readonly in descending order with given callback function `f`.
// The parameter `key` specifies the start entry for iterating. The `match` specifies whether
// starting iterating if the `key` is fully matched, or else using index searching iterating.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree[K, V]) IteratorDescFrom(key K, match bool, f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.doSearch(key)
	if match {
		if found {
			tree.doIteratorDesc(node, f)
		}
	} else {
		tree.doIteratorDesc(node, f)
	}
}

func (tree *RedBlackTree[K, V]) doIteratorDesc(node *RedBlackTreeNode[K, V], f func(key K, value V) bool) {
loop:
	if node == nil {
		return
	}
	if !f(node.Key, node.Value) {
		return
	}
	if node.left != nil {
		node = node.left
		for node.right != nil {
			node = node.right
		}
		goto loop
	}
	if node.parent != nil {
		old := node
		for node.parent != nil {
			node = node.parent
			if tree.getComparator()(old.Key, node.Key) >= 0 {
				goto loop
			}
		}
	}
}

// Clear removes all nodes from the tree.
func (tree *RedBlackTree[K, V]) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
}

// Replace the data of the tree with given `data`.
func (tree *RedBlackTree[K, V]) Replace(data map[K]V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// String returns a string representation of container.
func (tree *RedBlackTree[K, V]) String() string {
	if tree == nil {
		return ""
	}
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	str := ""
	if tree.size != 0 {
		tree.output(tree.root, "", true, &str)
	}
	return str
}

// Print prints the tree to stdout.
func (tree *RedBlackTree[K, V]) Print() {
	fmt.Println(tree.String())
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *RedBlackTree[K, V]) Search(key K) (value V, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.doSearch(key)
	if found {
		return node.Value, true
	}
	return
}

// Flip exchanges key-value of the tree to value-key.
// Note that you should guarantee the value is the same type as key,
// or else the comparator would panic.
//
// If the type of value is different with key, you pass the new `comparator`.
func (tree *RedBlackTree[K, V]) Flip(comparator func(v1, v2 V) int) *RedBlackTree[V, K] {
	t := (*RedBlackTree[V, K])(nil)
	t = NewRedBlackTree[V, K](comparator, tree.mu.IsSafe())
	tree.IteratorAsc(func(key K, value V) bool {
		t.doSet(value, key)
		return true
	})
	return t
}

func (tree *RedBlackTree[K, V]) output(node *RedBlackTreeNode[K, V], prefix string, isTail bool, str *string) {
	if node.right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		tree.output(node.right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += fmt.Sprintf("%v\n", node.Key)
	if node.left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		tree.output(node.left, newPrefix, true, str)
	}
}

// doSearch searches the tree with given `key` without mutex.
// It returns the node if found or otherwise nil.
func (tree *RedBlackTree[K, V]) doSearch(key K) (node *RedBlackTreeNode[K, V], found bool) {
	node = tree.root
	for node != nil {
		compare := tree.getComparator()(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			node = node.left
		case compare > 0:
			node = node.right
		}
	}
	return node, false
}

func (node *RedBlackTreeNode[K, V]) grandparent() *RedBlackTreeNode[K, V] {
	if node != nil && node.parent != nil {
		return node.parent.parent
	}
	return nil
}

func (node *RedBlackTreeNode[K, V]) uncle() *RedBlackTreeNode[K, V] {
	if node == nil || node.parent == nil || node.parent.parent == nil {
		return nil
	}
	return node.parent.sibling()
}

func (node *RedBlackTreeNode[K, V]) sibling() *RedBlackTreeNode[K, V] {
	if node == nil || node.parent == nil {
		return nil
	}
	if node == node.parent.left {
		return node.parent.right
	}
	return node.parent.left
}

func (tree *RedBlackTree[K, V]) rotateLeft(node *RedBlackTreeNode[K, V]) {
	right := node.right
	tree.replaceNode(node, right)
	node.right = right.left
	if right.left != nil {
		right.left.parent = node
	}
	right.left = node
	node.parent = right
}

func (tree *RedBlackTree[K, V]) rotateRight(node *RedBlackTreeNode[K, V]) {
	left := node.left
	tree.replaceNode(node, left)
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}
	left.right = node
	node.parent = left
}

func (tree *RedBlackTree[K, V]) replaceNode(old *RedBlackTreeNode[K, V], new *RedBlackTreeNode[K, V]) {
	if old.parent == nil {
		tree.root = new
	} else {
		if old == old.parent.left {
			old.parent.left = new
		} else {
			old.parent.right = new
		}
	}
	if new != nil {
		new.parent = old.parent
	}
}

func (tree *RedBlackTree[K, V]) insertCase1(node *RedBlackTreeNode[K, V]) {
	if node.parent == nil {
		node.color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *RedBlackTree[K, V]) insertCase2(node *RedBlackTreeNode[K, V]) {
	if tree.nodeColor(node.parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *RedBlackTree[K, V]) insertCase3(node *RedBlackTreeNode[K, V]) {
	uncle := node.uncle()
	if tree.nodeColor(uncle) == red {
		node.parent.color = black
		uncle.color = black
		node.grandparent().color = red
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *RedBlackTree[K, V]) insertCase4(node *RedBlackTreeNode[K, V]) {
	grandparent := node.grandparent()
	if node == node.parent.right && node.parent == grandparent.left {
		tree.rotateLeft(node.parent)
		node = node.left
	} else if node == node.parent.left && node.parent == grandparent.right {
		tree.rotateRight(node.parent)
		node = node.right
	}
	tree.insertCase5(node)
}

func (tree *RedBlackTree[K, V]) insertCase5(node *RedBlackTreeNode[K, V]) {
	node.parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == node.parent.left && node.parent == grandparent.left {
		tree.rotateRight(grandparent)
	} else if node == node.parent.right && node.parent == grandparent.right {
		tree.rotateLeft(grandparent)
	}
}

func (node *RedBlackTreeNode[K, V]) maximumNode() *RedBlackTreeNode[K, V] {
	if node == nil {
		return nil
	}
	for node.right != nil {
		return node.right
	}
	return node
}

func (tree *RedBlackTree[K, V]) deleteCase1(node *RedBlackTreeNode[K, V]) {
	if node.parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *RedBlackTree[K, V]) deleteCase2(node *RedBlackTreeNode[K, V]) {
	sibling := node.sibling()
	if tree.nodeColor(sibling) == red {
		node.parent.color = red
		sibling.color = black
		if node == node.parent.left {
			tree.rotateLeft(node.parent)
		} else {
			tree.rotateRight(node.parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *RedBlackTree[K, V]) deleteCase3(node *RedBlackTreeNode[K, V]) {
	sibling := node.sibling()
	if tree.nodeColor(node.parent) == black &&
		tree.nodeColor(sibling) == black &&
		tree.nodeColor(sibling.left) == black &&
		tree.nodeColor(sibling.right) == black {
		sibling.color = red
		tree.deleteCase1(node.parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *RedBlackTree[K, V]) deleteCase4(node *RedBlackTreeNode[K, V]) {
	sibling := node.sibling()
	if tree.nodeColor(node.parent) == red &&
		tree.nodeColor(sibling) == black &&
		tree.nodeColor(sibling.left) == black &&
		tree.nodeColor(sibling.right) == black {
		sibling.color = red
		node.parent.color = black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *RedBlackTree[K, V]) deleteCase5(node *RedBlackTreeNode[K, V]) {
	sibling := node.sibling()
	if node == node.parent.left &&
		tree.nodeColor(sibling) == black &&
		tree.nodeColor(sibling.left) == red &&
		tree.nodeColor(sibling.right) == black {
		sibling.color = red
		sibling.left.color = black
		tree.rotateRight(sibling)
	} else if node == node.parent.right &&
		tree.nodeColor(sibling) == black &&
		tree.nodeColor(sibling.right) == red &&
		tree.nodeColor(sibling.left) == black {
		sibling.color = red
		sibling.right.color = black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *RedBlackTree[K, V]) deleteCase6(node *RedBlackTreeNode[K, V]) {
	sibling := node.sibling()
	sibling.color = tree.nodeColor(node.parent)
	node.parent.color = black
	if node == node.parent.left && tree.nodeColor(sibling.right) == red {
		sibling.right.color = black
		tree.rotateLeft(node.parent)
	} else if tree.nodeColor(sibling.left) == red {
		sibling.left.color = black
		tree.rotateRight(node.parent)
	}
}

func (tree *RedBlackTree[K, V]) nodeColor(node *RedBlackTreeNode[K, V]) color {
	if node == nil {
		return black
	}
	return node.color
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree RedBlackTree[K, V]) MarshalJSON() (jsonBytes []byte, err error) {
	if tree.root == nil {
		return []byte("null"), nil
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('{')
	tree.Iterator(func(key K, value V) bool {
		valueBytes, valueJsonErr := json.Marshal(value)
		if valueJsonErr != nil {
			err = valueJsonErr
			return false
		}
		if buffer.Len() > 1 {
			buffer.WriteByte(',')
		}
		buffer.WriteString(fmt.Sprintf(`"%v":%s`, key, valueBytes))
		return true
	})
	buffer.WriteByte('}')
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (tree *RedBlackTree[K, V]) UnmarshalJSON(b []byte) error {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.comparator == nil {
		tree.comparator = comparator.ComparatorAny[K]
	}
	var data map[K]V
	if err := json.UnmarshalUseNumber(b, &data); err != nil {
		return err
	}
	for k, v := range data {
		tree.doSet(k, v)
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (tree *RedBlackTree[K, V]) UnmarshalValue(value V) (err error) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.comparator == nil {
		tree.comparator = comparator.ComparatorAny[K]
	}
	for k, v := range gconv.Map(value) {
		kt := conv.Convert[K](k)
		var vt V
		switch v.(type) {
		case string, []byte, json2.Number:
			var ok bool
			if vt, ok = v.(V); !ok {
				if err = json.UnmarshalUseNumber(gconv.Bytes(v), &vt); err != nil {
					return err
				}
			}
		default:
			vt, _ = v.(V)
		}
		tree.doSet(kt, vt)
	}
	return
}

// getComparator returns the comparator if it's previously set,
// or else it panics.
func (tree *RedBlackTree[K, V]) getComparator() func(a, b K) int {
	if tree.comparator == nil {
		panic("comparator is missing for tree")
	}
	return tree.comparator
}

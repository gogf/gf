// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree

import (
	"fmt"
	"github.com/gogf/gf/g/container/gvar"
	"github.com/gogf/gf/g/internal/rwmutex"
)

type color bool

const (
	black, red color = true, false
)

// RedBlackTree holds elements of the red-black tree.
type RedBlackTree struct {
	mu         *rwmutex.RWMutex
	root       *RedBlackTreeNode
	size       int
	comparator func(v1, v2 interface{}) int
}

// RedBlackTreeNode is a single element within the tree.
type RedBlackTreeNode struct {
	Key    interface{}
	Value  interface{}
	color  color
	left   *RedBlackTreeNode
	right  *RedBlackTreeNode
	parent *RedBlackTreeNode
}

// NewRedBlackTree instantiates a red-black tree with the custom comparator.
// The param <unsafe> used to specify whether using tree in un-concurrent-safety,
// which is false in default.
func NewRedBlackTree(comparator func(v1, v2 interface{}) int, unsafe...bool) *RedBlackTree {
	return &RedBlackTree {
		mu        : rwmutex.New(unsafe...),
		comparator: comparator,
	}
}

// NewRedBlackTreeFrom instantiates a red-black tree with the custom comparator and <data> map.
// The param <unsafe> used to specify whether using tree in un-concurrent-safety,
// which is false in default.
func NewRedBlackTreeFrom(comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}, unsafe...bool) *RedBlackTree {
	tree := NewRedBlackTree(comparator, unsafe...)
	for k, v := range data {
		tree.doSet(k, v)
	}
	return tree
}

// Clone returns a new tree with a copy of current tree.
func (tree *RedBlackTree) Clone(unsafe ...bool) *RedBlackTree {
	newTree := NewRedBlackTree(tree.comparator, !tree.mu.IsSafe())
	newTree.Sets(tree.Map())
	return newTree
}

// Set inserts key-value item into the tree.
func (tree *RedBlackTree) Set(key interface{}, value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.doSet(key, value)
}

// Sets batch sets key-values to the tree.
func (tree *RedBlackTree) Sets(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// doSet inserts key-value item into the tree without mutex.
func (tree *RedBlackTree) doSet(key interface{}, value interface{}) {
	insertedNode := (*RedBlackTreeNode)(nil)
	if tree.root == nil {
		// Assert key is of comparator's type for initial tree
		tree.comparator(key, key)
		tree.root    = &RedBlackTreeNode{Key: key, Value: value, color: red}
		insertedNode = tree.root
	} else {
		node := tree.root
		loop := true
		for loop {
			compare := tree.comparator(key, node.Key)
			switch {
				case compare == 0:
					//node.Key   = key
					node.Value = value
					return
				case compare < 0:
					if node.left == nil {
						node.left    = &RedBlackTreeNode{Key: key, Value: value, color: red}
						insertedNode = node.left
						loop         = false
					} else {
						node = node.left
					}
				case compare > 0:
					if node.right == nil {
						node.right   = &RedBlackTreeNode{Key: key, Value: value, color: red}
						insertedNode = node.right
						loop         = false
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

// Get searches the node in the tree by <key> and returns its value or nil if key is not found in tree.
func (tree *RedBlackTree) Get(key interface{}) (value interface{}) {
	value, _ = tree.Search(key)
	return
}

// doSetWithLockCheck checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given <key>,
// or else just return the existing value.
//
// When setting value, if <value> is type of <func() interface {}>,
// it will be executed with mutex.Lock of the hash map,
// and its return value will be set to the map with <key>.
//
// It returns value with given <key>.
func (tree *RedBlackTree) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if node := tree.doSearch(key); node != nil {
		return node.Value
	}
	if f, ok := value.(func() interface {}); ok {
		value = f()
	}
	tree.doSet(key, value)
	return value
}

// GetOrSet returns the value by key,
// or set value with given <value> if not exist and returns this value.
func (tree *RedBlackTree) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with return value of callback function <f> if not exist
// and returns this value.
func (tree *RedBlackTree) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock returns the value by key,
// or sets value with return value of callback function <f> if not exist
// and returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function <f>
// with mutex.Lock of the hash map.
func (tree *RedBlackTree) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

// GetVar returns a gvar.Var with the value by given <key>.
// The returned gvar.Var is un-concurrent safe.
func (tree *RedBlackTree) GetVar(key interface{}) *gvar.Var {
	return gvar.New(tree.Get(key), true)
}

// GetVarOrSet returns a gvar.Var with result from GetVarOrSet.
// The returned gvar.Var is un-concurrent safe.
func (tree *RedBlackTree) GetVarOrSet(key interface{}, value interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSet(key, value), true)
}

// GetVarOrSetFunc returns a gvar.Var with result from GetOrSetFunc.
// The returned gvar.Var is un-concurrent safe.
func (tree *RedBlackTree) GetVarOrSetFunc(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSetFunc(key, f), true)
}

// GetVarOrSetFuncLock returns a gvar.Var with result from GetOrSetFuncLock.
// The returned gvar.Var is un-concurrent safe.
func (tree *RedBlackTree) GetVarOrSetFuncLock(key interface{}, f func() interface{}) *gvar.Var {
	return gvar.New(tree.GetOrSetFuncLock(key, f), true)
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (tree *RedBlackTree) SetIfNotExist(key interface{}, value interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, then return true.
// It returns false if <key> exists, and <value> would be ignored.
func (tree *RedBlackTree) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function <f>, then return true.
// It returns false if <key> exists, and <value> would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function <f> with mutex.Lock of the hash map.
func (tree *RedBlackTree) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

// Contains checks whether <key> exists in the tree.
func (tree *RedBlackTree) Contains(key interface{}) bool {
	_, ok := tree.Search(key)
	return ok
}

// doRemove removes the node from the tree by <key> without mutex.
func (tree *RedBlackTree) doRemove(key interface{}) (value interface{}) {
	child := (*RedBlackTreeNode)(nil)
	node  := tree.doSearch(key)
	if node == nil {
		return
	}
	value = node.Value
	if node.left != nil && node.right != nil {
		p         := node.left.maximumNode()
		node.Key   = p.Key
		node.Value = p.Value
		node       = p
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

// Remove removes the node from the tree by <key>.
func (tree *RedBlackTree) Remove(key interface{}) (value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

// Removes batch deletes values of the tree by <keys>.
func (tree *RedBlackTree) Removes(keys []interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.doRemove(key)
	}
}

// IsEmpty returns true if tree does not contain any nodes.
func (tree *RedBlackTree) IsEmpty() bool {
	return tree.Size() == 0
}

// Size returns number of nodes in the tree.
func (tree *RedBlackTree) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.size
}

// Keys returns all keys in asc order.
func (tree *RedBlackTree) Keys() []interface{} {
	keys  := make([]interface{}, tree.Size())
	index := 0
	tree.IteratorAsc(func(key, value interface{}) bool {
		keys[index] = key
		index++
		return true
	})
	return keys
}

// Values returns all values in asc order based on the key.
func (tree *RedBlackTree) Values() []interface{} {
	values := make([]interface{}, tree.Size())
	index  := 0
	tree.IteratorAsc(func(key, value interface{}) bool {
		values[index] = value
		index++
		return true
	})
	return values
}

// Map returns all key-value items as map.
func (tree *RedBlackTree) Map() map[interface{}]interface{} {
	m := make(map[interface{}]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *RedBlackTree) Left() *RedBlackTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.leftNode()
	if tree.mu.IsSafe() {
		return &RedBlackTreeNode{
			Key   : node.Key,
			Value : node.Value,
		}
	}
	return node
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *RedBlackTree) Right() *RedBlackTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.rightNode()
	if tree.mu.IsSafe() {
		return &RedBlackTreeNode{
			Key   : node.Key,
			Value : node.Value,
		}
	}
	return node
}

// leftNode returns the left-most (min) node or nil if tree is empty.
func (tree *RedBlackTree) leftNode() *RedBlackTreeNode {
	p := (*RedBlackTreeNode)(nil)
	n := tree.root
	for n != nil {
		p = n
		n = n.left
	}
	return p
}

// rightNode returns the right-most (max) node or nil if tree is empty.
func (tree *RedBlackTree) rightNode() *RedBlackTreeNode {
	p := (*RedBlackTreeNode)(nil)
	n := tree.root
	for n != nil {
		p = n
		n = n.right
	}
	return p
}

// Floor Finds floor node of the input <key>, return the floor node or nil if no floor is found.
//
// Floor node is defined as the largest node that its key is smaller than or equal to the given <key>.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
func (tree *RedBlackTree) Floor(key interface{}) (floor *RedBlackTreeNode) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	found := false
	node  := tree.root
	for node != nil {
		compare := tree.comparator(key, node.Key)
		switch {
			case compare == 0:
				return node
			case compare < 0:
				node = node.left
			case compare > 0:
				floor, found = node, true
				node         = node.right
		}
	}
	if found {
		return floor
	}
	return nil
}

// Ceiling finds ceiling node of the input <key>, return the ceiling node or nil if no ceiling is found.
//
// Ceiling node is defined as the smallest node that its key is larger than or equal to the given <key>.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
func (tree *RedBlackTree) Ceiling(key interface{}) (ceiling *RedBlackTreeNode) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	found := false
	node  := tree.root
	for node != nil {
		compare := tree.comparator(key, node.Key)
		switch {
			case compare == 0:
				return node
			case compare < 0:
				ceiling, found = node, true
				node           = node.left
			case compare > 0:
				node = node.right
		}
	}
	if found {
		return ceiling
	}
	return nil
}

// Iterator is alias of IteratorAsc.
func (tree *RedBlackTree) Iterator(f func (key, value interface{}) bool) {
	tree.IteratorAsc(f)
}

// IteratorAsc iterates the tree in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorAsc(f func (key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.leftNode()
	if node == nil {
		return
	}
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
			if tree.comparator(old.Key, node.Key) <= 0 {
				goto loop
			}
		}
	}
}

// IteratorDesc iterates the tree in descending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorDesc(f func (key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.rightNode()
	if node == nil {
		return
	}
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
			if tree.comparator(old.Key, node.Key) >= 0 {
				goto loop
			}
		}
	}
}

// Clear removes all nodes from the tree.
func (tree *RedBlackTree) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
}

// String returns a string representation of container.
func (tree *RedBlackTree) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	str := "RedBlackTree\n"
	if tree.size != 0 {
		tree.output(tree.root, "", true, &str)
	}
	return str
}

// Print prints the tree to stdout.
func (tree *RedBlackTree) Print() {
	fmt.Println(tree.String())
}

// Search searches the tree with given <key>.
// Second return parameter <found> is true if key was found, otherwise false.
func (tree *RedBlackTree) Search(key interface{}) (value interface{}, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.doSearch(key)
	if node != nil {
		return node.Value, true
	}
	return nil, false
}

// Flip exchanges key-value of the tree to value-key.
// Note that you should guarantee the value is the same type as key,
// or else the comparator would panic.
//
// If the type of value is different with key, you pass the new <comparator>.
func (tree *RedBlackTree) Flip(comparator...func(v1, v2 interface{}) int) {
	t := (*RedBlackTree)(nil)
	if len(comparator) > 0 {
		t = NewRedBlackTree(comparator[0], !tree.mu.IsSafe())
	} else {
		t = NewRedBlackTree(tree.comparator, !tree.mu.IsSafe())
	}
	tree.IteratorAsc(func(key, value interface{}) bool {
		t.doSet(value, key)
		return true
	})
	tree.mu.Lock()
	tree.root = t.root
	tree.size = t.size
	tree.mu.Unlock()
}

func (tree *RedBlackTree) output(node *RedBlackTreeNode, prefix string, isTail bool, str *string) {
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

// doSearch searches the tree with given <key> without mutex.
// It returns the node if found or otherwise nil.
func (tree *RedBlackTree) doSearch(key interface{}) *RedBlackTreeNode {
	node := tree.root
	for node != nil {
		compare := tree.comparator(key, node.Key)
		switch {
			case compare == 0: return node
			case compare  < 0: node = node.left
			case compare  > 0: node = node.right
		}
	}
	return nil
}

func (node *RedBlackTreeNode) grandparent() *RedBlackTreeNode {
	if node != nil && node.parent != nil {
		return node.parent.parent
	}
	return nil
}

func (node *RedBlackTreeNode) uncle() *RedBlackTreeNode {
	if node == nil || node.parent == nil || node.parent.parent == nil {
		return nil
	}
	return node.parent.sibling()
}

func (node *RedBlackTreeNode) sibling() *RedBlackTreeNode {
	if node == nil || node.parent == nil {
		return nil
	}
	if node == node.parent.left {
		return node.parent.right
	}
	return node.parent.left
}

func (tree *RedBlackTree) rotateLeft(node *RedBlackTreeNode) {
	right := node.right
	tree.replaceNode(node, right)
	node.right = right.left
	if right.left != nil {
		right.left.parent = node
	}
	right.left  = node
	node.parent = right
}

func (tree *RedBlackTree) rotateRight(node *RedBlackTreeNode) {
	left := node.left
	tree.replaceNode(node, left)
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}
	left.right  = node
	node.parent = left
}

func (tree *RedBlackTree) replaceNode(old *RedBlackTreeNode, new *RedBlackTreeNode) {
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

func (tree *RedBlackTree) insertCase1(node *RedBlackTreeNode) {
	if node.parent == nil {
		node.color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *RedBlackTree) insertCase2(node *RedBlackTreeNode) {
	if tree.nodeColor(node.parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *RedBlackTree) insertCase3(node *RedBlackTreeNode) {
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

func (tree *RedBlackTree) insertCase4(node *RedBlackTreeNode) {
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

func (tree *RedBlackTree) insertCase5(node *RedBlackTreeNode) {
	node.parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == node.parent.left && node.parent == grandparent.left {
		tree.rotateRight(grandparent)
	} else if node == node.parent.right && node.parent == grandparent.right {
		tree.rotateLeft(grandparent)
	}
}

func (node *RedBlackTreeNode) maximumNode() *RedBlackTreeNode {
	if node == nil {
		return nil
	}
	for node.right != nil {
		return node.right
	}
	return node
}

func (tree *RedBlackTree) deleteCase1(node *RedBlackTreeNode) {
	if node.parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *RedBlackTree) deleteCase2(node *RedBlackTreeNode) {
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

func (tree *RedBlackTree) deleteCase3(node *RedBlackTreeNode) {
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

func (tree *RedBlackTree) deleteCase4(node *RedBlackTreeNode) {
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

func (tree *RedBlackTree) deleteCase5(node *RedBlackTreeNode) {
	sibling := node.sibling()
	if node == node.parent.left &&
		tree.nodeColor(sibling) == black &&
		tree.nodeColor(sibling.left) == red &&
		tree.nodeColor(sibling.right) == black {
		sibling.color      = red
		sibling.left.color = black
		tree.rotateRight(sibling)
	} else if node == node.parent.right &&
		tree.nodeColor(sibling) == black &&
		tree.nodeColor(sibling.right) == red &&
		tree.nodeColor(sibling.left) == black {
		sibling.color       = red
		sibling.right.color = black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *RedBlackTree) deleteCase6(node *RedBlackTreeNode) {
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

func (tree *RedBlackTree) nodeColor(node *RedBlackTreeNode) color {
	if node == nil {
		return black
	}
	return node.color
}
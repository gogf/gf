// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree

import (
	"fmt"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/rwmutex"
)

type color bool

const (
	black, red color = true, false
)

// RedBlackTree holds elements of the red-black tree.
type RedBlackTree struct {
	mu         rwmutex.RWMutex
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

// NewRedBlackTree instantiates a red-black tree with the custom key comparator.
// The parameter <safe> is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewRedBlackTree(comparator func(v1, v2 interface{}) int, safe ...bool) *RedBlackTree {
	return &RedBlackTree{
		mu:         rwmutex.Create(safe...),
		comparator: comparator,
	}
}

// NewRedBlackTreeFrom instantiates a red-black tree with the custom key comparator and <data> map.
// The parameter <safe> is used to specify whether using tree in concurrent-safety,
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
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.comparator = comparator
	if tree.size > 0 {
		data := make(map[interface{}]interface{}, tree.size)
		tree.doIteratorAsc(tree.leftNode(), func(key, value interface{}) bool {
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
func (tree *RedBlackTree) Clone() *RedBlackTree {
	newTree := NewRedBlackTree(tree.comparator, tree.mu.IsSafe())
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
		tree.getComparator()(key, key)
		tree.root = &RedBlackTreeNode{Key: key, Value: value, color: red}
		insertedNode = tree.root
	} else {
		node := tree.root
		loop := true
		for loop {
			compare := tree.getComparator()(key, node.Key)
			switch {
			case compare == 0:
				//node.Key   = key
				node.Value = value
				return
			case compare < 0:
				if node.left == nil {
					node.left = &RedBlackTreeNode{Key: key, Value: value, color: red}
					insertedNode = node.left
					loop = false
				} else {
					node = node.left
				}
			case compare > 0:
				if node.right == nil {
					node.right = &RedBlackTreeNode{Key: key, Value: value, color: red}
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
	if node, found := tree.doSearch(key); found {
		return node.Value
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value != nil {
		tree.doSet(key, value)
	}
	return value
}

// GetOrSet returns the value by key,
// or sets value with given <value> if it does not exist and then returns this value.
func (tree *RedBlackTree) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
func (tree *RedBlackTree) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
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

// SetIfNotExist sets <value> to the map if the <key> does not exist, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (tree *RedBlackTree) SetIfNotExist(key interface{}, value interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (tree *RedBlackTree) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function <f>, and then returns true.
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
	var (
		keys  = make([]interface{}, tree.Size())
		index = 0
	)
	tree.IteratorAsc(func(key, value interface{}) bool {
		keys[index] = key
		index++
		return true
	})
	return keys
}

// Values returns all values in asc order based on the key.
func (tree *RedBlackTree) Values() []interface{} {
	var (
		values = make([]interface{}, tree.Size())
		index  = 0
	)
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

// MapStrAny returns all key-value items as map[string]interface{}.
func (tree *RedBlackTree) MapStrAny() map[string]interface{} {
	m := make(map[string]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[gconv.String(key)] = value
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
			Key:   node.Key,
			Value: node.Value,
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
			Key:   node.Key,
			Value: node.Value,
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

// Floor Finds floor node of the input key, return the floor node or nil if no floor node is found.
// Second return parameter is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that its key is smaller than or equal to the given <key>.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
func (tree *RedBlackTree) Floor(key interface{}) (floor *RedBlackTreeNode, found bool) {
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
// Ceiling node is defined as the smallest node that its key is larger than or equal to the given <key>.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
func (tree *RedBlackTree) Ceiling(key interface{}) (ceiling *RedBlackTreeNode, found bool) {
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
func (tree *RedBlackTree) Iterator(f func(key, value interface{}) bool) {
	tree.IteratorAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
func (tree *RedBlackTree) IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorAsc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorAsc(tree.leftNode(), f)
}

// IteratorAscFrom iterates the tree readonly in ascending order with given callback function <f>.
// The parameter <key> specifies the start entry for iterating. The <match> specifies whether
// starting iterating if the <key> is fully matched, or else using index searching iterating.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
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

func (tree *RedBlackTree) doIteratorAsc(node *RedBlackTreeNode, f func(key, value interface{}) bool) {
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

// IteratorDesc iterates the tree readonly in descending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorDesc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorDesc(tree.rightNode(), f)
}

// IteratorDescFrom iterates the tree readonly in descending order with given callback function <f>.
// The parameter <key> specifies the start entry for iterating. The <match> specifies whether
// starting iterating if the <key> is fully matched, or else using index searching iterating.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *RedBlackTree) IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
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

func (tree *RedBlackTree) doIteratorDesc(node *RedBlackTreeNode, f func(key, value interface{}) bool) {
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
func (tree *RedBlackTree) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
}

// Replace the data of the tree with given <data>.
func (tree *RedBlackTree) Replace(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// String returns a string representation of container.
func (tree *RedBlackTree) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	str := ""
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
	node, found := tree.doSearch(key)
	if found {
		return node.Value, true
	}
	return nil, false
}

// Flip exchanges key-value of the tree to value-key.
// Note that you should guarantee the value is the same type as key,
// or else the comparator would panic.
//
// If the type of value is different with key, you pass the new <comparator>.
func (tree *RedBlackTree) Flip(comparator ...func(v1, v2 interface{}) int) {
	t := (*RedBlackTree)(nil)
	if len(comparator) > 0 {
		t = NewRedBlackTree(comparator[0], tree.mu.IsSafe())
	} else {
		t = NewRedBlackTree(tree.comparator, tree.mu.IsSafe())
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
func (tree *RedBlackTree) doSearch(key interface{}) (node *RedBlackTreeNode, found bool) {
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
	right.left = node
	node.parent = right
}

func (tree *RedBlackTree) rotateRight(node *RedBlackTreeNode) {
	left := node.left
	tree.replaceNode(node, left)
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}
	left.right = node
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

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree *RedBlackTree) MarshalJSON() ([]byte, error) {
	return json.Marshal(gconv.Map(tree.Map()))
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (tree *RedBlackTree) UnmarshalJSON(b []byte) error {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.comparator == nil {
		tree.comparator = gutil.ComparatorString
	}
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
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
	}
	for k, v := range gconv.Map(value) {
		tree.doSet(k, v)
	}
	return
}

// getComparator returns the comparator if it's previously set,
// or else it panics.
func (tree *RedBlackTree) getComparator() func(a, b interface{}) int {
	if tree.comparator == nil {
		panic("comparator is missing for tree")
	}
	return tree.comparator
}

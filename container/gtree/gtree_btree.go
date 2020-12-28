// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/internal/json"
	"strings"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/rwmutex"
)

// BTree holds elements of the B-tree.
type BTree struct {
	mu         rwmutex.RWMutex
	root       *BTreeNode
	comparator func(v1, v2 interface{}) int
	size       int // Total number of keys in the tree
	m          int // order (maximum number of children)
}

// BTreeNode is a single element within the tree.
type BTreeNode struct {
	Parent   *BTreeNode
	Entries  []*BTreeEntry // Contained keys in node
	Children []*BTreeNode  // Children nodes
}

// BTreeEntry represents the key-value pair contained within nodes.
type BTreeEntry struct {
	Key   interface{}
	Value interface{}
}

// NewBTree instantiates a B-tree with <m> (maximum number of children) and a custom key comparator.
// The parameter <safe> is used to specify whether using tree in concurrent-safety,
// which is false in default.
// Note that the <m> must be greater or equal than 3, or else it panics.
func NewBTree(m int, comparator func(v1, v2 interface{}) int, safe ...bool) *BTree {
	if m < 3 {
		panic("Invalid order, should be at least 3")
	}
	return &BTree{
		comparator: comparator,
		mu:         rwmutex.Create(safe...),
		m:          m,
	}
}

// NewBTreeFrom instantiates a B-tree with <m> (maximum number of children), a custom key comparator and data map.
// The parameter <safe> is used to specify whether using tree in concurrent-safety,
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

// doSet inserts key-value pair node into the tree.
// If key already exists, then its value is updated with the new value.
func (tree *BTree) doSet(key interface{}, value interface{}) {
	entry := &BTreeEntry{Key: key, Value: value}
	if tree.root == nil {
		tree.root = &BTreeNode{Entries: []*BTreeEntry{entry}, Children: []*BTreeNode{}}
		tree.size++
		return
	}

	if tree.insert(tree.root, entry) {
		tree.size++
	}
}

// Sets batch sets key-values to the tree.
func (tree *BTree) Sets(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// Get searches the node in the tree by <key> and returns its value or nil if key is not found in tree.
func (tree *BTree) Get(key interface{}) (value interface{}) {
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
func (tree *BTree) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if entry := tree.doSearch(key); entry != nil {
		return entry.Value
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
func (tree *BTree) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
func (tree *BTree) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
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
func (tree *BTree) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

// GetVar returns a gvar.Var with the value by given <key>.
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

// SetIfNotExist sets <value> to the map if the <key> does not exist, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (tree *BTree) SetIfNotExist(key interface{}, value interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (tree *BTree) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
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
func (tree *BTree) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

// Contains checks whether <key> exists in the tree.
func (tree *BTree) Contains(key interface{}) bool {
	_, ok := tree.Search(key)
	return ok
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *BTree) doRemove(key interface{}) (value interface{}) {
	node, index, found := tree.searchRecursively(tree.root, key)
	if found {
		value = node.Entries[index].Value
		tree.delete(node, index)
		tree.size--
	}
	return
}

// Remove removes the node from the tree by <key>.
func (tree *BTree) Remove(key interface{}) (value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

// Removes batch deletes values of the tree by <keys>.
func (tree *BTree) Removes(keys []interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.doRemove(key)
	}
}

// Empty returns true if tree does not contain any nodes
func (tree *BTree) IsEmpty() bool {
	return tree.Size() == 0
}

// Size returns number of nodes in the tree.
func (tree *BTree) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.size
}

// Keys returns all keys in asc order.
func (tree *BTree) Keys() []interface{} {
	keys := make([]interface{}, tree.Size())
	index := 0
	tree.IteratorAsc(func(key, value interface{}) bool {
		keys[index] = key
		index++
		return true
	})
	return keys
}

// Values returns all values in asc order based on the key.
func (tree *BTree) Values() []interface{} {
	values := make([]interface{}, tree.Size())
	index := 0
	tree.IteratorAsc(func(key, value interface{}) bool {
		values[index] = value
		index++
		return true
	})
	return values
}

// Map returns all key-value items as map.
func (tree *BTree) Map() map[interface{}]interface{} {
	m := make(map[interface{}]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

// MapStrAny returns all key-value items as map[string]interface{}.
func (tree *BTree) MapStrAny() map[string]interface{} {
	m := make(map[string]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[gconv.String(key)] = value
		return true
	})
	return m
}

// Clear removes all nodes from the tree.
func (tree *BTree) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
}

// Replace the data of the tree with given <data>.
func (tree *BTree) Replace(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
	for k, v := range data {
		tree.doSet(k, v)
	}
}

// Height returns the height of the tree.
func (tree *BTree) Height() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.root.height()
}

// Left returns the left-most (min) entry or nil if tree is empty.
func (tree *BTree) Left() *BTreeEntry {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.left(tree.root)
	return node.Entries[0]
}

// Right returns the right-most (max) entry or nil if tree is empty.
func (tree *BTree) Right() *BTreeEntry {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.right(tree.root)
	return node.Entries[len(node.Entries)-1]
}

// String returns a string representation of container (for debugging purposes)
func (tree *BTree) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var buffer bytes.Buffer
	if tree.size != 0 {
		tree.output(&buffer, tree.root, 0, true)
	}
	return buffer.String()
}

// Search searches the tree with given <key>.
// Second return parameter <found> is true if key was found, otherwise false.
func (tree *BTree) Search(key interface{}) (value interface{}, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, index, found := tree.searchRecursively(tree.root, key)
	if found {
		return node.Entries[index].Value, true
	}
	return nil, false
}

// Search searches the tree with given <key> without mutex.
// It returns the entry if found or otherwise nil.
func (tree *BTree) doSearch(key interface{}) *BTreeEntry {
	node, index, found := tree.searchRecursively(tree.root, key)
	if found {
		return node.Entries[index]
	}
	return nil
}

// Print prints the tree to stdout.
func (tree *BTree) Print() {
	fmt.Println(tree.String())
}

// Iterator is alias of IteratorAsc.
func (tree *BTree) Iterator(f func(key, value interface{}) bool) {
	tree.IteratorAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
func (tree *BTree) IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorAsc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.left(tree.root)
	if node == nil {
		return
	}
	tree.doIteratorAsc(node, node.Entries[0], 0, f)
}

// IteratorAscFrom iterates the tree readonly in ascending order with given callback function <f>.
// The parameter <key> specifies the start entry for iterating. The <match> specifies whether
// starting iterating if the <key> is fully matched, or else using index searching iterating.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, index, found := tree.searchRecursively(tree.root, key)
	if match {
		if found {
			tree.doIteratorAsc(node, node.Entries[index], index, f)
		}
	} else {
		tree.doIteratorAsc(node, node.Entries[index], index, f)
	}
}

func (tree *BTree) doIteratorAsc(node *BTreeNode, entry *BTreeEntry, index int, f func(key, value interface{}) bool) {
	first := true
loop:
	if entry == nil {
		return
	}
	if !f(entry.Key, entry.Value) {
		return
	}
	// Find current entry position in current node
	if !first {
		index, _ = tree.search(node, entry.Key)
	} else {
		first = false
	}
	// Try to go down to the child right of the current entry
	if index+1 < len(node.Children) {
		node = node.Children[index+1]
		// Try to go down to the child left of the current node
		for len(node.Children) > 0 {
			node = node.Children[0]
		}
		// Return the left-most entry
		entry = node.Entries[0]
		goto loop
	}
	// Above assures that we have reached a leaf node, so return the next entry in current node (if any)
	if index+1 < len(node.Entries) {
		entry = node.Entries[index+1]
		goto loop
	}
	// Reached leaf node and there are no entries to the right of the current entry, so go up to the parent
	for node.Parent != nil {
		node = node.Parent
		// Find next entry position in current node (note: search returns the first equal or bigger than entry)
		index, _ = tree.search(node, entry.Key)
		// Check that there is a next entry position in current node
		if index < len(node.Entries) {
			entry = node.Entries[index]
			goto loop
		}
	}
}

// IteratorDesc iterates the tree readonly in descending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorDesc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.right(tree.root)
	if node == nil {
		return
	}
	index := len(node.Entries) - 1
	entry := node.Entries[index]
	tree.doIteratorDesc(node, entry, index, f)
}

// IteratorDescFrom iterates the tree readonly in descending order with given callback function <f>.
// The parameter <key> specifies the start entry for iterating. The <match> specifies whether
// starting iterating if the <key> is fully matched, or else using index searching iterating.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *BTree) IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, index, found := tree.searchRecursively(tree.root, key)
	if match {
		if found {
			tree.doIteratorDesc(node, node.Entries[index], index, f)
		}
	} else {
		tree.doIteratorDesc(node, node.Entries[index], index, f)
	}
}

// IteratorDesc iterates the tree readonly in descending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *BTree) doIteratorDesc(node *BTreeNode, entry *BTreeEntry, index int, f func(key, value interface{}) bool) {
	first := true
loop:
	if entry == nil {
		return
	}
	if !f(entry.Key, entry.Value) {
		return
	}
	// Find current entry position in current node
	if !first {
		index, _ = tree.search(node, entry.Key)
	} else {
		first = false
	}
	// Try to go down to the child left of the current entry
	if index < len(node.Children) {
		node = node.Children[index]
		// Try to go down to the child right of the current node
		for len(node.Children) > 0 {
			node = node.Children[len(node.Children)-1]
		}
		// Return the right-most entry
		entry = node.Entries[len(node.Entries)-1]
		goto loop
	}
	// Above assures that we have reached a leaf node, so return the previous entry in current node (if any)
	if index-1 >= 0 {
		entry = node.Entries[index-1]
		goto loop
	}

	// Reached leaf node and there are no entries to the left of the current entry, so go up to the parent
	for node.Parent != nil {
		node = node.Parent
		// Find previous entry position in current node (note: search returns the first equal or bigger than entry)
		index, _ = tree.search(node, entry.Key)
		// Check that there is a previous entry position in current node
		if index-1 >= 0 {
			entry = node.Entries[index-1]
			goto loop
		}
	}
}

func (tree *BTree) output(buffer *bytes.Buffer, node *BTreeNode, level int, isTail bool) {
	for e := 0; e < len(node.Entries)+1; e++ {
		if e < len(node.Children) {
			tree.output(buffer, node.Children[e], level+1, true)
		}
		if e < len(node.Entries) {
			if _, err := buffer.WriteString(strings.Repeat("    ", level)); err != nil {
			}
			if _, err := buffer.WriteString(fmt.Sprintf("%v", node.Entries[e].Key) + "\n"); err != nil {
			}
		}
	}
}

func (node *BTreeNode) height() int {
	h := 0
	n := node
	for ; n != nil; n = n.Children[0] {
		h++
		if len(n.Children) == 0 {
			break
		}
	}
	return h
}

func (tree *BTree) isLeaf(node *BTreeNode) bool {
	return len(node.Children) == 0
}

//func (tree *BTree) isFull(node *BTreeNode) bool {
//	return len(node.Entries) == tree.maxEntries()
//}

func (tree *BTree) shouldSplit(node *BTreeNode) bool {
	return len(node.Entries) > tree.maxEntries()
}

func (tree *BTree) maxChildren() int {
	return tree.m
}

func (tree *BTree) minChildren() int {
	return (tree.m + 1) / 2 // ceil(m/2)
}

func (tree *BTree) maxEntries() int {
	return tree.maxChildren() - 1
}

func (tree *BTree) minEntries() int {
	return tree.minChildren() - 1
}

func (tree *BTree) middle() int {
	// "-1" to favor right nodes to have more keys when splitting
	return (tree.m - 1) / 2
}

// search searches only within the single node among its entries
func (tree *BTree) search(node *BTreeNode, key interface{}) (index int, found bool) {
	low, mid, high := 0, 0, len(node.Entries)-1
	for low <= high {
		mid = low + int((high-low)/2)
		compare := tree.getComparator()(key, node.Entries[mid].Key)
		switch {
		case compare > 0:
			low = mid + 1
		case compare < 0:
			high = mid - 1
		case compare == 0:
			return mid, true
		}
	}
	return low, false
}

// searchRecursively searches recursively down the tree starting at the startNode
func (tree *BTree) searchRecursively(startNode *BTreeNode, key interface{}) (node *BTreeNode, index int, found bool) {
	if tree.size == 0 {
		return nil, -1, false
	}
	node = startNode
	for {
		index, found = tree.search(node, key)
		if found {
			return node, index, true
		}
		if tree.isLeaf(node) {
			return node, index, false
		}
		node = node.Children[index]
	}
}

func (tree *BTree) insert(node *BTreeNode, entry *BTreeEntry) (inserted bool) {
	if tree.isLeaf(node) {
		return tree.insertIntoLeaf(node, entry)
	}
	return tree.insertIntoInternal(node, entry)
}

func (tree *BTree) insertIntoLeaf(node *BTreeNode, entry *BTreeEntry) (inserted bool) {
	insertPosition, found := tree.search(node, entry.Key)
	if found {
		node.Entries[insertPosition] = entry
		return false
	}
	// Insert entry's key in the middle of the node
	node.Entries = append(node.Entries, nil)
	copy(node.Entries[insertPosition+1:], node.Entries[insertPosition:])
	node.Entries[insertPosition] = entry
	tree.split(node)
	return true
}

func (tree *BTree) insertIntoInternal(node *BTreeNode, entry *BTreeEntry) (inserted bool) {
	insertPosition, found := tree.search(node, entry.Key)
	if found {
		node.Entries[insertPosition] = entry
		return false
	}
	return tree.insert(node.Children[insertPosition], entry)
}

func (tree *BTree) split(node *BTreeNode) {
	if !tree.shouldSplit(node) {
		return
	}

	if node == tree.root {
		tree.splitRoot()
		return
	}

	tree.splitNonRoot(node)
}

func (tree *BTree) splitNonRoot(node *BTreeNode) {
	middle := tree.middle()
	parent := node.Parent

	left := &BTreeNode{Entries: append([]*BTreeEntry(nil), node.Entries[:middle]...), Parent: parent}
	right := &BTreeNode{Entries: append([]*BTreeEntry(nil), node.Entries[middle+1:]...), Parent: parent}

	// Move children from the node to be split into left and right nodes
	if !tree.isLeaf(node) {
		left.Children = append([]*BTreeNode(nil), node.Children[:middle+1]...)
		right.Children = append([]*BTreeNode(nil), node.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	insertPosition, _ := tree.search(parent, node.Entries[middle].Key)

	// Insert middle key into parent
	parent.Entries = append(parent.Entries, nil)
	copy(parent.Entries[insertPosition+1:], parent.Entries[insertPosition:])
	parent.Entries[insertPosition] = node.Entries[middle]

	// Set child left of inserted key in parent to the created left node
	parent.Children[insertPosition] = left

	// Set child right of inserted key in parent to the created right node
	parent.Children = append(parent.Children, nil)
	copy(parent.Children[insertPosition+2:], parent.Children[insertPosition+1:])
	parent.Children[insertPosition+1] = right

	tree.split(parent)
}

func (tree *BTree) splitRoot() {
	middle := tree.middle()
	left := &BTreeNode{Entries: append([]*BTreeEntry(nil), tree.root.Entries[:middle]...)}
	right := &BTreeNode{Entries: append([]*BTreeEntry(nil), tree.root.Entries[middle+1:]...)}

	// Move children from the node to be split into left and right nodes
	if !tree.isLeaf(tree.root) {
		left.Children = append([]*BTreeNode(nil), tree.root.Children[:middle+1]...)
		right.Children = append([]*BTreeNode(nil), tree.root.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	// Root is a node with one entry and two children (left and right)
	newRoot := &BTreeNode{
		Entries:  []*BTreeEntry{tree.root.Entries[middle]},
		Children: []*BTreeNode{left, right},
	}

	left.Parent = newRoot
	right.Parent = newRoot
	tree.root = newRoot
}

func setParent(nodes []*BTreeNode, parent *BTreeNode) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

func (tree *BTree) left(node *BTreeNode) *BTreeNode {
	if tree.size == 0 {
		return nil
	}
	current := node
	for {
		if tree.isLeaf(current) {
			return current
		}
		current = current.Children[0]
	}
}

func (tree *BTree) right(node *BTreeNode) *BTreeNode {
	if tree.size == 0 {
		return nil
	}
	current := node
	for {
		if tree.isLeaf(current) {
			return current
		}
		current = current.Children[len(current.Children)-1]
	}
}

// leftSibling returns the node's left sibling and child index (in parent) if it exists, otherwise (nil,-1)
// key is any of keys in node (could even be deleted).
func (tree *BTree) leftSibling(node *BTreeNode, key interface{}) (*BTreeNode, int) {
	if node.Parent != nil {
		index, _ := tree.search(node.Parent, key)
		index--
		if index >= 0 && index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}
	return nil, -1
}

// rightSibling returns the node's right sibling and child index (in parent) if it exists, otherwise (nil,-1)
// key is any of keys in node (could even be deleted).
func (tree *BTree) rightSibling(node *BTreeNode, key interface{}) (*BTreeNode, int) {
	if node.Parent != nil {
		index, _ := tree.search(node.Parent, key)
		index++
		if index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}
	return nil, -1
}

// delete deletes an entry in node at entries' index
// ref.: https://en.wikipedia.org/wiki/B-tree#Deletion
func (tree *BTree) delete(node *BTreeNode, index int) {
	// deleting from a leaf node
	if tree.isLeaf(node) {
		deletedKey := node.Entries[index].Key
		tree.deleteEntry(node, index)
		tree.rebalance(node, deletedKey)
		if len(tree.root.Entries) == 0 {
			tree.root = nil
		}
		return
	}

	// deleting from an internal node
	leftLargestNode := tree.right(node.Children[index]) // largest node in the left sub-tree (assumed to exist)
	leftLargestEntryIndex := len(leftLargestNode.Entries) - 1
	node.Entries[index] = leftLargestNode.Entries[leftLargestEntryIndex]
	deletedKey := leftLargestNode.Entries[leftLargestEntryIndex].Key
	tree.deleteEntry(leftLargestNode, leftLargestEntryIndex)
	tree.rebalance(leftLargestNode, deletedKey)
}

// rebalance rebalances the tree after deletion if necessary and returns true, otherwise false.
// Note that we first delete the entry and then call rebalance, thus the passed deleted key as reference.
func (tree *BTree) rebalance(node *BTreeNode, deletedKey interface{}) {
	// check if rebalancing is needed
	if node == nil || len(node.Entries) >= tree.minEntries() {
		return
	}

	// try to borrow from left sibling
	leftSibling, leftSiblingIndex := tree.leftSibling(node, deletedKey)
	if leftSibling != nil && len(leftSibling.Entries) > tree.minEntries() {
		// rotate right
		node.Entries = append([]*BTreeEntry{node.Parent.Entries[leftSiblingIndex]}, node.Entries...) // prepend parent's separator entry to node's entries
		node.Parent.Entries[leftSiblingIndex] = leftSibling.Entries[len(leftSibling.Entries)-1]
		tree.deleteEntry(leftSibling, len(leftSibling.Entries)-1)
		if !tree.isLeaf(leftSibling) {
			leftSiblingRightMostChild := leftSibling.Children[len(leftSibling.Children)-1]
			leftSiblingRightMostChild.Parent = node
			node.Children = append([]*BTreeNode{leftSiblingRightMostChild}, node.Children...)
			tree.deleteChild(leftSibling, len(leftSibling.Children)-1)
		}
		return
	}

	// try to borrow from right sibling
	rightSibling, rightSiblingIndex := tree.rightSibling(node, deletedKey)
	if rightSibling != nil && len(rightSibling.Entries) > tree.minEntries() {
		// rotate left
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1]) // append parent's separator entry to node's entries
		node.Parent.Entries[rightSiblingIndex-1] = rightSibling.Entries[0]
		tree.deleteEntry(rightSibling, 0)
		if !tree.isLeaf(rightSibling) {
			rightSiblingLeftMostChild := rightSibling.Children[0]
			rightSiblingLeftMostChild.Parent = node
			node.Children = append(node.Children, rightSiblingLeftMostChild)
			tree.deleteChild(rightSibling, 0)
		}
		return
	}

	// merge with siblings
	if rightSibling != nil {
		// merge with right sibling
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Entries = append(node.Entries, rightSibling.Entries...)
		deletedKey = node.Parent.Entries[rightSiblingIndex-1].Key
		tree.deleteEntry(node.Parent, rightSiblingIndex-1)
		tree.appendChildren(node.Parent.Children[rightSiblingIndex], node)
		tree.deleteChild(node.Parent, rightSiblingIndex)
	} else if leftSibling != nil {
		// merge with left sibling
		entries := append([]*BTreeEntry(nil), leftSibling.Entries...)
		entries = append(entries, node.Parent.Entries[leftSiblingIndex])
		node.Entries = append(entries, node.Entries...)
		deletedKey = node.Parent.Entries[leftSiblingIndex].Key
		tree.deleteEntry(node.Parent, leftSiblingIndex)
		tree.prependChildren(node.Parent.Children[leftSiblingIndex], node)
		tree.deleteChild(node.Parent, leftSiblingIndex)
	}

	// make the merged node the root if its parent was the root and the root is empty
	if node.Parent == tree.root && len(tree.root.Entries) == 0 {
		tree.root = node
		node.Parent = nil
		return
	}

	// parent might underflow, so try to rebalance if necessary
	tree.rebalance(node.Parent, deletedKey)
}

func (tree *BTree) prependChildren(fromNode *BTreeNode, toNode *BTreeNode) {
	children := append([]*BTreeNode(nil), fromNode.Children...)
	toNode.Children = append(children, toNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (tree *BTree) appendChildren(fromNode *BTreeNode, toNode *BTreeNode) {
	toNode.Children = append(toNode.Children, fromNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (tree *BTree) deleteEntry(node *BTreeNode, index int) {
	copy(node.Entries[index:], node.Entries[index+1:])
	node.Entries[len(node.Entries)-1] = nil
	node.Entries = node.Entries[:len(node.Entries)-1]
}

func (tree *BTree) deleteChild(node *BTreeNode, index int) {
	if index >= len(node.Children) {
		return
	}
	copy(node.Children[index:], node.Children[index+1:])
	node.Children[len(node.Children)-1] = nil
	node.Children = node.Children[:len(node.Children)-1]
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree *BTree) MarshalJSON() ([]byte, error) {
	return json.Marshal(tree.Map())
}

// getComparator returns the comparator if it's previously set,
// or else it panics.
func (tree *BTree) getComparator() func(a, b interface{}) int {
	if tree.comparator == nil {
		panic("comparator is missing for tree")
	}
	return tree.comparator
}

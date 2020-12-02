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

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/rwmutex"
)

// AVLTree holds elements of the AVL tree.
type AVLTree struct {
	mu         rwmutex.RWMutex
	root       *AVLTreeNode
	comparator func(v1, v2 interface{}) int
	size       int
}

// AVLTreeNode is a single element within the tree.
type AVLTreeNode struct {
	Key      interface{}
	Value    interface{}
	parent   *AVLTreeNode
	children [2]*AVLTreeNode
	b        int8
}

// NewAVLTree instantiates an AVL tree with the custom key comparator.
// The parameter <safe> is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewAVLTree(comparator func(v1, v2 interface{}) int, safe ...bool) *AVLTree {
	return &AVLTree{
		mu:         rwmutex.Create(safe...),
		comparator: comparator,
	}
}

// NewAVLTreeFrom instantiates an AVL tree with the custom key comparator and data map.
// The parameter <safe> is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewAVLTreeFrom(comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}, safe ...bool) *AVLTree {
	tree := NewAVLTree(comparator, safe...)
	for k, v := range data {
		tree.put(k, v, nil, &tree.root)
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
	tree.put(key, value, nil, &tree.root)
}

// Sets batch sets key-values to the tree.
func (tree *AVLTree) Sets(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for key, value := range data {
		tree.put(key, value, nil, &tree.root)
	}
}

// Search searches the tree with given <key>.
// Second return parameter <found> is true if key was found, otherwise false.
func (tree *AVLTree) Search(key interface{}) (value interface{}, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if node, found := tree.doSearch(key); found {
		return node.Value, true
	}
	return nil, false
}

// doSearch searches the tree with given <key>.
// Second return parameter <found> is true if key was found, otherwise false.
func (tree *AVLTree) doSearch(key interface{}) (node *AVLTreeNode, found bool) {
	node = tree.root
	for node != nil {
		cmp := tree.getComparator()(key, node.Key)
		switch {
		case cmp == 0:
			return node, true
		case cmp < 0:
			node = node.children[0]
		case cmp > 0:
			node = node.children[1]
		}
	}
	return nil, false
}

// Get searches the node in the tree by <key> and returns its value or nil if key is not found in tree.
func (tree *AVLTree) Get(key interface{}) (value interface{}) {
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
func (tree *AVLTree) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if node, found := tree.doSearch(key); found {
		return node.Value
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value != nil {
		tree.put(key, value, nil, &tree.root)
	}
	return value
}

// GetOrSet returns the value by key,
// or sets value with given <value> if it does not exist and then returns this value.
func (tree *AVLTree) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
func (tree *AVLTree) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
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
func (tree *AVLTree) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

// GetVar returns a gvar.Var with the value by given <key>.
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

// SetIfNotExist sets <value> to the map if the <key> does not exist, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (tree *AVLTree) SetIfNotExist(key interface{}, value interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (tree *AVLTree) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
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
func (tree *AVLTree) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

// Contains checks whether <key> exists in the tree.
func (tree *AVLTree) Contains(key interface{}) bool {
	_, ok := tree.Search(key)
	return ok
}

// Remove removes the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *AVLTree) Remove(key interface{}) (value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	value, _ = tree.remove(key, &tree.root)
	return
}

// Removes batch deletes values of the tree by <keys>.
func (tree *AVLTree) Removes(keys []interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.remove(key, &tree.root)
	}
}

// IsEmpty returns true if tree does not contain any nodes.
func (tree *AVLTree) IsEmpty() bool {
	return tree.Size() == 0
}

// Size returns number of nodes in the tree.
func (tree *AVLTree) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.size
}

// Keys returns all keys in asc order.
func (tree *AVLTree) Keys() []interface{} {
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
func (tree *AVLTree) Values() []interface{} {
	values := make([]interface{}, tree.Size())
	index := 0
	tree.IteratorAsc(func(key, value interface{}) bool {
		values[index] = value
		index++
		return true
	})
	return values
}

// Left returns the minimum element of the AVL tree
// or nil if the tree is empty.
func (tree *AVLTree) Left() *AVLTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.bottom(0)
	if tree.mu.IsSafe() {
		return &AVLTreeNode{
			Key:   node.Key,
			Value: node.Value,
		}
	}
	return node
}

// Right returns the maximum element of the AVL tree
// or nil if the tree is empty.
func (tree *AVLTree) Right() *AVLTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.bottom(1)
	if tree.mu.IsSafe() {
		return &AVLTreeNode{
			Key:   node.Key,
			Value: node.Value,
		}
	}
	return node
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
	n := tree.root
	for n != nil {
		c := tree.getComparator()(key, n.Key)
		switch {
		case c == 0:
			return n, true
		case c < 0:
			n = n.children[0]
		case c > 0:
			floor, found = n, true
			n = n.children[1]
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
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree is smaller than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *AVLTree) Ceiling(key interface{}) (ceiling *AVLTreeNode, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	n := tree.root
	for n != nil {
		c := tree.getComparator()(key, n.Key)
		switch {
		case c == 0:
			return n, true
		case c > 0:
			n = n.children[1]
		case c < 0:
			ceiling, found = n, true
			n = n.children[0]
		}
	}
	if found {
		return
	}
	return nil, false
}

// Clear removes all nodes from the tree.
func (tree *AVLTree) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
}

// Replace the data of the tree with given <data>.
func (tree *AVLTree) Replace(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
	for key, value := range data {
		tree.put(key, value, nil, &tree.root)
	}
}

// String returns a string representation of container
func (tree *AVLTree) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	str := ""
	if tree.size != 0 {
		output(tree.root, "", true, &str)
	}
	return str
}

// Print prints the tree to stdout.
func (tree *AVLTree) Print() {
	fmt.Println(tree.String())
}

// Map returns all key-value items as map.
func (tree *AVLTree) Map() map[interface{}]interface{} {
	m := make(map[interface{}]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

// MapStrAny returns all key-value items as map[string]interface{}.
func (tree *AVLTree) MapStrAny() map[string]interface{} {
	m := make(map[string]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[gconv.String(key)] = value
		return true
	})
	return m
}

// Flip exchanges key-value of the tree to value-key.
// Note that you should guarantee the value is the same type as key,
// or else the comparator would panic.
//
// If the type of value is different with key, you pass the new <comparator>.
func (tree *AVLTree) Flip(comparator ...func(v1, v2 interface{}) int) {
	t := (*AVLTree)(nil)
	if len(comparator) > 0 {
		t = NewAVLTree(comparator[0], tree.mu.IsSafe())
	} else {
		t = NewAVLTree(tree.comparator, tree.mu.IsSafe())
	}
	tree.IteratorAsc(func(key, value interface{}) bool {
		t.put(value, key, nil, &t.root)
		return true
	})
	tree.mu.Lock()
	tree.root = t.root
	tree.size = t.size
	tree.mu.Unlock()
}

// Iterator is alias of IteratorAsc.
func (tree *AVLTree) Iterator(f func(key, value interface{}) bool) {
	tree.IteratorAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
func (tree *AVLTree) IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// IteratorAsc iterates the tree readonly in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *AVLTree) IteratorAsc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorAsc(tree.bottom(0), f)
}

// IteratorAscFrom iterates the tree readonly in ascending order with given callback function <f>.
// The parameter <key> specifies the start entry for iterating. The <match> specifies whether
// starting iterating if the <key> is fully matched, or else using index searching iterating.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *AVLTree) IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
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

func (tree *AVLTree) doIteratorAsc(node *AVLTreeNode, f func(key, value interface{}) bool) {
	for node != nil {
		if !f(node.Key, node.Value) {
			return
		}
		node = node.Next()
	}
}

// IteratorDesc iterates the tree readonly in descending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *AVLTree) IteratorDesc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorDesc(tree.bottom(1), f)
}

// IteratorDescFrom iterates the tree readonly in descending order with given callback function <f>.
// The parameter <key> specifies the start entry for iterating. The <match> specifies whether
// starting iterating if the <key> is fully matched, or else using index searching iterating.
// If <f> returns true, then it continues iterating; or false to stop.
func (tree *AVLTree) IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
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

func (tree *AVLTree) doIteratorDesc(node *AVLTreeNode, f func(key, value interface{}) bool) {
	for node != nil {
		if !f(node.Key, node.Value) {
			return
		}
		node = node.Prev()
	}
}

func (tree *AVLTree) put(key interface{}, value interface{}, p *AVLTreeNode, qp **AVLTreeNode) bool {
	q := *qp
	if q == nil {
		tree.size++
		*qp = &AVLTreeNode{Key: key, Value: value, parent: p}
		return true
	}

	c := tree.getComparator()(key, q.Key)
	if c == 0 {
		q.Key = key
		q.Value = value
		return false
	}

	if c < 0 {
		c = -1
	} else {
		c = 1
	}
	a := (c + 1) / 2
	if tree.put(key, value, q, &q.children[a]) {
		return putFix(int8(c), qp)
	}
	return false
}

func (tree *AVLTree) remove(key interface{}, qp **AVLTreeNode) (value interface{}, fix bool) {
	q := *qp
	if q == nil {
		return nil, false
	}

	c := tree.getComparator()(key, q.Key)
	if c == 0 {
		tree.size--
		value = q.Value
		fix = true
		if q.children[1] == nil {
			if q.children[0] != nil {
				q.children[0].parent = q.parent
			}
			*qp = q.children[0]
			return
		}
		if removeMin(&q.children[1], &q.Key, &q.Value) {
			return value, removeFix(-1, qp)
		}
		return
	}

	if c < 0 {
		c = -1
	} else {
		c = 1
	}
	a := (c + 1) / 2
	value, fix = tree.remove(key, &q.children[a])
	if fix {
		return value, removeFix(int8(-c), qp)
	}
	return value, false
}

func removeMin(qp **AVLTreeNode, minKey *interface{}, minVal *interface{}) bool {
	q := *qp
	if q.children[0] == nil {
		*minKey = q.Key
		*minVal = q.Value
		if q.children[1] != nil {
			q.children[1].parent = q.parent
		}
		*qp = q.children[1]
		return true
	}
	fix := removeMin(&q.children[0], minKey, minVal)
	if fix {
		return removeFix(1, qp)
	}
	return false
}

func putFix(c int8, t **AVLTreeNode) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return true
	}

	if s.b == -c {
		s.b = 0
		return false
	}

	if s.children[(c+1)/2].b == c {
		s = singleRotate(c, s)
	} else {
		s = doubleRotate(c, s)
	}
	*t = s
	return false
}

func removeFix(c int8, t **AVLTreeNode) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return false
	}

	if s.b == -c {
		s.b = 0
		return true
	}

	a := (c + 1) / 2
	if s.children[a].b == 0 {
		s = rotate(c, s)
		s.b = -c
		*t = s
		return false
	}

	if s.children[a].b == c {
		s = singleRotate(c, s)
	} else {
		s = doubleRotate(c, s)
	}
	*t = s
	return true
}

func singleRotate(c int8, s *AVLTreeNode) *AVLTreeNode {
	s.b = 0
	s = rotate(c, s)
	s.b = 0
	return s
}

func doubleRotate(c int8, s *AVLTreeNode) *AVLTreeNode {
	a := (c + 1) / 2
	r := s.children[a]
	s.children[a] = rotate(-c, s.children[a])
	p := rotate(c, s)

	switch {
	default:
		s.b = 0
		r.b = 0
	case p.b == c:
		s.b = -c
		r.b = 0
	case p.b == -c:
		s.b = 0
		r.b = c
	}

	p.b = 0
	return p
}

func rotate(c int8, s *AVLTreeNode) *AVLTreeNode {
	a := (c + 1) / 2
	r := s.children[a]
	s.children[a] = r.children[a^1]
	if s.children[a] != nil {
		s.children[a].parent = s
	}
	r.children[a^1] = s
	r.parent = s.parent
	s.parent = r
	return r
}

func (tree *AVLTree) bottom(d int) *AVLTreeNode {
	n := tree.root
	if n == nil {
		return nil
	}

	for c := n.children[d]; c != nil; c = n.children[d] {
		n = c
	}
	return n
}

// Prev returns the previous element in an inorder
// walk of the AVL tree.
func (node *AVLTreeNode) Prev() *AVLTreeNode {
	return node.walk1(0)
}

// Next returns the next element in an inorder
// walk of the AVL tree.
func (node *AVLTreeNode) Next() *AVLTreeNode {
	return node.walk1(1)
}

func (node *AVLTreeNode) walk1(a int) *AVLTreeNode {
	if node == nil {
		return nil
	}
	n := node
	if n.children[a] != nil {
		n = n.children[a]
		for n.children[a^1] != nil {
			n = n.children[a^1]
		}
		return n
	}

	p := n.parent
	for p != nil && p.children[a] == n {
		n = p
		p = p.parent
	}
	return p
}

func output(node *AVLTreeNode, prefix string, isTail bool, str *string) {
	if node.children[1] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.children[1], newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += fmt.Sprintf("%v\n", node.Key)
	if node.children[0] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.children[0], newPrefix, true, str)
	}
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree *AVLTree) MarshalJSON() ([]byte, error) {
	return json.Marshal(tree.Map())
}

// getComparator returns the comparator if it's previously set,
// or else it panics.
func (tree *AVLTree) getComparator() func(a, b interface{}) int {
	if tree.comparator == nil {
		panic("comparator is missing for tree")
	}
	return tree.comparator
}

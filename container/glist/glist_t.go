// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package glist

import (
	"bytes"
	"container/list"

	"github.com/gogf/gf/v2/internal/deepcopy"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/util/gconv"
)

// TElement is an element of a linked list.
type TElement[T any] struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *TElement[T]

	// The list to which this element belongs.
	list *TList[T]

	// The value stored with this element.
	Value T
}

// Next returns the next list element or nil.
func (e *TElement[T]) Next() *TElement[T] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *TElement[T]) Prev() *TElement[T] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// TList is a doubly linked list containing a concurrent-safe/unsafe switch.
// The switch should be set when its initialization and cannot be changed then.

type TList[T any] struct {
	mu   rwmutex.RWMutex
	root TElement[T] // sentinel list element, only &root, root.prev, and root.next are used
	len  int         // current list length excluding (this) sentinel element
}

// NewT creates and returns a new empty doubly linked list.
func NewT[T any](safe ...bool) *TList[T] {
	l := &TList[T]{
		mu: rwmutex.Create(safe...),
	}
	return l.init()
}

// NewTFrom creates and returns a list from a copy of given slice `array`.
// The parameter `safe` is used to specify whether using list in concurrent-safety,
// which is false in default.
func NewTFrom[T any](array []T, safe ...bool) *TList[T] {
	l := NewT[T](safe...)
	for _, v := range array {
		l.insertValue(v, l.root.prev)
	}
	return l
}

// PushFront inserts a new element `e` with value `v` at the front of list `l` and returns `e`.
func (l *TList[T]) PushFront(v T) (e *TElement[T]) {
	l.mu.Lock()
	l.lazyInit()
	e = l.insertValue(v, &l.root)
	l.mu.Unlock()
	return
}

// PushBack inserts a new element `e` with value `v` at the back of list `l` and returns `e`.
func (l *TList[T]) PushBack(v T) (e *TElement[T]) {
	l.mu.Lock()
	l.lazyInit()
	e = l.insertValue(v, l.root.prev)
	l.mu.Unlock()
	return
}

// PushFronts inserts multiple new elements with values `values` at the front of list `l`.
func (l *TList[T]) PushFronts(values []T) {
	l.mu.Lock()
	l.lazyInit()
	for _, v := range values {
		l.insertValue(v, &l.root)
	}
	l.mu.Unlock()
}

// PushBacks inserts multiple new elements with values `values` at the back of list `l`.
func (l *TList[T]) PushBacks(values []T) {
	l.mu.Lock()
	l.lazyInit()
	for _, v := range values {
		l.insertValue(v, l.root.prev)
	}
	l.mu.Unlock()
}

// PopBack removes the element from back of `l` and returns the value of the element.
func (l *TList[T]) PopBack() (value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	if l.len == 0 {
		return
	}
	return l.remove(l.root.prev)
}

// PopFront removes the element from front of `l` and returns the value of the element.
func (l *TList[T]) PopFront() (value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	if l.len == 0 {
		return
	}
	return l.remove(l.root.next)
}

// PopBacks removes `max` elements from back of `l`
// and returns values of the removed elements as slice.
func (l *TList[T]) PopBacks(max int) (values []T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()

	length := l.len
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]T, length)
		for i := 0; i < length; i++ {
			values[i] = l.remove(l.root.prev)
		}
	}
	return
}

// PopFronts removes `max` elements from front of `l`
// and returns values of the removed elements as slice.
func (l *TList[T]) PopFronts(max int) (values []T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()

	length := l.len
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]T, length)
		for i := 0; i < length; i++ {
			values[i] = l.remove(l.root.next)
		}
	}
	return
}

// PopBackAll removes all elements from back of `l`
// and returns values of the removed elements as slice.
func (l *TList[T]) PopBackAll() []T {
	return l.PopBacks(-1)
}

// PopFrontAll removes all elements from front of `l`
// and returns values of the removed elements as slice.
func (l *TList[T]) PopFrontAll() []T {
	return l.PopFronts(-1)
}

// FrontAll copies and returns values of all elements from front of `l` as slice.
func (l *TList[T]) FrontAll() (values []T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	length := l.len
	if length > 0 {
		values = make([]T, length)
		for i, e := 0, l.front(); i < length; i, e = i+1, e.Next() {
			values[i] = e.Value
		}
	}
	return
}

// BackAll copies and returns values of all elements from back of `l` as slice.
func (l *TList[T]) BackAll() (values []T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	length := l.len
	if length > 0 {
		values = make([]T, length)
		for i, e := 0, l.back(); i < length; i, e = i+1, e.Prev() {
			values[i] = e.Value
		}
	}
	return
}

// FrontValue returns value of the first element of `l` or zero value of T if the list is empty.
func (l *TList[T]) FrontValue() (value T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	if e := l.front(); e != nil {
		value = e.Value
	}
	return
}

// BackValue returns value of the last element of `l` or zero value of T if the list is empty.
func (l *TList[T]) BackValue() (value T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	if e := l.back(); e != nil {
		value = e.Value
	}
	return
}

// Front returns the first element of list `l` or nil if the list is empty.
func (l *TList[T]) Front() (e *TElement[T]) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	l.lazyInit()

	e = l.front()
	return
}

// Back returns the last element of list `l` or nil if the list is empty.
func (l *TList[T]) Back() (e *TElement[T]) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()

	e = l.back()
	return
}

// Len returns the number of elements of list `l`.
// The complexity is O(1).
func (l *TList[T]) Len() (length int) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	l.lazyInit()

	length = l.len
	return
}

// Size is alias of Len.
func (l *TList[T]) Size() int {
	return l.Len()
}

// MoveBefore moves element `e` to its new position before `p`.
// If `e` or `p` is not an element of `l`, or `e` == `p`, the list is not modified.
// The element and `p` must not be nil.
func (l *TList[T]) MoveBefore(e, p *TElement[T]) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.lazyInit()

	if e.list != l || e == p || p.list != l {
		return
	}
	l.move(e, p.prev)
}

// MoveAfter moves element `e` to its new position after `p`.
// If `e` or `p` is not an element of `l`, or `e` == `p`, the list is not modified.
// The element and `p` must not be nil.
func (l *TList[T]) MoveAfter(e, p *TElement[T]) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	if e.list != l || e == p || p.list != l {
		return
	}
	l.move(e, p)
}

// MoveToFront moves element `e` to the front of list `l`.
// If `e` is not an element of `l`, the list is not modified.
// The element must not be nil.
func (l *TList[T]) MoveToFront(e *TElement[T]) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()

	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves element `e` to the back of list `l`.
// If `e` is not an element of `l`, the list is not modified.
// The element must not be nil.
func (l *TList[T]) MoveToBack(e *TElement[T]) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()

	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, l.root.prev)
}

// PushBackList inserts a copy of an other list at the back of list `l`.
// The lists `l` and `other` may be the same, but they must not be nil.
func (l *TList[T]) PushBackList(other *TList[T]) {
	if l != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	l.lazyInit()

	for i, e := other.len, other.front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of an other list at the front of list `l`.
// The lists `l` and `other` may be the same, but they must not be nil.
func (l *TList[T]) PushFrontList(other *TList[T]) {
	if l != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	l.lazyInit()

	for i, e := other.len, other.back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

// InsertAfter inserts a new element `e` with value `v` immediately after `p` and returns `e`.
// If `p` is not an element of `l`, the list is not modified.
// The `p` must not be nil.
func (l *TList[T]) InsertAfter(p *TElement[T], v T) (e *TElement[T]) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.lazyInit()
	if p.list != l {
		return nil
	}
	e = l.insertValue(v, p)
	return
}

// InsertBefore inserts a new element `e` with value `v` immediately before `p` and returns `e`.
// If `p` is not an element of `l`, the list is not modified.
// The `p` must not be nil.
func (l *TList[T]) InsertBefore(p *TElement[T], v T) (e *TElement[T]) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.lazyInit()
	if p.list != l {
		return nil
	}
	e = l.insertValue(v, p.prev)
	return
}

// Remove removes `e` from `l` if `e` is an element of list `l`.
// It returns the element value e.Value.
// The element must not be nil.
func (l *TList[T]) Remove(e *TElement[T]) (value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	return l.remove(e)
}

// Removes removes multiple elements `es` from `l` if `es` are elements of list `l`.
func (l *TList[T]) Removes(es []*TElement[T]) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	for _, e := range es {
		l.remove(e)
	}
}

// RemoveAll removes all elements from list `l`.
func (l *TList[T]) RemoveAll() {
	l.mu.Lock()
	l.init()
	l.mu.Unlock()
}

// Clear is alias of RemoveAll.
func (l *TList[T]) Clear() {
	l.RemoveAll()
}

// ToList converts TList[T] to list.List
func (l *TList[T]) ToList() *list.List {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.toList()
}

// toList converts TList[T] to list.List
func (l *TList[T]) toList() *list.List {
	l.lazyInit()

	nl := list.New()

	for e := l.front(); e != nil; e = e.Next() {
		nl.PushBack(e.Value)
	}
	return nl
}

// AppendList append list.List to the end
func (l *TList[T]) AppendList(nl *list.List) {
	if nl.Len() == 0 {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.appendList(nl)
}

// appendList append list.List to the end
func (l *TList[T]) appendList(nl *list.List) {
	if nl.Len() == 0 {
		return
	}

	l.lazyInit()

	for e := nl.Front(); e != nil; e = e.Next() {
		if v, ok := e.Value.(T); ok {
			l.insertValue(v, l.root.prev)
		}
	}
}

// AssignList assigns list.List to now TList[T].
// It will clear TList[T] first, and append the list.List.
// Note: Elements in nl that are not assignable to T are silently skipped.
// Returns the number of skipped (incompatible) elements.
func (l *TList[T]) AssignList(nl *list.List) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.assignList(nl)
}

// assignList assigns list.List to now TList[T].
// It will clear TList[T] first, and append the list.List.
// Returns the number of skipped (incompatible) elements.
func (l *TList[T]) assignList(nl *list.List) int {
	l.init()
	if nl.Len() == 0 {
		return 0
	}
	skipped := 0
	for e := nl.Front(); e != nil; e = e.Next() {
		if v, ok := e.Value.(T); ok {
			l.insertValue(v, l.root.prev)
		} else {
			skipped++
		}
	}
	return skipped
}

// RLockFunc locks reading with given callback function `f` within RWMutex.RLock.
func (l *TList[T]) RLockFunc(f func(list *list.List)) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f(l.toList())
}

// LockFunc locks writing with given callback function `f` within RWMutex.Lock.
func (l *TList[T]) LockFunc(f func(list *list.List)) {
	l.mu.Lock()
	defer l.mu.Unlock()

	nl := l.toList()
	f(nl)
	l.assignList(nl)
}

// Iterator is alias of IteratorAsc.
func (l *TList[T]) Iterator(f func(e *TElement[T]) bool) {
	l.IteratorAsc(f)
}

// IteratorAsc iterates the list readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (l *TList[T]) IteratorAsc(f func(e *TElement[T]) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	length := l.len
	if length > 0 {
		for i, e := 0, l.front(); i < length; i, e = i+1, e.Next() {
			if !f(e) {
				break
			}
		}
	}
}

// IteratorDesc iterates the list readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (l *TList[T]) IteratorDesc(f func(e *TElement[T]) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	length := l.len
	if length > 0 {
		for i, e := 0, l.back(); i < length; i, e = i+1, e.Prev() {
			if !f(e) {
				break
			}
		}
	}
}

// Join joins list elements with a string `glue`.
func (l *TList[T]) Join(glue string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()

	buffer := bytes.NewBuffer(nil)
	length := l.len
	if length > 0 {
		for i, e := 0, l.front(); i < length; i, e = i+1, e.Next() {
			buffer.WriteString(gconv.String(e.Value))
			if i != length-1 {
				buffer.WriteString(glue)
			}
		}
	}
	return buffer.String()
}

// String returns current list as a string.
func (l *TList[T]) String() string {
	if l == nil {
		return ""
	}
	return "[" + l.Join(",") + "]"
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (l *TList[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.FrontAll())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (l *TList[T]) UnmarshalJSON(b []byte) error {
	var array []T
	if err := json.UnmarshalUseNumber(b, &array); err != nil {
		return err
	}
	l.init()
	l.PushBacks(array)
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for list.
func (l *TList[T]) UnmarshalValue(value any) (err error) {
	var array []T
	switch value.(type) {
	case string, []byte:
		err = json.UnmarshalUseNumber(gconv.Bytes(value), &array)
	default:
		anyArray := gconv.SliceAny(value)
		if err = gconv.Scan(anyArray, &array); err != nil {
			return
		}
	}
	l.init()
	l.PushBacks(array)
	return err
}

// DeepCopy implements interface for deep copy of current type.
func (l *TList[T]) DeepCopy() any {
	if l == nil {
		return nil
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	l.lazyInit()

	var (
		length  = l.len
		valuesT = make([]T, length)
	)
	if length > 0 {
		for i, e := 0, l.front(); i < length; i, e = i+1, e.Next() {
			valuesT[i] = deepcopy.Copy(e.Value).(T)
		}
	}
	return NewTFrom(valuesT, l.mu.IsSafe())
}

// Init initializes or clears list l.
func (l *TList[T]) init() *TList[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// lazyInit lazily initializes a zero List value.
func (l *TList[T]) lazyInit() {
	if l.root.next == nil {
		l.init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *TList[T]) insert(e, at *TElement[T]) *TElement[T] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *TList[T]) insertValue(v T, at *TElement[T]) *TElement[T] {
	return l.insert(&TElement[T]{Value: v}, at)
}

// remove removes e from its list, decrements l.len
func (l *TList[T]) remove(e *TElement[T]) (val T) {
	if e.list != l {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--

	return e.Value
}

// move moves e to next to at.
func (l *TList[T]) move(e, at *TElement[T]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// front returns the first element of list l or nil if the list is empty.
func (l *TList[T]) front() *TElement[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// back returns the last element of list l or nil if the list is empty.
func (l *TList[T]) back() *TElement[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

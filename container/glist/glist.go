// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with l file,
// You can obtain one at https://github.com/gogf/gf.
//

// Package glist provides a concurrent-safe/unsafe doubly linked list.
package glist

import (
	"bytes"
	"container/list"
	"encoding/json"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/internal/rwmutex"
)

type (
	List struct {
		mu   rwmutex.RWMutex
		list list.List
	}

	Element = list.Element
)

// New creates and returns a new empty doubly linked list.
func New(safe ...bool) *List {
	return &List{
		mu:   rwmutex.Create(safe...),
		list: list.List{},
	}
}

// NewFrom creates and returns a list from a copy of given slice <array>.
// The parameter <safe> is used to specify whether using list in concurrent-safety,
// which is false in default.
func NewFrom(array []interface{}, safe ...bool) *List {
	l := list.List{}
	for _, v := range array {
		l.PushBack(v)
	}
	return &List{
		mu:   rwmutex.Create(safe...),
		list: l,
	}
}

// PushFront inserts a new element <e> with value <v> at the front of list <l> and returns <e>.
func (l *List) PushFront(v interface{}) (e *Element) {
	l.mu.Lock()
	e = l.list.PushFront(v)
	l.mu.Unlock()
	return
}

// PushBack inserts a new element <e> with value <v> at the back of list <l> and returns <e>.
func (l *List) PushBack(v interface{}) (e *Element) {
	l.mu.Lock()
	e = l.list.PushBack(v)
	l.mu.Unlock()
	return
}

// PushFronts inserts multiple new elements with values <values> at the front of list <l>.
func (l *List) PushFronts(values []interface{}) {
	l.mu.Lock()
	for _, v := range values {
		l.list.PushFront(v)
	}
	l.mu.Unlock()
}

// PushBacks inserts multiple new elements with values <values> at the back of list <l>.
func (l *List) PushBacks(values []interface{}) {
	l.mu.Lock()
	for _, v := range values {
		l.list.PushBack(v)
	}
	l.mu.Unlock()
}

// PopBack removes the element from back of <l> and returns the value of the element.
func (l *List) PopBack() (value interface{}) {
	l.mu.Lock()
	if e := l.list.Back(); e != nil {
		value = l.list.Remove(e)
	}
	l.mu.Unlock()
	return
}

// PopFront removes the element from front of <l> and returns the value of the element.
func (l *List) PopFront() (value interface{}) {
	l.mu.Lock()
	if e := l.list.Front(); e != nil {
		value = l.list.Remove(e)
	}
	l.mu.Unlock()
	return
}

// PopBacks removes <max> elements from back of <l>
// and returns values of the removed elements as slice.
func (l *List) PopBacks(max int) (values []interface{}) {
	l.mu.Lock()
	length := l.list.Len()
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]interface{}, length)
		for i := 0; i < length; i++ {
			values[i] = l.list.Remove(l.list.Back())
		}
	}
	l.mu.Unlock()
	return
}

// PopFronts removes <max> elements from front of <l>
// and returns values of the removed elements as slice.
func (l *List) PopFronts(max int) (values []interface{}) {
	l.mu.Lock()
	length := l.list.Len()
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]interface{}, length)
		for i := 0; i < length; i++ {
			values[i] = l.list.Remove(l.list.Front())
		}
	}
	l.mu.Unlock()
	return
}

// PopBackAll removes all elements from back of <l>
// and returns values of the removed elements as slice.
func (l *List) PopBackAll() []interface{} {
	return l.PopBacks(-1)
}

// PopFrontAll removes all elements from front of <l>
// and returns values of the removed elements as slice.
func (l *List) PopFrontAll() []interface{} {
	return l.PopFronts(-1)
}

// FrontAll copies and returns values of all elements from front of <l> as slice.
func (l *List) FrontAll() (values []interface{}) {
	l.mu.RLock()
	length := l.list.Len()
	if length > 0 {
		values = make([]interface{}, length)
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			values[i] = e.Value
		}
	}
	l.mu.RUnlock()
	return
}

// BackAll copies and returns values of all elements from back of <l> as slice.
func (l *List) BackAll() (values []interface{}) {
	l.mu.RLock()
	length := l.list.Len()
	if length > 0 {
		values = make([]interface{}, length)
		for i, e := 0, l.list.Back(); i < length; i, e = i+1, e.Prev() {
			values[i] = e.Value
		}
	}
	l.mu.RUnlock()
	return
}

// FrontValue returns value of the first element of <l> or nil if the list is empty.
func (l *List) FrontValue() (value interface{}) {
	l.mu.RLock()
	if e := l.list.Front(); e != nil {
		value = e.Value
	}
	l.mu.RUnlock()
	return
}

// BackValue returns value of the last element of <l> or nil if the list is empty.
func (l *List) BackValue() (value interface{}) {
	l.mu.RLock()
	if e := l.list.Back(); e != nil {
		value = e.Value
	}
	l.mu.RUnlock()
	return
}

// Front returns the first element of list <l> or nil if the list is empty.
func (l *List) Front() (e *Element) {
	l.mu.RLock()
	e = l.list.Front()
	l.mu.RUnlock()
	return
}

// Back returns the last element of list <l> or nil if the list is empty.
func (l *List) Back() (e *Element) {
	l.mu.RLock()
	e = l.list.Back()
	l.mu.RUnlock()
	return
}

// Len returns the number of elements of list <l>.
// The complexity is O(1).
func (l *List) Len() (length int) {
	l.mu.RLock()
	length = l.list.Len()
	l.mu.RUnlock()
	return
}

// Size is alias of Len.
func (l *List) Size() int {
	return l.Len()
}

// MoveBefore moves element <e> to its new position before <p>.
// If <e> or <p> is not an element of <l>, or <e> == <p>, the list is not modified.
// The element and <p> must not be nil.
func (l *List) MoveBefore(e, p *Element) {
	l.mu.Lock()
	l.list.MoveBefore(e, p)
	l.mu.Unlock()
}

// MoveAfter moves element <e> to its new position after <p>.
// If <e> or <p> is not an element of <l>, or <e> == <p>, the list is not modified.
// The element and <p> must not be nil.
func (l *List) MoveAfter(e, p *Element) {
	l.mu.Lock()
	l.list.MoveAfter(e, p)
	l.mu.Unlock()
}

// MoveToFront moves element <e> to the front of list <l>.
// If <e> is not an element of <l>, the list is not modified.
// The element must not be nil.
func (l *List) MoveToFront(e *Element) {
	l.mu.Lock()
	l.list.MoveToFront(e)
	l.mu.Unlock()
}

// MoveToBack moves element <e> to the back of list <l>.
// If <e> is not an element of <l>, the list is not modified.
// The element must not be nil.
func (l *List) MoveToBack(e *Element) {
	l.mu.Lock()
	l.list.MoveToBack(e)
	l.mu.Unlock()
}

// PushBackList inserts a copy of an other list at the back of list <l>.
// The lists <l> and <other> may be the same, but they must not be nil.
func (l *List) PushBackList(other *List) {
	if l != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	l.mu.Lock()
	l.list.PushBackList(&other.list)
	l.mu.Unlock()
}

// PushFrontList inserts a copy of an other list at the front of list <l>.
// The lists <l> and <other> may be the same, but they must not be nil.
func (l *List) PushFrontList(other *List) {
	if l != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	l.mu.Lock()
	l.list.PushFrontList(&other.list)
	l.mu.Unlock()
}

// InsertAfter inserts a new element <e> with value <v> immediately after <p> and returns <e>.
// If <p> is not an element of <l>, the list is not modified.
// The <p> must not be nil.
func (l *List) InsertAfter(p *Element, v interface{}) (e *Element) {
	l.mu.Lock()
	e = l.list.InsertAfter(v, p)
	l.mu.Unlock()
	return
}

// InsertBefore inserts a new element <e> with value <v> immediately before <p> and returns <e>.
// If <p> is not an element of <l>, the list is not modified.
// The <p> must not be nil.
func (l *List) InsertBefore(p *Element, v interface{}) (e *Element) {
	l.mu.Lock()
	e = l.list.InsertBefore(v, p)
	l.mu.Unlock()
	return
}

// Remove removes <e> from <l> if <e> is an element of list <l>.
// It returns the element value e.Value.
// The element must not be nil.
func (l *List) Remove(e *Element) (value interface{}) {
	l.mu.Lock()
	value = l.list.Remove(e)
	l.mu.Unlock()
	return
}

// Removes removes multiple elements <es> from <l> if <es> are elements of list <l>.
func (l *List) Removes(es []*Element) {
	l.mu.Lock()
	for _, e := range es {
		l.list.Remove(e)
	}
	l.mu.Unlock()
	return
}

// RemoveAll removes all elements from list <l>.
func (l *List) RemoveAll() {
	l.mu.Lock()
	l.list = list.List{}
	l.mu.Unlock()
}

// See RemoveAll().
func (l *List) Clear() {
	l.RemoveAll()
}

// RLockFunc locks reading with given callback function <f> within RWMutex.RLock.
func (l *List) RLockFunc(f func(list *list.List)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	f(&l.list)
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (l *List) LockFunc(f func(list *list.List)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f(&l.list)
}

// Iterator is alias of IteratorAsc.
func (l *List) Iterator(f func(e *Element) bool) {
	l.IteratorAsc(f)
}

// IteratorAsc iterates the list in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (l *List) IteratorAsc(f func(e *Element) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	length := l.list.Len()
	if length > 0 {
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			if !f(e) {
				break
			}
		}
	}
}

// IteratorDesc iterates the list in descending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (l *List) IteratorDesc(f func(e *Element) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	length := l.list.Len()
	if length > 0 {
		for i, e := 0, l.list.Back(); i < length; i, e = i+1, e.Prev() {
			if !f(e) {
				break
			}
		}
	}
}

// Join joins list elements with a string <glue>.
func (l *List) Join(glue string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	buffer := bytes.NewBuffer(nil)
	length := l.list.Len()
	if length > 0 {
		s := ""
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			s = gconv.String(e.Value)
			if gstr.IsNumeric(s) {
				buffer.WriteString(s)
			} else {
				buffer.WriteString(`"` + gstr.QuoteMeta(s, `"\`) + `"`)
			}
			if i != length-1 {
				buffer.WriteString(glue)
			}
		}
	}
	return buffer.String()
}

// String returns current list as a string.
func (l *List) String() string {
	return "[" + l.Join(",") + "]"
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (l *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.FrontAll())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (l *List) UnmarshalJSON(b []byte) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	var array []interface{}
	if err := json.Unmarshal(b, &array); err != nil {
		return err
	}
	l.PushBacks(array)
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for list.
func (l *List) UnmarshalValue(value interface{}) (err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var array []interface{}
	switch value.(type) {
	case string, []byte:
		err = json.Unmarshal(gconv.Bytes(value), &array)
	default:
		array = gconv.SliceAny(value)
	}
	l.PushBacks(array)
	return err
}

<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

// 并发安全的双向链表.
package glist

import (
	"sync"
	"container/list"
)

// 变长双向链表
type List struct {
	mu   sync.RWMutex
	list *list.List
}

// 获得一个变长链表指针
func New() *List {
	return &List{list: list.New()}
}

// 往链表头入栈数据项
func (this *List) PushFront(v interface{}) *list.Element {
	this.mu.Lock()
	e := this.list.PushFront(v)
	this.mu.Unlock()
	return e
}

// 往链表尾入栈数据项
func (this *List) PushBack(v interface{}) *list.Element {
	this.mu.Lock()
	r := this.list.PushBack(v)
	this.mu.Unlock()
	return r
}

// 在list 中元素mark之后插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (this *List) InsertAfter(v interface{}, mark *list.Element) *list.Element {
    this.mu.Lock()
    r := this.list.InsertAfter(v, mark)
    this.mu.Unlock()
    return r
}

// 在list 中元素mark之前插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (this *List) InsertBefore(v interface{}, mark *list.Element) *list.Element {
    this.mu.Lock()
    r := this.list.InsertBefore(v, mark)
    this.mu.Unlock()
    return r
}


// 批量往链表头入栈数据项
func (this *List) BatchPushFront(vs []interface{}) {
	this.mu.Lock()
	for _, item := range vs {
		this.list.PushFront(item)
	}
	this.mu.Unlock()
}

// 从链表尾端出栈数据项(删除)
func (this *List) PopBack() interface{} {
	this.mu.Lock()
	if elem := this.list.Back(); elem != nil {
		item := this.list.Remove(elem)
		this.mu.Unlock()
		return item
	}
	this.mu.Unlock()
	return nil
}

// 从链表头端出栈数据项(删除)
func (this *List) PopFront() interface{} {
	this.mu.Lock()
	if elem := this.list.Front(); elem != nil {
		item := this.list.Remove(elem)
		this.mu.Unlock()
		return item
	}
	this.mu.Unlock()
	return nil
}

// 批量从链表尾端出栈数据项(删除)
func (this *List) BatchPopBack(max int) []interface{} {
	this.mu.Lock()
	count := this.list.Len()
	if count == 0 {
		this.mu.Unlock()
		return []interface{}{}
	}

	if count > max {
		count = max
	}
	items := make([]interface{}, count)
	for i := 0; i < count; i++ {
		items[i] = this.list.Remove(this.list.Back())
	}
	this.mu.Unlock()
	return items
}

// 批量从链表头端出栈数据项(删除)
func (this *List) BatchPopFront(max int) []interface{} {
	this.mu.Lock()
	count := this.list.Len()
	if count == 0 {
		this.mu.Unlock()
		return []interface{}{}
	}

	if count > max {
		count = max
	}
	items := make([]interface{}, count)
	for i := 0; i < count; i++ {
		items[i] = this.list.Remove(this.list.Front())
	}
	this.mu.Unlock()
	return items
}

// 批量从链表尾端依次获取所有数据(删除)
func (this *List) PopBackAll() []interface{} {
	this.mu.Lock()
	count := this.list.Len()
	if count == 0 {
		this.mu.Unlock()
		return []interface{}{}
	}
	items := make([]interface{}, count)
	for i := 0; i < count; i++ {
		items[i] = this.list.Remove(this.list.Back())
	}
	this.mu.Unlock()
	return items
}

// 批量从链表头端依次获取所有数据(删除)
func (this *List) PopFrontAll() []interface{} {
    this.mu.Lock()
    count := this.list.Len()
    if count == 0 {
        this.mu.Unlock()
        return []interface{}{}
    }
    items := make([]interface{}, count)
    for i := 0; i < count; i++ {
        items[i] = this.list.Remove(this.list.Front())
    }
    this.mu.Unlock()
    return items
}

// 删除数据项
func (this *List) Remove(e *list.Element) interface{} {
	this.mu.Lock()
	r := this.list.Remove(e)
	this.mu.Unlock()
	return r
}

// 删除所有数据项
func (this *List) RemoveAll() {
	this.mu.Lock()
	this.list = list.New()
	this.mu.Unlock()
}

// 从链表头获取所有数据(不删除)
func (this *List) FrontAll() []interface{} {
	this.mu.RLock()
	count := this.list.Len()
	if count == 0 {
        this.mu.RUnlock()
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for e := this.list.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}
    this.mu.RUnlock()
	return items
}

// 从链表尾获取所有数据(不删除)
func (this *List) BackAll() []interface{} {
	this.mu.RLock()
	count := this.list.Len()
	if count == 0 {
        this.mu.RUnlock()
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for e := this.list.Back(); e != nil; e = e.Prev() {
		items = append(items, e.Value)
	}
    this.mu.RUnlock()
	return items
}

// 获取链表头值(不删除)
func (this *List) FrontItem() interface{} {
	this.mu.RLock()
	if f := this.list.Front(); f != nil {
		this.mu.RUnlock()
		return f.Value
	}

	this.mu.RUnlock()
	return nil
}

// 获取链表尾值(不删除)
func (this *List) BackItem() interface{} {
    this.mu.RLock()
    if f := this.list.Back(); f != nil {
        this.mu.RUnlock()
        return f.Value
    }

    this.mu.RUnlock()
    return nil
}

// 获取表头指针
func (this *List) Front() *list.Element {
    this.mu.RLock()
    r := this.list.Front()
    this.mu.RUnlock()
    return r
}

// 获取表位指针
func (this *List) Back() *list.Element {
    this.mu.RLock()
    r := this.list.Back()
    this.mu.RUnlock()
    return r
}

// 获取链表长度
func (this *List) Len() int {
	this.mu.RLock()
    length := this.list.Len()
	this.mu.RUnlock()
	return length
}
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with l file,
// You can obtain one at https://github.com/gogf/gf.
//

// Package glist provides a concurrent-safe/unsafe doubly linked list.
package glist

import (
    "container/list"
    "github.com/gogf/gf/g/internal/rwmutex"
)

type (
	List struct {
	    mu   *rwmutex.RWMutex
	    list *list.List
	}

	Element = list.Element
)

// New creates and returns a new empty doubly linked list.
func New(unsafe...bool) *List {
	return &List {
	    mu   : rwmutex.New(unsafe...),
	    list : list.New(),
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
        tempe := (*Element)(nil)
        values = make([]interface{}, length)
        for i := 0; i < length; i++ {
            tempe     = l.list.Back()
            values[i] = l.list.Remove(tempe)
        }
    }
    l.mu.Unlock()
    return
}

// PopFronts removes <max> elements from front of <l>
// and returns values of the removed elements as slice.
func (l *List) PopFronts(max int) (values []interface{}) {
    l.mu.RLock()
    length := l.list.Len()
    if length > 0 {
        if max > 0 && max < length {
            length = max
        }
        tempe := (*Element)(nil)
        values = make([]interface{}, length)
        for i := 0; i < length; i++ {
            tempe     = l.list.Front()
            values[i] = l.list.Remove(tempe)
        }
    }
    l.mu.RUnlock()
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
        for i, e := 0, l.list.Front(); i < length; i, e = i + 1, e.Next() {
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
        for i, e := 0, l.list.Back(); i < length; i, e = i + 1, e.Prev() {
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

// Alias of Len.
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
    l.list.PushBackList(other.list)
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
    l.list.PushFrontList(other.list)
    l.mu.Unlock()
}

// InsertAfter inserts a new element <e> with value <v> immediately after <p> and returns <e>.
// If <p> is not an element of <l>, the list is not modified.
// The <p> must not be nil.
func (l *List) InsertAfter(v interface{}, p *Element) (e *Element) {
    l.mu.Lock()
    e = l.list.InsertAfter(v, p)
    l.mu.Unlock()
    return
}

// InsertBefore inserts a new element <e> with value <v> immediately before <p> and returns <e>.
// If <p> is not an element of <l>, the list is not modified.
// The <p> must not be nil.
func (l *List) InsertBefore(v interface{}, p *Element) (e *Element) {
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
    l.list = list.New()
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
    f(l.list)
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (l *List) LockFunc(f func(list *list.List)) {
    l.mu.Lock()
    defer l.mu.Unlock()
    f(l.list)
}

// Iterator is alias of IteratorAsc.
func (l *List) Iterator(f func (e *Element) bool) {
	l.IteratorAsc(f)
}

// IteratorAsc iterates the list in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (l *List) IteratorAsc(f func (e *Element) bool) {
	l.mu.RLock()
	length := l.list.Len()
	if length > 0 {
		for i, e := 0, l.list.Front(); i < length; i, e = i + 1, e.Next() {
			if !f(e) {
				break
			}
		}
	}
	l.mu.RUnlock()
}

// IteratorDesc iterates the list in descending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (l *List) IteratorDesc(f func (e *Element) bool) {
	l.mu.RLock()
	length := l.list.Len()
	if length > 0 {
		for i, e := 0, l.list.Back(); i < length; i, e = i + 1, e.Prev() {
			if !f(e) {
				break
			}
		}
	}
	l.mu.RUnlock()
}
>>>>>>> upstream/master

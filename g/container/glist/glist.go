// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with l file,
// You can obtain one at https://gitee.com/johng/gf.
//

// Package glist provides a concurrent-safe(alternative) doubly linked list.
// 并发安全的双向链表.
package glist

import (
	"container/list"
	"gitee.com/johng/gf/g/container/internal/rwmutex"
)

// 变长双向链表
type List struct {
	mu   *rwmutex.RWMutex
	list *list.List
}

// 获得一个变长链表指针
func New(safe...bool) *List {
	return &List {
	    mu   : rwmutex.New(safe...),
		list : list.New(),
    }
}

// 往链表头入栈数据项
func (l *List) PushFront(v interface{}) *list.Element {
	l.mu.Lock()
	e := l.list.PushFront(v)
	l.mu.Unlock()
	return e
}

// 往链表尾入栈数据项
func (l *List) PushBack(v interface{}) *list.Element {
	l.mu.Lock()
	r := l.list.PushBack(v)
	l.mu.Unlock()
	return r
}

// 在list 中元素mark之后插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (l *List) InsertAfter(v interface{}, mark *list.Element) *list.Element {
    l.mu.Lock()
    r := l.list.InsertAfter(v, mark)
    l.mu.Unlock()
    return r
}

// 在list 中元素mark之前插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (l *List) InsertBefore(v interface{}, mark *list.Element) *list.Element {
    l.mu.Lock()
    r := l.list.InsertBefore(v, mark)
    l.mu.Unlock()
    return r
}


// 批量往链表头入栈数据项
func (l *List) BatchPushFront(vs []interface{}) {
	l.mu.Lock()
	for _, item := range vs {
		l.list.PushFront(item)
	}
	l.mu.Unlock()
}

// 从链表尾端出栈数据项(删除)
func (l *List) PopBack() interface{} {
	l.mu.Lock()
	if elem := l.list.Back(); elem != nil {
		item := l.list.Remove(elem)
		l.mu.Unlock()
		return item
	}
	l.mu.Unlock()
	return nil
}

// 从链表头端出栈数据项(删除)
func (l *List) PopFront() interface{} {
	l.mu.Lock()
	if elem := l.list.Front(); elem != nil {
		item := l.list.Remove(elem)
		l.mu.Unlock()
		return item
	}
	l.mu.Unlock()
	return nil
}

// 批量从链表尾端出栈数据项(删除)
func (l *List) BatchPopBack(max int) []interface{} {
	l.mu.Lock()
	count := l.list.Len()
	if count == 0 {
		l.mu.Unlock()
		return []interface{}{}
	}

	if count > max {
		count = max
	}
	items := make([]interface{}, count)
	for i := 0; i < count; i++ {
		items[i] = l.list.Remove(l.list.Back())
	}
	l.mu.Unlock()
	return items
}

// 批量从链表头端出栈数据项(删除)
func (l *List) BatchPopFront(max int) []interface{} {
	l.mu.Lock()
	count := l.list.Len()
	if count == 0 {
		l.mu.Unlock()
		return []interface{}{}
	}

	if count > max {
		count = max
	}
	items := make([]interface{}, count)
	for i := 0; i < count; i++ {
		items[i] = l.list.Remove(l.list.Front())
	}
	l.mu.Unlock()
	return items
}

// 批量从链表尾端依次获取所有数据(删除)
func (l *List) PopBackAll() []interface{} {
	l.mu.Lock()
	count := l.list.Len()
	if count == 0 {
		l.mu.Unlock()
		return []interface{}{}
	}
	items := make([]interface{}, count)
	for i := 0; i < count; i++ {
		items[i] = l.list.Remove(l.list.Back())
	}
	l.mu.Unlock()
	return items
}

// 批量从链表头端依次获取所有数据(删除)
func (l *List) PopFrontAll() []interface{} {
    l.mu.Lock()
    count := l.list.Len()
    if count == 0 {
        l.mu.Unlock()
        return []interface{}{}
    }
    items := make([]interface{}, count)
    for i := 0; i < count; i++ {
        items[i] = l.list.Remove(l.list.Front())
    }
    l.mu.Unlock()
    return items
}

// 删除数据项
func (l *List) Remove(e *list.Element) interface{} {
	l.mu.Lock()
	r := l.list.Remove(e)
	l.mu.Unlock()
	return r
}

// 删除所有数据项
func (l *List) RemoveAll() {
	l.mu.Lock()
	l.list = list.New()
	l.mu.Unlock()
}

// 从链表头获取所有数据(不删除)
func (l *List) FrontAll() []interface{} {
	l.mu.RLock()
	count := l.list.Len()
	if count == 0 {
        l.mu.RUnlock()
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for e := l.list.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}
    l.mu.RUnlock()
	return items
}

// 从链表尾获取所有数据(不删除)
func (l *List) BackAll() []interface{} {
	l.mu.RLock()
	count := l.list.Len()
	if count == 0 {
        l.mu.RUnlock()
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for e := l.list.Back(); e != nil; e = e.Prev() {
		items = append(items, e.Value)
	}
    l.mu.RUnlock()
	return items
}

// 获取链表头值(不删除)
func (l *List) FrontItem() interface{} {
	l.mu.RLock()
	if f := l.list.Front(); f != nil {
		l.mu.RUnlock()
		return f.Value
	}

	l.mu.RUnlock()
	return nil
}

// 获取链表尾值(不删除)
func (l *List) BackItem() interface{} {
    l.mu.RLock()
    if f := l.list.Back(); f != nil {
        l.mu.RUnlock()
        return f.Value
    }

    l.mu.RUnlock()
    return nil
}

// 获取表头指针
func (l *List) Front() *list.Element {
    l.mu.RLock()
    r := l.list.Front()
    l.mu.RUnlock()
    return r
}

// 获取表位指针
func (l *List) Back() *list.Element {
    l.mu.RLock()
    r := l.list.Back()
    l.mu.RUnlock()
    return r
}

// 获取链表长度
func (l *List) Len() int {
	l.mu.RLock()
    length := l.list.Len()
	l.mu.RUnlock()
	return length
}

// 读锁操作
func (l *List) RLockFunc(f func(l *list.List)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	f(l.list)
}

// 写锁操作
func (l *List) LockFunc(f func(l *list.List)) {
    l.mu.Lock()
    defer l.mu.Unlock()
    f(l.list)
}

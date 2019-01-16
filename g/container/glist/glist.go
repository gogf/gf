// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with l file,
// You can obtain one at https://gitee.com/johng/gf.
//

// Package glist provides a concurrent-safe(alternative) doubly linked list.
//
// 并发安全的双向链表.
package glist

import (
    "container/list"
    "gitee.com/johng/gf/g/internal/rwmutex"
)

// 变长双向链表
type List struct {
    mu     *rwmutex.RWMutex
    list   *list.List
}

type Element = list.Element


// 获得一个变长链表指针
func New(unsafe...bool) *List {
	return &List {
	    mu   : rwmutex.New(unsafe...),
	    list : list.New(),
    }
}

// 往链表头入栈数据项
func (l *List) PushFront(v interface{}) (e *Element) {
    l.mu.Lock()
    e = l.list.PushFront(v)
    l.mu.Unlock()
    return
}

// 往链表尾入栈数据项
func (l *List) PushBack(v interface{}) (e *Element) {
    l.mu.Lock()
    e = l.list.PushBack(v)
    l.mu.Unlock()
    return
}

// 批量往链表头入栈数据项
func (l *List) BatchPushFront(values []interface{}) {
    l.mu.Lock()
	for _, v := range values {
        l.list.PushFront(v)
	}
    l.mu.Unlock()
}

// 批量往链表尾入栈数据项
func (l *List) BatchPushBack(values []interface{}) {
    l.mu.Lock()
    for _, v := range values {
        l.list.PushBack(v)
    }
    l.mu.Unlock()
}

// 从链表尾端出栈数据项(删除)
func (l *List) PopBack() (value interface{}) {
    l.mu.Lock()
	if e := l.list.Back(); e != nil {
        value = l.list.Remove(e)
	}
    l.mu.Unlock()
	return
}

// 从链表头端出栈数据项(删除)
func (l *List) PopFront() (value interface{}) {
    l.mu.Lock()
    if e := l.list.Front(); e != nil {
        value = l.list.Remove(e)
    }
    l.mu.Unlock()
    return
}

// 批量从链表尾端出栈数据项(删除)
func (l *List) BatchPopBack(max int) (values []interface{}) {
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

// 批量从链表头端出栈数据项(删除)
func (l *List) BatchPopFront(max int) (values []interface{}) {
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

// 批量从链表尾端依次获取所有数据(删除)
func (l *List) PopBackAll() []interface{} {
	return l.BatchPopBack(-1)
}

// 批量从链表头端依次获取所有数据(删除)
func (l *List) PopFrontAll() []interface{} {
    return l.BatchPopFront(-1)
}

// 从链表头获取所有数据(不删除)
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

// 从链表尾获取所有数据(不删除)
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

// 获取链表头值(不删除)
func (l *List) FrontItem() (value interface{}) {
    l.mu.RLock()
    if e := l.list.Front(); e != nil {
        value = e.Value
    }
    l.mu.RUnlock()
    return
}

// 获取链表尾值(不删除)
func (l *List) BackItem() (value interface{}) {
    l.mu.RLock()
    if e := l.list.Back(); e != nil {
        value = e.Value
    }
    l.mu.RUnlock()
    return
}

// 获取表头指针
func (l *List) Front() (e *Element) {
    l.mu.RLock()
    e = l.list.Front()
    l.mu.RUnlock()
    return
}

// 获取表位指针
func (l *List) Back() (e *Element) {
    l.mu.RLock()
    e = l.list.Back()
    l.mu.RUnlock()
    return
}

// 获取链表长度
func (l *List) Len() (length int) {
    l.mu.RLock()
    length = l.list.Len()
    l.mu.RUnlock()
	return
}

func (l *List) MoveBefore(e, p *Element) {
    l.mu.Lock()
    l.list.MoveBefore(e, p)
    l.mu.Unlock()
}

func (l *List) MoveAfter(e, p *Element) {
    l.mu.Lock()
    l.list.MoveAfter(e, p)
    l.mu.Unlock()
}

func (l *List) MoveToFront(e *Element) {
    l.mu.Lock()
    l.list.MoveToFront(e)
    l.mu.Unlock()
}

func (l *List) MoveToBack(e *Element) {
    l.mu.Lock()
    l.list.MoveToBack(e)
    l.mu.Unlock()
}

func (l *List) PushBackList(other *List) {
    if l != other {
        other.mu.RLock()
        defer other.mu.RUnlock()
    }
    l.mu.Lock()
    l.list.PushBackList(other.list)
    l.mu.Unlock()
}

func (l *List) PushFrontList(other *List) {
    if l != other {
        other.mu.RLock()
        defer other.mu.RUnlock()
    }
    l.mu.Lock()
    l.list.PushFrontList(other.list)
    l.mu.Unlock()
}

// 在list中元素项p之后插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (l *List) InsertAfter(v interface{}, p *Element) (e *Element) {
    l.mu.Lock()
    e = l.list.InsertAfter(v, p)
    l.mu.Unlock()
    return
}

// 在list中元素项p之前插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (l *List) InsertBefore(v interface{}, p *Element) (e *Element) {
    l.mu.Lock()
    e = l.list.InsertBefore(v, p)
    l.mu.Unlock()
    return
}

// 删除数据项e, 并返回删除项的元素项
func (l *List) Remove(e *Element) (value interface{}) {
    l.mu.Lock()
    value = l.list.Remove(e)
    l.mu.Unlock()
    return
}

// 批量删除数据项
func (l *List) BatchRemove(es []*Element) {
    l.mu.Lock()
    for _, e := range es {
        l.list.Remove(e)
    }
    l.mu.Unlock()
    return
}

// 删除所有数据项
func (l *List) RemoveAll() {
    l.mu.Lock()
    l.list = list.New()
    l.mu.Unlock()
}

// 读锁操作
func (l *List) RLockFunc(f func(list *list.List)) {
    l.mu.RLock()
    defer l.mu.RUnlock()
    f(l.list)
}

// 写锁操作
func (l *List) LockFunc(f func(list *list.List)) {
    l.mu.Lock()
    defer l.mu.Unlock()
    f(l.list)
}
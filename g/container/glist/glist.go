// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with l file,
// You can obtain one at https://gitee.com/johng/gf.
//

// Package glist provides a concurrent-safe(alternative) doubly linked list/并发安全的双向链表.
package glist

import (
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/container/internal/rwmutex"
)

// 变长双向链表
type List struct {
    mu     *rwmutex.RWMutex
    root   *Element
    length *gtype.Int
}

// 获得一个变长链表指针
func New(safe...bool) *List {
	l := &List{
	    mu     : rwmutex.New(safe...),
	    length : gtype.NewInt(),
    }
	l.root      = newElement(nil, l, safe...)
	l.root.list = l
    l.root.next = l.root
    l.root.prev = l.root
    return l
}

// 往链表头入栈数据项
func (l *List) PushFront(v interface{}) *Element {
    return l.InsertAfter(v, l.root)
}

// 往链表尾入栈数据项
func (l *List) PushBack(v interface{}) *Element {
    return l.InsertBefore(v, l.root)
}

// 批量往链表头入栈数据项
func (l *List) BatchPushFront(values []interface{}) {
    l.mu.Lock()
    defer l.mu.Unlock()
	for _, v := range values {
        l.InsertAfter(v, l.root)
	}
}

// 批量往链表尾入栈数据项
func (l *List) BatchPushBack(values []interface{}) {
    l.mu.Lock()
    defer l.mu.Unlock()
    for _, v := range values {
        l.InsertBefore(v, l.root)
    }
}

// 从链表尾端出栈数据项(删除)
func (l *List) PopBack() interface{} {
	if e := l.Back(); e != nil {
		if o := l.Remove(e); o != nil {
		    return o.Value()
        }
	}
	return nil
}

// 从链表头端出栈数据项(删除)
func (l *List) PopFront() interface{} {
	if e := l.Front(); e != nil {
        if o := l.Remove(e); o != nil {
            return o.Value()
        }
	}
	return nil
}

// 批量从链表尾端出栈数据项(删除)
func (l *List) BatchPopBack(max int) []interface{} {
	count := l.Len()
	if count == 0 {
		return []interface{}{}
	}
	if count > max {
		count = max
	}
	items := make([]interface{}, count)
	for i := 0; i < count; i++ {
		items[i] = l.PopBack()
	}
	return items
}

// 批量从链表头端出栈数据项(删除)
func (l *List) BatchPopFront(max int) []interface{} {
	count := l.Len()
	if count == 0 {
		return []interface{}{}
	}
	if count > max {
		count = max
	}
	items := make([]interface{}, count)
	for i := 0; i < count; i++ {
		items[i] = l.PopFront()
	}
	return items
}

// 批量从链表尾端依次获取所有数据(删除)
func (l *List) PopBackAll() []interface{} {
	return l.BatchPopFront(l.Len())
}

// 批量从链表头端依次获取所有数据(删除)
func (l *List) PopFrontAll() []interface{} {
    return l.BatchPopFront(l.Len())
}

// 从链表头获取所有数据(不删除)
func (l *List) FrontAll() []interface{} {
	count := l.Len()
	if count == 0 {
		return nil
	}
	items := make([]interface{}, 0, count)
	for e := l.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value())
	}
	return items
}

// 从链表尾获取所有数据(不删除)
func (l *List) BackAll() []interface{} {
	count := l.Len()
	if count == 0 {
		return nil
	}
	items := make([]interface{}, 0, count)
	for e := l.Back(); e != nil; e = e.Prev() {
		items = append(items, e.Value())
	}
	return items
}

// 获取链表头值(不删除)
func (l *List) FrontItem() interface{} {
	if e := l.Front(); e != nil {
		return e.Value()
	}
	return nil
}

// 获取链表尾值(不删除)
func (l *List) BackItem() interface{} {
    if e := l.Back(); e != nil {
        return e.Value()
    }
    return nil
}

// 获取表头指针
func (l *List) Front() *Element {
    if l.length.Val() == 0 {
        return nil
    }
    return l.root.getNext()
}

// 获取表位指针
func (l *List) Back() *Element {
    if l.length.Val() == 0 {
        return nil
    }
    return l.root.getPrev()
}

// 获取链表长度
func (l *List) Len() int {
	return l.length.Val()
}

func (l *List) MoveBefore(e, p *Element) {
    if e.getList() != l || p.getList() != l || e == p {
        return
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    l.doInsertElementBefore(l.doRemove(e), p)
}

func (l *List) MoveAfter(e, p *Element) {
    if e.getList() != l || p.getList() != l || e == p {
        return
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    l.doInsertElementAfter(l.doRemove(e), p)
}

func (l *List) MoveToFront(e *Element) {
    if e.getList() != l {
        return
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    l.doInsertElementAfter(l.doRemove(e), l.root)
}

func (l *List) MoveToBack(e *Element) {
    if e.getList() != l {
        return
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    l.doInsertElementBefore(l.doRemove(e), l.root)
}

func (l *List) PushBackList(other *List) {
    if other.Len() == 0 {
        return
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    for i, e := other.Len(), other.Front(); i > 0; i, e = i - 1, e.Next() {
        l.doInsertBefore(e.Value(), l.root)
    }
}

func (l *List) PushFrontList(other *List) {
    if other.Len() == 0 {
        return
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    for i, e := other.Len(), other.Back(); i > 0; i, e = i - 1, e.Prev() {
        l.doInsertAfter(e.Value(), l.root)
    }
}


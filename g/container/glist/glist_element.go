// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with l file,
// You can obtain one at https://gitee.com/johng/gf.
//

package glist

import (
    "gitee.com/johng/gf/g/container/internal/rwmutex"
)

// 链表元素项
type Element struct {
	mu    *rwmutex.RWMutex
    list  *List
    prev  *Element
    next  *Element
    value interface{}
}

// 创建一个并发安全的列表元素项
func newElement(value interface{}, list *List, safe...bool) *Element {
	return &Element {
	    mu    : rwmutex.New(safe...),
		value : value,
	    list  : list,
    }
}

// 获得元素项值
func (e *Element) Value() (v interface{}) {
	e.mu.RLock()
	v = e.value
	e.mu.RUnlock()
	return
}

// 获得下一个元素项(遍历使用)
func (e *Element) Next() *Element {
    e.mu.RLock()
    if p := e.next; p != e.list.root {
        e.mu.RUnlock()
        return p
    }
    e.mu.RUnlock()
    return nil
}

// 获得前一个元素项(遍历使用)
func (e *Element) Prev() *Element {
    e.mu.RLock()
    if p := e.prev; p != e.list.root {
        e.mu.RUnlock()
        return p
    }
    e.mu.RUnlock()
    return nil
}

// 只读锁操作
func (e *Element) RLockFunc(f func(e *Element)) {
    e.mu.RLock()
    defer e.mu.RUnlock()
    f(e)
}

// 写锁操作
func (e *Element) LockFunc(f func(e *Element)) {
    e.mu.Lock()
    defer e.mu.Unlock()
    f(e)
}

func (e *Element) setPrev(prev *Element) (old *Element) {
    e.mu.Lock()
    old    = e.prev
    e.prev = prev
    e.mu.Unlock()
    return
}

func (e *Element) setNext(next *Element) (old *Element) {
    e.mu.Lock()
    old    = e.next
    e.next = next
    e.mu.Unlock()
    return
}

func (e *Element) setList(list *List) {
    e.mu.Lock()
    e.list = list
    e.mu.Unlock()
}

// 获得前一个元素项(内部并发安全使用)
func (e *Element) getPrev() (prev *Element) {
    e.mu.RLock()
    prev = e.prev
    e.mu.RUnlock()
    return
}

// 获得下一个元素项(内部并发安全使用)
func (e *Element) getNext() (next *Element) {
    e.mu.RLock()
    next = e.next
    e.mu.RUnlock()
    return
}

func (e *Element) getList() (list *List) {
    e.mu.RLock()
    list = e.list
    e.mu.RUnlock()
    return
}
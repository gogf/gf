// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with l file,
// You can obtain one at https://gitee.com/johng/gf.

package glist

import "gitee.com/johng/gf/g/container/internal/rwmutex"

// 在list中元素项p之后插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (l *List) InsertAfter(v interface{}, p *Element) *Element {
    if p.getList() != l {
        return nil
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.doInsertAfter(v, p)
}

// 在list中元素项p之前插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (l *List) InsertBefore(v interface{}, p *Element) *Element {
    if p.getList() != l {
        return nil
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.doInsertBefore(v, p)
}

// 在元素项p后添加元素项e, 注意这里的p和e都需要加锁，以保证并发安全性
func (l *List) InsertElementAfter(e, p *Element) *Element {
    if p.getList() != l {
        return nil
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.doInsertElementAfter(e, p)
}

// 在元素项p前添加元素项e, 注意这里的p和e都需要加锁，以保证并发安全性
func (l *List) InsertElementBefore(e, p *Element) *Element {
    if p.getList() != l {
        return nil
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.doInsertElementBefore(e, p)
}

func (l *List) doInsertAfter(v interface{}, p *Element) *Element {
    n := p.getNext()
    e := &Element {
        mu    : rwmutex.New(l.mu.IsSafe()),
        value : v,
        prev  : p,
        next  : n,
        list  : l,
    }
    p.setNext(e)
    n.setPrev(e)
    l.length.Add(1)
    return e
}

func (l *List) doInsertBefore(v interface{}, p *Element) *Element {
    n := p.getPrev()
    e := &Element {
        mu    : rwmutex.New(l.mu.IsSafe()),
        value : v,
        prev  : n,
        next  : p,
        list  : l,
    }
    p.setPrev(e)
    n.setNext(e)
    l.length.Add(1)
    return e
}

func (l *List) doInsertElementAfter(e, p *Element) *Element {
    o := p.setNext(e)
    o.setPrev(e)
    e.mu.Lock()
    e.prev = p
    e.next = o
    e.list = l
    e.mu.Unlock()
    l.length.Add(1)
    return e
}

func (l *List) doInsertElementBefore(e, p *Element) *Element {
    o := p.setPrev(e)
    o.setNext(e)
    e.mu.Lock()
    e.prev = o
    e.next = p
    e.list = l
    e.mu.Unlock()
    l.length.Add(1)
    return e
}
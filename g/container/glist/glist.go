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
)

// 变长双向链表
type List struct {
    safe   bool
    root   *Element
    length *gtype.Int
}

// 获得一个变长链表指针
func New(safe...bool) *List {
	l := &List{
	    length : gtype.NewInt(),
    }
	l.root      = newElement(nil, safe...)
	l.root.list = l
    l.root.next = l.root
    l.root.prev = l.root
    if len(safe) > 0 {
        l.safe = safe[0]
    } else {
        l.safe = true
    }
    return l
}

// 往链表头入栈数据项
func (l *List) PushFront(v interface{}) *Element {
    return l.insertValue(v, l.root)
}

// 往链表尾入栈数据项
func (l *List) PushBack(v interface{}) *Element {
    return l.insertValue(v, l.root.prev)
}

// 在list 中元素mark之后插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (l *List) InsertAfter(v interface{}, p *Element) *Element {
    if p.checkList(l) == false {
        return nil
    }
    return l.insertValue(v, p)
}

// 在list 中元素mark之前插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (l *List) InsertBefore(v interface{}, p *Element) *Element {
    if p.checkList(l) == false {
        return nil
    }
    return l.insertValue(v, p.getPrev())
}

// 批量往链表头入栈数据项
func (l *List) BatchPushFront(vs []interface{}) {
	for _, item := range vs {
		l.PushFront(item)
	}
}

// 从链表尾端出栈数据项(删除)
func (l *List) PopBack() interface{} {
	if e := l.Back(); e != nil {
		return l.Remove(e)
	}
	return nil
}

// 从链表头端出栈数据项(删除)
func (l *List) PopFront() interface{} {
	if e := l.Front(); e != nil {
		return l.Remove(e)
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

// 删除数据项e
func (l *List) Remove(e *Element) interface{} {
    if e.list == l {
        l.remove(e)
    }
    return e.Value()
}

// 删除所有数据项
func (l *List) RemoveAll() {
    l.length.Set(0)
    l.root.setNext(l.root)
    l.root.setPrev(l.root)
}

// 从链表头获取所有数据(不删除)
func (l *List) FrontAll() []interface{} {
	count := l.Len()
	if count == 0 {
		return []interface{}{}
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
		return []interface{}{}
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
    return l.root.next
}

// 获取表位指针
func (l *List) Back() *Element {
    if l.length.Val() == 0 {
        return nil
    }
    return l.root.prev
}

// 获取链表长度
func (l *List) Len() int {
	return l.length.Val()
}

// 内部使用insertValue而非insert, 是为去掉insert方法中对元素项e的并发安全保护，提高并发执行效率
func (l *List) MoveBefore(e, p *Element) {
    if e.checkList(l) == false || e == p || p.checkList(l) == false {
        return
    }
    l.insert(l.remove(e), p.getPrev())
}

func (l *List) MoveAfter(e, p *Element) {
    if e.checkList(l) == false || e == p || p.checkList(l) == false {
        return
    }
    l.insert(l.remove(e), p)
}

func (l *List) MoveToFront(e *Element) {
    if e.checkList(l) == false || l.root.next == e {
        return
    }
    l.insert(l.remove(e), l.root)
}

func (l *List) MoveToBack(e *Element) {
    if e.checkList(l) == false || l.root.prev == e {
        return
    }
    l.insert(l.remove(e), l.root.prev)
}

func (l *List) PushBackList(other *List) {
    for i, e := other.Len(), other.Front(); i > 0; i, e = i - 1, e.Next() {
        l.insertValue(e.Value(), l.root.prev)
    }
}

func (l *List) PushFrontList(other *List) {
    for i, e := other.Len(), other.Back(); i > 0; i, e = i - 1, e.Prev() {
        l.insertValue(e.Value(), l.root)
    }
}

// 在元素项p后添加数值value, 自动创建元素项
func (l *List) insertValue(value interface{}, p *Element) *Element {
    return l.insert(newElement(value, l.safe), p)
}

// 在元素项p后添加**新增**元素项e, 注意这里的e不需要加锁
func (l *List) insert(e, p *Element) *Element {
    n := p.setNext(e)
    n.setPrev(e)
    e.LockFunc(func(e *Element) {
        e.prev = p
        e.next = n
        e.list = l
    })
    l.length.Add(1)
    return e
}

// 从列表中删除元素项e
func (l *List) remove(e *Element) *Element {
    e.RLockFunc(func(e *Element) {
        e.prev.setNext(e.next)
        e.next.setPrev(e.prev)
    })
    e.LockFunc(func(e *Element) {
        e.next = nil
        e.prev = nil
        e.list = nil
    })
    l.length.Add(-1)
    return e
}

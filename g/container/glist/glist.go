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

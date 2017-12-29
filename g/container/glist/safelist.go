// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//
//
package glist

import (
	"container/list"
	"sync"
)

// 变长双向链表
type SafeList struct {
	sync.RWMutex
	L *list.List
}

// 获得一个变长链表指针
func NewSafeList() *SafeList {
	return &SafeList{L: list.New()}
}

// 往链表头入栈数据项
func (this *SafeList) PushFront(v interface{}) *list.Element {
	this.Lock()
	e := this.L.PushFront(v)
	this.Unlock()
	return e
}

// 往链表尾入栈数据项
func (this *SafeList) PushBack(v interface{}) *list.Element {
	this.Lock()
	r := this.L.PushBack(v)
	this.Unlock()
	return r
}

// 在list 中元素mark之后插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (this *SafeList) InsertAfter(v interface{}, mark *list.Element) *list.Element {
    this.Lock()
    r := this.L.InsertAfter(v, mark)
    this.Unlock()
    return r
}

// 在list 中元素mark之前插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
func (this *SafeList) InsertBefore(v interface{}, mark *list.Element) *list.Element {
    this.Lock()
    r := this.L.InsertBefore(v, mark)
    this.Unlock()
    return r
}


// 批量往链表头入栈数据项
func (this *SafeList) BatchPushFront(vs []interface{}) {
	this.Lock()
	for _, item := range vs {
		this.L.PushFront(item)
	}
	this.Unlock()
}

// 从链表尾端出栈数据项(删除)
func (this *SafeList) PopBack() interface{} {
	this.Lock()
	if elem := this.L.Back(); elem != nil {
		item := this.L.Remove(elem)
		this.Unlock()
		return item
	}
	this.Unlock()
	return nil
}

// 批量从链表尾端出栈数据项(删除)
func (this *SafeList) BatchPopBack(max int) []interface{} {
	this.Lock()
	count := this.L.Len()
	if count == 0 {
		this.Unlock()
		return []interface{}{}
	}

	if count > max {
		count = max
	}
	items := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		item := this.L.Remove(this.L.Back())
		items = append(items, item)
	}
	this.Unlock()
	return items
}

// 批量从链表尾端依次获取所有数据
func (this *SafeList) PopBackAll() []interface{} {
	this.Lock()

	count := this.L.Len()
	if count == 0 {
		this.Unlock()
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		item := this.L.Remove(this.L.Back())
		items = append(items, item)
	}

	this.Unlock()
	return items
}

// 删除数据项
func (this *SafeList) Remove(e *list.Element) interface{} {
	this.Lock()
	r := this.L.Remove(e)
	this.Unlock()
	return r
}

// 删除所有数据项
func (this *SafeList) RemoveAll() {
	this.Lock()
	this.L = list.New()
	this.Unlock()
}

// 从链表头获取所有数据(不删除)
func (this *SafeList) FrontAll() []interface{} {
	this.RLock()
	count := this.L.Len()
	if count == 0 {
        this.RUnlock()
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for e := this.L.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}
    this.RUnlock()
	return items
}

// 从链表尾获取所有数据(不删除)
func (this *SafeList) BackAll() []interface{} {
	this.RLock()
	count := this.L.Len()
	if count == 0 {
        this.RUnlock()
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for e := this.L.Back(); e != nil; e = e.Prev() {
		items = append(items, e.Value)
	}
    this.RUnlock()
	return items
}

// 获取链表头值(不删除)
func (this *SafeList) FrontItem() interface{} {
	this.RLock()
	if f := this.L.Front(); f != nil {
		this.RUnlock()
		return f.Value
	}

	this.RUnlock()
	return nil
}

// 获取链表尾值(不删除)
func (this *SafeList) BackItem() interface{} {
    this.RLock()
    if f := this.L.Back(); f != nil {
        this.RUnlock()
        return f.Value
    }

    this.RUnlock()
    return nil
}

// 获取表头指针
func (this *SafeList) Front() *list.Element {
    this.RLock()
    r := this.L.Front()
    this.RUnlock()
    return r
}

// 获取表位指针
func (this *SafeList) Back() *list.Element {
    this.RLock()
    r := this.L.Back()
    this.RUnlock()
    return r
}

// 获取链表长度
func (this *SafeList) Len() int {
	this.RLock()
    length := this.L.Len()
	this.RUnlock()
	return length
}

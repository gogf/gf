package glist

import (
	"container/list"
	"sync"
)

// 变长链表
type SafeList struct {
	m sync.RWMutex
	L *list.List
}

// 获得一个变长链表指针
func NewSafeList() *SafeList {
	return &SafeList{L: list.New()}
}

// 往链表头入栈数据项
func (this *SafeList) PushFront(v interface{}) *list.Element {
	this.m.Lock()
	e := this.L.PushFront(v)
	this.m.Unlock()
	return e
}

// 批量往链表头入栈数据项
func (this *SafeList) BatchPushFront(vs []interface{}) {
	this.m.Lock()
	for _, item := range vs {
		this.L.PushFront(item)
	}
	this.m.Unlock()
}

// 从链表尾端出栈数据项
func (this *SafeList) PopBack() interface{} {
	this.m.Lock()
	if elem := this.L.Back(); elem != nil {
		item := this.L.Remove(elem)
		this.m.Unlock()
		return item
	}
	this.m.Unlock()
	return nil
}

// 批量从链表尾端出栈数据项
func (this *SafeList) BatchPopBack(max int) []interface{} {
	this.m.Lock()
	count := this.L.Len()
	if count == 0 {
		this.m.Unlock()
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
	this.m.Unlock()
	return items
}

// 批量从链表尾端依次获取所有数据
func (this *SafeList) PopBackAll() []interface{} {
	this.m.Lock()

	count := this.L.Len()
	if count == 0 {
		this.m.Unlock()
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		item := this.L.Remove(this.L.Back())
		items = append(items, item)
	}

	this.m.Unlock()
	return items
}

// (查找并)删除数据项
func (this *SafeList) Remove(e *list.Element) interface{} {
	this.m.Lock()
	defer this.m.Unlock()
	return this.L.Remove(e)
}

// 删除所有数据项
func (this *SafeList) RemoveAll() {
	this.m.Lock()
	this.L = list.New()
	this.m.Unlock()
}

// 从链表头获取所有数据(不删除)
func (this *SafeList) FrontAll() []interface{} {
	this.m.RLock()
	defer this.m.RUnlock()

	count := this.L.Len()
	if count == 0 {
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for e := this.L.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}
	return items
}

// 从链表尾获取所有数据(不删除)
func (this *SafeList) BackAll() []interface{} {
	this.m.RLock()
	count := this.L.Len()
	if count == 0 {
        this.m.RUnlock()
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for e := this.L.Back(); e != nil; e = e.Prev() {
		items = append(items, e.Value)
	}
    this.m.RUnlock()
	return items
}

// 获取链表头值(不删除)
func (this *SafeList) Front() interface{} {
	this.m.RLock()
	if f := this.L.Front(); f != nil {
		this.m.RUnlock()
		return f.Value
	}

	this.m.RUnlock()
	return nil
}

// 获取链表长度
func (this *SafeList) Len() int {
	this.m.RLock()
    length := this.L.Len()
	this.m.RUnlock()
	return length
}


// 固定长度的链表
type SafeListLimited struct {
	maxSize int
	SL      *SafeList
}

func NewSafeListLimited(maxSize int) *SafeListLimited {
	return &SafeListLimited{SL: NewSafeList(), maxSize: maxSize}
}

func (this *SafeListLimited) PopBack() interface{} {
	return this.SL.PopBack()
}

func (this *SafeListLimited) PopBackBy(max int) []interface{} {
	return this.SL.BatchPopBack(max)
}

func (this *SafeListLimited) PushFront(v interface{}) bool {
	if this.SL.Len() >= this.maxSize {
		return false
	}

	this.SL.PushFront(v)
	return true
}

func (this *SafeListLimited) PushFrontBatch(vs []interface{}) bool {
	if this.SL.Len() >= this.maxSize {
		return false
	}

	this.SL.BatchPushFront(vs)
	return true
}

func (this *SafeListLimited) PushFrontViolently(v interface{}) bool {
	this.SL.PushFront(v)
	if this.SL.Len() > this.maxSize {
		this.SL.PopBack()
	}

	return true
}

func (this *SafeListLimited) RemoveAll() {
	this.SL.RemoveAll()
}

func (this *SafeListLimited) Front() interface{} {
	return this.SL.Front()
}

func (this *SafeListLimited) FrontAll() []interface{} {
	return this.SL.FrontAll()
}

func (this *SafeListLimited) Len() int {
	return this.SL.Len()
}

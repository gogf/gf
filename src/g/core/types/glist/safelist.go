package glist

import (
	"container/list"
	"sync"
)

// 变长链表
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

// 批量往链表头入栈数据项
func (this *SafeList) BatchPushFront(vs []interface{}) {
	this.Lock()
	for _, item := range vs {
		this.L.PushFront(item)
	}
	this.Unlock()
}

// 从链表尾端出栈数据项
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

// 批量从链表尾端出栈数据项
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

// (查找并)删除数据项
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
func (this *SafeList) Front() interface{} {
	this.RLock()
	if f := this.L.Front(); f != nil {
		this.RUnlock()
		return f.Value
	}

	this.RUnlock()
	return nil
}

// 获取链表长度
func (this *SafeList) Len() int {
	this.RLock()
    length := this.L.Len()
	this.RUnlock()
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

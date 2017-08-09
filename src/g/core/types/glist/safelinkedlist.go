package glist

import (
	"container/list"
	"sync"
)

type SafeLinkedList struct {
	sync.RWMutex
	L *list.List
}

func NewSafeLinkedList() *SafeLinkedList {
	return &SafeLinkedList{L: list.New()}
}

func (this *SafeLinkedList) PopBack(max int) []interface{} {
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

func (this *SafeLinkedList) PushFront(v interface{}) *list.Element {
	this.Lock()
	r := this.L.PushFront(v)
	this.Unlock()
	return r
}

func (this *SafeLinkedList) Front() *list.Element {
	this.RLock()
	r := this.L.Front()
	this.RUnlock()
	return r
}

// TODO 异步rcu锁
func (this *SafeLinkedList) Len() int {
	this.RLock()
	r := this.L.Len()
	this.RUnlock()
	return r
}

// SafeLinkedList with Limited Size
type SafeLinkedListLimited struct {
	MaxSize int
	SL      *SafeLinkedList
}

func NewSafeLinkedListLimited(maxSize int) *SafeLinkedListLimited {
	return &SafeLinkedListLimited{SL: NewSafeLinkedList(), MaxSize: maxSize}
}

func (this *SafeLinkedListLimited) PopBack(max int) []interface{} {
	return this.SL.PopBack(max)
}

func (this *SafeLinkedListLimited) PushFront(v interface{}) bool {
	if this.SL.Len() >= this.MaxSize {
		return false
	}

	this.SL.PushFront(v)
	return true
}

func (this *SafeLinkedListLimited) Front() *list.Element {
	return this.SL.Front()
}

func (this *SafeLinkedListLimited) Len() int {
	return this.SL.Len()
}

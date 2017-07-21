package list

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
	defer this.Unlock()

	count := this.L.Len()
	if count == 0 {
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

	return items
}

func (this *SafeLinkedList) PushFront(v interface{}) *list.Element {
	this.Lock()
	defer this.Unlock()
	return this.L.PushFront(v)
}

func (this *SafeLinkedList) Front() *list.Element {
	this.RLock()
	defer this.RUnlock()
	return this.L.Front()
}

// TODO 异步rcu锁
func (this *SafeLinkedList) Len() int {
	this.RLock()
	defer this.RUnlock()
	return this.L.Len()
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

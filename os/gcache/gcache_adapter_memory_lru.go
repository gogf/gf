// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"sync"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/container/gmap"
)

// memoryLru holds LRU info.
// It uses list.List from stdlib for its underlying doubly linked list.
type memoryLru struct {
	mu   sync.RWMutex // Mutex to guarantee concurrent safety.
	cap  int          // LRU cap.
	data *gmap.Map    // Key mapping to the item of the list.
	list *glist.List  // Key list.
}

// newMemoryLru creates and returns a new LRU manager.
func newMemoryLru(cap int) *memoryLru {
	lru := &memoryLru{
		cap:  cap,
		data: gmap.New(false),
		list: glist.New(false),
	}
	return lru
}

// Remove deletes the `key` FROM `lru`.
func (l *memoryLru) Remove(keys ...any) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, key := range keys {
		if v := l.data.Remove(key); v != nil {
			l.list.Remove(v.(*glist.Element))
		}
	}
}

// SaveAndEvict saves the keys into LRU, evicts and returns the spare keys.
func (l *memoryLru) SaveAndEvict(keys ...any) (evictedKeys []any) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	evictedKeys = make([]any, 0)
	for _, key := range keys {
		if evictedKey := l.doSaveAndEvict(key); evictedKey != nil {
			evictedKeys = append(evictedKeys, evictedKey)
		}
	}
	return
}

func (l *memoryLru) doSaveAndEvict(key any) (evictedKey any) {
	var element *glist.Element
	if v := l.data.Get(key); v != nil {
		element = v.(*glist.Element)
		if element.Prev() == nil {
			// It this element is already on top of list,
			// it ignores the element moving.
			return
		}
		l.list.Remove(element)
	}

	// pushes the active key to top of list.
	element = l.list.PushFront(key)
	l.data.Set(key, element)
	// evict the spare key from list.
	if l.data.Size() <= l.cap {
		return
	}

	if evictedKey = l.list.PopBack(); evictedKey != nil {
		l.data.Remove(evictedKey)
	}
	return
}

// Clear deletes all keys.
func (l *memoryLru) Clear() {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.data.Clear()
	l.list.Clear()
}

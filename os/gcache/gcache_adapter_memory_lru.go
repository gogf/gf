// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/os/gtimer"
)

// LRU cache object.
// It uses list.List from stdlib for its underlying doubly linked list.
type adapterMemoryLru struct {
	cache   *AdapterMemory // Parent cache object.
	data    *gmap.Map      // Key mapping to the item of the list.
	list    *glist.List    // Key list.
	rawList *glist.List    // History for key adding.
	closed  *gtype.Bool    // Closed or not.
}

// newMemCacheLru creates and returns a new LRU object.
func newMemCacheLru(cache *AdapterMemory) *adapterMemoryLru {
	lru := &adapterMemoryLru{
		cache:   cache,
		data:    gmap.New(true),
		list:    glist.New(true),
		rawList: glist.New(true),
		closed:  gtype.NewBool(),
	}
	return lru
}

// Close closes the LRU object.
func (lru *adapterMemoryLru) Close() {
	lru.closed.Set(true)
}

// Remove deletes the `key` FROM `lru`.
func (lru *adapterMemoryLru) Remove(key interface{}) {
	if v := lru.data.Get(key); v != nil {
		lru.data.Remove(key)
		lru.list.Remove(v.(*glist.Element))
	}
}

// Size returns the size of `lru`.
func (lru *adapterMemoryLru) Size() int {
	return lru.data.Size()
}

// Push pushes `key` to the tail of `lru`.
func (lru *adapterMemoryLru) Push(key interface{}) {
	lru.rawList.PushBack(key)
}

// Pop deletes and returns the key from tail of `lru`.
func (lru *adapterMemoryLru) Pop() interface{} {
	if v := lru.list.PopBack(); v != nil {
		lru.data.Remove(v)
		return v
	}
	return nil
}

// SyncAndClear synchronizes the keys from `rawList` to `list` and `data`
// using Least Recently Used algorithm.
func (lru *adapterMemoryLru) SyncAndClear(ctx context.Context) {
	if lru.closed.Val() {
		gtimer.Exit()
		return
	}
	// Data synchronization.
	var alreadyExistItem interface{}
	for {
		if rawListItem := lru.rawList.PopFront(); rawListItem != nil {
			// Deleting the key from list.
			if alreadyExistItem = lru.data.Get(rawListItem); alreadyExistItem != nil {
				lru.list.Remove(alreadyExistItem.(*glist.Element))
			}
			// Pushing key to the head of the list
			// and setting its list item to hash table for quick indexing.
			lru.data.Set(rawListItem, lru.list.PushFront(rawListItem))
		} else {
			break
		}
	}
	// Data cleaning up.
	for clearLength := lru.Size() - lru.cache.cap; clearLength > 0; clearLength-- {
		if topKey := lru.Pop(); topKey != nil {
			lru.cache.clearByKey(topKey, true)
		}
	}
}

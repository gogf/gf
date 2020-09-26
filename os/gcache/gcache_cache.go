// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/util/gconv"
	"time"
)

// Cache struct.
type Cache struct {
	Adapter // Adapter for cache features.
}

// New creates and returns a new cache object using default memory adapter.
// Note that the LRU feature is only available using memory adapter.
func New(lruCap ...int) *Cache {
	memAdapter := newAdapterMemory(lruCap...)
	c := &Cache{
		Adapter: memAdapter,
	}
	// Here may be a "timer leak" if adapter is manually changed from memory adapter.
	// Do not worry about this, as adapter is less changed and it dose nothing if it's not used.
	gtimer.AddSingleton(time.Second, memAdapter.syncEventAndClearExpired)
	return c
}

// SetAdapter changes the adapter for this cache.
// Be very note that, this setting function is not concurrent-safe, which means you should not call
// this setting function concurrently in multiple goroutines.
func (c *Cache) SetAdapter(adapter Adapter) {
	c.Adapter = adapter
}

// Contains returns true if <key> exists in the cache, or else returns false.
func (c *Cache) Contains(key interface{}) bool {
	return c.Get(key) != nil
}

// GetVar retrieves and returns the value of <key> as gvar.Var.
func (c *Cache) GetVar(key interface{}) *gvar.Var {
	return gvar.New(c.Get(key))
}

// Removes deletes <keys> in the cache.
// Deprecated, use Remove instead.
func (c *Cache) Removes(keys []interface{}) {
	c.Remove(keys...)
}

// KeyStrings returns all keys in the cache as string slice.
func (c *Cache) KeyStrings() []string {
	return gconv.Strings(c.Keys())
}

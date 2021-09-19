// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/util/gconv"
	"time"
)

// Cache struct.
type Cache struct {
	localAdapter
}

// localAdapter is alias of Adapter, for embedded attribute purpose only.
type localAdapter = Adapter

// New creates and returns a new cache object using default memory adapter.
// Note that the LRU feature is only available using memory adapter.
func New(lruCap ...int) *Cache {
	memAdapter := NewAdapterMemory(lruCap...)
	c := &Cache{
		localAdapter: memAdapter,
	}
	// Here may be a "timer leak" if adapter is manually changed from memory adapter.
	// Do not worry about this, as adapter is less changed, and it does nothing if it's not used.
	gtimer.AddSingleton(time.Second, memAdapter.syncEventAndClearExpired)
	return c
}

// NewWithAdapter creates and returns a Cache object with given Adapter implements.
func NewWithAdapter(adapter Adapter) *Cache {
	return &Cache{
		localAdapter: adapter,
	}
}

// SetAdapter changes the adapter for this cache.
// Be very note that, this setting function is not concurrent-safe, which means you should not call
// this setting function concurrently in multiple goroutines.
func (c *Cache) SetAdapter(adapter Adapter) {
	c.localAdapter = adapter
}

// Removes deletes `keys` in the cache.
func (c *Cache) Removes(ctx context.Context, keys []interface{}) error {
	_, err := c.Remove(ctx, keys...)
	return err
}

// KeyStrings returns all keys in the cache as string slice.
func (c *Cache) KeyStrings(ctx context.Context) ([]string, error) {
	keys, err := c.Keys(ctx)
	if err != nil {
		return nil, err
	}
	return gconv.Strings(keys), nil
}

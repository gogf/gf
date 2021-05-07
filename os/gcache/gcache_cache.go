// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/util/gconv"
	"time"
)

// Cache struct.
type Cache struct {
	adapter Adapter         // Adapter for cache features.
	ctx     context.Context // Context for operations.
}

// New creates and returns a new cache object using default memory adapter.
// Note that the LRU feature is only available using memory adapter.
func New(lruCap ...int) *Cache {
	memAdapter := newAdapterMemory(lruCap...)
	c := &Cache{
		adapter: memAdapter,
	}
	// Here may be a "timer leak" if adapter is manually changed from memory adapter.
	// Do not worry about this, as adapter is less changed and it dose nothing if it's not used.
	gtimer.AddSingleton(time.Second, memAdapter.syncEventAndClearExpired)
	return c
}

// Clone returns a shallow copy of current object.
func (c *Cache) Clone() *Cache {
	return &Cache{
		adapter: c.adapter,
		ctx:     c.ctx,
	}
}

// Ctx is a chaining function, which shallowly clones current object and sets the context
// for next operation.
func (c *Cache) Ctx(ctx context.Context) *Cache {
	newCache := c.Clone()
	newCache.ctx = ctx
	return newCache
}

// SetAdapter changes the adapter for this cache.
// Be very note that, this setting function is not concurrent-safe, which means you should not call
// this setting function concurrently in multiple goroutines.
func (c *Cache) SetAdapter(adapter Adapter) {
	c.adapter = adapter
}

// GetVar retrieves and returns the value of <key> as gvar.Var.
func (c *Cache) GetVar(key interface{}) (*gvar.Var, error) {
	v, err := c.Get(key)
	return gvar.New(v), err
}

// Removes deletes <keys> in the cache.
// Deprecated, use Remove instead.
func (c *Cache) Removes(keys []interface{}) error {
	_, err := c.Remove(keys...)
	return err
}

// KeyStrings returns all keys in the cache as string slice.
func (c *Cache) KeyStrings() ([]string, error) {
	keys, err := c.Keys()
	if err != nil {
		return nil, err
	}
	return gconv.Strings(keys), nil
}

// KeyStrings returns all keys in the cache as string slice.
func (c *Cache) getCtx() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

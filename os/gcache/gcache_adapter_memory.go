// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"
	"math"
	"time"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
)

// AdapterMemory is an adapter implements using memory.
type AdapterMemory struct {
	data        *memoryData        // data is the underlying cache data which is stored in a hash table.
	expireTimes *memoryExpireTimes // expireTimes is the expiring key to its timestamp mapping, which is used for quick indexing and deleting.
	expireSets  *memoryExpireSets  // expireSets is the expiring timestamp to its key set mapping, which is used for quick indexing and deleting.
	lru         *memoryLru         // lru is the LRU manager, which is enabled when attribute cap > 0.
	eventList   *glist.List        // eventList is the asynchronous event list for internal data synchronization.
	closed      *gtype.Bool        // closed controls the cache closed or not.
}

// Internal event item.
type adapterMemoryEvent struct {
	k any   // Key.
	e int64 // Expire time in milliseconds.
}

const (
	// defaultMaxExpire is the default expire time for no expiring items.
	// It equals to math.MaxInt64/1000000.
	defaultMaxExpire = 9223372036854
)

// NewAdapterMemory creates and returns a new adapter_memory cache object.
func NewAdapterMemory() *AdapterMemory {
	return doNewAdapterMemory()
}

// NewAdapterMemoryLru creates and returns a new adapter_memory cache object with LRU.
func NewAdapterMemoryLru(cap int) *AdapterMemory {
	c := doNewAdapterMemory()
	c.lru = newMemoryLru(cap)
	return c
}

// doNewAdapterMemory creates and returns a new adapter_memory cache object.
func doNewAdapterMemory() *AdapterMemory {
	c := &AdapterMemory{
		data:        newMemoryData(),
		expireTimes: newMemoryExpireTimes(),
		expireSets:  newMemoryExpireSets(),
		eventList:   glist.New(true),
		closed:      gtype.NewBool(),
	}
	// Here may be a "timer leak" if adapter is manually changed from adapter_memory adapter.
	// Do not worry about this, as adapter is less changed, and it does nothing if it's not used.
	gtimer.AddSingleton(context.Background(), time.Second, c.syncEventAndClearExpired)
	return c
}

// Set sets cache with `key`-`value` pair, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory) Set(ctx context.Context, key any, value any, duration time.Duration) error {
	defer c.handleLruKey(ctx, key)
	expireTime := c.getInternalExpire(duration)
	c.data.Set(key, memoryDataItem{
		v: value,
		e: expireTime,
	})
	c.eventList.PushBack(&adapterMemoryEvent{
		k: key,
		e: expireTime,
	})
	return nil
}

// SetMap batch sets cache with key-value pairs by `data` map, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory) SetMap(ctx context.Context, data map[any]any, duration time.Duration) error {
	var (
		expireTime = c.getInternalExpire(duration)
		err        = c.data.SetMap(data, expireTime)
	)
	if err != nil {
		return err
	}
	for k := range data {
		c.eventList.PushBack(&adapterMemoryEvent{
			k: k,
			e: expireTime,
		})
	}
	if c.lru != nil {
		for key := range data {
			c.handleLruKey(ctx, key)
		}
	}
	return nil
}

// SetIfNotExist sets cache with `key`-`value` pair which is expired after `duration`
// if `key` does not exist in the cache. It returns true the `key` does not exist in the
// cache, and it sets `value` successfully to the cache, or else it returns false.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory) SetIfNotExist(ctx context.Context, key any, value any, duration time.Duration) (bool, error) {
	defer c.handleLruKey(ctx, key)
	isContained, err := c.Contains(ctx, key)
	if err != nil {
		return false, err
	}
	if !isContained {
		if _, err = c.doSetWithLockCheck(ctx, key, value, duration); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// SetIfNotExistFunc sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// The parameter `value` can be type of `func() any`, but it does nothing if its
// result is nil.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory) SetIfNotExistFunc(ctx context.Context, key any, f Func, duration time.Duration) (bool, error) {
	defer c.handleLruKey(ctx, key)
	isContained, err := c.Contains(ctx, key)
	if err != nil {
		return false, err
	}
	if !isContained {
		value, err := f(ctx)
		if err != nil {
			return false, err
		}
		if _, err = c.doSetWithLockCheck(ctx, key, value, duration); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// SetIfNotExistFuncLock sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
//
// Note that it differs from function `SetIfNotExistFunc` is that the function `f` is executed within
// writing mutex lock for concurrent safety purpose.
func (c *AdapterMemory) SetIfNotExistFuncLock(ctx context.Context, key any, f Func, duration time.Duration) (bool, error) {
	defer c.handleLruKey(ctx, key)
	isContained, err := c.Contains(ctx, key)
	if err != nil {
		return false, err
	}
	if !isContained {
		if _, err = c.doSetWithLockCheck(ctx, key, f, duration); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Get retrieves and returns the associated value of given `key`.
// It returns nil if it does not exist, or its value is nil, or it's expired.
// If you would like to check if the `key` exists in the cache, it's better using function Contains.
func (c *AdapterMemory) Get(ctx context.Context, key any) (*gvar.Var, error) {
	item, ok := c.data.Get(key)
	if ok && !item.IsExpired() {
		c.handleLruKey(ctx, key)
		return gvar.New(item.v), nil
	}
	return nil, nil
}

// GetOrSet retrieves and returns the value of `key`, or sets `key`-`value` pair and
// returns `value` if `key` does not exist in the cache. The key-value pair expires
// after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (c *AdapterMemory) GetOrSet(ctx context.Context, key any, value any, duration time.Duration) (*gvar.Var, error) {
	defer c.handleLruKey(ctx, key)
	v, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return c.doSetWithLockCheck(ctx, key, value, duration)
	}
	return v, nil
}

// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (c *AdapterMemory) GetOrSetFunc(ctx context.Context, key any, f Func, duration time.Duration) (*gvar.Var, error) {
	defer c.handleLruKey(ctx, key)
	v, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		value, err := f(ctx)
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, nil
		}
		return c.doSetWithLockCheck(ctx, key, value, duration)
	}
	return v, nil
}

// GetOrSetFuncLock retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
//
// Note that it differs from function `GetOrSetFunc` is that the function `f` is executed within
// writing mutex lock for concurrent safety purpose.
func (c *AdapterMemory) GetOrSetFuncLock(ctx context.Context, key any, f Func, duration time.Duration) (*gvar.Var, error) {
	defer c.handleLruKey(ctx, key)
	v, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return c.doSetWithLockCheck(ctx, key, f, duration)
	}
	return v, nil
}

// Contains checks and returns true if `key` exists in the cache, or else returns false.
func (c *AdapterMemory) Contains(ctx context.Context, key any) (bool, error) {
	v, err := c.Get(ctx, key)
	if err != nil {
		return false, err
	}
	return v != nil, nil
}

// GetExpire retrieves and returns the expiration of `key` in the cache.
//
// Note that,
// It returns 0 if the `key` does not expire.
// It returns -1 if the `key` does not exist in the cache.
func (c *AdapterMemory) GetExpire(ctx context.Context, key any) (time.Duration, error) {
	if item, ok := c.data.Get(key); ok {
		c.handleLruKey(ctx, key)
		return time.Duration(item.e-gtime.TimestampMilli()) * time.Millisecond, nil
	}
	return -1, nil
}

// Remove deletes one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the last deleted item.
func (c *AdapterMemory) Remove(ctx context.Context, keys ...any) (*gvar.Var, error) {
	defer c.lru.Remove(keys...)
	return c.doRemove(ctx, keys...)
}

func (c *AdapterMemory) doRemove(_ context.Context, keys ...any) (*gvar.Var, error) {
	var removedKeys []any
	removedKeys, value, err := c.data.Remove(keys...)
	if err != nil {
		return nil, err
	}
	for _, key := range removedKeys {
		c.eventList.PushBack(&adapterMemoryEvent{
			k: key,
			e: gtime.TimestampMilli() - 1000,
		})
	}
	return gvar.New(value), nil
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func (c *AdapterMemory) Update(ctx context.Context, key any, value any) (oldValue *gvar.Var, exist bool, err error) {
	v, exist, err := c.data.Update(key, value)
	if exist {
		c.handleLruKey(ctx, key)
	}
	return gvar.New(v), exist, err
}

// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func (c *AdapterMemory) UpdateExpire(ctx context.Context, key any, duration time.Duration) (oldDuration time.Duration, err error) {
	newExpireTime := c.getInternalExpire(duration)
	oldDuration, err = c.data.UpdateExpire(key, newExpireTime)
	if err != nil {
		return
	}
	if oldDuration != -1 {
		c.eventList.PushBack(&adapterMemoryEvent{
			k: key,
			e: newExpireTime,
		})
		c.handleLruKey(ctx, key)
	}
	return
}

// Size returns the size of the cache.
func (c *AdapterMemory) Size(ctx context.Context) (size int, err error) {
	return c.data.Size()
}

// Data returns a copy of all key-value pairs in the cache as map type.
func (c *AdapterMemory) Data(ctx context.Context) (map[any]any, error) {
	return c.data.Data()
}

// Keys returns all keys in the cache as slice.
func (c *AdapterMemory) Keys(ctx context.Context) ([]any, error) {
	return c.data.Keys()
}

// Values returns all values in the cache as slice.
func (c *AdapterMemory) Values(ctx context.Context) ([]any, error) {
	return c.data.Values()
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
func (c *AdapterMemory) Clear(ctx context.Context) error {
	c.data.Clear()
	c.lru.Clear()
	return nil
}

// Close closes the cache.
func (c *AdapterMemory) Close(ctx context.Context) error {
	c.closed.Set(true)
	return nil
}

// doSetWithLockCheck sets cache with `key`-`value` pair if `key` does not exist in the
// cache, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// The parameter `value` can be type of <func() any>, but it does nothing if the
// function result is nil.
//
// It doubly checks the `key` whether exists in the cache using mutex writing lock
// before setting it to the cache.
func (c *AdapterMemory) doSetWithLockCheck(ctx context.Context, key any, value any, duration time.Duration) (result *gvar.Var, err error) {
	expireTimestamp := c.getInternalExpire(duration)
	v, err := c.data.SetWithLock(ctx, key, value, expireTimestamp)
	c.eventList.PushBack(&adapterMemoryEvent{k: key, e: expireTimestamp})
	return gvar.New(v), err
}

// getInternalExpire converts and returns the expiration time with given expired duration in milliseconds.
func (c *AdapterMemory) getInternalExpire(duration time.Duration) int64 {
	if duration == 0 {
		return defaultMaxExpire
	}
	return gtime.TimestampMilli() + duration.Nanoseconds()/1000000
}

// makeExpireKey groups the `expire` in milliseconds to its according seconds.
func (c *AdapterMemory) makeExpireKey(expire int64) int64 {
	return int64(math.Ceil(float64(expire/1000)+1) * 1000)
}

// syncEventAndClearExpired does the asynchronous task loop:
//  1. Asynchronously process the data in the event list,
//     and synchronize the results to the `expireTimes` and `expireSets` properties.
//  2. Clean up the expired key-value pair data.
func (c *AdapterMemory) syncEventAndClearExpired(ctx context.Context) {
	if c.closed.Val() {
		gtimer.Exit()
		return
	}
	var (
		event         *adapterMemoryEvent
		oldExpireTime int64
		newExpireTime int64
	)
	// ================================
	// Data expiration synchronization.
	// ================================
	for {
		v := c.eventList.PopFront()
		if v == nil {
			break
		}
		event = v.(*adapterMemoryEvent)
		// Fetching the old expire set.
		oldExpireTime = c.expireTimes.Get(event.k)
		// Calculating the new expiration time set.
		newExpireTime = c.makeExpireKey(event.e)
		// Expiration changed for this key.
		if newExpireTime != oldExpireTime {
			c.expireSets.GetOrNew(newExpireTime).Add(event.k)
			if oldExpireTime != 0 {
				c.expireSets.GetOrNew(oldExpireTime).Remove(event.k)
			}
			// Updating the expired time for `event.k`.
			c.expireTimes.Set(event.k, newExpireTime)
		}
	}
	// =================================
	// Data expiration auto cleaning up.
	// =================================
	var (
		expireSet  *gset.Set
		expireTime int64
		currentEk  = c.makeExpireKey(gtime.TimestampMilli())
	)
	// auto removing expiring key set for latest seconds.
	for i := int64(1); i <= 5; i++ {
		expireTime = currentEk - i*1000
		if expireSet = c.expireSets.Get(expireTime); expireSet != nil {
			// Iterating the set to delete all keys in it.
			expireSet.Iterator(func(key any) bool {
				c.deleteExpiredKey(key)
				// remove auto expired key for lru.
				c.lru.Remove(key)
				return true
			})
			// Deleting the set after all of its keys are deleted.
			c.expireSets.Delete(expireTime)
		}
	}
}

func (c *AdapterMemory) handleLruKey(ctx context.Context, keys ...any) {
	if c.lru == nil {
		return
	}
	if evictedKeys := c.lru.SaveAndEvict(keys...); len(evictedKeys) > 0 {
		_, _ = c.doRemove(ctx, evictedKeys...)
		return
	}
}

// clearByKey deletes the key-value pair with given `key`.
// The parameter `force` specifies whether doing this deleting forcibly.
func (c *AdapterMemory) deleteExpiredKey(key any) {
	// Doubly check before really deleting it from cache.
	c.data.Delete(key)
	// Deleting its expiration time from `expireTimes`.
	c.expireTimes.Delete(key)
}

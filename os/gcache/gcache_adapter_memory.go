// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"
	"github.com/gogf/gf/container/gvar"
	"math"
	"time"

	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gtimer"
)

// AdapterMemory is an adapter implements using memory.
type AdapterMemory struct {
	// cap limits the size of the cache pool.
	// If the size of the cache exceeds the cap,
	// the cache expiration process performs according to the LRU algorithm.
	// It is 0 in default which means no limits.
	cap         int
	data        *adapterMemoryData        // data is the underlying cache data which is stored in a hash table.
	expireTimes *adapterMemoryExpireTimes // expireTimes is the expiring key to its timestamp mapping, which is used for quick indexing and deleting.
	expireSets  *adapterMemoryExpireSets  // expireSets is the expiring timestamp to its key set mapping, which is used for quick indexing and deleting.
	lru         *adapterMemoryLru         // lru is the LRU manager, which is enabled when attribute cap > 0.
	lruGetList  *glist.List               // lruGetList is the LRU history according to Get function.
	eventList   *glist.List               // eventList is the asynchronous event list for internal data synchronization.
	closed      *gtype.Bool               // closed controls the cache closed or not.
}

// Internal cache item.
type adapterMemoryItem struct {
	v interface{} // Value.
	e int64       // Expire timestamp in milliseconds.
}

// Internal event item.
type adapterMemoryEvent struct {
	k interface{} // Key.
	e int64       // Expire time in milliseconds.
}

const (
	// defaultMaxExpire is the default expire time for no expiring items.
	// It equals to math.MaxInt64/1000000.
	defaultMaxExpire = 9223372036854
)

// NewAdapterMemory creates and returns a new memory cache object.
func NewAdapterMemory(lruCap ...int) Adapter {
	c := &AdapterMemory{
		data:        newAdapterMemoryData(),
		lruGetList:  glist.New(true),
		expireTimes: newAdapterMemoryExpireTimes(),
		expireSets:  newAdapterMemoryExpireSets(),
		eventList:   glist.New(true),
		closed:      gtype.NewBool(),
	}
	if len(lruCap) > 0 {
		c.cap = lruCap[0]
		c.lru = newMemCacheLru(c)
	}
	return c
}

// Set sets cache with `key`-`value` pair, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory) Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error {
	expireTime := c.getInternalExpire(duration)
	c.data.Set(key, adapterMemoryItem{
		v: value,
		e: expireTime,
	})
	c.eventList.PushBack(&adapterMemoryEvent{
		k: key,
		e: expireTime,
	})
	return nil
}

// Sets batch sets cache with key-value pairs by `data` map, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory) Sets(ctx context.Context, data map[interface{}]interface{}, duration time.Duration) error {
	var (
		expireTime = c.getInternalExpire(duration)
		err        = c.data.Sets(data, expireTime)
	)
	if err != nil {
		return err
	}
	for k, _ := range data {
		c.eventList.PushBack(&adapterMemoryEvent{
			k: k,
			e: expireTime,
		})
	}
	return nil
}

// SetIfNotExist sets cache with `key`-`value` pair which is expired after `duration`
// if `key` does not exist in the cache. It returns true the `key` does not exist in the
// cache, and it sets `value` successfully to the cache, or else it returns false.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory) SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (bool, error) {
	isContained, err := c.Contains(ctx, key)
	if err != nil {
		return false, err
	}
	if !isContained {
		if _, err = c.doSetWithLockCheck(key, value, duration); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// SetIfNotExistFunc sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// The parameter `value` can be type of `func() interface{}`, but it does nothing if its
// result is nil.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory) SetIfNotExistFunc(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (bool, error) {
	isContained, err := c.Contains(ctx, key)
	if err != nil {
		return false, err
	}
	if !isContained {
		value, err := f()
		if err != nil {
			return false, err
		}
		if _, err = c.doSetWithLockCheck(key, value, duration); err != nil {
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
func (c *AdapterMemory) SetIfNotExistFuncLock(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (bool, error) {
	isContained, err := c.Contains(ctx, key)
	if err != nil {
		return false, err
	}
	if !isContained {
		if _, err = c.doSetWithLockCheck(key, f, duration); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Get retrieves and returns the associated value of given `key`.
// It returns nil if it does not exist, or its value is nil, or it's expired.
// If you would like to check if the `key` exists in the cache, it's better using function Contains.
func (c *AdapterMemory) Get(ctx context.Context, key interface{}) (*gvar.Var, error) {
	item, ok := c.data.Get(key)
	if ok && !item.IsExpired() {
		// Adding to LRU history if LRU feature is enabled.
		if c.cap > 0 {
			c.lruGetList.PushBack(key)
		}
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
func (c *AdapterMemory) GetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (*gvar.Var, error) {
	v, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return c.doSetWithLockCheck(key, value, duration)
	} else {
		return v, nil
	}
}

// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (c *AdapterMemory) GetOrSetFunc(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (*gvar.Var, error) {
	v, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		value, err := f()
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, nil
		}
		return c.doSetWithLockCheck(key, value, duration)
	} else {
		return v, nil
	}
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
func (c *AdapterMemory) GetOrSetFuncLock(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (*gvar.Var, error) {
	v, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return c.doSetWithLockCheck(key, f, duration)
	} else {
		return v, nil
	}
}

// Contains checks and returns true if `key` exists in the cache, or else returns false.
func (c *AdapterMemory) Contains(ctx context.Context, key interface{}) (bool, error) {
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
func (c *AdapterMemory) GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	if item, ok := c.data.Get(key); ok {
		return time.Duration(item.e-gtime.TimestampMilli()) * time.Millisecond, nil
	}
	return -1, nil
}

// Remove deletes one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the last deleted item.
func (c *AdapterMemory) Remove(ctx context.Context, keys ...interface{}) (*gvar.Var, error) {
	var removedKeys []interface{}
	removedKeys, value, err := c.data.Remove(keys...)
	if err != nil {
		return nil, err
	}
	for _, key := range removedKeys {
		c.eventList.PushBack(&adapterMemoryEvent{
			k: key,
			e: gtime.TimestampMilli() - 1000000,
		})
	}
	return gvar.New(value), nil
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func (c *AdapterMemory) Update(ctx context.Context, key interface{}, value interface{}) (oldValue *gvar.Var, exist bool, err error) {
	v, exist, err := c.data.Update(key, value)
	return gvar.New(v), exist, err
}

// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func (c *AdapterMemory) UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
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
	}
	return
}

// Size returns the size of the cache.
func (c *AdapterMemory) Size(ctx context.Context) (size int, err error) {
	return c.data.Size()
}

// Data returns a copy of all key-value pairs in the cache as map type.
func (c *AdapterMemory) Data(ctx context.Context) (map[interface{}]interface{}, error) {
	return c.data.Data()
}

// Keys returns all keys in the cache as slice.
func (c *AdapterMemory) Keys(ctx context.Context) ([]interface{}, error) {
	return c.data.Keys()
}

// Values returns all values in the cache as slice.
func (c *AdapterMemory) Values(ctx context.Context) ([]interface{}, error) {
	return c.data.Values()
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
func (c *AdapterMemory) Clear(ctx context.Context) error {
	return c.data.Clear()
}

// Close closes the cache.
func (c *AdapterMemory) Close(ctx context.Context) error {
	if c.cap > 0 {
		c.lru.Close()
	}
	c.closed.Set(true)
	return nil
}

// doSetWithLockCheck sets cache with `key`-`value` pair if `key` does not exist in the
// cache, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// The parameter `value` can be type of <func() interface{}>, but it does nothing if the
// function result is nil.
//
// It doubly checks the `key` whether exists in the cache using mutex writing lock
// before setting it to the cache.
func (c *AdapterMemory) doSetWithLockCheck(key interface{}, value interface{}, duration time.Duration) (result *gvar.Var, err error) {
	expireTimestamp := c.getInternalExpire(duration)
	v, err := c.data.SetWithLock(key, value, expireTimestamp)
	c.eventList.PushBack(&adapterMemoryEvent{k: key, e: expireTimestamp})
	return gvar.New(v), err
}

// getInternalExpire converts and returns the expiration time with given expired duration in milliseconds.
func (c *AdapterMemory) getInternalExpire(duration time.Duration) int64 {
	if duration == 0 {
		return defaultMaxExpire
	} else {
		return gtime.TimestampMilli() + duration.Nanoseconds()/1000000
	}
}

// makeExpireKey groups the `expire` in milliseconds to its according seconds.
func (c *AdapterMemory) makeExpireKey(expire int64) int64 {
	return int64(math.Ceil(float64(expire/1000)+1) * 1000)
}

// syncEventAndClearExpired does the asynchronous task loop:
// 1. Asynchronously process the data in the event list,
//    and synchronize the results to the `expireTimes` and `expireSets` properties.
// 2. Clean up the expired key-value pair data.
func (c *AdapterMemory) syncEventAndClearExpired() {
	if c.closed.Val() {
		gtimer.Exit()
		return
	}
	var (
		event         *adapterMemoryEvent
		oldExpireTime int64
		newExpireTime int64
	)
	// ========================
	// Data Synchronization.
	// ========================
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
		if newExpireTime != oldExpireTime {
			c.expireSets.GetOrNew(newExpireTime).Add(event.k)
			if oldExpireTime != 0 {
				c.expireSets.GetOrNew(oldExpireTime).Remove(event.k)
			}
			// Updating the expired time for <event.k>.
			c.expireTimes.Set(event.k, newExpireTime)
		}
		// Adding the key the LRU history by writing operations.
		if c.cap > 0 {
			c.lru.Push(event.k)
		}
	}
	// Processing expired keys from LRU.
	if c.cap > 0 && c.lruGetList.Len() > 0 {
		for {
			if v := c.lruGetList.PopFront(); v != nil {
				c.lru.Push(v)
			} else {
				break
			}
		}
	}
	// ========================
	// Data Cleaning up.
	// ========================
	var (
		expireSet *gset.Set
		ek        = c.makeExpireKey(gtime.TimestampMilli())
		eks       = []int64{ek - 1000, ek - 2000, ek - 3000, ek - 4000, ek - 5000}
	)
	for _, expireTime := range eks {
		if expireSet = c.expireSets.Get(expireTime); expireSet != nil {
			// Iterating the set to delete all keys in it.
			expireSet.Iterator(func(key interface{}) bool {
				c.clearByKey(key)
				return true
			})
			// Deleting the set after all of its keys are deleted.
			c.expireSets.Delete(expireTime)
		}
	}
}

// clearByKey deletes the key-value pair with given `key`.
// The parameter `force` specifies whether doing this deleting forcibly.
func (c *AdapterMemory) clearByKey(key interface{}, force ...bool) {
	// Doubly check before really deleting it from cache.
	c.data.DeleteWithDoubleCheck(key, force...)

	// Deleting its expiration time from `expireTimes`.
	c.expireTimes.Delete(key)

	// Deleting it from LRU.
	if c.cap > 0 {
		c.lru.Remove(key)
	}
}

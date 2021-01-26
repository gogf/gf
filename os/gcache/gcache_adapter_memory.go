// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"math"
	"time"

	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gtimer"
)

// Internal cache object.
type adapterMemory struct {
	// cap limits the size of the cache pool.
	// If the size of the cache exceeds the cap,
	// the cache expiration process performs according to the LRU algorithm.
	// It is 0 in default which means no limits.
	cap         int
	data        *adapterMemoryData        // data is the underlying cache data which is stored in a hash table.
	expireTimes *adapterMemoryExpireTimes // expireTimes is the expiring key to its timestamp mapping, which is used for quick indexing and deleting.
	expireSets  *adapterMemoryExpireSets  // expireSets is the expiring timestamp to its key set mapping, which is used for quick indexing and deleting.
	lru         *adapterMemoryLru         // lru is the LRU manager, which is enabled when attribute cap > 0.
	lruGetList  *glist.List               // lruGetList is the LRU history according with Get function.
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

// newAdapterMemory creates and returns a new memory cache object.
func newAdapterMemory(lruCap ...int) *adapterMemory {
	c := &adapterMemory{
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

// Set sets cache with <key>-<value> pair, which is expired after <duration>.
//
// It does not expire if <duration> == 0.
// It deletes the <key> if <duration> < 0.
func (c *adapterMemory) Set(key interface{}, value interface{}, duration time.Duration) error {
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

// Update updates the value of <key> without changing its expiration and returns the old value.
// The returned value <exist> is false if the <key> does not exist in the cache.
//
// It deletes the <key> if given <value> is nil.
// It does nothing if <key> does not exist in the cache.
func (c *adapterMemory) Update(key interface{}, value interface{}) (oldValue interface{}, exist bool, err error) {
	return c.data.Update(key, value)
}

// UpdateExpire updates the expiration of <key> and returns the old expiration duration value.
//
// It returns -1 and does nothing if the <key> does not exist in the cache.
// It deletes the <key> if <duration> < 0.
func (c *adapterMemory) UpdateExpire(key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
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

// GetExpire retrieves and returns the expiration of <key> in the cache.
//
// It returns 0 if the <key> does not expire.
// It returns -1 if the <key> does not exist in the cache.
func (c *adapterMemory) GetExpire(key interface{}) (time.Duration, error) {
	if item, ok := c.data.Get(key); ok {
		return time.Duration(item.e-gtime.TimestampMilli()) * time.Millisecond, nil
	}
	return -1, nil
}

// SetIfNotExist sets cache with <key>-<value> pair which is expired after <duration>
// if <key> does not exist in the cache. It returns true the <key> dose not exist in the
// cache and it sets <value> successfully to the cache, or else it returns false.
// The parameter <value> can be type of <func() interface{}>, but it dose nothing if its
// result is nil.
//
// It does not expire if <duration> == 0.
// It deletes the <key> if <duration> < 0 or given <value> is nil.
func (c *adapterMemory) SetIfNotExist(key interface{}, value interface{}, duration time.Duration) (bool, error) {
	isContained, err := c.Contains(key)
	if err != nil {
		return false, err
	}
	if !isContained {
		_, err := c.doSetWithLockCheck(key, value, duration)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Sets batch sets cache with key-value pairs by <data>, which is expired after <duration>.
//
// It does not expire if <duration> == 0.
// It deletes the keys of <data> if <duration> < 0 or given <value> is nil.
func (c *adapterMemory) Sets(data map[interface{}]interface{}, duration time.Duration) error {
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

// Get retrieves and returns the associated value of given <key>.
// It returns nil if it does not exist or its value is nil.
func (c *adapterMemory) Get(key interface{}) (interface{}, error) {
	item, ok := c.data.Get(key)
	if ok && !item.IsExpired() {
		// Adding to LRU history if LRU feature is enabled.
		if c.cap > 0 {
			c.lruGetList.PushBack(key)
		}
		return item.v, nil
	}
	return nil, nil
}

// GetOrSet retrieves and returns the value of <key>, or sets <key>-<value> pair and
// returns <value> if <key> does not exist in the cache. The key-value pair expires
// after <duration>.
//
// It does not expire if <duration> == 0.
// It deletes the <key> if <duration> < 0 or given <value> is nil, but it does nothing
// if <value> is a function and the function result is nil.
func (c *adapterMemory) GetOrSet(key interface{}, value interface{}, duration time.Duration) (interface{}, error) {
	v, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return c.doSetWithLockCheck(key, value, duration)
	} else {
		return v, nil
	}
}

// GetOrSetFunc retrieves and returns the value of <key>, or sets <key> with result of
// function <f> and returns its result if <key> does not exist in the cache. The key-value
// pair expires after <duration>.
//
// It does not expire if <duration> == 0.
// It deletes the <key> if <duration> < 0 or given <value> is nil, but it does nothing
// if <value> is a function and the function result is nil.
func (c *adapterMemory) GetOrSetFunc(key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	v, err := c.Get(key)
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

// GetOrSetFuncLock retrieves and returns the value of <key>, or sets <key> with result of
// function <f> and returns its result if <key> does not exist in the cache. The key-value
// pair expires after <duration>.
//
// It does not expire if <duration> == 0.
// It does nothing if function <f> returns nil.
//
// Note that the function <f> should be executed within writing mutex lock for concurrent
// safety purpose.
func (c *adapterMemory) GetOrSetFuncLock(key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	v, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return c.doSetWithLockCheck(key, f, duration)
	} else {
		return v, nil
	}
}

// Contains returns true if <key> exists in the cache, or else returns false.
func (c *adapterMemory) Contains(key interface{}) (bool, error) {
	v, err := c.Get(key)
	if err != nil {
		return false, err
	}
	return v != nil, nil
}

// Remove deletes the one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the deleted last item.
func (c *adapterMemory) Remove(keys ...interface{}) (value interface{}, err error) {
	var removedKeys []interface{}
	removedKeys, value, err = c.data.Remove(keys...)
	if err != nil {
		return
	}
	for _, key := range removedKeys {
		c.eventList.PushBack(&adapterMemoryEvent{
			k: key,
			e: gtime.TimestampMilli() - 1000000,
		})
	}
	return
}

// Data returns a copy of all key-value pairs in the cache as map type.
func (c *adapterMemory) Data() (map[interface{}]interface{}, error) {
	return c.data.Data()
}

// Keys returns all keys in the cache as slice.
func (c *adapterMemory) Keys() ([]interface{}, error) {
	return c.data.Keys()
}

// Values returns all values in the cache as slice.
func (c *adapterMemory) Values() ([]interface{}, error) {
	return c.data.Values()
}

// Size returns the size of the cache.
func (c *adapterMemory) Size() (size int, err error) {
	return c.data.Size()
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
func (c *adapterMemory) Clear() error {
	return c.data.Clear()
}

// Close closes the cache.
func (c *adapterMemory) Close() error {
	if c.cap > 0 {
		c.lru.Close()
	}
	c.closed.Set(true)
	return nil
}

// doSetWithLockCheck sets cache with <key>-<value> pair if <key> does not exist in the
// cache, which is expired after <duration>.
//
// It does not expire if <duration> == 0.
// The parameter <value> can be type of <func() interface{}>, but it dose nothing if the
// function result is nil.
//
// It doubly checks the <key> whether exists in the cache using mutex writing lock
// before setting it to the cache.
func (c *adapterMemory) doSetWithLockCheck(key interface{}, value interface{}, duration time.Duration) (result interface{}, err error) {
	expireTimestamp := c.getInternalExpire(duration)
	result, err = c.data.SetWithLock(key, value, expireTimestamp)
	c.eventList.PushBack(&adapterMemoryEvent{k: key, e: expireTimestamp})
	return
}

// getInternalExpire converts and returns the expire time with given expired duration in milliseconds.
func (c *adapterMemory) getInternalExpire(duration time.Duration) int64 {
	if duration == 0 {
		return defaultMaxExpire
	} else {
		return gtime.TimestampMilli() + duration.Nanoseconds()/1000000
	}
}

// makeExpireKey groups the <expire> in milliseconds to its according seconds.
func (c *adapterMemory) makeExpireKey(expire int64) int64 {
	return int64(math.Ceil(float64(expire/1000)+1) * 1000)
}

// syncEventAndClearExpired does the asynchronous task loop:
// 1. Asynchronously process the data in the event list,
//    and synchronize the results to the <expireTimes> and <expireSets> properties.
// 2. Clean up the expired key-value pair data.
func (c *adapterMemory) syncEventAndClearExpired() {
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
		// Calculating the new expire set.
		newExpireTime = c.makeExpireKey(event.e)
		if newExpireTime != oldExpireTime {
			c.expireSets.GetOrNew(newExpireTime).Add(event.k)
			if oldExpireTime != 0 {
				c.expireSets.GetOrNew(oldExpireTime).Remove(event.k)
			}
			// Updating the expire time for <event.k>.
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

// clearByKey deletes the key-value pair with given <key>.
// The parameter <force> specifies whether doing this deleting forcibly.
func (c *adapterMemory) clearByKey(key interface{}, force ...bool) {
	// Doubly check before really deleting it from cache.
	c.data.DeleteWithDoubleCheck(key, force...)

	// Deleting its expire time from <expireTimes>.
	c.expireTimes.Delete(key)

	// Deleting it from LRU.
	if c.cap > 0 {
		c.lru.Remove(key)
	}
}

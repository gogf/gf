// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"math"
	"sync"
	"time"

	"github.com/gogf/gf/g/container/glist"
	"github.com/gogf/gf/g/container/gset"
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/os/gtimer"
	"github.com/gogf/gf/g/util/gconv"
)

// Internal cache object.
type memCache struct {
	dataMu       sync.RWMutex
	expireTimeMu sync.RWMutex
	expireSetMu  sync.RWMutex

	// <cap> limits the size of the cache pool.
	// If the size of the cache exceeds the <cap>,
	// the cache expiration process is performed according to the LRU algorithm.
	// It is 0 in default which means no limits.
	cap         int
	data        map[interface{}]memCacheItem // Underlying cache data which is stored in a hash table.
	expireTimes map[interface{}]int64        // Expiring key mapping to its timestamp, which is used for quick indexing and deleting.
	expireSets  map[int64]*gset.Set          // Expiring timestamp mapping to its key set, which is used for quick indexing and deleting.

	lru        *memCacheLru // LRU object, which is enabled when <cap> > 0.
	lruGetList *glist.List  // LRU history according with Get function.
	eventList  *glist.List  // Asynchronous event list for internal data synchronization.
	closed     *gtype.Bool  // Is this cache closed or not.
}

// Internal cache item.
type memCacheItem struct {
	v interface{} // Value.
	e int64       // Expire time in milliseconds.
}

// Internal event item.
type memCacheEvent struct {
	k interface{} // Key.
	e int64       // Expire time in milliseconds.
}

const (
	// Default expire time for no expiring items.
	// It equals to math.MaxInt64/1000000.
	gDEFAULT_MAX_EXPIRE = 9223372036854
)

// newMemCache creates and returns a new memory cache object.
func newMemCache(lruCap ...int) *memCache {
	c := &memCache{
		lruGetList:  glist.New(),
		data:        make(map[interface{}]memCacheItem),
		expireTimes: make(map[interface{}]int64),
		expireSets:  make(map[int64]*gset.Set),
		eventList:   glist.New(),
		closed:      gtype.NewBool(),
	}
	if len(lruCap) > 0 {
		c.cap = lruCap[0]
		c.lru = newMemCacheLru(c)
	}
	return c
}

// makeExpireKey groups the <expire> in milliseconds to its according seconds.
func (c *memCache) makeExpireKey(expire int64) int64 {
	return int64(math.Ceil(float64(expire/1000)+1) * 1000)
}

// getExpireSet returns the expire set for given <expire> in seconds.
func (c *memCache) getExpireSet(expire int64) (expireSet *gset.Set) {
	c.expireSetMu.RLock()
	expireSet, _ = c.expireSets[expire]
	c.expireSetMu.RUnlock()
	return
}

// getOrNewExpireSet returns the expire set for given <expire> in seconds.
// It creates and returns a new set for <expire> if it does not exist.
func (c *memCache) getOrNewExpireSet(expire int64) (expireSet *gset.Set) {
	if expireSet = c.getExpireSet(expire); expireSet == nil {
		expireSet = gset.New()
		c.expireSetMu.Lock()
		if es, ok := c.expireSets[expire]; ok {
			expireSet = es
		} else {
			c.expireSets[expire] = expireSet
		}
		c.expireSetMu.Unlock()
	}
	return
}

// getMilliExpire converts parameter <duration> to int type in milliseconds.
//
// Note that there's some performance cost in type assertion here, but it's valuable.
func (c *memCache) getMilliExpire(duration interface{}) int {
	if d, ok := duration.(time.Duration); ok {
		return int(d.Nanoseconds() / 1000000)
	} else {
		return duration.(int)
	}
}

// Set sets cache with <key>-<value> pair, which is expired after <duration>.
//
// The parameter <duration> can be either type of int or time.Duration.
// If <duration> is type of int, it means <duration> milliseconds.
// If <duration> <=0 means it does not expire.
func (c *memCache) Set(key interface{}, value interface{}, duration interface{}) {
	expire := c.getMilliExpire(duration)
	expireTime := c.getInternalExpire(expire)
	c.dataMu.Lock()
	c.data[key] = memCacheItem{v: value, e: expireTime}
	c.dataMu.Unlock()
	c.eventList.PushBack(&memCacheEvent{k: key, e: expireTime})
}

// doSetWithLockCheck sets cache with <key>-<value> pair if <key> does not exist in the cache,
// which is expired after <duration>.
//
// The parameter <duration> can be either type of int or time.Duration.
// If <duration> is type of int, it means <duration> milliseconds.
// If <duration> <=0 means it does not expire.
//
// It doubly checks the <key> whether exists in the cache using mutex writing lock
// before setting it to the cache.
func (c *memCache) doSetWithLockCheck(key interface{}, value interface{}, duration interface{}) interface{} {
	expire := c.getMilliExpire(duration)
	expireTimestamp := c.getInternalExpire(expire)
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	if v, ok := c.data[key]; ok && !v.IsExpired() {
		return v.v
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value == nil {
		return nil
	}
	c.data[key] = memCacheItem{v: value, e: expireTimestamp}
	c.eventList.PushBack(&memCacheEvent{k: key, e: expireTimestamp})
	return value
}

// getInternalExpire returns the expire time with given expire duration in milliseconds.
func (c *memCache) getInternalExpire(expire int) int64 {
	if expire != 0 {
		return gtime.Millisecond() + int64(expire)
	} else {
		return gDEFAULT_MAX_EXPIRE
	}
}

// SetIfNotExist sets cache with <key>-<value> pair if <key> does not exist in the cache,
// which is expired after <duration>.
//
// The parameter <duration> can be either type of int or time.Duration.
// If <duration> is type of int, it means <duration> milliseconds.
// If <duration> <=0 means it does not expire.
func (c *memCache) SetIfNotExist(key interface{}, value interface{}, duration interface{}) bool {
	expire := c.getMilliExpire(duration)
	if !c.Contains(key) {
		c.doSetWithLockCheck(key, value, expire)
		return true
	}
	return false
}

// Sets batch sets cache with key-value pairs by <data>, which is expired after <duration>.
//
// The parameter <duration> can be either type of int or time.Duration.
// If <duration> is type of int, it means <duration> milliseconds.
// If <duration> <=0 means it does not expire.
func (c *memCache) Sets(data map[interface{}]interface{}, duration interface{}) {
	expire := c.getMilliExpire(duration)
	expireTime := c.getInternalExpire(expire)
	for k, v := range data {
		c.dataMu.Lock()
		c.data[k] = memCacheItem{v: v, e: expireTime}
		c.dataMu.Unlock()
		c.eventList.PushBack(&memCacheEvent{k: k, e: expireTime})
	}
}

// Get returns the value of <key>.
// It returns nil if it does not exist or its value is nil.
func (c *memCache) Get(key interface{}) interface{} {
	c.dataMu.RLock()
	item, ok := c.data[key]
	c.dataMu.RUnlock()
	if ok && !item.IsExpired() {
		// Adding to LRU history if LRU feature is enbaled.
		if c.cap > 0 {
			c.lruGetList.PushBack(key)
		}
		return item.v
	}
	return nil
}

// GetOrSet returns the value of <key>,
// or sets <key>-<value> pair and returns <value> if <key> does not exist in the cache.
// The key-value pair expires after <duration>.
//
// The parameter <duration> can be either type of int or time.Duration.
// If <duration> is type of int, it means <duration> milliseconds.
// If <duration> <=0 means it does not expire.
func (c *memCache) GetOrSet(key interface{}, value interface{}, duration interface{}) interface{} {
	if v := c.Get(key); v == nil {
		return c.doSetWithLockCheck(key, value, duration)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value of <key>,
// or sets <key> with result of function <f> and returns its result
// if <key> does not exist in the cache.
// The key-value pair expires after <duration>.
//
// The parameter <duration> can be either type of int or time.Duration.
// If <duration> is type of int, it means <duration> milliseconds.
// If <duration> <=0 means it does not expire.
func (c *memCache) GetOrSetFunc(key interface{}, f func() interface{}, duration interface{}) interface{} {
	if v := c.Get(key); v == nil {
		return c.doSetWithLockCheck(key, f(), duration)
	} else {
		return v
	}
}

// GetOrSetFuncLock returns the value of <key>,
// or sets <key> with result of function <f> and returns its result
// if <key> does not exist in the cache.
// The key-value pair expires after <duration>.
//
// The parameter <duration> can be either type of int or time.Duration.
// If <duration> is type of int, it means <duration> milliseconds.
// If <duration> <=0 means it does not expire.
//
// Note that the function <f> is executed within writing mutex lock.
func (c *memCache) GetOrSetFuncLock(key interface{}, f func() interface{}, duration interface{}) interface{} {
	if v := c.Get(key); v == nil {
		return c.doSetWithLockCheck(key, f, duration)
	} else {
		return v
	}
}

// Contains returns true if <key> exists in the cache, or else returns false.
func (c *memCache) Contains(key interface{}) bool {
	return c.Get(key) != nil
}

// Remove deletes the <key> in the cache, and returns its value.
func (c *memCache) Remove(key interface{}) (value interface{}) {
	c.dataMu.RLock()
	item, ok := c.data[key]
	c.dataMu.RUnlock()
	if ok {
		value = item.v
		c.dataMu.Lock()
		delete(c.data, key)
		c.dataMu.Unlock()
		c.eventList.PushBack(&memCacheEvent{k: key, e: gtime.Millisecond() - 1000})
	}
	return
}

// Removes deletes <keys> in the cache.
func (c *memCache) Removes(keys []interface{}) {
	for _, key := range keys {
		c.Remove(key)
	}
}

// Data returns a copy of all key-value pairs in the cache as map type.
func (c *memCache) Data() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	c.dataMu.RLock()
	for k, v := range c.data {
		if !v.IsExpired() {
			m[k] = v.v
		}
	}
	c.dataMu.RUnlock()
	return m
}

// Keys returns all keys in the cache as slice.
func (c *memCache) Keys() []interface{} {
	keys := make([]interface{}, 0)
	c.dataMu.RLock()
	for k, v := range c.data {
		if !v.IsExpired() {
			keys = append(keys, k)
		}
	}
	c.dataMu.RUnlock()
	return keys
}

// KeyStrings returns all keys in the cache as string slice.
func (c *memCache) KeyStrings() []string {
	return gconv.Strings(c.Keys())
}

// Values returns all values in the cache as slice.
func (c *memCache) Values() []interface{} {
	values := make([]interface{}, 0)
	c.dataMu.RLock()
	for _, v := range c.data {
		if !v.IsExpired() {
			values = append(values, v.v)
		}
	}
	c.dataMu.RUnlock()
	return values
}

// Size returns the size of the cache.
func (c *memCache) Size() (size int) {
	c.dataMu.RLock()
	size = len(c.data)
	c.dataMu.RUnlock()
	return
}

// Close closes the cache.
func (c *memCache) Close() {
	if c.cap > 0 {
		c.lru.Close()
	}
	c.closed.Set(true)
}

// Asynchronous task loop:
// 1. asynchronously process the data in the event list,
//    and synchronize the results to the <expireTimes> and <expireSets> properties.
// 2. clean up the expired key-value pair data.
func (c *memCache) syncEventAndClearExpired() {
	event := (*memCacheEvent)(nil)
	oldExpireTime := int64(0)
	newExpireTime := int64(0)
	if c.closed.Val() {
		gtimer.Exit()
		return
	}
	// ========================
	// Data Synchronization.
	// ========================
	for {
		v := c.eventList.PopFront()
		if v == nil {
			break
		}
		event = v.(*memCacheEvent)
		// Fetching the old expire set.
		c.expireTimeMu.RLock()
		oldExpireTime = c.expireTimes[event.k]
		c.expireTimeMu.RUnlock()
		// Calculating the new expire set.
		newExpireTime = c.makeExpireKey(event.e)
		if newExpireTime != oldExpireTime {
			c.getOrNewExpireSet(newExpireTime).Add(event.k)
			if oldExpireTime != 0 {
				c.getOrNewExpireSet(oldExpireTime).Remove(event.k)
			}
			// Updating the expire time for <event.k>.
			c.expireTimeMu.Lock()
			c.expireTimes[event.k] = newExpireTime
			c.expireTimeMu.Unlock()
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
	ek := c.makeExpireKey(gtime.Millisecond())
	eks := []int64{ek - 1000, ek - 2000, ek - 3000, ek - 4000, ek - 5000}
	for _, expireTime := range eks {
		if expireSet := c.getExpireSet(expireTime); expireSet != nil {
			// Iterating the set to delete all keys in it.
			expireSet.Iterator(func(key interface{}) bool {
				c.clearByKey(key)
				return true
			})
			// Deleting the set after all of its keys are deleted.
			c.expireSetMu.Lock()
			delete(c.expireSets, expireTime)
			c.expireSetMu.Unlock()
		}
	}
}

// clearByKey deletes the key-value pair with given <key>.
// The parameter <force> specifies whether doing this deleting forcely.
func (c *memCache) clearByKey(key interface{}, force ...bool) {
	c.dataMu.Lock()
	// Doubly check before really deleting it from cache.
	if item, ok := c.data[key]; (ok && item.IsExpired()) || (len(force) > 0 && force[0]) {
		delete(c.data, key)
	}
	c.dataMu.Unlock()

	// Deleting its expire time from <expireTimes>.
	c.expireTimeMu.Lock()
	delete(c.expireTimes, key)
	c.expireTimeMu.Unlock()

	// Deleting it from LRU.
	if c.cap > 0 {
		c.lru.Remove(key)
	}
}

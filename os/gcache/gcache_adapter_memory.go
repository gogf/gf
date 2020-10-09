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

	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gtimer"
)

// Internal cache object.
type adapterMemory struct {
	// dataMu ensures the concurrent safety of underlying data map.
	dataMu sync.RWMutex

	// expireTimeMu ensures the concurrent safety of expireTimes map.
	expireTimeMu sync.RWMutex

	// expireSetMu ensures the concurrent safety of expireSets map.
	expireSetMu sync.RWMutex

	// cap limits the size of the cache pool.
	// If the size of the cache exceeds the cap,
	// the cache expiration process performs according to the LRU algorithm.
	// It is 0 in default which means no limits.
	cap int

	// data is the underlying cache data which is stored in a hash table.
	data map[interface{}]adapterMemoryItem

	// expireTimes is the expiring key to its timestamp mapping,
	// which is used for quick indexing and deleting.
	expireTimes map[interface{}]int64

	// expireSets is the expiring timestamp to its key set mapping,
	// which is used for quick indexing and deleting.
	expireSets map[int64]*gset.Set

	// lru is the LRU manager, which is enabled when attribute cap > 0.
	lru *adapterMemoryLru

	// lruGetList is the LRU history according with Get function.
	lruGetList *glist.List

	// eventList is the asynchronous event list for internal data synchronization.
	eventList *glist.List

	// closed controls the cache closed or not.
	closed *gtype.Bool
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
	// gDEFAULT_MAX_EXPIRE is the default expire time for no expiring items.
	// It equals to math.MaxInt64/1000000.
	gDEFAULT_MAX_EXPIRE = 9223372036854
)

// newAdapterMemory creates and returns a new memory cache object.
func newAdapterMemory(lruCap ...int) *adapterMemory {
	c := &adapterMemory{
		lruGetList:  glist.New(true),
		data:        make(map[interface{}]adapterMemoryItem),
		expireTimes: make(map[interface{}]int64),
		expireSets:  make(map[int64]*gset.Set),
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
	c.dataMu.Lock()
	c.data[key] = adapterMemoryItem{
		v: value,
		e: expireTime,
	}
	c.dataMu.Unlock()
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
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	if item, ok := c.data[key]; ok {
		c.data[key] = adapterMemoryItem{
			v: value,
			e: item.e,
		}
		return item.v, true, nil
	}
	return nil, false, nil
}

// UpdateExpire updates the expiration of <key> and returns the old expiration duration value.
//
// It returns -1 and does nothing if the <key> does not exist in the cache.
// It deletes the <key> if <duration> < 0.
func (c *adapterMemory) UpdateExpire(key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	newExpireTime := c.getInternalExpire(duration)
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	if item, ok := c.data[key]; ok {
		c.data[key] = adapterMemoryItem{
			v: item.v,
			e: newExpireTime,
		}
		c.eventList.PushBack(&adapterMemoryEvent{
			k: key,
			e: newExpireTime,
		})
		return time.Duration(item.e-gtime.TimestampMilli()) * time.Millisecond, nil
	}
	return -1, nil
}

// GetExpire retrieves and returns the expiration of <key> in the cache.
//
// It returns 0 if the <key> does not expire.
// It returns -1 if the <key> does not exist in the cache.
func (c *adapterMemory) GetExpire(key interface{}) (time.Duration, error) {
	c.dataMu.RLock()
	defer c.dataMu.RUnlock()
	if item, ok := c.data[key]; ok {
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
	expireTime := c.getInternalExpire(duration)
	for k, v := range data {
		c.dataMu.Lock()
		c.data[k] = adapterMemoryItem{
			v: v,
			e: expireTime,
		}
		c.dataMu.Unlock()
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
	c.dataMu.RLock()
	item, ok := c.data[key]
	c.dataMu.RUnlock()
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
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	for _, key := range keys {
		item, ok := c.data[key]
		if ok {
			value = item.v
			delete(c.data, key)
			c.eventList.PushBack(&adapterMemoryEvent{
				k: key,
				e: gtime.TimestampMilli() - 1000,
			})
		}
	}
	return value, nil
}

// Data returns a copy of all key-value pairs in the cache as map type.
func (c *adapterMemory) Data() (map[interface{}]interface{}, error) {
	m := make(map[interface{}]interface{})
	c.dataMu.RLock()
	for k, v := range c.data {
		if !v.IsExpired() {
			m[k] = v.v
		}
	}
	c.dataMu.RUnlock()
	return m, nil
}

// Keys returns all keys in the cache as slice.
func (c *adapterMemory) Keys() ([]interface{}, error) {
	keys := make([]interface{}, 0)
	c.dataMu.RLock()
	for k, v := range c.data {
		if !v.IsExpired() {
			keys = append(keys, k)
		}
	}
	c.dataMu.RUnlock()
	return keys, nil
}

// Values returns all values in the cache as slice.
func (c *adapterMemory) Values() ([]interface{}, error) {
	values := make([]interface{}, 0)
	c.dataMu.RLock()
	for _, v := range c.data {
		if !v.IsExpired() {
			values = append(values, v.v)
		}
	}
	c.dataMu.RUnlock()
	return values, nil
}

// Size returns the size of the cache.
func (c *adapterMemory) Size() (size int, err error) {
	c.dataMu.RLock()
	size = len(c.data)
	c.dataMu.RUnlock()
	return size, nil
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
func (c *adapterMemory) Clear() error {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	c.data = make(map[interface{}]adapterMemoryItem)
	return nil
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
func (c *adapterMemory) doSetWithLockCheck(key interface{}, value interface{}, duration time.Duration) (interface{}, error) {
	expireTimestamp := c.getInternalExpire(duration)
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	if v, ok := c.data[key]; ok && !v.IsExpired() {
		return v.v, nil
	}
	if f, ok := value.(func() (interface{}, error)); ok {
		v, err := f()
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, nil
		} else {
			value = v
		}
	}
	c.data[key] = adapterMemoryItem{v: value, e: expireTimestamp}
	c.eventList.PushBack(&adapterMemoryEvent{k: key, e: expireTimestamp})
	return value, nil
}

// getInternalExpire converts and returns the expire time with given expired duration in milliseconds.
func (c *adapterMemory) getInternalExpire(duration time.Duration) int64 {
	if duration == 0 {
		return gDEFAULT_MAX_EXPIRE
	} else {
		return gtime.TimestampMilli() + duration.Nanoseconds()/1000000
	}
}

// makeExpireKey groups the <expire> in milliseconds to its according seconds.
func (c *adapterMemory) makeExpireKey(expire int64) int64 {
	return int64(math.Ceil(float64(expire/1000)+1) * 1000)
}

// getExpireSet returns the expire set for given <expire> in seconds.
func (c *adapterMemory) getExpireSet(expire int64) (expireSet *gset.Set) {
	c.expireSetMu.RLock()
	expireSet, _ = c.expireSets[expire]
	c.expireSetMu.RUnlock()
	return
}

// getOrNewExpireSet returns the expire set for given <expire> in seconds.
// It creates and returns a new set for <expire> if it does not exist.
func (c *adapterMemory) getOrNewExpireSet(expire int64) (expireSet *gset.Set) {
	if expireSet = c.getExpireSet(expire); expireSet == nil {
		expireSet = gset.New(true)
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
	var (
		expireSet *gset.Set
		ek        = c.makeExpireKey(gtime.TimestampMilli())
		eks       = []int64{ek - 1000, ek - 2000, ek - 3000, ek - 4000, ek - 5000}
	)
	for _, expireTime := range eks {
		if expireSet = c.getExpireSet(expireTime); expireSet != nil {
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
// The parameter <force> specifies whether doing this deleting forcibly.
func (c *adapterMemory) clearByKey(key interface{}, force ...bool) {
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

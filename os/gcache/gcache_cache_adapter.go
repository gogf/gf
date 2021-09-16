// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"time"
)

// Set sets cache with `key`-`value` pair, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *Cache) Set(key interface{}, value interface{}, duration time.Duration) error {
	return c.adapter.Set(c.getCtx(), key, value, duration)
}

// Sets batch sets cache with key-value pairs by `data` map, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *Cache) Sets(data map[interface{}]interface{}, duration time.Duration) error {
	return c.adapter.Sets(c.getCtx(), data, duration)
}

// SetIfNotExist sets cache with `key`-`value` pair which is expired after `duration`
// if `key` does not exist in the cache. It returns true the `key` does not exist in the
// cache, and it sets `value` successfully to the cache, or else it returns false.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *Cache) SetIfNotExist(key interface{}, value interface{}, duration time.Duration) (bool, error) {
	return c.adapter.SetIfNotExist(c.getCtx(), key, value, duration)
}

// SetIfNotExistFunc sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// The parameter `value` can be type of `func() interface{}`, but it does nothing if its
// result is nil.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *Cache) SetIfNotExistFunc(key interface{}, f func() (interface{}, error), duration time.Duration) (bool, error) {
	return c.adapter.SetIfNotExistFunc(c.getCtx(), key, f, duration)
}

// SetIfNotExistFuncLock sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
//
// Note that it differs from function `SetIfNotExistFunc` is that the function `f` is executed within
// writing mutex lock for concurrent safety purpose.
func (c *Cache) SetIfNotExistFuncLock(key interface{}, f func() (interface{}, error), duration time.Duration) (bool, error) {
	return c.adapter.SetIfNotExistFuncLock(c.getCtx(), key, f, duration)
}

// Get retrieves and returns the associated value of given `key`.
// It returns nil if it does not exist, or its value is nil, or it's expired.
// If you would like to check if the `key` exists in the cache, it's better using function Contains.
//
// It is suggested using GetVar instead for compatibility of different adapters purpose.
func (c *Cache) Get(key interface{}) (interface{}, error) {
	return c.adapter.Get(c.getCtx(), key)
}

// GetOrSet retrieves and returns the value of `key`, or sets `key`-`value` pair and
// returns `value` if `key` does not exist in the cache. The key-value pair expires
// after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
//
// It is suggested using GetVarOrSet instead for compatibility of different adapters purpose.
func (c *Cache) GetOrSet(key interface{}, value interface{}, duration time.Duration) (interface{}, error) {
	return c.adapter.GetOrSet(c.getCtx(), key, value, duration)
}

// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
//
// It is suggested using GetVarOrSetFunc instead for compatibility of different adapters purpose.
func (c *Cache) GetOrSetFunc(key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	return c.adapter.GetOrSetFunc(c.getCtx(), key, f, duration)
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
//
// It is suggested using GetVarOrSetFuncLock instead for compatibility of different adapters purpose.
func (c *Cache) GetOrSetFuncLock(key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	return c.adapter.GetOrSetFuncLock(c.getCtx(), key, f, duration)
}

// Contains checks and returns true if `key` exists in the cache, or else returns false.
func (c *Cache) Contains(key interface{}) (bool, error) {
	return c.adapter.Contains(c.getCtx(), key)
}

// GetExpire retrieves and returns the expiration of `key` in the cache.
//
// Note that,
// It returns 0 if the `key` does not expire.
// It returns -1 if the `key` does not exist in the cache.
func (c *Cache) GetExpire(key interface{}) (time.Duration, error) {
	return c.adapter.GetExpire(c.getCtx(), key)
}

// Remove deletes one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the last deleted item.
//
// It is suggested using RemoveVar instead for compatibility of different adapters purpose.
func (c *Cache) Remove(keys ...interface{}) (value interface{}, err error) {
	return c.adapter.Remove(c.getCtx(), keys...)
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func (c *Cache) Update(key interface{}, value interface{}) (oldValue interface{}, exist bool, err error) {
	return c.adapter.Update(c.getCtx(), key, value)
}

// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func (c *Cache) UpdateExpire(key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	return c.adapter.UpdateExpire(c.getCtx(), key, duration)
}

// Size returns the number of items in the cache.
func (c *Cache) Size() (size int, err error) {
	return c.adapter.Size(c.getCtx())
}

// Data returns a copy of all key-value pairs in the cache as map type.
// Note that this function may lead lots of memory usage, you can implement this function
// if necessary.
func (c *Cache) Data() (map[interface{}]interface{}, error) {
	return c.adapter.Data(c.getCtx())
}

// Keys returns all keys in the cache as slice.
func (c *Cache) Keys() ([]interface{}, error) {
	return c.adapter.Keys(c.getCtx())
}

// Values returns all values in the cache as slice.
func (c *Cache) Values() ([]interface{}, error) {
	return c.adapter.Values(c.getCtx())
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
func (c *Cache) Clear() error {
	return c.adapter.Clear(c.getCtx())
}

// Close closes the cache if necessary.
func (c *Cache) Close() error {
	return c.adapter.Close(c.getCtx())
}

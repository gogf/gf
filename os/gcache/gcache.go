// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcache provides kinds of cache management for process.
//
// It provides a concurrent-safe in-memory cache adapter for process in default.
package gcache

import (
	"context"
	"github.com/gogf/gf/container/gvar"
	"time"
)

// Default cache object.
var defaultCache = New()

// Ctx is a chaining function, which shallowly clones current object and sets the context
// for next operation.
func Ctx(ctx context.Context) *Cache {
	return defaultCache.Ctx(ctx)
}

// Set sets cache with `key`-`value` pair, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func Set(key interface{}, value interface{}, duration time.Duration) error {
	return defaultCache.Set(key, value, duration)
}

// Sets batch sets cache with key-value pairs by `data` map, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func Sets(data map[interface{}]interface{}, duration time.Duration) error {
	return defaultCache.Sets(data, duration)
}

// SetIfNotExist sets cache with `key`-`value` pair which is expired after `duration`
// if `key` does not exist in the cache. It returns true the `key` does not exist in the
// cache, and it sets `value` successfully to the cache, or else it returns false.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func SetIfNotExist(key interface{}, value interface{}, duration time.Duration) (bool, error) {
	return defaultCache.SetIfNotExist(key, value, duration)
}

// SetIfNotExistFunc sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// The parameter `value` can be type of `func() interface{}`, but it does nothing if its
// result is nil.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func SetIfNotExistFunc(key interface{}, f func() (interface{}, error), duration time.Duration) (bool, error) {
	return defaultCache.SetIfNotExistFunc(key, f, duration)
}

// SetIfNotExistFuncLock sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
//
// Note that it differs from function `SetIfNotExistFunc` is that the function `f` is executed within
// writing mutex lock for concurrent safety purpose.
func SetIfNotExistFuncLock(key interface{}, f func() (interface{}, error), duration time.Duration) (bool, error) {
	return defaultCache.SetIfNotExistFuncLock(key, f, duration)
}

// Get retrieves and returns the associated value of given `key`.
// It returns nil if it does not exist, or its value is nil, or it's expired.
// If you would like to check if the `key` exists in the cache, it's better using function Contains.
func Get(key interface{}) (interface{}, error) {
	return defaultCache.Get(key)
}

// GetVar retrieves and returns the value of `key` as gvar.Var.
func GetVar(key interface{}) (*gvar.Var, error) {
	return defaultCache.GetVar(key)
}

// GetOrSet retrieves and returns the value of `key`, or sets `key`-`value` pair and
// returns `value` if `key` does not exist in the cache. The key-value pair expires
// after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func GetOrSet(key interface{}, value interface{}, duration time.Duration) (interface{}, error) {
	return defaultCache.GetOrSet(key, value, duration)
}

// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func GetOrSetFunc(key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	return defaultCache.GetOrSetFunc(key, f, duration)
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
func GetOrSetFuncLock(key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	return defaultCache.GetOrSetFuncLock(key, f, duration)
}

// Contains checks and returns true if `key` exists in the cache, or else returns false.
func Contains(key interface{}) (bool, error) {
	return defaultCache.Contains(key)
}

// GetExpire retrieves and returns the expiration of `key` in the cache.
//
// Note that,
// It returns 0 if the `key` does not expire.
// It returns -1 if the `key` does not exist in the cache.
func GetExpire(key interface{}) (time.Duration, error) {
	return defaultCache.GetExpire(key)
}

// Remove deletes one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the last deleted item.
func Remove(keys ...interface{}) (value interface{}, err error) {
	return defaultCache.Remove(keys...)
}

// Removes deletes `keys` in the cache.
func Removes(keys []interface{}) error {
	return defaultCache.Removes(keys)
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func Update(key interface{}, value interface{}) (oldValue interface{}, exist bool, err error) {
	return defaultCache.Update(key, value)
}

// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func UpdateExpire(key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	return defaultCache.UpdateExpire(key, duration)
}

// Size returns the number of items in the cache.
func Size() (int, error) {
	return defaultCache.Size()
}

// Data returns a copy of all key-value pairs in the cache as map type.
// Note that this function may lead lots of memory usage, you can implement this function
// if necessary.
func Data() (map[interface{}]interface{}, error) {
	return defaultCache.Data()
}

// Keys returns all keys in the cache as slice.
func Keys() ([]interface{}, error) {
	return defaultCache.Keys()
}

// KeyStrings returns all keys in the cache as string slice.
func KeyStrings() ([]string, error) {
	return defaultCache.KeyStrings()
}

// Values returns all values in the cache as slice.
func Values() ([]interface{}, error) {
	return defaultCache.Values()
}

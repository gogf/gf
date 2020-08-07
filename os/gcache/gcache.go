// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcache provides high performance and concurrent-safe in-memory cache for process.
package gcache

import (
	"github.com/gogf/gf/container/gvar"
	"time"
)

// Default cache object.
var cache = New()

// Set sets cache with <key>-<value> pair, which is expired after <duration>.
// It does not expire if <duration> == 0.
func Set(key interface{}, value interface{}, duration time.Duration) {
	cache.Set(key, value, duration)
}

// SetVar edit cache with <key>-<value> , but will not change expire.
// If key is not exist return false. Nothing will change.
func SetVar(key interface{}, value interface{}) bool {
	return cache.SetVar(key, value)
}

// SetIfNotExist sets cache with <key>-<value> pair if <key> does not exist in the cache,
// which is expired after <duration>. It does not expire if <duration> == 0.
func SetIfNotExist(key interface{}, value interface{}, duration time.Duration) bool {
	return cache.SetIfNotExist(key, value, duration)
}

// Sets batch sets cache with key-value pairs by <data>, which is expired after <duration>.
//
// It does not expire if <duration> == 0.
func Sets(data map[interface{}]interface{}, duration time.Duration) {
	cache.Sets(data, duration)
}

// Get returns the value of <key>.
// It returns nil if it does not exist or its value is nil.
func Get(key interface{}) interface{} {
	return cache.Get(key)
}

// GetVar retrieves and returns the value of <key> as gvar.Var.
func GetVar(key interface{}) *gvar.Var {
	return cache.GetVar(key)
}

// GetOrSet returns the value of <key>,
// or sets <key>-<value> pair and returns <value> if <key> does not exist in the cache.
// The key-value pair expires after <duration>.
//
// It does not expire if <duration> == 0.
func GetOrSet(key interface{}, value interface{}, duration time.Duration) interface{} {
	return cache.GetOrSet(key, value, duration)
}

// GetOrSetFunc returns the value of <key>, or sets <key> with result of function <f>
// and returns its result if <key> does not exist in the cache. The key-value pair expires
// after <duration>. It does not expire if <duration> == 0.
func GetOrSetFunc(key interface{}, f func() interface{}, duration time.Duration) interface{} {
	return cache.GetOrSetFunc(key, f, duration)
}

// GetOrSetFuncLock returns the value of <key>, or sets <key> with result of function <f>
// and returns its result if <key> does not exist in the cache. The key-value pair expires
// after <duration>. It does not expire if <duration> == 0.
//
// Note that the function <f> is executed within writing mutex lock.
func GetOrSetFuncLock(key interface{}, f func() interface{}, duration time.Duration) interface{} {
	return cache.GetOrSetFuncLock(key, f, duration)
}

// Contains returns true if <key> exists in the cache, or else returns false.
func Contains(key interface{}) bool {
	return cache.Contains(key)
}

// Remove deletes the one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the deleted last item.
func Remove(keys ...interface{}) (value interface{}) {
	return cache.Remove(keys...)
}

// Removes deletes <keys> in the cache.
// Deprecated, use Remove instead.
func Removes(keys []interface{}) {
	cache.Removes(keys)
}

// Data returns a copy of all key-value pairs in the cache as map type.
func Data() map[interface{}]interface{} {
	return cache.Data()
}

// Keys returns all keys in the cache as slice.
func Keys() []interface{} {
	return cache.Keys()
}

// KeyStrings returns all keys in the cache as string slice.
func KeyStrings() []string {
	return cache.KeyStrings()
}

// Values returns all values in the cache as slice.
func Values() []interface{} {
	return cache.Values()
}

// Size returns the size of the cache.
func Size() int {
	return cache.Size()
}

// GetExpire returns the expire time with given expired duration in milliseconds.
func GetExpire(key interface{}) (int64, bool) {
	return cache.GetExpire(key)
}

// SetExpire set cache expired after <duration>.
func SetExpire(key interface{}, duration time.Duration) {
	cache.SetExpire(key, duration)
}

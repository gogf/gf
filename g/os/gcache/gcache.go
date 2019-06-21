// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcache provides high performance and concurrent-safe in-memory cache for process.
package gcache

// Default cache object.
var cache = New()

// Set sets cache with <key>-<value> pair, which is expired after <expire> milliseconds.
// If <expire> <=0 means it does not expire.
func Set(key interface{}, value interface{}, expire int) {
	cache.Set(key, value, expire)
}

// SetIfNotExist sets cache with <key>-<value> pair if <key> does not exist in the cache,
// which is expired after <expire> milliseconds.
// If <expire> <=0 means it does not expire.
func SetIfNotExist(key interface{}, value interface{}, expire int) bool {
	return cache.SetIfNotExist(key, value, expire)
}

// Sets batch sets cache with key-value pairs by <data>, which is expired after <expire> milliseconds.
// If <expire> <=0 means it does not expire.
func Sets(data map[interface{}]interface{}, expire int) {
	cache.Sets(data, expire)
}

// Get returns the value of <key>.
// It returns nil if it does not exist or its value is nil.
func Get(key interface{}) interface{} {
	return cache.Get(key)
}

// GetOrSet returns the value of <key>,
// or sets <key>-<value> pair and returns <value> if <key> does not exist in the cache.
// The key-value pair expires after <expire> milliseconds.
// If <expire> <=0 means it does not expire.
func GetOrSet(key interface{}, value interface{}, expire int) interface{} {
	return cache.GetOrSet(key, value, expire)
}

// GetOrSetFunc returns the value of <key>,
// or sets <key> with result of function <f> and returns its result
// if <key> does not exist in the cache.
// The key-value pair expires after <expire> milliseconds.
// If <expire> <=0 means it does not expire.
func GetOrSetFunc(key interface{}, f func() interface{}, expire int) interface{} {
	return cache.GetOrSetFunc(key, f, expire)
}

// GetOrSetFuncLock returns the value of <key>,
// or sets <key> with result of function <f> and returns its result
// if <key> does not exist in the cache.
// The key-value pair expires after <expire> milliseconds.
// If <expire> <=0 means it does not expire.
//
// Note that the function <f> is executed within writing mutex lock.
func GetOrSetFuncLock(key interface{}, f func() interface{}, expire int) interface{} {
	return cache.GetOrSetFuncLock(key, f, expire)
}

// Contains returns true if <key> exists in the cache, or else returns false.
func Contains(key interface{}) bool {
	return cache.Contains(key)
}

// Remove deletes the <key> in the cache, and returns its value.
func Remove(key interface{}) interface{} {
	return cache.Remove(key)
}

// Removes deletes <keys> in the cache.
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

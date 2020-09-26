// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"time"
)

// Adapter is the adapter for cache features implements.
type Adapter interface {
	// Set sets cache with <key>-<value> pair, which is expired after <duration>.
	//
	// It does not expire if <duration> == 0.
	// It deletes the <key> if <duration> < 0.
	Set(key interface{}, value interface{}, duration time.Duration)

	// Sets batch sets cache with key-value pairs by <data>, which is expired after <duration>.
	//
	// It does not expire if <duration> == 0.
	// It deletes the keys of <data> if <duration> < 0 or given <value> is nil.
	Sets(data map[interface{}]interface{}, duration time.Duration)

	// SetIfNotExist sets cache with <key>-<value> pair which is expired after <duration>
	// if <key> does not exist in the cache.
	// The parameter <value> can be type of <func() interface{}>, but it dose nothing if its
	// result is nil.
	//
	// It does not expire if <duration> == 0.
	// It deletes the <key> if <duration> < 0 or given <value> is nil.
	SetIfNotExist(key interface{}, value interface{}, duration time.Duration) bool

	// Get retrieves and returns the associated value of given <key>.
	// It returns nil if it does not exist or its value is nil.
	Get(key interface{}) interface{}

	// GetOrSet retrieves and returns the value of <key>, or sets <key>-<value> pair and
	// returns <value> if <key> does not exist in the cache. The key-value pair expires
	// after <duration>.
	//
	// It does not expire if <duration> == 0.
	// It deletes the <key> if <duration> < 0 or given <value> is nil, but it does nothing
	// if <value> is a function and the function result is nil.
	GetOrSet(key interface{}, value interface{}, duration time.Duration) interface{}

	// GetOrSetFunc retrieves and returns the value of <key>, or sets <key> with result of
	// function <f> and returns its result if <key> does not exist in the cache. The key-value
	// pair expires after <duration>.
	//
	// It does not expire if <duration> == 0.
	// It deletes the <key> if <duration> < 0 or given <value> is nil, but it does nothing
	// if <value> is a function and the function result is nil.
	GetOrSetFunc(key interface{}, f func() interface{}, duration time.Duration) interface{}

	// GetOrSetFuncLock retrieves and returns the value of <key>, or sets <key> with result of
	// function <f> and returns its result if <key> does not exist in the cache. The key-value
	// pair expires after <duration>.
	//
	// It does not expire if <duration> == 0.
	// It does nothing if function <f> returns nil.
	//
	// Note that the function <f> should be executed within writing mutex lock for concurrent
	// safety purpose.
	GetOrSetFuncLock(key interface{}, f func() interface{}, duration time.Duration) interface{}

	// Remove deletes one or more keys from cache, and returns its value.
	// If multiple keys are given, it returns the value of the last deleted item.
	Remove(keys ...interface{}) (value interface{})

	// Update updates the value of <key> without changing its expiration and returns the old value.
	// The returned value <exist> is false if the <key> does not exist in the cache.
	Update(key interface{}, value interface{}) (oldValue interface{}, exist bool)

	// UpdateExpire updates the expiration of <key> and returns the old expiration duration value.
	// It returns -1 if the <key> does not exist in the cache.
	UpdateExpire(key interface{}, duration time.Duration) (oldDuration time.Duration)

	// Contains checks and returns whether given <key> exists in the cache.
	Contains(key interface{}) bool

	// GetExpire retrieves and returns the expiration of <key> in the cache.
	// It returns -1 if the <key> does not exist in the cache.
	GetExpire(key interface{}) time.Duration

	// Size returns the number of items in the cache.
	Size() (size int)

	// Data returns a copy of all key-value pairs in the cache as map type.
	// Note that this function may leads lots of memory usage, you can implement this function
	// if necessary.
	Data() map[interface{}]interface{}

	// Keys returns all keys in the cache as slice.
	Keys() []interface{}

	// Values returns all values in the cache as slice.
	Values() []interface{}

	// Clear clears all data of the cache.
	// Note that this function is sensitive and should be carefully used.
	Clear() error

	// Close closes the cache if necessary.
	Close() error
}

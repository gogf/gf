package gcache_test

import (
	"context"
	"github.com/gogf/gf/v2/os/gcache"

	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"time"
)

// Test cache is the gcache adapter implements using Redis server.
type T_cache struct {
	cache *gcache.AdapterMemory
}

// NewCache creates and returns a new cache memory cache object.
func NewCache(redis *gcache.AdapterMemory) gcache.Adapter {
	return &T_cache{
		cache: redis,
	}
}

// Set sets cache with `key`-`value` pair, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (r *T_cache) Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error {
	fmt.Println("111111111")
	return r.cache.Set(ctx, key, value, duration)
}

// SetMap batch sets cache with key-value pairs by `data` map, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (r *T_cache) SetMap(ctx context.Context, data map[interface{}]interface{}, duration time.Duration) error {
	return r.cache.SetMap(ctx, data, duration)
}

// SetIfNotExist sets cache with `key`-`value` pair which is expired after `duration`
// if `key` does not exist in the cache. It returns true the `key` does not exist in the
// cache, and it sets `value` successfully to the cache, or else it returns false.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (r *T_cache) SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (ok bool, err error) {
	return r.cache.SetIfNotExist(ctx, key, value, duration)
}

// SetIfNotExistFunc sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// The parameter `value` can be type of `func() interface{}`, but it does nothing if its
// result is nil.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (r *T_cache) SetIfNotExistFunc(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (ok bool, err error) {
	return r.cache.SetIfNotExistFunc(ctx, key, f, duration)
}

// SetIfNotExistFuncLock sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
//
// Note that it differs from function `SetIfNotExistFunc` is that the function `f` is executed within
// writing mutex lock for concurrent safety purpose.
func (r *T_cache) SetIfNotExistFuncLock(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (ok bool, err error) {
	return r.cache.SetIfNotExistFuncLock(ctx, key, f, duration)
}

// Get retrieves and returns the associated value of given `key`.
// It returns nil if it does not exist, or its value is nil, or it's expired.
// If you would like to check if the `key` exists in the cache, it's better using function Contains.
func (r *T_cache) Get(ctx context.Context, key interface{}) (*gvar.Var, error) {
	return r.cache.Get(ctx, key)
}

// GetOrSet retrieves and returns the value of `key`, or sets `key`-`value` pair and
// returns `value` if `key` does not exist in the cache. The key-value pair expires
// after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (r *T_cache) GetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (result *gvar.Var, err error) {
	return r.cache.GetOrSet(ctx, key, value, duration)
}

// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (r *T_cache) GetOrSetFunc(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (result *gvar.Var, err error) {
	return r.cache.GetOrSetFunc(ctx, key, f, duration)
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
func (r *T_cache) GetOrSetFuncLock(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (result *gvar.Var, err error) {
	return r.cache.GetOrSetFuncLock(ctx, key, f, duration)
}

// Contains checks and returns true if `key` exists in the cache, or else returns false.
func (r *T_cache) Contains(ctx context.Context, key interface{}) (bool, error) {
	return r.cache.Contains(ctx, key)
}

// Size returns the number of items in the cache.
func (r *T_cache) Size(ctx context.Context) (size int, err error) {
	return r.cache.Size(ctx)
}

// Data returns a copy of all key-value pairs in the cache as map type.
// Note that this function may lead lots of memory usage, you can implement this function
// if necessary.
func (r *T_cache) Data(ctx context.Context) (data map[interface{}]interface{}, err error) {
	return r.cache.Data(ctx)
}

// Keys returns all keys in the cache as slice.
func (r *T_cache) Keys(ctx context.Context) (keys []interface{}, err error) {
	return r.cache.Keys(ctx)
}

// Values returns all values in the cache as slice.
func (r *T_cache) Values(ctx context.Context) (values []interface{}, err error) {
	return r.cache.Values(ctx)
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func (r *T_cache) Update(ctx context.Context, key interface{}, value interface{}) (oldValue *gvar.Var, exist bool, err error) {
	return r.cache.Update(ctx, key, value)
}

// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func (r *T_cache) UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	return r.cache.UpdateExpire(ctx, key, duration)
}

// GetExpire retrieves and returns the expiration of `key` in the cache.
//
// Note that,
// It returns 0 if the `key` does not expire.
// It returns -1 if the `key` does not exist in the cache.
func (r *T_cache) GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	return r.cache.GetExpire(ctx, key)
}

// Remove deletes one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the last deleted item.
func (r *T_cache) Remove(ctx context.Context, keys ...interface{}) (lastValue *gvar.Var, err error) {
	return r.cache.Remove(ctx, keys)
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
func (r *T_cache) Clear(ctx context.Context) error {
	return r.cache.Clear(ctx)
}

// Close closes the cache if necessary.
func (r *T_cache) Close(ctx context.Context) error {
	return r.cache.Close(ctx)
}

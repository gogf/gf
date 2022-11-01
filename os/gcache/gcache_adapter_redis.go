// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
)

// AdapterRedis is the gcache adapter implements using Redis server.
type AdapterRedis struct {
	redis *gredis.Redis
}

// NewAdapterRedis creates and returns a new memory cache object.
func NewAdapterRedis(redis *gredis.Redis) Adapter {
	return &AdapterRedis{
		redis: redis,
	}
}

// Set sets cache with `key`-`value` pair, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *AdapterRedis) Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (err error) {
	if value == nil || duration < 0 {
		_, err = c.redis.Do(ctx, "DEL", key)
	} else {
		if duration == 0 {
			_, err = c.redis.Do(ctx, "SET", key, value)
		} else {
			_, err = c.redis.Do(ctx, "SETEX", key, uint64(duration.Seconds()), value)
		}
	}
	return err
}

// SetMap batch sets cache with key-value pairs by `data` map, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *AdapterRedis) SetMap(ctx context.Context, data map[interface{}]interface{}, duration time.Duration) error {
	if len(data) == 0 {
		return nil
	}
	// DEL.
	if duration < 0 {
		var (
			index = 0
			keys  = make([]interface{}, len(data))
		)
		for k := range data {
			keys[index] = k
			index += 1
		}
		_, err := c.redis.Do(ctx, "DEL", keys...)
		if err != nil {
			return err
		}
	}
	if duration == 0 {
		var (
			index     = 0
			keyValues = make([]interface{}, len(data)*2)
		)
		for k, v := range data {
			keyValues[index] = k
			keyValues[index+1] = v
			index += 2
		}
		_, err := c.redis.Do(ctx, "MSET", keyValues...)
		if err != nil {
			return err
		}
	}
	if duration > 0 {
		var err error
		for k, v := range data {
			if err = c.Set(ctx, k, v, duration); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetIfNotExist sets cache with `key`-`value` pair which is expired after `duration`
// if `key` does not exist in the cache. It returns true the `key` does not exist in the
// cache, and it sets `value` successfully to the cache, or else it returns false.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *AdapterRedis) SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (bool, error) {
	var (
		v   *gvar.Var
		err error
	)
	// Execute the function and retrieve the result.
	f, ok := value.(Func)
	if !ok {
		// Compatible with raw function value.
		f, ok = value.(func(ctx context.Context) (value interface{}, err error))
	}
	if ok {
		if value, err = f(ctx); err != nil {
			return false, err
		}
	}
	// DEL.
	if duration < 0 || value == nil {
		if v, err = c.redis.Do(ctx, "DEL", key, value); err != nil {
			return false, err
		}
		if v.Int() == 1 {
			return true, err
		} else {
			return false, err
		}
	}
	if v, err = c.redis.Do(ctx, "SETNX", key, value); err != nil {
		return false, err
	}
	if v.Int() > 0 && duration > 0 {
		// Set the expiration.
		_, err = c.redis.Do(ctx, "EXPIRE", key, uint64(duration.Seconds()))
		if err != nil {
			return false, err
		}
		return true, err
	}
	return false, err
}

// SetIfNotExistFunc sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// The parameter `value` can be type of `func() interface{}`, but it does nothing if its
// result is nil.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *AdapterRedis) SetIfNotExistFunc(ctx context.Context, key interface{}, f Func, duration time.Duration) (ok bool, err error) {
	value, err := f(ctx)
	if err != nil {
		return false, err
	}
	return c.SetIfNotExist(ctx, key, value, duration)
}

// SetIfNotExistFuncLock sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
//
// Note that it differs from function `SetIfNotExistFunc` is that the function `f` is executed within
// writing mutex lock for concurrent safety purpose.
func (c *AdapterRedis) SetIfNotExistFuncLock(ctx context.Context, key interface{}, f Func, duration time.Duration) (ok bool, err error) {
	value, err := f(ctx)
	if err != nil {
		return false, err
	}
	return c.SetIfNotExist(ctx, key, value, duration)
}

// Get retrieves and returns the associated value of given <key>.
// It returns nil if it does not exist or its value is nil.
func (c *AdapterRedis) Get(ctx context.Context, key interface{}) (*gvar.Var, error) {
	return c.redis.Do(ctx, "GET", key)
}

// GetOrSet retrieves and returns the value of `key`, or sets `key`-`value` pair and
// returns `value` if `key` does not exist in the cache. The key-value pair expires
// after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (c *AdapterRedis) GetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (result *gvar.Var, err error) {
	result, err = c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if result.IsNil() {
		return gvar.New(value), c.Set(ctx, key, value, duration)
	}
	return
}

// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (c *AdapterRedis) GetOrSetFunc(ctx context.Context, key interface{}, f Func, duration time.Duration) (result *gvar.Var, err error) {
	v, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v.IsNil() {
		value, err := f(ctx)
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, nil
		}
		return gvar.New(value), c.Set(ctx, key, value, duration)
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
func (c *AdapterRedis) GetOrSetFuncLock(ctx context.Context, key interface{}, f Func, duration time.Duration) (result *gvar.Var, err error) {
	return c.GetOrSetFunc(ctx, key, f, duration)
}

// Contains checks and returns true if `key` exists in the cache, or else returns false.
func (c *AdapterRedis) Contains(ctx context.Context, key interface{}) (bool, error) {
	v, err := c.redis.Do(ctx, "EXISTS", key)
	if err != nil {
		return false, err
	}
	return v.Bool(), nil
}

// Size returns the number of items in the cache.
func (c *AdapterRedis) Size(ctx context.Context) (size int, err error) {
	v, err := c.redis.Do(ctx, "DBSIZE")
	if err != nil {
		return 0, err
	}
	return v.Int(), nil
}

// Data returns a copy of all key-value pairs in the cache as map type.
// Note that this function may lead lots of memory usage, you can implement this function
// if necessary.
func (c *AdapterRedis) Data(ctx context.Context) (map[interface{}]interface{}, error) {
	// Keys.
	v, err := c.redis.Do(ctx, "KEYS", "*")
	if err != nil {
		return nil, err
	}
	keys := v.Slice()
	// Values.
	v, err = c.redis.Do(ctx, "MGET", keys...)
	if err != nil {
		return nil, err
	}
	values := v.Slice()
	// Compose keys and values.
	data := make(map[interface{}]interface{})
	for i := 0; i < len(keys); i++ {
		data[keys[i]] = values[i]
	}
	return data, nil
}

// Keys returns all keys in the cache as slice.
func (c *AdapterRedis) Keys(ctx context.Context) ([]interface{}, error) {
	v, err := c.redis.Do(ctx, "KEYS", "*")
	if err != nil {
		return nil, err
	}
	return v.Slice(), nil
}

// Values returns all values in the cache as slice.
func (c *AdapterRedis) Values(ctx context.Context) ([]interface{}, error) {
	// Keys.
	v, err := c.redis.Do(ctx, "KEYS", "*")
	if err != nil {
		return nil, err
	}
	keys := v.Slice()
	// Values.
	v, err = c.redis.Do(ctx, "MGET", keys...)
	if err != nil {
		return nil, err
	}
	return v.Slice(), nil
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func (c *AdapterRedis) Update(ctx context.Context, key interface{}, value interface{}) (oldValue *gvar.Var, exist bool, err error) {
	var (
		v           *gvar.Var
		oldDuration time.Duration
	)
	// TTL.
	v, err = c.redis.Do(ctx, "TTL", key)
	if err != nil {
		return
	}
	oldDuration = v.Duration()
	if oldDuration == -2 {
		// It does not exist.
		return
	}
	// Check existence.
	v, err = c.redis.Do(ctx, "GET", key)
	if err != nil {
		return
	}
	oldValue = v
	// DEL.
	if value == nil {
		_, err = c.redis.Do(ctx, "DEL", key)
		if err != nil {
			return
		}
		return
	}
	// Update the value.
	if oldDuration == -1 {
		_, err = c.redis.Do(ctx, "SET", key, value)
	} else {
		oldDuration *= time.Second
		_, err = c.redis.Do(ctx, "SETEX", key, uint64(oldDuration.Seconds()), value)
	}
	return oldValue, true, err
}

// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func (c *AdapterRedis) UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	var (
		v *gvar.Var
	)
	// TTL.
	v, err = c.redis.Do(ctx, "TTL", key)
	if err != nil {
		return
	}
	oldDuration = v.Duration()
	if oldDuration == -2 {
		// It does not exist.
		oldDuration = -1
		return
	}
	oldDuration *= time.Second
	// DEL.
	if duration < 0 {
		_, err = c.redis.Do(ctx, "DEL", key)
		return
	}
	// Update the expire.
	if duration > 0 {
		_, err = c.redis.Do(ctx, "EXPIRE", key, uint64(duration.Seconds()))
	}
	// No expire.
	if duration == 0 {
		v, err = c.redis.Do(ctx, "GET", key)
		if err != nil {
			return
		}
		_, err = c.redis.Do(ctx, "SET", key, v.Val())
	}
	return
}

// GetExpire retrieves and returns the expiration of `key` in the cache.
//
// Note that,
// It returns 0 if the `key` does not expire.
// It returns -1 if the `key` does not exist in the cache.
func (c *AdapterRedis) GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	v, err := c.redis.Do(ctx, "TTL", key)
	if err != nil {
		return 0, err
	}
	switch v.Int() {
	case -1:
		return 0, nil
	case -2:
		return -1, nil
	default:
		return v.Duration() * time.Second, nil
	}
}

// Remove deletes the one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the deleted last item.
func (c *AdapterRedis) Remove(ctx context.Context, keys ...interface{}) (lastValue *gvar.Var, err error) {
	if len(keys) == 0 {
		return nil, nil
	}
	// Retrieves the last key value.
	if lastValue, err = c.redis.Do(ctx, "GET", keys[len(keys)-1]); err != nil {
		return nil, err
	}
	// Deletes all given keys.
	_, err = c.redis.Do(ctx, "DEL", keys...)
	return
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
// It uses `FLUSHDB` command in redis server, which might be disabled in server.
func (c *AdapterRedis) Clear(ctx context.Context) (err error) {
	// The "FLUSHDB" may not be available.
	_, err = c.redis.Do(ctx, "FLUSHDB")
	return
}

// Close closes the cache.
func (c *AdapterRedis) Close(ctx context.Context) error {
	// It does nothing.
	return nil
}

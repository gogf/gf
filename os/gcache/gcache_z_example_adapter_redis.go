package gcache

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"time"
)

// Redis is the gcache adapter implements using Redis server.
type Redis struct {
	redis *gredis.Redis
}

// NewRedis creates and returns a new redis memory cache object.
func NewRedis(redis *gredis.Redis) Adapter {
	return &Redis{
		redis: redis,
	}
}


// Set sets cache with `key`-`value` pair, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (r Redis) Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error {
	var err error
	if value == nil || duration < 0 {
		_, err = r.redis.Do(ctx,"DEL", key)
	} else {
		if duration == 0 {
			_, err = r.redis.Do(ctx,"SET", key, value)
		} else {
			_, err = r.redis.Do(ctx,"SETEX", key, uint64(duration.Seconds()), value)
		}
	}
	return err
}

// SetMap batch sets cache with key-value pairs by `data` map, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (r Redis) SetMap(ctx context.Context, data map[interface{}]interface{}, duration time.Duration) error {
	panic("implement me")
}


// SetIfNotExist sets cache with `key`-`value` pair which is expired after `duration`
// if `key` does not exist in the cache. It returns true the `key` does not exist in the
// cache, and it sets `value` successfully to the cache, or else it returns false.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (r Redis) SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (ok bool, err error) {
	panic("implement me")
}


// SetIfNotExistFunc sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// The parameter `value` can be type of `func() interface{}`, but it does nothing if its
// result is nil.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (r Redis) SetIfNotExistFunc(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (ok bool, err error) {
	panic("implement me")
}

// SetIfNotExistFuncLock sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
//
// Note that it differs from function `SetIfNotExistFunc` is that the function `f` is executed within
// writing mutex lock for concurrent safety purpose.
func (r Redis) SetIfNotExistFuncLock(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (ok bool, err error) {
	panic("implement me")
}

// Get retrieves and returns the associated value of given `key`.
// It returns nil if it does not exist, or its value is nil, or it's expired.
// If you would like to check if the `key` exists in the cache, it's better using function Contains.
func (r Redis) Get(ctx context.Context, key interface{}) (*gvar.Var, error) {
	v, err := r.redis.Do(ctx,"GET", key)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// GetOrSet retrieves and returns the value of `key`, or sets `key`-`value` pair and
// returns `value` if `key` does not exist in the cache. The key-value pair expires
// after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (r Redis) GetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (result *gvar.Var, err error) {
	panic("implement me")
}

// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (r Redis) GetOrSetFunc(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (result *gvar.Var, err error) {
	panic("implement me")
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
func (r Redis) GetOrSetFuncLock(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (result *gvar.Var, err error) {
	panic("implement me")
}

// Contains checks and returns true if `key` exists in the cache, or else returns false.
func (r Redis) Contains(ctx context.Context, key interface{}) (bool, error) {
	panic("implement me")
}

// Size returns the number of items in the cache.
func (r Redis) Size(ctx context.Context) (size int, err error) {
	panic("implement me")
}

// Data returns a copy of all key-value pairs in the cache as map type.
// Note that this function may lead lots of memory usage, you can implement this function
// if necessary.
func (r Redis) Data(ctx context.Context) (data map[interface{}]interface{}, err error) {
	panic("implement me")
}

// Keys returns all keys in the cache as slice.
func (r Redis) Keys(ctx context.Context) (keys []interface{}, err error) {
	panic("implement me")
}

// Values returns all values in the cache as slice.
func (r Redis) Values(ctx context.Context) (values []interface{}, err error) {
	panic("implement me")
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func (r Redis) Update(ctx context.Context, key interface{}, value interface{}) (oldValue *gvar.Var, exist bool, err error) {
	panic("implement me")
}


// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func  (r Redis) UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	return defaultCache.UpdateExpire(ctx, key, duration)
}

// GetExpire retrieves and returns the expiration of `key` in the cache.
//
// Note that,
// It returns 0 if the `key` does not expire.
// It returns -1 if the `key` does not exist in the cache.
func (r Redis) GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	panic("implement me")
}

// Remove deletes one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the last deleted item.
func (r Redis) Remove(ctx context.Context, keys ...interface{}) (lastValue *gvar.Var, err error) {
	panic("implement me")
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
func (r Redis) Clear(ctx context.Context) error {
	panic("implement me")
}

// Close closes the cache if necessary.
func (r Redis) Close(ctx context.Context) error {
	panic("implement me")
}




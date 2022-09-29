// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"time"
)

type RedisGroupExpire struct {
	redis *Redis
}

func (r *Redis) Expire() *RedisGroupExpire {
	return &RedisGroupExpire{
		redis: r,
	}
}

// Expire set a timeout on key.
// After the timeout has expired, the key will automatically be deleted.
//
// https://redis.io/commands/expire/
func (r *RedisGroupExpire) Expire(ctx context.Context, key string, seconds time.Duration) (bool, error) {
	v, err := r.redis.Do(ctx, "EXPIRE", key, seconds.Seconds())
	return v.Bool(), err
}

// ExpireAt has the same effect and semantic as EXPIRE, but instead of specifying the number of
// seconds representing the TTL (time to live), it takes an absolute Unix timestamp (seconds since
// January 1, 1970).
// A timestamp in the past will delete the key immediately.
//
// https://redis.io/commands/expireat/
func (r *RedisGroupExpire) ExpireAt(ctx context.Context, key string, time time.Time) (bool, error) {
	v, err := r.redis.Do(ctx, "EXPIREAT", key, time)
	return v.Bool(), err
}

// TTL return the remaining time to live of a key that has a timeout.
// This introspection capability allows a Redis client to check how many seconds a given key
// will continue to be part of the dataset.
// In Redis 2.6 or older the command returns -1 if the key does not exist or if the key exist but has
// no associated expire.
// Starting with Redis 2.8 the return value in case of error changed:
// The command returns -2 if the key does not exist.
// The command returns -1 if the key exists but has no associated expire.
// See also the PTTL command that returns the same information with milliseconds resolution (Only
// available in Redis 2.6 or greater).
//
// https://redis.io/commands/ttl/
func (r *RedisGroupExpire) TTL(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "TTL", key)
	return v.Int64(), err
}

// PErsist remove the existing timeout on key, turning the key from volatile (a key with an expire set)
// to persistent (a key that will never expire as no timeout is associated).
//
// https://redis.io/commands/persist/
func (r *RedisGroupExpire) PErsist(ctx context.Context, key string) (bool, error) {
	v, err := r.redis.Do(ctx, "PERSIST", key)
	return v.Bool(), err
}

// PExpire works exactly like EXPIRE but the time to live of the key is specified in milliseconds
// instead of seconds.
//
// https://redis.io/commands/pexpire/
func (r *RedisGroupExpire) PExpire(ctx context.Context, key string, time time.Duration, options string) (bool, error) {
	v, err := r.redis.Do(ctx, "PEXPIRE", key, time.Milliseconds(), options)
	return v.Bool(), err
}

// PExpireAt has the same effect and semantic as EXPIREAT, but the Unix time at which the key will
// expire is specified in milliseconds instead of seconds.
//
// https://redis.io/commands/pexpireat/
func (r *RedisGroupExpire) PExpireAt(ctx context.Context, key string, time time.Time) (bool, error) {
	v, err := r.redis.Do(ctx, "PEXPIREAT", key, time)
	return v.Bool(), err
}

// PTTL like TTL this command returns the remaining time to live of a key that has an expire set,
// with the sole difference that TTL returns the amount of remaining time in seconds while PTTL
// returns it in milliseconds.
//
// In Redis 2.6 or older the command returns -1 if the key does not exist or if the key exist but has
// no associated expire.
//
//  https://redis.io/commands/pttl/
func (r *RedisGroupExpire) PTTL(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "PTTL", key)
	return v.Int64(), err
}

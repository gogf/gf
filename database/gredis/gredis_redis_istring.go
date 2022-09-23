// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
)

type RedisString struct {
	redis *Redis
}

func (r *Redis) String() *RedisString {
	return &RedisString{
		redis: r,
	}
}

// Set key to hold the string value. If key already holds a value, it is overwritten, regardless of its type.
// Any previous time to live associated with the key is discarded on successful SET operation.
//
// https://redis.io/commands/set/
func (r *RedisString) Set(ctx context.Context, key string, value interface{}) (string, error) {
	v, err := r.redis.Do(ctx, "SET", key, value)
	return v.String(), err
}

func (r *RedisString) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
}

func (r *RedisString) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
}

// Get the value of key. If the key does not exist the special value nil is returned.
// An error is returned if the value stored at key is not a string, because GET only handles string values.
//
// https://redis.io/commands/get/
func (r *RedisString) Get(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "GET", key)
}

// GetDel gets the value of key and delete the key.
// This command is similar to GET, except for the fact that it also deletes the key on success
// (if and only if the key's value type is a string).
//
// https://redis.io/commands/getdel/
func (r *RedisString) GetDel(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "GETDEL", key)
}

func (r *RedisString) GetEX(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "GETDEL", key)
}

func (r *RedisString) GetSet(ctx context.Context, key string, value interface{}) (string, error) {}

func (r *RedisString) StrLen(ctx context.Context, key string) (int64, error) {}

// Append appends the value at the end of the string, if key already exists and is a string.
// If key does not exist it is created and set as an empty string,
// so APPEND will be similar to SET in this special case.
//
// https://redis.io/commands/append/
func (r *RedisString) Append(ctx context.Context, key string, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "APPEND", key, value)
	return v.Int64(), err
}

func (r *RedisString) SetRange(ctx context.Context, key string, offset int64, value string) (int64, error) {
}

func (r *RedisString) GetRange(ctx context.Context, key string, start, end int64) (string, error) {}

func (r *RedisString) Incr(ctx context.Context, key string) (int64, error) {}

func (r *RedisString) IncrBy(ctx context.Context, key string, value int64) (int64, error) {}

func (r *RedisString) IncrByFloat(ctx context.Context, key string, value float64) (float64, error) {}

// Decr decrements the number stored at key by one.
// If the key does not exist, it is set to 0 before performing the operation.
// An error is returned if the key contains a value of the wrong type or contains a string
// that can not be represented as integer. This operation is limited to 64 bits signed integers.
//
// https://redis.io/commands/decr/
func (r *RedisString) Decr(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "DECR", key)
	return v.Int64(), err
}

// DecrBy decrements the number stored at key by decrement.
// If the key does not exist, it is set to 0 before performing the operation.
// An error is returned if the key contains a value of the wrong type or contains a string
// that can not be represented as integer. This operation is limited to 64 bits signed integers.
//
// https://redis.io/commands/decrby/
func (r *RedisString) DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	v, err := r.redis.Do(ctx, "DECRBY", key, decrement)
	return v.Int64(), err
}

func (r *RedisString) MSet(ctx context.Context, pairs ...interface{}) (string, error) {}

func (r *RedisString) MSetNX(ctx context.Context, pairs ...interface{}) (bool, error) {}

func (r *RedisString) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {}

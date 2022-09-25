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

// RedisGroupList is the redis group list object.
type RedisGroupList struct {
	redis *Redis
}

// List is the redis list object.
func (r *Redis) List() *RedisGroupList {
	return &RedisGroupList{
		redis: r,
	}
}

// LPush inserts all the specified values at the head of the list stored at key.
func (r *RedisGroupList) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "LPUSH", key, values)
	return v.Int64(), err
}

// LPushX inserts value at the head of the list stored at key, only if key exists and holds a list.
func (r *RedisGroupList) LPushX(ctx context.Context, key string, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "LPushX", key, value)
	return v.Int64(), err
}

// RPush inserts all the specified values at the tail of the list stored at key.
func (r *RedisGroupList) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "RPush", key, values)
	return v.Int64(), err
}

// RPushX inserts value at the tail of the list stored at key, only if key exists and holds a list.
func (r *RedisGroupList) RPushX(ctx context.Context, key string, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "RPushX", key, value)
	return v.Int64(), err
}

// LPop removes and returns the first element of the list stored at key.
func (r *RedisGroupList) LPop(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "LPop", key)
}

// RPop removes and returns the last element of the list stored at key.
func (r *RedisGroupList) RPop(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "RPop", key)
}

// RPopLPush removes the last element in list source, appends it to the front of list destination and returns it.
func (RedisGroupList) RPopLPush(ctx context.Context, source, destination string) (string, error) {
	panic("implement me")
}

// LRem removes the first count occurrences of elements equal to value from the list stored at key.
func (RedisGroupList) LRem(ctx context.Context, key string, count int64, value string) (int64, error) {
	panic("implement me")
}

// LLen returns the length of the list stored at key.
func (r *RedisGroupList) LLen(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "LLen", key)
	return v.Int64(), err
}

// LIndex returns the element at index in the list stored at key.
func (RedisGroupList) LIndex(ctx context.Context, key string, index int64) (string, error) {
	panic("implement me")
}

// LInsert inserts value in the list stored at key either before or after the reference value pivot.
func (RedisGroupList) LInsert(ctx context.Context, key, op string, pivot, value string) (int64, error) {
	panic("implement me")
}

// LSet sets the list element at index to value.
func (RedisGroupList) LSet(ctx context.Context, key string, index int64, value string) (string, error) {
	panic("implement me")
}

// LRange returns the specified elements of the list stored at key.
func (RedisGroupList) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	panic("implement me")
}

// LTrim removes the first count occurrences of elements equal to value from the list stored at key.
func (RedisGroupList) LTrim(ctx context.Context, key string, start, stop int64) (string, error) {
	panic("implement me")
}

// BLPop removes and returns the first element of the list stored at key, or blocks until one is available.
func (RedisGroupList) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	panic("implement me")
}

// BRPop removes and returns the last element of the list stored at key, or blocks until one is available.
func (RedisGroupList) BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	panic("implement me")
}

// BRPopLPush removes the last element in list source, appends it to the front of list destination and returns it.
func (RedisGroupList) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error) {
	panic("implement me")
}

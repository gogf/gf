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

// LPush inserts all the specified values at the head of the list stored at key
// Insert all the specified values at the head of the list stored at key.
// If key does not exist, it is created as empty list before performing the push operations.
// When key holds a value that is not a list, an error is returned.
//
// It is possible to push multiple elements using a single command call just specifying multiple arguments at
// the end of the command. Elements are inserted one after the other to the head of the list,
// from the leftmost element to the rightmost element.
// So for instance the command `LPUSH mylist a b c` will result into a list containing c as first element,
// b as second element and a as third element
//
// https://redis.io/commands/lpush/
func (r *RedisGroupList) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "LPUSH", key, values)
	return v.Int64(), err
}

// LPushX inserts value at the head of the list stored at key, only if key exists and holds a list.
// Inserts specified values at the head of the list stored at key, only if key already exists and holds a list.
// In contrary to LPUSH, no operation will be performed when key does not yet exist.
// Return Integer reply: the length of the list after the push operation.
//
// https://redis.io/commands/lpushx
func (r *RedisGroupList) LPushX(ctx context.Context, key string, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "LPushX", key, value)
	return v.Int64(), err
}

// RPush inserts all the specified values at the tail of the list stored at key.
// Insert all the specified values at the tail of the list stored at key.
// If key does not exist, it is created as empty list before performing the push operation.
// When key holds a value that is not a list, an error is returned.
//
// It is possible to push multiple elements using a single command call just specifying multiple arguments at
// the end of the command. Elements are inserted one after the other to the tail of the list,
// from the leftmost element to the rightmost element.
// So for instance the command RPUSH mylist a b c will result into a list containing a as first element,
// b as second element and c as third element.
//
// https://redis.io/commands/rpush
func (r *RedisGroupList) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "RPush", key, values)
	return v.Int64(), err
}

// RPushX inserts value at the tail of the list stored at key, only if key exists and holds a list.
// Inserts specified values at the tail of the list stored at key, only if key already exists and holds a list.
// In contrary to RPUSH, no operation will be performed when key does not yet exist.
//
// Return Integer reply: the length of the list after the push operation.
//
// https://redis.io/commands/rpushx
func (r *RedisGroupList) RPushX(ctx context.Context, key string, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "RPushX", key, value)
	return v.Int64(), err
}

// LPop removes and returns the first element of the list stored at key.
// Removes and returns the first elements of the list stored at key.
//
// By default, the command pops a single element from the beginning of the list.
// When provided with the optional count argument, the reply will consist of up to count elements,
// depending on the list's length.
//
// Return
// When called without the count argument:
// Bulk string reply: the value of the first element, or nil when key does not exist.
//
// When called with the count argument:
// Array reply: list of popped elements, or nil when key does not exist.
//
// https://redis.io/commands/lpop
func (r *RedisGroupList) LPop(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "LPop", key)
}

// RPop removes and returns the last element of the list stored at key.
// Removes and returns the last elements of the list stored at key.
// By default, the command pops a single element from the end of the list.
// When provided with the optional count argument, the reply will consist of up to count elements,
// depending on the list's length.
//
// Return
// When called without the count argument:
// Bulk string reply: the value of the last element, or nil when key does not exist.
//
// When called with the count argument:
// Array reply: list of popped elements, or nil when key does not exist.
//
// https://redis.io/commands/rpop
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
// Returns the length of the list stored at key.
// If key does not exist, it is interpreted as an empty list and 0 is returned.
// An error is returned when the value stored at key is not a list.
//
// Return
// Integer reply: the length of the list at key.
//
// https://redis.io/commands/llen
func (r *RedisGroupList) LLen(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "LLen", key)
	return v.Int64(), err
}

// LIndex returns the element at index in the list stored at key.
// Returns the element at index index in the list stored at key.
// The index is zero-based, so 0 means the first element, 1 the second element and so on.
// Negative indices can be used to designate elements starting at the tail of the list. Here,
// -1 means the last element, -2 means the penultimate and so forth.
//
// When the value at key is not a list, an error is returned.
//
// Return
// Bulk string reply: the requested element, or nil when index is out of range.
//
// https://redis.io/commands/lindex
func (r *RedisGroupList) LIndex(ctx context.Context, key string, index int64) (*gvar.Var, error) {
	return r.redis.Do(ctx, "LINDEX", key, index)
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

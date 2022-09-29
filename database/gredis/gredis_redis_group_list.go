// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
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

// LPush insert all the specified values at the head of the list stored at key
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

// LPushX insert value at the head of the list stored at key, only if key exists and holds a list.
// Inserts specified values at the head of the list stored at key, only if key already exists and holds a list.
// In contrary to LPUSH, no operation will be performed when key does not yet exist.
// Return Integer reply: the length of the list after the push operation.
//
// https://redis.io/commands/lpushx
func (r *RedisGroupList) LPushX(ctx context.Context, key string, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "LPUSHX", key, value)
	return v.Int64(), err
}

// RPush insert all the specified values at the tail of the list stored at key.
// Insert all the specified values at the tail of the list stored at key.
// If key does not exist, it is created as empty list before performing the push operation.
// When key holds a value that is not a list, an error is returned.
// It is possible to push multiple elements using a single command call just specifying multiple
// arguments at  the end of the command. Elements are inserted one after the other to the tail of the
// list, from the leftmost element to the rightmost element.
// So for instance the command RPUSH mylist a b c will result into a list containing a as first element,
// b as second element and c as third element.
//
// https://redis.io/commands/rpush
func (r *RedisGroupList) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "RPUSH", key, values)
	return v.Int64(), err
}

// RPushX insert value at the tail of the list stored at key, only if key exists and holds a list.
// Inserts specified values at the tail of the list stored at key, only if key already exists and
// holds a list.
// In contrary to RPUSH, no operation will be performed when key does not yet exist.
//
// Return Integer reply: the length of the list after the push operation.
//
// https://redis.io/commands/rpushx
func (r *RedisGroupList) RPushX(ctx context.Context, key string, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "RPUSHX", key, value)
	return v.Int64(), err
}

// LPop remove and returns the first element of the list stored at key.
// Removes and returns the first elements of the list stored at key.
//
// By default, the command pops a single element from the beginning of the list.
// When provided with the optional count argument, the reply will consist of up to count elements,
// depending on the list's length.
//
// Return When called without the count argument:
// Bulk string reply: the value of the first element, or nil when key does not exist.
//
// When called with the count argument:
// Array reply: list of popped elements, or nil when key does not exist.
//
// https://redis.io/commands/lpop
func (r *RedisGroupList) LPop(ctx context.Context, key string, count int) ([]string, error) {
	v, err := r.redis.Do(ctx, "LPOP", key, count)
	return gconv.SliceStr(v), err
}

// RPop remove and returns the last element of the list stored at key.
// Removes and returns the last elements of the list stored at key.
// By default, the command pops a single element from the end of the list.
// When provided with the optional count argument, the reply will consist of up to count elements,
// depending on the list's length.
// Return When called without the count argument:
// Bulk string reply: the value of the last element, or nil when key does not exist.
// When called with the count argument:
// Array reply: list of popped elements, or nil when key does not exist.
//
// https://redis.io/commands/rpop
func (r *RedisGroupList) RPop(ctx context.Context, key string) (string, error) {
	v, err := r.redis.Do(ctx, "RPOP", key)
	return v.String(), err
}

// RPopLPush remove the last element in list source, appends it to the front of list destination and
// returns it.
//
// https://redis.io/commands/rpoplpush/
func (r *RedisGroupList) RPopLPush(ctx context.Context, source, destination string) ([]string, error) {
	v, err := r.redis.Do(ctx, "RPOPLPUSH", source, destination)
	return gconv.SliceStr(v), err
}

// LRem remove the first count occurrences of elements equal to value from the list stored at key.
//
// https://redis.io/commands/lrem/
func (r *RedisGroupList) LRem(ctx context.Context, key string, count int64, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "LREM", key, count, value)
	return v.Int64(), err
}

// LLen return the length of the list stored at key.
// Returns the length of the list stored at key.
// If key does not exist, it is interpreted as an empty list and 0 is returned.
// An error is returned when the value stored at key is not a list.
//
// Return
// Integer reply: the length of the list at key.
//
// https://redis.io/commands/llen
func (r *RedisGroupList) LLen(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "LLEN", key)
	return v.Int64(), err
}

// LIndex return the element at index in the list stored at key.
// Returns the element at index index in the list stored at key.
// The index is zero-based, so 0 means the first element, 1 the second element and so on.
// Negative indices can be used to designate elements starting at the tail of the list.
// Here, -1 means the last element, -2 means the penultimate and so forth.
// When the value at key is not a list, an error is returned.
// Return
// Bulk string reply: the requested element, or nil when index is out of range.
//
// https://redis.io/commands/lindex
func (r *RedisGroupList) LIndex(ctx context.Context, key string, index int64) (*gvar.Var, error) {
	return r.redis.Do(ctx, "LINDEX", key, index)
}

// LInsert insert element in the list stored at key either before or after the reference value pivot.
// When key does not exist, it is considered an empty list and no operation is performed.
// An error is returned when key exists but does not hold a list value.
//
// https://redis.io/commands/linsert/
func (r *RedisGroupList) LInsert(ctx context.Context, key, op string, pivot, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "LINSERT", key, op, pivot, value)
	return v.Int64(), err
}

// LSet set the list element at index to element.
// For more information on the index argument, see LINDEX.
// An error is returned for out of range indexes.
//
// https://redis.io/commands/lset/
func (r *RedisGroupList) LSet(ctx context.Context, key string, index int64, value string) (string, error) {
	v, err := r.redis.Do(ctx, "LSET", key, index, value)
	return v.String(), err
}

// LRange return the specified elements of the list stored at key.
// The offsets start and stop are zero-based indexes, with 0 being the first element of the list (the
// head of the list), 1 being the next element and so on.
//
// These offsets can also be negative numbers indicating offsets starting at the end of the list.
// For example, -1 is the last element of the list, -2 the penultimate, and so on.
//
// https://redis.io/commands/lrange/
func (r *RedisGroupList) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	v, err := r.redis.Do(ctx, "LRANGE", key, start, stop)
	return gconv.SliceStr(v), err
}

// LTrim trim an existing list so that it will contain only the specified range of elements
// specified. Both start and stop are zero-based indexes, where 0 is the first element of the list
// (the head), 1 the next element and so on.
// For example: LTRIM foobar 0 2 will modify the list stored at foobar so that only the first three
// elements of the list will remain.
// start and end can also be negative numbers indicating offsets from the end of the list, where -1
// is the last element of the list, -2 the penultimate element and so on.
// Out of range indexes will not produce an error: if start is larger than the end of the list, or
// start > end, the result will be an empty list (which causes key to be removed). If end is larger
// than the end of the list, Redis will treat it like the last element of the list.
//
//
// https://redis.io/commands/ltrim/
func (r *RedisGroupList) LTrim(ctx context.Context, key string, start, stop int64) (string, error) {
	v, err := r.redis.Do(ctx, "LTRIM", key, start, stop)
	return v.String(), err
}

// BLPop is a blocking list pop primitive.
// It is the blocking version of LPOP because it blocks the connection when there are no elements to
// pop from any of the given lists.
// An element is popped from the head of the first list that is non-empty, with the given keys being
// checked in the order that they are given.
//
// https://redis.io/commands/blpop/
func (r *RedisGroupList) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	v, err := r.redis.Do(ctx, "BLPOP", keys, timeout.Seconds())
	return gconv.SliceStr(v), err
}

// BRPop is a blocking list pop primitive.
// It is the blocking version of RPOP because it blocks the connection when there are no elements to
// pop from any of the given lists. An element is popped from the tail of the first list that is
// non-empty, with the given keys being checked in the order that they are given.
//
// https://redis.io/commands/brpop/
func (r *RedisGroupList) BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	v, err := r.redis.Do(ctx, "BRPOP", keys, timeout.Seconds())
	return gconv.SliceStr(v), err
}

// BRPopLPush is the blocking variant of RPOPLPUSH.
// When source contains elements, this command behaves exactly like RPOPLPUSH. When used inside a
// MULTI/EXEC block,
// this command behaves exactly like RPOPLPUSH. When source is empty, Redis will block the connection
// until another // client pushes to it or until timeout is reached.
// A timeout of zero can be used to block indefinitely.
//
// https://redis.io/commands/brpoplpush/
func (r *RedisGroupList) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error) {
	v, err := r.redis.Do(ctx, "BRPOPLPUSH", source, destination, timeout.Seconds())
	return v.String(), err
}

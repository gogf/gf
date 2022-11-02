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
	"github.com/gogf/gf/v2/util/gconv"
)

// RedisGroupString is the function group manager for string operations.
type RedisGroupString struct {
	redis *Redis
}

// GroupString is the redis group object for string operations.
func (r *Redis) GroupString() RedisGroupString {
	return RedisGroupString{
		redis: r,
	}
}

// TTLOption provides extra option for TTL related functions.
type TTLOption struct {
	EX      *int64 // EX seconds -- Set the specified expire time, in seconds.
	PX      *int64 // PX milliseconds -- Set the specified expire time, in milliseconds.
	EXAT    *int64 // EXAT timestamp-seconds -- Set the specified Unix time at which the key will expire, in seconds.
	PXAT    *int64 // PXAT timestamp-milliseconds -- Set the specified Unix time at which the key will expire, in milliseconds.
	KeepTTL bool   // Retain the time to live associated with the key.
}

// SetOption provides extra option for Set function.
type SetOption struct {
	TTLOption
	NX bool // Only set the key if it does not already exist.
	XX bool // Only set the key if it already exists.

	// Return the old string stored at key, or nil if key did not exist.
	// An error is returned and SET aborted if the value stored at key is not a string.
	Get bool
}

// Set key to hold the string value. If key already holds a value, it is overwritten,
// regardless of its type.
// Any previous time to live associated with the key is discarded on successful SET operation.
//
// https://redis.io/commands/set/
func (r RedisGroupString) Set(ctx context.Context, key string, value interface{}, option ...SetOption) (*gvar.Var, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	return r.redis.Do(ctx, "Set", mustMergeOptionToArgs(
		[]interface{}{key, value}, usedOption,
	)...)
}

// SetNX sets key to hold string value if key does not exist.
// In that case, it is equal to SET.
// When key already holds a value, no operation is performed.
// SetNX is short for "SET if Not exists".
//
// It returns:
// true:  if the all the keys were set.
// false: if no key was set (at least one key already existed).
//
// https://redis.io/commands/setnx/
func (r RedisGroupString) SetNX(ctx context.Context, key string, value interface{}) (bool, error) {
	v, err := r.redis.Do(ctx, "SetNX", key, value)
	return v.Bool(), err
}

// SetEX sets key to hold the string value and set key to timeout after a given number of seconds.
// This command is equivalent to executing the following commands:
//
//	SET myKey value
//	EXPIRE myKey seconds
//
// SetEX is atomic, and can be reproduced by using the previous two commands inside an MULTI / EXEC block.
// It is provided as a faster alternative to the given sequence of operations, because this operation is very
// common when Redis is used as a cache.
//
// An error is returned when seconds invalid.
//
// https://redis.io/commands/setex/
func (r RedisGroupString) SetEX(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	_, err := r.redis.Do(ctx, "SetEX", key, int64(ttl.Seconds()), value)
	return err
}

// Get the value of key. If the key does not exist the special value nil is returned.
// An error is returned if the value stored at key is not a string, because GET only handles string values.
//
// https://redis.io/commands/get/
func (r RedisGroupString) Get(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "Get", key)
}

// GetDel gets the value of key and delete the key.
// This command is similar to GET, except for the fact that it also deletes the key on success
// (if and only if the key's value type is a string).
//
// https://redis.io/commands/getdel/
func (r RedisGroupString) GetDel(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "GetDel", key)
}

// GetEXOption provides extra option for GetEx function.
type GetEXOption struct {
	TTLOption
	Persist bool // Persist -- Remove the time to live associated with the key.
}

// GetEX is similar to GET, but is a write command with additional options.
//
// https://redis.io/commands/getex/
func (r RedisGroupString) GetEX(ctx context.Context, key string, option ...GetEXOption) (*gvar.Var, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	return r.redis.Do(ctx, "GetEX", mustMergeOptionToArgs(
		[]interface{}{key}, usedOption,
	)...)
}

// GetSet atomically sets key to value and returns the old value stored at key.
// Returns an error when key exists but does not hold a string value. Any previous time to live associated with
// the key is discarded on successful SET operation.
//
// https://redis.io/commands/getset/
func (r RedisGroupString) GetSet(ctx context.Context, key string, value interface{}) (*gvar.Var, error) {
	return r.redis.Do(ctx, "GetSet", key, value)
}

// StrLen returns the length of the string value stored at key.
// An error is returned when key holds a non-string value.
//
// It returns the length of the string at key, or 0 when key does not exist.
//
// https://redis.io/commands/strlen/
func (r RedisGroupString) StrLen(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "StrLen", key)
	return v.Int64(), err
}

// Append appends the value at the end of the string, if key already exists and is a string.
// If key does not exist it is created and set as an empty string,
// so APPEND will be similar to SET in this special case.
//
// https://redis.io/commands/append/
func (r RedisGroupString) Append(ctx context.Context, key string, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "Append", key, value)
	return v.Int64(), err
}

// SetRange overwrites part of the string stored at key, starting at the specified offset, for the entire length
// of value. If the offset is larger than the current length of the string at key, the string is padded with
// zero-bytes to make offset fit. Non-existing keys are considered as empty strings, so this command will
// make sure it holds a string large enough to be able to set value at offset.
//
// It returns the length of the string after it was modified by the command.
//
// https://redis.io/commands/setrange/
func (r RedisGroupString) SetRange(ctx context.Context, key string, offset int64, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "SetRange", key, offset, value)
	return v.Int64(), err
}

// GetRange returns the substring of the string value stored at key,
// determined by the offsets start and end (both are inclusive). Negative offsets can be used in order to provide
// an offset starting from the end of the string. So -1 means the last character, -2 the penultimate and so forth.
//
// The function handles out of range requests by limiting the resulting range to the actual length of the string.
//
// https://redis.io/commands/getrange/
func (r RedisGroupString) GetRange(ctx context.Context, key string, start, end int64) (string, error) {
	v, err := r.redis.Do(ctx, "GetRange", key, start, end)
	return v.String(), err
}

// Incr increments the number stored at key by one.
// If the key does not exist, it is set to 0 before performing the operation.
// An error is returned if the key contains a value of the wrong type or contains a string that can not be
// represented as integer. This operation is limited to 64 bits signed integers.
//
// https://redis.io/commands/incr/
func (r RedisGroupString) Incr(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "Incr", key)
	return v.Int64(), err
}

// IncrBy increments the number stored at key by increment. If the key does not exist, it is set to 0 before
// performing the operation.
//
// An error is returned if the key contains a value of the wrong type or contains a
// string that can not be represented as integer. This operation is limited to 64 bits signed integers.
//
// https://redis.io/commands/incrby/
func (r RedisGroupString) IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	v, err := r.redis.Do(ctx, "IncrBy", key, increment)
	return v.Int64(), err
}

// IncrByFloat increments the string representing a floating point number stored at key by the specified increment.
//
// https://redis.io/commands/incrbyfloat/
func (r RedisGroupString) IncrByFloat(ctx context.Context, key string, increment float64) (float64, error) {
	v, err := r.redis.Do(ctx, "IncrByFloat", key, increment)
	return v.Float64(), err
}

// Decr decrements the number stored at key by one.
//
// https://redis.io/commands/decr/
func (r RedisGroupString) Decr(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "Decr", key)
	return v.Int64(), err
}

// DecrBy decrements the number stored at key by decrement.
//
// https://redis.io/commands/decrby/
func (r RedisGroupString) DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	v, err := r.redis.Do(ctx, "DecrBy", key, decrement)
	return v.Int64(), err
}

// MSet sets the given keys to their respective values.
// MSet replaces existing values with new values, just as regular SET.
// See MSetNX if you don't want to overwrite existing values.
//
// MSet is atomic, so all given keys are set at once. It is not possible for clients to see that some keys
// were updated while others are unchanged.
//
// https://redis.io/commands/mset/
func (r RedisGroupString) MSet(ctx context.Context, keyValueMap map[string]interface{}) error {
	var args []interface{}
	for k, v := range keyValueMap {
		args = append(args, k, v)
	}
	_, err := r.redis.Do(ctx, "MSet", args...)
	return err
}

// MSetNX sets the given keys to their respective values.
//
// It returns:
// true:  if the all the keys were set.
// false: if no key was set (at least one key already existed).
func (r RedisGroupString) MSetNX(ctx context.Context, keyValueMap map[string]interface{}) (bool, error) {
	var args []interface{}
	for k, v := range keyValueMap {
		args = append(args, k, v)
	}
	v, err := r.redis.Do(ctx, "MSetNX", args...)
	return v.Bool(), err
}

// MGet returns the values of all specified keys.
//
// https://redis.io/commands/mget/
func (r RedisGroupString) MGet(ctx context.Context, keys ...string) (map[string]*gvar.Var, error) {
	var result = make(map[string]*gvar.Var)
	v, err := r.redis.Do(ctx, "MGet", gconv.Interfaces(keys)...)
	if err == nil {
		values := v.Vars()
		for i, key := range keys {
			result[key] = values[i]
		}
	}
	return result, err
}

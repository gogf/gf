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

// String returns the function group manager for string operations.
func (r *Redis) String() *RedisGroupString {
	return &RedisGroupString{
		redis: r,
	}
}

// SetOption provides extra option for Set function.
type SetOption struct {
	TTLOption
	NX      bool // Only set the key if it does not already exist.
	XX      bool // Only set the key if it already exists.
	KEEPTTL bool // Retain the time to live associated with the key.

	// Return the old string stored at key, or nil if key did not exist.
	// An error is returned and SET aborted if the value stored at key is not a string.
	GET bool
}

// Set key to hold the string value. If key already holds a value, it is overwritten, regardless of its type.
// Any previous time to live associated with the key is discarded on successful SET operation.
//
// https://redis.io/commands/set/
func (r *RedisGroupString) Set(ctx context.Context, key string, value interface{}, option ...SetOption) (*gvar.Var, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	return r.redis.Do(ctx, "SET", mustMergeOptionToArgs(
		[]interface{}{key, value}, usedOption,
	)...)
}

// SetNX sets key to hold string value if key does not exist. In that case, it is equal to SET.
// When key already holds a value, no operation is performed. SetNX is short for "SET if Not exists".
//
// It returns:
// true:  if the all the keys were set.
// false: if no key was set (at least one key already existed).
//
// https://redis.io/commands/setnx/
func (r *RedisGroupString) SetNX(ctx context.Context, key string, value interface{}) (bool, error) {
	v, err := r.redis.Do(ctx, "SETNX", key, value)
	return v.Bool(), err
}

// SetEX sets key to hold the string value and set key to timeout after a given number of seconds.
// This command is equivalent to executing the following commands:
//
//     SET mykey value
//     EXPIRE mykey seconds
//
// SetEX is atomic, and can be reproduced by using the previous two commands inside an MULTI / EXEC block.
// It is provided as a faster alternative to the given sequence of operations, because this operation is very
// common when Redis is used as a cache.
//
// An error is returned when seconds invalid.
//
// https://redis.io/commands/setex/
func (r *RedisGroupString) SetEX(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	_, err := r.redis.Do(ctx, "SETEX", key, int64(ttl.Seconds()), value)
	return err
}

// Get the value of key. If the key does not exist the special value nil is returned.
// An error is returned if the value stored at key is not a string, because GET only handles string values.
//
// https://redis.io/commands/get/
func (r *RedisGroupString) Get(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "GET", key)
}

// GetDel gets the value of key and delete the key.
// This command is similar to GET, except for the fact that it also deletes the key on success
// (if and only if the key's value type is a string).
//
// https://redis.io/commands/getdel/
func (r *RedisGroupString) GetDel(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "GETDEL", key)
}

// GetEXOption provides extra option for GetEx function.
type GetEXOption struct {
	TTLOption
	PERSIST bool // PERSIST -- Remove the time to live associated with the key.
}

// GetEX is similar to GET, but is a write command with additional options.
//
// https://redis.io/commands/getex/
func (r *RedisGroupString) GetEX(ctx context.Context, key string, option ...GetEXOption) (*gvar.Var, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	return r.redis.Do(ctx, "GETDEL", mustMergeOptionToArgs(
		[]interface{}{key}, usedOption,
	)...)
}

// GetSet atomically sets key to value and returns the old value stored at key.
// Returns an error when key exists but does not hold a string value. Any previous time to live associated with
// the key is discarded on successful SET operation.
//
// https://redis.io/commands/getset/
func (r *RedisGroupString) GetSet(ctx context.Context, key string, value interface{}) (*gvar.Var, error) {
	return r.redis.Do(ctx, "GETSET", key, value)
}

// StrLen returns the length of the string value stored at key.
// An error is returned when key holds a non-string value.
//
// It returns the length of the string at key, or 0 when key does not exist.
//
// https://redis.io/commands/strlen/
func (r *RedisGroupString) StrLen(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "STRLEN", key)
	return v.Int64(), err
}

// Append appends the value at the end of the string, if key already exists and is a string.
// If key does not exist it is created and set as an empty string,
// so APPEND will be similar to SET in this special case.
//
// https://redis.io/commands/append/
func (r *RedisGroupString) Append(ctx context.Context, key string, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "APPEND", key, value)
	return v.Int64(), err
}

// SetRange overwrites part of the string stored at key, starting at the specified offset, for the entire length
// of value. If the offset is larger than the current length of the string at key, the string is padded with
// zero-bytes to make offset fit. Non-existing keys are considered as empty strings, so this command will
// make sure it holds a string large enough to be able to set value at offset.
//
// Note that the maximum offset that you can set is 2^29 -1 (536870911), as Redis Strings are limited to 512 megabytes.
// If you need to grow beyond this size, you can use multiple keys.
//
// It returns the length of the string after it was modified by the command.
//
// https://redis.io/commands/setrange/
func (r *RedisGroupString) SetRange(ctx context.Context, key string, offset int64, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "SETRANGE", key, offset, value)
	return v.Int64(), err
}

// GetRange returns the substring of the string value stored at key,
// determined by the offsets start and end (both are inclusive). Negative offsets can be used in order to provide
// an offset starting from the end of the string. So -1 means the last character, -2 the penultimate and so forth.
//
// The function handles out of range requests by limiting the resulting range to the actual length of the string.
//
// https://redis.io/commands/getrange/
func (r *RedisGroupString) GetRange(ctx context.Context, key string, start, end int64) (string, error) {
	v, err := r.redis.Do(ctx, "GETRANGE", key, start, end)
	return v.String(), err
}

// Incr increments the number stored at key by one.
// If the key does not exist, it is set to 0 before performing the operation.
// An error is returned if the key contains a value of the wrong type or contains a string that can not be
// represented as integer. This operation is limited to 64 bits signed integers.
//
// Note: this is a string operation because Redis does not have a dedicated integer type.
// The string stored at the key is interpreted as a base-10 64 bits signed integer to execute the operation.
//
// Redis stores integers in their integer representation, so for string values that actually hold an integer,
// there is no overhead for storing the string representation of the integer.
//
// https://redis.io/commands/incr/
func (r *RedisGroupString) Incr(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "INCR", key)
	return v.Int64(), err
}

// IncrBy increments the number stored at key by increment. If the key does not exist, it is set to 0 before
// performing the operation. An error is returned if the key contains a value of the wrong type or contains a
// string that can not be represented as integer. This operation is limited to 64 bits signed integers.
//
// See Incr for extra information on increment/decrement operations.
//
// https://redis.io/commands/incrby/
func (r *RedisGroupString) IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	v, err := r.redis.Do(ctx, "INCRBY", key, increment)
	return v.Int64(), err
}

// IncrByFloat increments the string representing a floating point number stored at key by the specified increment.
// By using a negative increment value, the result is that the value stored at the key is decremented
// (by the obvious properties of addition). If the key does not exist, it is set to 0 before performing the operation.
// An error is returned if one of the following conditions occur:
//
// The key contains a value of the wrong type (not a string).
// The current key content or the specified increment are not parsable as a double precision floating point number.
// If the command is successful the new incremented value is stored as the new value of the key (replacing the old one),
// and returned to the caller as a string.
//
// Both the value already contained in the string key and the increment argument can be optionally provided in
// exponential notation, however the value computed after the increment is stored consistently in the same format,
// that is, an integer number followed (if needed) by a dot, and a variable number of digits representing the decimal
// part of the number. Trailing zeroes are always removed.
//
// The precision of the output is fixed at 17 digits after the decimal point regardless of the actual internal
// precision of the computation.
//
// https://redis.io/commands/incrbyfloat/
func (r *RedisGroupString) IncrByFloat(ctx context.Context, key string, increment float64) (float64, error) {
	v, err := r.redis.Do(ctx, "INCRBYFLOAT", key, increment)
	return v.Float64(), err
}

// Decr decrements the number stored at key by one.
// If the key does not exist, it is set to 0 before performing the operation.
// An error is returned if the key contains a value of the wrong type or contains a string
// that can not be represented as integer. This operation is limited to 64 bits signed integers.
//
// https://redis.io/commands/decr/
func (r *RedisGroupString) Decr(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "DECR", key)
	return v.Int64(), err
}

// DecrBy decrements the number stored at key by decrement.
// If the key does not exist, it is set to 0 before performing the operation.
// An error is returned if the key contains a value of the wrong type or contains a string
// that can not be represented as integer. This operation is limited to 64 bits signed integers.
//
// https://redis.io/commands/decrby/
func (r *RedisGroupString) DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	v, err := r.redis.Do(ctx, "DECRBY", key, decrement)
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
func (r *RedisGroupString) MSet(ctx context.Context, keyValues ...interface{}) error {
	_, err := r.redis.Do(ctx, "MSET", keyValues...)
	return err
}

// MSetNX sets the given keys to their respective values.
// MSetNX will not perform any operation at all even if just a single key already exists.
//
// Because of this semantic MSetNX can be used in order to set different keys representing different fields of
// a unique logic object in a way that ensures that either all the fields or none at all are set.
//
// MSetNX is atomic, so all given keys are set at once. It is not possible for clients to see that some keys
// were updated while others are unchanged.
//
// It returns:
// true:  if the all the keys were set.
// false: if no key was set (at least one key already existed).
func (r *RedisGroupString) MSetNX(ctx context.Context, keyValues ...interface{}) (bool, error) {
	v, err := r.redis.Do(ctx, "MSETNX", keyValues...)
	return v.Bool(), err
}

// MGet returns the values of all specified keys. For every key that does not hold a string value or does not exist,
// the special value nil is returned. Because of this, the operation never fails.
//
// https://redis.io/commands/mget/
func (r *RedisGroupString) MGet(ctx context.Context, keys ...string) ([]*gvar.Var, error) {
	v, err := r.redis.Do(ctx, "MGET", gconv.Interfaces(keys)...)
	return v.Vars(), err
}

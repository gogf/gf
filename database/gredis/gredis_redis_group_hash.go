// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/util/gconv"
)

type RedisGroupHash struct {
	redis *Redis
}

func (r *Redis) Hash() *RedisGroupHash {
	return &RedisGroupHash{
		redis: r,
	}
}

// HSet sets field in the hash stored at key to value.
// If key does not exist, a new key holding a hash is created.
// If field already exists in the hash, it is overwritten.
//
// It returns the number of fields that were added.
//
// https://redis.io/commands/hset/
func (r *RedisGroupHash) HSet(ctx context.Context, key string, fields map[string]interface{}) (int64, error) {
	var s = []interface{}{key}
	for k, v := range fields {
		s = append(s, k, v)
	}
	v, err := r.redis.Do(ctx, "HSET", s...)
	return v.Int64(), err
}

// HSetNX sets field in the hash stored at key to value, only if field does not yet exist.
// If key does not exist, a new key holding a hash is created.
// If field already exists, this operation has no effect.
//
// It returns:
// - 1 if field is a new field in the hash and value was set.
// - 0 if field already exists in the hash and no operation was performed.
//
// https://redis.io/commands/hsetnx/
func (r *RedisGroupHash) HSetNX(ctx context.Context, key, field string, value interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "HSETNX", key, field, value)
	return v.Int64(), err
}

// HGet returns the value associated with field in the hash stored at key.
//
// It returns the value associated with field, or nil when field is not present in the hash or key does not exist.
//
// https://redis.io/commands/hget/
func (r *RedisGroupHash) HGet(ctx context.Context, key, field string) (string, error) {
	v, err := r.redis.Do(ctx, "HGET", key, field)
	return v.String(), err
}

// HStrLen Returns the string length of the value associated with field in the hash stored at key.
// If the key or the field do not exist, 0 is returned.
//
// It returns the string length of the value associated with field,
// or zero when field is not present in the hash or key does not exist at all.
//
// https://redis.io/commands/hstrlen/
func (r *RedisGroupHash) HStrLen(ctx context.Context, key, field string) (int64, error) {
	v, err := r.redis.Do(ctx, "HSTRLEN", key, field)
	return v.Int64(), err
}

// HExists returns if field is an existing field in the hash stored at key.
//
// It returns:
// - 1 if the hash contains field.
// - 0 if the hash does not contain field, or key does not exist.
//
// https://redis.io/commands/hexists/
func (r *RedisGroupHash) HExists(ctx context.Context, key, field string) (int64, error) {
	v, err := r.redis.Do(ctx, "HEXISTS", key, field)
	return v.Int64(), err
}

// HDel removes the specified fields from the hash stored at key.
// Specified fields that do not exist within this hash are ignored.
// If key does not exist, it is treated as an empty hash and this command returns 0.
//
// It returns the number of fields that were removed from the hash, not including specified but non-existing fields.
//
// https://redis.io/commands/hdel/
func (r *RedisGroupHash) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	v, err := r.redis.Do(ctx, "HDEL", append([]interface{}{key}, gconv.Interfaces(fields)...)...)
	return v.Int64(), err
}

// HLen returns the number of fields contained in the hash stored at key.
//
// https://redis.io/commands/hlen/
func (r *RedisGroupHash) HLen(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "HLEN", key)
	return v.Int64(), err
}

// HIncrBy increments the number stored at field in the hash stored at key by increment.
// If key does not exist, a new key holding a hash is created.
// If field does not exist the value is set to 0 before the operation is performed.
//
// The range of values supported by HINCRBY is limited to 64-bit signed integers.
//
// It returns the value at field after the increment operation.
//
// https://redis.io/commands/hincrby/
func (r *RedisGroupHash) HIncrBy(ctx context.Context, key, field string, increment int64) (int64, error) {
	v, err := r.redis.Do(ctx, "HINCRBY", key, field, increment)
	return v.Int64(), err
}

// HIncrByFloat increments the specified field of a hash stored at key, and representing a floating
// point number, by the specified increment. If the increment value is negative, the result is to
// have the hash field value decremented instead of incremented. If the field does not exist, it is
// set to 0 before performing the operation.
// An error is returned if one of the following conditions occur:
//
// The field contains a value of the wrong type (not a string).
// The current field content or the specified increment are not parsable as a double precision
// floating point number.
// The exact behavior of this command is identical to the one of the INCRBYFLOAT command,
// please refer to the documentation of INCRBYFLOAT for further information.
//
// It returns the value of field after the increment.
//
// https://redis.io/commands/hincrbyfloat/
func (r *RedisGroupHash) HIncrByFloat(ctx context.Context, key, field string, increment float64) (float64, error) {
	v, err := r.redis.Do(ctx, "HINCRBYFLOAT", key, field, increment)
	return v.Float64(), err
}

// HMSet sets the specified fields to their respective values in the hash stored at key.
// This command overwrites any specified fields already existing in the hash.
// If key does not exist, a new key holding a hash is created.
//
// https://redis.io/commands/hmset/
func (r *RedisGroupHash) HMSet(ctx context.Context, key string, fields map[string]interface{}) error {
	var s = []interface{}{key}
	for k, v := range fields {
		s = append(s, k, v)
	}
	_, err := r.redis.Do(ctx, "HMSET", s...)
	return err
}

// HMGet return  the values associated with the specified fields in the hash stored at key.
// For every field that does not exist in the hash, a nil value is returned.
// Because non-existing keys are treated as empty hashes, running HMGET against a non-existing key
// will return a list of nil values.
//
// https://redis.io/commands/hmget/
func (r *RedisGroupHash) HMGet(ctx context.Context, key string, fields ...string) ([]*gvar.Var, error) {
	v, err := r.redis.Do(ctx, "HMGET", append([]interface{}{key}, gconv.Interfaces(fields)...)...)
	return v.Vars(), err
}

// HKeys returns all field names in the hash stored at key.
//
// https://redis.io/commands/hkeys/
func (r *RedisGroupHash) HKeys(ctx context.Context, key string) ([]string, error) {
	v, err := r.redis.Do(ctx, "HKEYS", key)
	return v.Strings(), err
}

// HVals return all values in the hash stored at key.
//
// https://redis.io/commands/hvals/
func (r *RedisGroupHash) HVals(ctx context.Context, key string) ([]*gvar.Var, error) {
	v, err := r.redis.Do(ctx, "HVALS", key)
	return v.Vars(), err
}

// HGetAll returns all fields and values of the hash stored at key.
// In the returned value, every field name is followed by its value,
// so the length of the reply is twice the size of the hash.
//
// https://redis.io/commands/hgetall/
func (r *RedisGroupHash) HGetAll(ctx context.Context, key string) (map[string]*gvar.Var, error) {
	v, err := r.redis.Do(ctx, "HGETALL", key)
	return v.MapStrVar(), err
}

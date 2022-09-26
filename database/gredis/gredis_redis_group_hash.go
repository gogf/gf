// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
)

type RedisGroupHash struct {
	redis *Redis
}

func (r *Redis) Hash() *RedisGroupHash {
	return &RedisGroupHash{
		redis: r,
	}
}

// HSet  Sets field in the hash stored at key to value. If key does not exist, a new key holding a hash is created. If field already exists in the hash, it is overwritten.
//
// https://redis.io/commands/hset/
func (r *RedisGroupHash) HSet(ctx context.Context, key, field, value string) (int64, error) {
	v, err := r.redis.Do(ctx, "HSET", key, field, value)
	return v.Int64(), err
}

// HSetNX Sets field in the hash stored at key to value, only if field does not yet exist.
// If key does not exist, a new key holding a hash is created. If field already exists, this operation has no effect.
//
// https://redis.io/commands/hsetnx/
func (r *RedisGroupHash) HSetNX(ctx context.Context, key, field, value string) (bool, error) {
	v, err := r.redis.Do(ctx, "HSETNX", key, field, value)
	return v.Bool(), err
}

// HGet Returns the value associated with field in the hash stored at key.
//
// https://redis.io/commands/hget/
func (r *RedisGroupHash) HGet(ctx context.Context, key, field string) (string, error) {
	v, err := r.redis.Do(ctx, "HGET", key, field)
	return v.String(), err
}

// HExists Returns if field is an existing field in the hash stored at key.
//
// https://redis.io/commands/hexists/
func (r *RedisGroupHash) HExists(ctx context.Context, key, field string) (bool, error) {
	v, err := r.redis.Do(ctx, "HEXISTS", key, field)
	return v.Bool(), err
}

// HDel Removes the specified fields from the hash stored at key. Specified fields that do not exist within this hash are ignored. If key does not exist, it is treated as an empty hash and this command returns 0.
//
// https://redis.io/commands/hdel/
func (r *RedisGroupHash) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	v, err := r.redis.Do(ctx, "HDEL", key, fields)
	return v.Int64(), err
}

// HLen Returns the number of fields contained in the hash stored at key.
//
// https://redis.io/commands/hlen/
func (r *RedisGroupHash) HLen(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "HLEN", key)
	return v.Int64(), err
}

// HIncrBy Increments the number stored at field in the hash stored at key by increment. If key does not exist, a new key holding a hash is created. If field does not exist the value is set to 0 before the operation is performed.
//
//The range of values supported by HINCRBY is limited to 64 bit signed integers.
//
// https://redis.io/commands/hincrby/
func (r *RedisGroupHash) HIncrBy(ctx context.Context, key, field string, value int64) (int64, error) {
	v, err := r.redis.Do(ctx, "HINCRBY", key, field, value)
	return v.Int64(), err
}

// HIncrByFloat Increment the specified field of a hash stored at key, and representing a floating point number, by the specified increment. If the increment value is negative, the result is to have the hash field value decremented instead of incremented. If the field does not exist, it is set to 0 before performing the operation. An error is returned if one of the following conditions occur:
//
//The field contains a value of the wrong type (not a string).
//The current field content or the specified increment are not parsable as a double precision floating point number.
//The exact behavior of this command is identical to the one of the INCRBYFLOAT command, please refer to the documentation of INCRBYFLOAT for further information.
//
// https://redis.io/commands/hincrbyfloat/
func (r *RedisGroupHash) HIncrByFloat(ctx context.Context, key, field string, value float64) (float64, error) {
	v, err := r.redis.Do(ctx, "HINCRBYFLOAT", key, field, value)
	return v.Float64(), err
}

// HMSet Sets the specified fields to their respective values in the hash stored at key. This command overwrites any specified fields already existing in the hash. If key does not exist, a new key holding a hash is created.
//
// https://redis.io/commands/hmset/
func (r *RedisGroupHash) HMSet(ctx context.Context, key string, fields ...map[string]interface{}) (bool, error) {
	v, err := r.redis.Do(ctx, "HMSET", key, fields)
	return v.Bool(), err
}

// HMGet Returns the values associated with the specified fields in the hash stored at key.
//
//For every field that does not exist in the hash, a nil value is returned. Because non-existing keys are treated as empty hashes, running HMGET against a non-existing key will return a list of nil values.
//
// https://redis.io/commands/hmget/
func (r *RedisGroupHash) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	v, err := r.redis.Do(ctx, "HMGET", key, fields)
	return v.Slice(), err
}

// HKeys Returns all field names in the hash stored at key.
//
// https://redis.io/commands/hkeys/
func (r *RedisGroupHash) HKeys(ctx context.Context, key string) ([]string, error) {
	v, err := r.redis.Do(ctx, "HKEYS", key)
	return v.Strings(), err
}

// HVals Returns all values in the hash stored at key.
//
// https://redis.io/commands/hvals/
func (r *RedisGroupHash) HVals(ctx context.Context, key string) ([]string, error) {
	v, err := r.redis.Do(ctx, "HVALS", key)
	return v.Strings(), err
}

// HGetAll  Returns all fields and values of the hash stored at key. In the returned value, every field name is followed by its value, so the length of the reply is twice the size of the hash.
//
// https://redis.io/commands/hgetall/
func (r *RedisGroupHash) HGetAll(ctx context.Context, key string) (map[string]interface{}, error) {
	v, err := r.redis.Do(ctx, "HGETALL", key)
	return v.MapStrAny(), err
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
)

type RedisGroupSet struct {
	redis *Redis
}

func (r *Redis) Set() *RedisGroupSet {
	return &RedisGroupSet{
		redis: r,
	}
}

// SAdd Add the specified members to the set stored at key.
// Specified members that are already a member of this set are ignored. If key does not exist, a new set is created before adding the specified members.
//
// An error is returned when the value stored at key is not a set.
//
// https://redis.io/commands/sadd/
func (r *RedisGroupSet) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "SADD", key, members)
	return v.Int64(), err
}

// SIsMember Returns if member is a member of the set stored at key.
//
// https://redis.io/commands/sismember/
func (r *RedisGroupSet) SIsMember(ctx context.Context, key string, member string) (bool, error) {
	v, err := r.redis.Do(ctx, "SISMEMBER", key, member)
	return v.Bool(), err
}

// SPop Removes and returns one or more random members from the set value store at key.
//
// This operation is similar to SRANDMEMBER, that returns one or more random elements from a set but does not remove it.
// By default, the command pops a single member from the set. When provided with the optional count argument, the reply will consist of up to count members, depending on the set's cardinality.
//
// https://redis.io/commands/spop/
func (r *RedisGroupSet) SPop(ctx context.Context, key string) (string, error) {
	v, err := r.redis.Do(ctx, "SPOP", key)
	return v.String(), err
}

// SRandMember When called with just the key argument, return a random element from the set value stored at key.
// If the provided count argument is positive, return an array of distinct elements.
// The array's length is either count or the set's cardinality (SCARD), whichever is lower.
// If called with a negative count, the behavior changes and the command is allowed to return the same element multiple times. In this case, the number of returned elements is the absolute value of the specified count.
//
// https://redis.io/commands/srandmember/
func (r *RedisGroupSet) SRandMember(ctx context.Context, key string, count int) (string, error) {
	v, err := r.redis.Do(ctx, "SRANDMEMBER", key, count)
	return v.String(), err
}

// SRem Remove the specified members from the set stored at key.
// Specified members that are not a member of this set are ignored.
// If key does not exist, it is treated as an empty set and this command returns 0.
//
// An error is returned when the value stored at key is not a set.
//
// https://redis.io/commands/srem/
func (r *RedisGroupSet) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "SREM", key, members)
	return v.Int64(), err
}

// SMove Move member from the set at source to the set at destination.
// This operation is atomic. In every given moment the element will appear to be a member of source or destination for other clients.
// If the source set does not exist or does not contain the specified element, no operation is performed and 0 is returned. Otherwise, the element is removed from the source set and added to the destination set.
// When the specified element already exists in the destination set, it is only removed from the source set.
//
// An error is returned if source or destination does not hold a set value.
//
// https://redis.io/commands/smove/
func (r *RedisGroupSet) SMove(ctx context.Context, source, destination, member string) (bool, error) {
	v, err := r.redis.Do(ctx, "SMOVE", source, destination, member)
	return v.Bool(), err
}

// SCard Returns the set cardinality (number of elements) of the set stored at key.
//
// https://redis.io/commands/scard/
func (r *RedisGroupSet) SCard(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "SCARD", key)
	return v.Int64(), err
}

// SMembers Returns all the members of the set value stored at key.
// This has the same effect as running SINTER with one argument key.
//
// https://redis.io/commands/smembers/
func (r *RedisGroupSet) SMembers(ctx context.Context, key string) ([]string, error) {
	v, err := r.redis.Do(ctx, "SMEMBERS", key)
	return v.Strings(), err
}

// SInter Returns the members of the set resulting from the intersection of all the given sets.
//
// https://redis.io/commands/sinter/
func (r *RedisGroupSet) SInter(ctx context.Context, keys ...string) ([]string, error) {
	v, err := r.redis.Do(ctx, "SINTER", keys)
	return v.Strings(), err
}

// SInterStore This command is equal to SINTER, but instead of returning the resulting set, it is stored in destination.
//
// If destination already exists, it is overwritten.
//
// https://redis.io/commands/sinterstore/
func (r *RedisGroupSet) SInterStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	v, err := r.redis.Do(ctx, "SINTERSTORE", destination, keys)
	return v.Int64(), err
}

// SUnion Returns the members of the set resulting from the union of all the given sets.
//
// https://redis.io/commands/sunion/
func (r *RedisGroupSet) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	v, err := r.redis.Do(ctx, "SUNION", keys)
	return v.Strings(), err
}

// SUnionStore This command is equal to SUNION, but instead of returning the resulting set, it is stored in destination.
//
//  If destination already exists, it is overwritten.
//
// https://redis.io/commands/sunionstore/
func (r *RedisGroupSet) SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	v, err := r.redis.Do(ctx, "SUNIONSTORE", destination, keys)
	return v.Int64(), err
}

// SDiff  Returns the members of the set resulting from the difference between the first set and all the successive sets.
//
// https://redis.io/commands/sdiff/
func (r *RedisGroupSet) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	v, err := r.redis.Do(ctx, "SDIFF", keys)
	return v.Strings(), err
}

// SDiffStore This command is equal to SDIFF, but instead of returning the resulting set, it is stored in destination.
//
//If destination already exists, it is overwritten.
//
// https://redis.io/commands/sdiffstore/
func (r *RedisGroupSet) SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	v, err := r.redis.Do(ctx, "SDIFFSTORE", destination, keys)
	return v.Int64(), err
}

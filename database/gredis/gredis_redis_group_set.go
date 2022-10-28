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

// RedisGroupSet provides set functions for redis.
type RedisGroupSet struct {
	redis *Redis
}

// GroupSet creates and returns RedisGroupSet.
func (r *Redis) GroupSet() RedisGroupSet {
	return RedisGroupSet{
		redis: r,
	}
}

// SAdd adds the specified members to the set stored at key.
// Specified members that are already a member of this set are ignored.
// If key does not exist, a new set is created before adding the specified members.
//
// An error is returned when the value stored at key is not a set.
//
// It returns the number of elements that were added to the set,
// not including all the elements already present in the set.
//
// https://redis.io/commands/sadd/
func (r RedisGroupSet) SAdd(ctx context.Context, key string, member interface{}, members ...interface{}) (int64, error) {
	var s = []interface{}{key}
	s = append(s, member)
	s = append(s, members...)
	v, err := r.redis.Do(ctx, "SAdd", s...)
	return v.Int64(), err
}

// SIsMember returns if member is a member of the set stored at key.
//
// It returns:
// - 1 if the element is a member of the set.
// - 0 if the element is not a member of the set, or if key does not exist.
//
// https://redis.io/commands/sismember/
func (r RedisGroupSet) SIsMember(ctx context.Context, key string, member interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "SIsMember", key, member)
	return v.Int64(), err
}

// SPop removes and returns one or more random members from the set value store at key.
//
// This operation is similar to SRandMember, that returns one or more random elements from a set but
// does not remove it.
// By default, the command pops a single member from the set. When provided with the optional count
// argument, the reply will consist of up to count members, depending on the set's cardinality.
//
// It returns:
// - When called without the count argument:
//   Bulk string reply: the removed member, or nil when key does not exist.
// - When called with the count argument:
//   Array reply: the removed members, or an empty array when key does not exist.
//
// https://redis.io/commands/spop/
func (r RedisGroupSet) SPop(ctx context.Context, key string, count ...int) (*gvar.Var, error) {
	var s = []interface{}{key}
	s = append(s, gconv.Interfaces(count)...)
	v, err := r.redis.Do(ctx, "SPop", s...)
	return v, err
}

// SRandMember called with just the key argument, return a random element from the set value stored
// at key.
// If the provided count argument is positive, return an array of distinct elements.
// The array's length is either count or the set's cardinality (SCard), whichever is lower.
// If called with a negative count, the behavior changes and the command is allowed to return the
// same element multiple times. In this case, the number of returned elements is the absolute value
// of the specified count.
//
// It returns:
// - Bulk string reply: without the additional count argument, the command returns a Bulk Reply with the
//   randomly selected element, or nil when key does not exist.
// - Array reply: when the additional count argument is passed, the command returns an array of elements,
//   or an empty array when key does not exist.
//
// https://redis.io/commands/srandmember/
func (r RedisGroupSet) SRandMember(ctx context.Context, key string, count ...int) (*gvar.Var, error) {
	var s = []interface{}{key}
	s = append(s, gconv.Interfaces(count)...)
	v, err := r.redis.Do(ctx, "SRandMember", s...)
	return v, err
}

// SRem removes the specified members from the set stored at key.
// Specified members that are not a member of this set are ignored.
// If key does not exist, it is treated as an empty set and this command returns 0.
//
// An error is returned when the value stored at key is not a set.
//
// It returns the number of members that were removed from the set, not including non existing members.
//
// https://redis.io/commands/srem/
func (r RedisGroupSet) SRem(ctx context.Context, key string, member interface{}, members ...interface{}) (int64, error) {
	var s = []interface{}{key}
	s = append(s, member)
	s = append(s, members...)
	v, err := r.redis.Do(ctx, "SRem", s...)
	return v.Int64(), err
}

// SMove moves member from the set at source to the set at destination.
// This operation is atomic. In every given moment the element will appear to be a member of source or
// destination for other clients.
// If the source set does not exist or does not contain the specified element, no operation is performed and 0
// is returned. Otherwise, the element is removed from the source set and added to the destination set.
// When the specified element already exists in the destination set, it is only removed from the source set.
//
// An error is returned if source or destination does not hold a set value.
//
// It returns:
// - 1 if the element is moved.
// - 0 if the element is not a member of source and no operation was performed.
//
// https://redis.io/commands/smove/
func (r RedisGroupSet) SMove(ctx context.Context, source, destination string, member interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "SMove", source, destination, member)
	return v.Int64(), err
}

// SCard returns the set cardinality (number of elements) of the set stored at key.
//
// It returns the cardinality (number of elements) of the set, or 0 if key does not exist.
//
// https://redis.io/commands/scard/
func (r RedisGroupSet) SCard(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "SCard", key)
	return v.Int64(), err
}

// SMembers returns all the members of the set value stored at key.
// This has the same effect as running SINTER with one argument key.
//
// It returns all elements of the set.
//
// https://redis.io/commands/smembers/
func (r RedisGroupSet) SMembers(ctx context.Context, key string) ([]*gvar.Var, error) {
	v, err := r.redis.Do(ctx, "SMembers", key)
	return v.Vars(), err
}

// SMIsMember returns whether each member is a member of the set stored at key.
//
// For every member, 1 is returned if the value is a member of the set, or 0 if the element is not a member of
// the set or if key does not exist.
//
// It returns list representing the membership of the given elements, in the same order as they are requested.
//
// https://redis.io/commands/smismember/
func (r RedisGroupSet) SMIsMember(ctx context.Context, key, member interface{}, members ...interface{}) ([]int, error) {
	var s = []interface{}{key, member}
	s = append(s, members...)
	v, err := r.redis.Do(ctx, "SMIsMember", s...)
	return v.Ints(), err
}

// SInter returns the members of the set resulting from the intersection of all the given sets.
//
// It returns list with members of the resulting set.
//
// https://redis.io/commands/sinter/
func (r RedisGroupSet) SInter(ctx context.Context, key string, keys ...string) ([]*gvar.Var, error) {
	var s = []interface{}{key}
	s = append(s, gconv.Interfaces(keys)...)
	v, err := r.redis.Do(ctx, "SInter", s...)
	return v.Vars(), err
}

// SInterStore is equal to SInter, but instead of returning the resulting set, it is stored in
// destination.
//
// If destination already exists, it is overwritten.
//
// It returns the number of elements in the resulting set.
//
// https://redis.io/commands/sinterstore/
func (r RedisGroupSet) SInterStore(ctx context.Context, destination string, key string, keys ...string) (int64, error) {
	var s = []interface{}{destination, key}
	s = append(s, gconv.Interfaces(keys)...)
	v, err := r.redis.Do(ctx, "SInterStore", s...)
	return v.Int64(), err
}

// SUnion returns the members of the set resulting from the union of all the given sets.
//
// It returns list with members of the resulting set.
//
// https://redis.io/commands/sunion/
func (r RedisGroupSet) SUnion(ctx context.Context, key string, keys ...string) ([]*gvar.Var, error) {
	var s = []interface{}{key}
	s = append(s, gconv.Interfaces(keys)...)
	v, err := r.redis.Do(ctx, "SUnion", s...)
	return v.Vars(), err
}

// SUnionStore is equal to SUnion, but instead of returning the resulting set, it is stored in destination.
//
//  If destination already exists, it is overwritten.
//
// It returns the number of elements in the resulting set.
//
// https://redis.io/commands/sunionstore/
func (r RedisGroupSet) SUnionStore(ctx context.Context, destination, key string, keys ...string) (int64, error) {
	var s = []interface{}{destination, key}
	s = append(s, gconv.Interfaces(keys)...)
	v, err := r.redis.Do(ctx, "SUnionStore", s...)
	return v.Int64(), err
}

// SDiff returns the members of the set resulting from the difference between the first set and all the
// successive sets.
//
// It returns list with members of the resulting set.
//
// https://redis.io/commands/sdiff/
func (r RedisGroupSet) SDiff(ctx context.Context, key string, keys ...string) ([]*gvar.Var, error) {
	var s = []interface{}{key}
	s = append(s, gconv.Interfaces(keys)...)
	v, err := r.redis.Do(ctx, "SDiff", s...)
	return v.Vars(), err
}

// SDiffStore is equal to SDiff, but instead of returning the resulting set, it is stored in destination.
//
// If destination already exists, it is overwritten.
//
// It returns the number of elements in the resulting set.
//
// https://redis.io/commands/sdiffstore/
func (r RedisGroupSet) SDiffStore(ctx context.Context, destination string, key string, keys ...string) (int64, error) {
	var s = []interface{}{destination, key}
	s = append(s, gconv.Interfaces(keys)...)
	v, err := r.redis.Do(ctx, "SDiffStore", s...)
	return v.Int64(), err
}

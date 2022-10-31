// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
)

// RedisGroupSortedSet provides sorted set functions for redis.
type RedisGroupSortedSet struct {
	redis *Redis
}

// GroupSortedSet creates and returns RedisGroupSortedSet.
func (r *Redis) GroupSortedSet() RedisGroupSortedSet {
	return RedisGroupSortedSet{
		redis: r,
	}
}

// ZAddOption provides options for function ZAdd.
type ZAddOption struct {
	XX bool // Only update elements that already exist. Don't add new elements.
	NX bool // Only add new elements. Don't update already existing elements.
	// Only update existing elements if the new score is less than the current score.
	// This flag doesn't prevent adding new elements.
	LT bool

	// Only update existing elements if the new score is greater than the current score.
	// This flag doesn't prevent adding new elements.
	GT bool

	// Modify the return value from the number of new elements added, to the total number of elements changed (CH is an abbreviation of changed).
	// Changed elements are new elements added and elements already existing for which the score was updated.
	// So elements specified in the command line having the same score as they had in the past are not counted.
	// Note: normally the return value of ZAdd only counts the number of new elements added.
	CH bool

	// When this option is specified ZAdd acts like ZIncrBy. Only one score-element pair can be specified in this mode.
	INCR bool
}

// ZAddMember is element struct for set.
type ZAddMember struct {
	Score  float64
	Member interface{}
}

// ZAdd adds all the specified members with the specified scores to the sorted set stored at key.
// It is possible to specify multiple score / member pairs.
// If a specified member is already a member of the sorted set, the score is updated and the element reinserted
// at the right position to ensure the correct ordering.
//
// If key does not exist, a new sorted set with the specified members as sole members is created, like if the
// sorted set was empty. If the key exists but does not hold a sorted set, an error is returned.
//
// The score values should be the string representation of a double precision floating point number. +inf and
// -inf values are valid values as well.
//
// It returns:
// - When used without optional arguments, the number of elements added to the sorted set (excluding score updates).
// - If the CH option is specified, the number of elements that were changed (added or updated).
//
// If the INCR option is specified, the return value will be Bulk string reply:
// - The new score of member (a double precision floating point number) represented as string, or nil if the operation
//   was aborted (when called with either the XX or the NX option).
//
// https://redis.io/commands/zadd/
func (r RedisGroupSortedSet) ZAdd(ctx context.Context, key string, option *ZAddOption, member ZAddMember, members ...ZAddMember) (*gvar.Var, error) {
	s := mustMergeOptionToArgs(
		[]interface{}{key}, option,
	)
	s = append(s, member.Score, member.Member)
	for _, item := range members {
		s = append(s, item.Score, item.Member)
	}
	v, err := r.redis.Do(ctx, "ZAdd", s...)
	return v, err
}

// ZScore Returns the score of member in the sorted set at key.
//
// If member does not exist in the sorted set, or key does not exist, nil is returned.
//
// It returns the score of member (a double precision floating point number), represented as string.
//
// https://redis.io/commands/zscore/
func (r RedisGroupSortedSet) ZScore(ctx context.Context, key string, member interface{}) (float64, error) {
	v, err := r.redis.Do(ctx, "ZScore", key, member)
	return v.Float64(), err
}

// ZIncrBy increments the score of member in the sorted set stored at key by increment.
// If member does not exist in the sorted set, it is added with increment as its score (as if its previous score
// was 0.0). If key does not exist, a new sorted set with the specified member as its sole member is created.
//
// An error is returned when key exists but does not hold a sorted set.
//
// The score value should be the string representation of a numeric value, and accepts double precision floating
// point numbers. It is possible to provide a negative value to decrement the score.
//
// It returns the new score of member (a double precision floating point number).
//
// https://redis.io/commands/zincrby/
func (r RedisGroupSortedSet) ZIncrBy(ctx context.Context, key string, increment float64, member interface{}) (float64, error) {
	v, err := r.redis.Do(ctx, "ZIncrBy", key, increment, member)
	return v.Float64(), err
}

// ZCard returns the sorted set cardinality (number of elements) of the sorted set stored at key.
//
// It returns the cardinality (number of elements) of the sorted set, or 0 if key does not exist.
//
// https://redis.io/commands/zcard/
func (r RedisGroupSortedSet) ZCard(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "ZCard", key)
	return v.Int64(), err
}

// ZCount returns the number of elements in the sorted set at key with a score between min and max.
//
// The min and max arguments have the same semantic as described for ZRangeByScore.
//
// Note: the command has a complexity of just O(log(N)) because it uses elements ranks (see ZRANK) to get an
// idea of the range. Because of this there is no need to do a work proportional to the size of the range.
//
// It returns the number of elements in the specified score range.
//
// https://redis.io/commands/zcount/
func (r RedisGroupSortedSet) ZCount(ctx context.Context, key string, min, max string) (int64, error) {
	v, err := r.redis.Do(ctx, "ZCount", key, min, max)
	return v.Int64(), err
}

// ZRangeOption provides extra option for ZRange function.
type ZRangeOption struct {
	ByScore bool
	ByLex   bool
	// The optional REV argument reverses the ordering, so elements are ordered from highest to lowest score,
	// and score ties are resolved by reverse lexicographical ordering.
	Rev   bool
	Limit *ZRangeOptionLimit
	// The optional WithScores argument supplements the command's reply with the scores of elements returned.
	WithScores bool
}

// ZRangeOptionLimit provides LIMIT argument for ZRange function.
// The optional LIMIT argument can be used to obtain a sub-range from the matching elements
// (similar to SELECT LIMIT offset, count in SQL). A negative `Count` returns all elements from the `Offset`.
type ZRangeOptionLimit struct {
	Offset *int
	Count  *int
}

// ZRange return the specified range of elements in the sorted set stored at <key>.
//
// ZRange can perform different types of range queries: by index (rank), by the score, or by lexicographical
// order.
//
// https://redis.io/commands/zrange/
func (r RedisGroupSortedSet) ZRange(ctx context.Context, key string, start, stop int64, option ...ZRangeOption) ([]*gvar.Var, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	v, err := r.redis.Do(ctx, "ZRange", mustMergeOptionToArgs(
		[]interface{}{key, start, stop}, usedOption,
	)...)
	return v.Vars(), err
}

// ZRevRangeOption provides options for function ZRevRange.
type ZRevRangeOption struct {
	WithScores bool
}

// ZRevRange returns the specified range of elements in the sorted set stored at key.
// The elements are considered to be ordered from the highest to the lowest score.
// Descending lexicographical order is used for elements with equal score.
//
// Apart from the reversed ordering, ZRevRange is similar to ZRange.
//
// It returns list of elements in the specified range (optionally with their scores).
//
// https://redis.io/commands/zrevrange/
func (r RedisGroupSortedSet) ZRevRange(ctx context.Context, key string, start, stop int64, option ...ZRevRangeOption) (*gvar.Var, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	return r.redis.Do(ctx, "ZRevRange", mustMergeOptionToArgs(
		[]interface{}{key, start, stop}, usedOption,
	)...)
}

// ZRank returns the rank of member in the sorted set stored at key, with the scores ordered from low to high.
// The rank (or index) is 0-based, which means that the member with the lowest score has rank 0.
//
// Use ZRevRank to get the rank of an element with the scores ordered from high to low.
//
// It returns:
// - If member exists in the sorted set, Integer reply: the rank of member.
// - If member does not exist in the sorted set or key does not exist, Bulk string reply: nil.
//
// https://redis.io/commands/zrank/
func (r RedisGroupSortedSet) ZRank(ctx context.Context, key string, member interface{}) (*gvar.Var, error) {
	v, err := r.redis.Do(ctx, "ZRank", key, member)
	return v, err
}

// ZRevRank returns the rank of member in the sorted set stored at key, with the scores ordered from high to low.
// The rank (or index) is 0-based, which means that the member with the highest score has rank 0.
//
// Use ZRank to get the rank of an element with the scores ordered from low to high.
//
// It returns:
// - If member exists in the sorted set, Integer reply: the rank of member.
// - If member does not exist in the sorted set or key does not exist, Bulk string reply: nil.
//
// https://redis.io/commands/zrevrank/
func (r RedisGroupSortedSet) ZRevRank(ctx context.Context, key string, member interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "ZRevRank", key, member)
	return v.Int64(), err
}

// ZRem remove the specified members from the sorted set stored at key.
// Non-existing members are ignored.
//
// An error is returned when key exists and does not hold a sorted set.
//
// It returns the number of members removed from the sorted set, not including non existing members.
//
// https://redis.io/commands/zrem/
func (r RedisGroupSortedSet) ZRem(ctx context.Context, key string, member interface{}, members ...interface{}) (int64, error) {
	var s = []interface{}{key}
	s = append(s, member)
	s = append(s, members...)
	v, err := r.redis.Do(ctx, "ZRem", s...)
	return v.Int64(), err
}

// ZRemRangeByRank removes all elements in the sorted set stored at key with rank between start and stop.
// Both start and stop are 0 -based indexes with 0 being the element with the lowest score.
//
// These indexes can be negative numbers, where they indicate offsets starting at the element with the highest
// score. For example: -1 is the element with the highest score, -2 the element with the second-highest score
// and so forth.
//
// It returns the number of elements removed.
//
// https://redis.io/commands/zremrangebyrank/
func (r RedisGroupSortedSet) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	v, err := r.redis.Do(ctx, "ZRemRangeByRank", key, start, stop)
	return v.Int64(), err
}

// ZRemRangeByScore removes all elements in the sorted set stored at key with a score between min and max
// (inclusive).
//
// It returns the number of elements removed.
//
// https://redis.io/commands/zremrangebyscore/
func (r RedisGroupSortedSet) ZRemRangeByScore(ctx context.Context, key string, min, max string) (int64, error) {
	v, err := r.redis.Do(ctx, "ZRemRangeByScore", key, min, max)
	return v.Int64(), err
}

// ZRemRangeByLex removes all elements in the sorted set stored at key between the
// lexicographical range specified by min and max.
//
// The meaning of min and max are the same of the ZRangeByLex command.
// Similarly, this command actually removes the same elements that ZRangeByLex would return if called with the
// same min and max arguments.
//
// It returns the number of elements removed.
//
// https://redis.io/commands/zremrangebylex/
func (r RedisGroupSortedSet) ZRemRangeByLex(ctx context.Context, key string, min, max string) (int64, error) {
	v, err := r.redis.Do(ctx, "ZRemRangeByLex", key, min, max)
	return v.Int64(), err
}

// ZLexCount all the elements in a sorted set are inserted with the same score,
// in order to force lexicographical ordering, this command returns the number of elements in the sorted
// set at key with a value between min and max.
//
// The min and max arguments have the same meaning as described for ZRangeByLex.
//
// Note: the command has a complexity of just O(log(N)) because it uses elements ranks (see ZRank) to get an
// idea of the range. Because of this there is no need to do a work proportional to the size of the range.
//
// It returns the number of elements in the specified score range.
//
// https://redis.io/commands/zlexcount/
func (r RedisGroupSortedSet) ZLexCount(ctx context.Context, key, min, max string) (int64, error) {
	v, err := r.redis.Do(ctx, "ZLexCount", key, min, max)
	return v.Int64(), err
}

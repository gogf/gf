// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"

	"github.com/gogf/gf/v2/util/gconv"
)

// RedisGroupSortedSet is the redis group object for sorted set operations.
type RedisGroupSortedSet struct {
	redis *Redis
}

// GroupSortedSet creates and returns a redis group object for sorted set operations.
func (r *Redis) GroupSortedSet() RedisGroupSortedSet {
	return RedisGroupSortedSet{
		redis: r,
	}
}

// ZAddOption is the option for function ZAdd.
type ZAddOption struct {
	NX   bool // Only update elements that already exist. Don't add new elements.
	XX   bool // Only add new elements. Don't update already existing elements.
	GT   bool // Only update existing elements if the new score is greater than the current score.
	LT   bool // Only update existing elements if the new score is less than the current score.
	CH   bool // Modify the return value from the number of new elements added, to the total number of elements changed.
	INCR bool // When this option is specified ZAdd acts like ZIncrBy.
}

// ZAddItem contains the item for function ZAdd.
type ZAddItem struct {
	Score  float64
	Member string
}

// ZAdd add all the specified members with the specified scores to the sorted set stored at key.
// It is possible to specifyZAdd multiple score / member pairs.
// If a specified member is already a member of the sorted set, the score is updated and the element reinserted
// at the right position to ensure the correct ordering.
//
// If key does not exist, a new sorted set with the specified members as sole members is created, like if the
// sorted set was empty. If the key exists but does not hold a sorted set, an error is returned.
//
// The score values should be the string representation of a double precision floating point number. +inf and
// -inf values are valid values as well.
//
// https://redis.io/commands/zadd/
func (r RedisGroupSortedSet) ZAdd(ctx context.Context, key string, items []ZAddItem, option ...ZAddOption) (int64, error) {
	var (
		usedOption interface{}
		args       = make([]interface{}, 0)
	)
	if len(option) > 0 {
		usedOption = option[0]
	}
	args = mustMergeOptionToArgs(
		[]interface{}{key}, usedOption,
	)
	for _, item := range items {
		args = append(args, item.Score, item.Member)
	}
	v, err := r.redis.Do(ctx, "ZRange", args...)
	return v.Int64(), err
}

// ZScore Returns the score of member in the sorted set at key.
//
// If member does not exist in the sorted set, or key does not exist, nil is returned.
//
// https://redis.io/commands/zscore/
func (r RedisGroupSortedSet) ZScore(ctx context.Context, key string, member string) (float64, error) {
	v, err := r.redis.Do(ctx, "ZScore", key, member)
	return v.Float64(), err
}

// ZIncrBy increment the score of member in the sorted set stored at key by increment.
// If member does not exist in the sorted set, it is added with increment as its score (as if its previous score
// was 0.0). If key does not exist, a new sorted set with the specified member as its sole member is created.
//
// An error is returned when key exists but does not hold a sorted set.
//
// The score value should be the string representation of a numeric value, and accepts double precision floating
// point numbers. It is possible to provide a negative value to decrement the score.
//
// https://redis.io/commands/zincrby/
func (r RedisGroupSortedSet) ZIncrBy(ctx context.Context, key string, value float64, member string) (float64, error) {
	v, err := r.redis.Do(ctx, "ZIncrBy", key, value, member)
	return v.Float64(), err
}

// ZCard return the sorted set cardinality (number of elements) of the sorted set stored at key.
//
// https://redis.io/commands/zcard/
func (r RedisGroupSortedSet) ZCard(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "ZCard", key)
	return v.Int64(), err
}

// ZCount return the number of elements in the sorted set at key with a score between min and max.
//
// The min and max arguments have the same semantic as described for ZRangeByScore.
//
// Note: the command has a complexity of just O(log(N)) because it uses elements ranks (see ZRank) to get an
// idea of the range. Because of this there is no need to do a work proportional to the size of the range.
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
	Offset int
	Count  int
}

// ZRange return the specified range of elements in the sorted set stored at <key>.
//
// ZRange can perform different types of range queries: by index (rank), by the score, or by lexicographical
// order.
//
// https://redis.io/commands/zrange/
func (r RedisGroupSortedSet) ZRange(ctx context.Context, key string, start, stop int64, option ...ZRangeOption) ([]string, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	v, err := r.redis.Do(ctx, "ZRange", mustMergeOptionToArgs(
		[]interface{}{key, start, stop}, usedOption,
	)...)
	return v.Strings(), err
}

// ZRevRange return the specified range of elements in the sorted set stored at key.
// The elements are considered to be ordered from the highest to the lowest score.
// Descending lexicographical order is used for elements with equal score.
//
// Apart from the reversed ordering, ZRevRange is similar to ZRange.
//
// https://redis.io/commands/zrevrange/
func (r RedisGroupSortedSet) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	v, err := r.redis.Do(ctx, "ZRevRange", key, start, stop)
	return gconv.SliceStr(v), err
}

// ZRank return the rank of member in the sorted set stored at key, with the scores ordered from low to high.
// The rank (or index) is 0-based, which means that the member with the lowest score has rank 0.
//
// Use ZRevRank to get the rank of an element with the scores ordered from high to low.
//
// https://redis.io/commands/zrank/
func (r RedisGroupSortedSet) ZRank(ctx context.Context, key, member string) (int64, error) {
	v, err := r.redis.Do(ctx, "ZRank", key, member)
	return v.Int64(), err
}

// ZRevRank return the rank of member in the sorted set stored at key, with the scores ordered from high to low.
// The rank (or index) is 0-based, which means that the member with the highest score has rank 0.
//
// Use ZRank to get the rank of an element with the scores ordered from low to high.
//
// https://redis.io/commands/zrevrank/
func (r RedisGroupSortedSet) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	v, err := r.redis.Do(ctx, "ZRevRank", key, member)
	return v.Int64(), err
}

// ZRem remove the specified members from the sorted set stored at key.
// Non existing members are ignored.
//
// An error is returned when key exists and does not hold a sorted set.
//
// https://redis.io/commands/zrem/
func (r RedisGroupSortedSet) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "ZRem", key, members)
	return v.Int64(), err
}

// ZRemRangeByRank remove all elements in the sorted set stored at key with rank between start and stop.
// Both start and stop are 0 -based indexes with 0 being the element with the lowest score.
// These indexes can be negative numbers, where they indicate offsets starting at the element with the highest
// score. For example: -1 is the element with the highest score, -2 the element with the second highest score
// and so forth.
//
// https://redis.io/commands/zremrangebyrank/
func (r RedisGroupSortedSet) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	v, err := r.redis.Do(ctx, "ZRemRangeByRank", key, start, stop)
	return v.Int64(), err
}

// ZRemRangeByScore remove all elements in the sorted set stored at key with a score between min and max
// (inclusive).
//
// https://redis.io/commands/zremrangebyscore/
func (r RedisGroupSortedSet) ZRemRangeByScore(ctx context.Context, key string, min, max string) (int64, error) {
	v, err := r.redis.Do(ctx, "ZRemRangeByScore", key, min, max)
	return v.Int64(), err
}

// ZRemRangeByLex all the elements in a sorted set are inserted with the same score, in order to force
// lexicographical ordering, this command removes all elements in the sorted set stored at key between the
// lexicographical range specified by min and max.
//
// The meaning of min and max are the same of the ZRANGEBYLEX command.
// Similarly, this command actually removes the same elements that ZRANGEBYLEX would return if called with the
// same min and max arguments.
func (r RedisGroupSortedSet) ZRemRangeByLex(ctx context.Context, key string, min, max string) (int64, error) {
	v, err := r.redis.Do(ctx, "ZRemRangeByLex", key, min, max)
	return v.Int64(), err
}

// ZLexCount all the elements in a sorted set are inserted with the same score,
// in order to force lexicographical ordering, this command returns the number of elements in the sorted
// set at key with a value between min and max.
//
// The min and max arguments have the same meaning as described for ZRANGEBYLEX.
//
// Note: the command has a complexity of just O(log(N)) because it uses elements ranks (see ZRANK) to get an
// idea of the range. Because of this there is no need to do a work proportional to the size of the range.
//
// https://redis.io/commands/zlexcount/
func (r RedisGroupSortedSet) ZLexCount(ctx context.Context, key, min, max string) (int64, error) {
	v, err := r.redis.Do(ctx, "ZLexCount", key, min, max)
	return v.Int64(), err
}

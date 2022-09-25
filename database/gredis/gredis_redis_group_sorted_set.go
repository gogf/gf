// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
)

type RedisGroupSortedSet struct {
	redis *Redis
}

func (r *Redis) SortedSet() *RedisGroupSortedSet {
	return &RedisGroupSortedSet{
		redis: r,
	}
}

func (RedisGroupSortedSet) ZAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZScore(ctx context.Context, key string, member string) (float64, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZIncrBy(ctx context.Context, key string, value float64, member string) (float64, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZCard(ctx context.Context, key string) (int64, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZCount(ctx context.Context, key string, min, max string) (int64, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZRank(ctx context.Context, key, member string) (int64, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZRemRangeByScore(ctx context.Context, key string, min, max string) (int64, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZRemRangeByLex(ctx context.Context, key string, min, max string) (int64, error) {
	panic("implement me")
}

func (RedisGroupSortedSet) ZLexCount(ctx context.Context, key, min, max string) (int64, error) {
	panic("implement me")
}

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

func (RedisGroupSet) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	panic("implement me")
}

func (RedisGroupSet) SIsMember(ctx context.Context, key string, member string) (bool, error) {
	panic("implement me")
}

func (RedisGroupSet) SPop(ctx context.Context, key string) (string, error) {
	panic("implement me")
}

func (RedisGroupSet) SRandMember(ctx context.Context, key string) (string, error) {
	panic("implement me")
}

func (RedisGroupSet) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	panic("implement me")
}

func (RedisGroupSet) SMove(ctx context.Context, source, destination, member string) (bool, error) {
	panic("implement me")
}

func (RedisGroupSet) SCard(ctx context.Context, key string) (int64, error) {
	panic("implement me")
}

func (RedisGroupSet) SMembers(ctx context.Context, key string) ([]string, error) {
	panic("implement me")
}

func (RedisGroupSet) SInter(ctx context.Context, keys ...string) ([]string, error) {
	panic("implement me")
}

func (RedisGroupSet) SInterStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	panic("implement me")
}

func (RedisGroupSet) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	panic("implement me")
}

func (RedisGroupSet) SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	panic("implement me")
}

func (RedisGroupSet) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	panic("implement me")
}

func (RedisGroupSet) SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	panic("implement me")
}

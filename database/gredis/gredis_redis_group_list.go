// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"time"
)

type RedisGroupList struct {
	redis *Redis
}

func (r *Redis) List() *RedisGroupList {
	return &RedisGroupList{
		redis: r,
	}
}

func (RedisGroupList) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	panic("implement me")
}

func (RedisGroupList) LPushX(ctx context.Context, key string, value string) (int64, error) {
	panic("implement me")
}

func (RedisGroupList) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	panic("implement me")
}

func (RedisGroupList) RPushX(ctx context.Context, key string, value string) (int64, error) {
	panic("implement me")
}

func (RedisGroupList) LPop(ctx context.Context, key string) (string, error) {
	panic("implement me")
}

func (RedisGroupList) RPop(ctx context.Context, key string) (string, error) {
	panic("implement me")
}

func (RedisGroupList) RPopLPush(ctx context.Context, source, destination string) (string, error) {
	panic("implement me")
}

func (RedisGroupList) LRem(ctx context.Context, key string, count int64, value string) (int64, error) {
	panic("implement me")
}

func (RedisGroupList) LLen(ctx context.Context, key string) (int64, error) {
	panic("implement me")
}

func (RedisGroupList) LIndex(ctx context.Context, key string, index int64) (string, error) {
	panic("implement me")
}

func (RedisGroupList) LInsert(ctx context.Context, key, op string, pivot, value string) (int64, error) {
	panic("implement me")
}

func (RedisGroupList) LSet(ctx context.Context, key string, index int64, value string) (string, error) {
	panic("implement me")
}

func (RedisGroupList) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	panic("implement me")
}

func (RedisGroupList) LTrim(ctx context.Context, key string, start, stop int64) (string, error) {
	panic("implement me")
}

func (RedisGroupList) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	panic("implement me")
}

func (RedisGroupList) BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	panic("implement me")
}

func (RedisGroupList) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error) {
	panic("implement me")
}

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

func (RedisGroupHash) HSet(ctx context.Context, key, field, value string) (int64, error) {
	panic("implement me")
}

func (RedisGroupHash) HSetNX(ctx context.Context, key, field, value string) (bool, error) {
	panic("implement me")
}

func (RedisGroupHash) HGet(ctx context.Context, key, field string) (string, error) {
	panic("implement me")
}

func (RedisGroupHash) HExists(ctx context.Context, key, field string) (bool, error) {
	panic("implement me")
}

func (RedisGroupHash) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	panic("implement me")
}

func (RedisGroupHash) HLen(ctx context.Context, key string) (int64, error) {
	panic("implement me")
}

func (RedisGroupHash) HIncrBy(ctx context.Context, key, field string, value int64) (int64, error) {
	panic("implement me")
}

func (RedisGroupHash) HIncrByFloat(ctx context.Context, key, field string, value float64) (float64, error) {
	panic("implement me")
}

func (RedisGroupHash) HMSet(ctx context.Context, key string, fields map[string]string) (bool, error) {
	panic("implement me")
}

func (RedisGroupHash) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	panic("implement me")
}

func (RedisGroupHash) HKeys(ctx context.Context, key string) ([]string, error) {
	panic("implement me")
}

func (RedisGroupHash) HVals(ctx context.Context, key string) ([]string, error) {
	panic("implement me")
}

func (RedisGroupHash) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	panic("implement me")
}

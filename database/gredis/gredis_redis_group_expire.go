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

type RedisGroupExpire struct {
	redis *Redis
}

func (r *Redis) Expire() *RedisGroupExpire {
	return &RedisGroupExpire{
		redis: r,
	}
}

func (RedisGroupExpire) Expire(ctx context.Context, key string, seconds time.Duration) (bool, error) {
	panic("implement me")
}

func (RedisGroupExpire) ExpireAt(ctx context.Context, key string, time time.Time) (bool, error) {
	panic("implement me")
}

func (RedisGroupExpire) TTL(ctx context.Context, key string) (time.Duration, error) {
	panic("implement me")
}

func (RedisGroupExpire) PErsist(ctx context.Context, key string, time time.Duration) (bool, error) {
	panic("implement me")
}

func (RedisGroupExpire) PExpire(ctx context.Context, key string, time time.Duration) (bool, error) {
	panic("implement me")
}

func (RedisGroupExpire) PExpireAt(ctx context.Context, key string, time time.Time) (bool, error) {
	panic("implement me")
}

func (RedisGroupExpire) PTTL(ctx context.Context, key string) (time.Duration, error) {
	panic("implement me")
}

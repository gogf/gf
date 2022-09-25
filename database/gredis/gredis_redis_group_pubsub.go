// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
)

type RedisGroupPubSub struct {
	redis *Redis
}

func (r *Redis) PubSub() *RedisGroupPubSub {
	return &RedisGroupPubSub{
		redis: r,
	}
}

func (RedisGroupPubSub) Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	panic("implement me")
}

func (RedisGroupPubSub) Subscribe(ctx context.Context, channels ...string) (interface{}, error) {
	panic("implement me")
}

func (RedisGroupPubSub) PSubscribe(ctx context.Context, channels ...string) (interface{}, error) {
	panic("implement me")
}

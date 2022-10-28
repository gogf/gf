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

// RedisGroupPubSub is the redis group object for pub/sub operations.
type RedisGroupPubSub struct {
	redis *Redis
}

// GroupPubSub creates and returns a redis group object for pub/sub operations.
func (r *Redis) GroupPubSub() RedisGroupPubSub {
	return RedisGroupPubSub{
		redis: r,
	}
}

// Publish post a message to the given channel.
//
// In a Redis Cluster clients can publish to every node. The cluster makes sure that published
// messages are forwarded as needed, so clients can subscribe to any channel by connecting to any one
// of the nodes.
//
// https://redis.io/commands/publish/
func (r RedisGroupPubSub) Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "Publish", channel, message)
	return v.Int64(), err
}

// Subscribe the client to the specified channels.
//
// https://redis.io/commands/subscribe/
func (r RedisGroupPubSub) Subscribe(ctx context.Context, channels ...string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "Subscribe", channels)
}

// PSubscribe the client to the given patterns.
//
// https://redis.io/commands/psubscribe/
func (r RedisGroupPubSub) PSubscribe(ctx context.Context, channels ...string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "PSubscribe", channels)
}

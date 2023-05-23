// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis

import (
	"context"

	"github.com/gogf/gf/v2/database/gredis"
)

// GroupPubSub provides pub/sub functions for redis.
type GroupPubSub struct {
	redis *Redis
}

// GroupPubSub creates and returns GroupPubSub.
func (r *Redis) GroupPubSub() gredis.IGroupPubSub {
	return GroupPubSub{
		redis: r,
	}
}

// Publish posts a message to the given channel.
//
// In a Redis Cluster clients can publish to every node. The cluster makes sure that published
// messages are forwarded as needed, so clients can subscribe to any channel by connecting to any one
// of the nodes.
//
// It returns the number of clients that received the message.
// Note that in a Redis Cluster, only clients that are connected to the same node as the publishing client
// are included in the count.
//
// https://redis.io/commands/publish/
func (r GroupPubSub) Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "Publish", channel, message)
	return v.Int64(), err
}

// Subscribe subscribes the client to the specified channels.
//
// https://redis.io/commands/subscribe/
func (r GroupPubSub) Subscribe(
	ctx context.Context, channel string, channels ...string,
) (gredis.Conn, []*gredis.Subscription, error) {
	conn, err := r.redis.Conn(ctx)
	if err != nil {
		return nil, nil, err
	}
	subs, err := conn.Subscribe(ctx, channel, channels...)
	if err != nil {
		return conn, nil, err
	}
	return conn, subs, nil
}

// PSubscribe subscribes the client to the given patterns.
//
// Supported glob-style patterns:
// - h?llo subscribes to hello, hallo and hxllo
// - h*llo subscribes to hllo and heeeello
// - h[ae]llo subscribes to hello and hallo, but not hillo
//
// Use \ to escape special characters if you want to match them verbatim.
//
// https://redis.io/commands/psubscribe/
func (r GroupPubSub) PSubscribe(
	ctx context.Context, pattern string, patterns ...string,
) (gredis.Conn, []*gredis.Subscription, error) {
	conn, err := r.redis.Conn(ctx)
	if err != nil {
		return nil, nil, err
	}
	subs, err := conn.PSubscribe(ctx, pattern, patterns...)
	if err != nil {
		return conn, nil, err
	}
	return conn, subs, nil
}

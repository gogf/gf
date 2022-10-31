// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"fmt"
)

// RedisGroupPubSub provides pub/sub functions for redis.
type RedisGroupPubSub struct {
	redis *Redis
}

// Message received as result of a PUBLISH command issued by another client.
type Message struct {
	Channel      string
	Pattern      string
	Payload      string
	PayloadSlice []string
}

// Subscription received after a successful subscription to channel.
type Subscription struct {
	Kind    string // Can be "subscribe", "unsubscribe", "psubscribe" or "punsubscribe".
	Channel string // Channel name we have subscribed to.
	Count   int    // Number of channels we are currently subscribed to.
}

// String converts current object to a readable string.
func (m *Subscription) String() string {
	return fmt.Sprintf("%s: %s", m.Kind, m.Channel)
}

// GroupPubSub creates and returns RedisGroupPubSub.
func (r *Redis) GroupPubSub() RedisGroupPubSub {
	return RedisGroupPubSub{
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
func (r RedisGroupPubSub) Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	v, err := r.redis.Do(ctx, "Publish", channel, message)
	return v.Int64(), err
}

// Subscribe subscribes the client to the specified channels.
//
// https://redis.io/commands/subscribe/
func (r RedisGroupPubSub) Subscribe(ctx context.Context, channel string, channels ...string) (*RedisConn, []*Subscription, error) {
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
func (r RedisGroupPubSub) PSubscribe(ctx context.Context, pattern string, patterns ...string) (*RedisConn, []*Subscription, error) {
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

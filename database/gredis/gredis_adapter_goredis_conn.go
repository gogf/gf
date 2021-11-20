// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type localAdapterGoRedisConn struct {
	ps    *redis.PubSub
	redis *AdapterGoRedis
}

// Do sends a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (c *localAdapterGoRedisConn) Do(ctx context.Context, command string, args ...interface{}) (reply *gvar.Var, err error) {
	switch gstr.ToLower(command) {
	case `subscribe`:
		c.ps = c.redis.client.Subscribe(ctx, gconv.Strings(args)...)

	case `psubscribe`:
		c.ps = c.redis.client.PSubscribe(ctx, gconv.Strings(args)...)

	case `unsubscribe`:
		if c.ps != nil {
			err = c.ps.Unsubscribe(ctx, gconv.Strings(args)...)
		}

	case `punsubscribe`:
		if c.ps != nil {
			err = c.ps.PUnsubscribe(ctx, gconv.Strings(args)...)
		}

	default:
		arguments := make([]interface{}, len(args)+1)
		copy(arguments, []interface{}{command})
		copy(arguments[1:], args)
		reply, err = c.resultToVar(
			c.redis.client.Do(ctx, arguments...).Result(),
		)
	}

	return
}

// Receive receives a single reply as gvar.Var from the Redis server.
func (c *localAdapterGoRedisConn) Receive(ctx context.Context) (*gvar.Var, error) {
	if c.ps != nil {
		return c.resultToVar(c.ps.Receive(ctx))
	}
	return nil, nil
}

// Close closes current PubSub or puts the connection back to connection pool.
func (c *localAdapterGoRedisConn) Close(ctx context.Context) error {
	if c.ps != nil {
		return c.ps.Close()
	}
	return nil
}

// resultToVar converts redis operation result to gvar.Var.
func (c *localAdapterGoRedisConn) resultToVar(result interface{}, err error) (*gvar.Var, error) {
	if err == redis.Nil {
		err = nil
	}
	if err == nil {
		switch v := result.(type) {
		case []byte:
			return gvar.New(string(v)), err

		case []interface{}:
			return gvar.New(gconv.Strings(v)), err

		case *redis.Message:
			result = &Message{
				Channel:      v.Channel,
				Pattern:      v.Pattern,
				Payload:      v.Payload,
				PayloadSlice: v.PayloadSlice,
			}

		case *redis.Subscription:
			result = &Subscription{
				Kind:    v.Kind,
				Channel: v.Channel,
				Count:   v.Count,
			}
		}
	}
	return gvar.New(result), err
}

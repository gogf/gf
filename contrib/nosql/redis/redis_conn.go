// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/database/gredis"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type localAdapterGoRedisConn struct {
	ps    *redis.PubSub
	redis *Redis
}

// Do send a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (c *localAdapterGoRedisConn) Do(ctx context.Context, command string, args ...interface{}) (reply *gvar.Var, err error) {
	argStrSlice := gconv.Strings(args)
	switch gstr.ToLower(command) {
	case `subscribe`:
		c.ps = c.redis.client.Subscribe(ctx, argStrSlice...)

	case `psubscribe`:
		c.ps = c.redis.client.PSubscribe(ctx, argStrSlice...)

	case `unsubscribe`:
		if c.ps != nil {
			err = c.ps.Unsubscribe(ctx, argStrSlice...)
			if err != nil {
				err = gerror.Wrapf(err, `Redis PubSub Unsubscribe failed with arguments "%v"`, argStrSlice)
			}
		}

	case `punsubscribe`:
		if c.ps != nil {
			err = c.ps.PUnsubscribe(ctx, argStrSlice...)
			if err != nil {
				err = gerror.Wrapf(err, `Redis PubSub PUnsubscribe failed with arguments "%v"`, argStrSlice)
			}
		}

	default:
		arguments := make([]interface{}, len(args)+1)
		copy(arguments, []interface{}{command})
		copy(arguments[1:], args)
		reply, err = c.resultToVar(c.redis.client.Do(ctx, arguments...).Result())
		if err != nil {
			err = gerror.Wrapf(err, `Redis Client Do failed with arguments "%v"`, arguments)
		}
	}
	return
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
			result = &gredis.Message{
				Channel:      v.Channel,
				Pattern:      v.Pattern,
				Payload:      v.Payload,
				PayloadSlice: v.PayloadSlice,
			}

		case *redis.Subscription:
			result = &gredis.Subscription{
				Kind:    v.Kind,
				Channel: v.Channel,
				Count:   v.Count,
			}
		}
	}

	return gvar.New(result), err
}

// Receive receives a single reply as gvar.Var from the Redis server.
func (c *localAdapterGoRedisConn) Receive(ctx context.Context) (*gvar.Var, error) {
	if c.ps != nil {
		v, err := c.resultToVar(c.ps.Receive(ctx))
		if err != nil {
			err = gerror.Wrapf(err, `Redis PubSub Receive failed`)
		}
		return v, err
	}
	return nil, nil
}

// Close closes current PubSub or puts the connection back to connection pool.
func (c *localAdapterGoRedisConn) Close(ctx context.Context) (err error) {
	if c.ps != nil {
		err = c.ps.Close()
		if err != nil {
			err = gerror.Wrapf(err, `Redis PubSub Close failed`)
		}
	}
	return
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gtime"
	"reflect"
)

// RedisConn is a connection of redis client.
type RedisConn struct {
	conn  Conn
	redis *Redis
}

// Do sends a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (c *RedisConn) Do(ctx context.Context, command string, args ...interface{}) (reply *gvar.Var, err error) {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
	)
	for k, v := range args {
		reflectValue = reflect.ValueOf(v)
		reflectKind = reflectValue.Kind()
		if reflectKind == reflect.Ptr {
			reflectValue = reflectValue.Elem()
			reflectKind = reflectValue.Kind()
		}
		switch reflectKind {
		case
			reflect.Struct,
			reflect.Map,
			reflect.Slice,
			reflect.Array:
			// Ignore slice type of: []byte.
			if _, ok := v.([]byte); !ok {
				if args[k], err = json.Marshal(v); err != nil {
					return nil, err
				}
			}
		}
	}
	timestampMilli1 := gtime.TimestampMilli()
	reply, err = c.conn.Do(ctx, command, args...)
	timestampMilli2 := gtime.TimestampMilli()

	// Tracing.
	c.addTracingItem(ctx, &tracingItem{
		err:       err,
		command:   command,
		args:      args,
		costMilli: timestampMilli2 - timestampMilli1,
	})
	return
}

// Receive receives a single reply as gvar.Var from the Redis server.
func (c *RedisConn) Receive(ctx context.Context) (*gvar.Var, error) {
	return c.conn.Receive(ctx)
}

// Close puts the connection back to connection pool.
func (c *RedisConn) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

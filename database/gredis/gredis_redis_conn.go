// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"github.com/gogf/gf/container/gvar"
)

type RedisConn struct {
	conn Conn
}

func (c *RedisConn) Do(ctx context.Context, command string, args ...interface{}) (reply *gvar.Var, err error) {
	args, option := parseArgsToArgsAndOption(args)
	return c.conn.Do(ctx, command, args, option)
}

// Receive receives a single reply as gvar.Var from the Redis server.
func (c *RedisConn) Receive(ctx context.Context, option ...*Option) (*gvar.Var, error) {
	var (
		usedOption *Option
	)
	if len(option) > 0 {
		usedOption = option[0]
	} else {
		usedOption = defaultOption()
	}
	return c.conn.Receive(ctx, usedOption)
}

func (c *RedisConn) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

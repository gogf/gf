// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
)

// Redis client.
type Redis struct {
	adapter Adapter
}

const (
	errorNilRedis = `the Redis object is nil`
)

func (r *Redis) SetAdapter(adapter Adapter) {
	if r == nil {
		return
	}
	r.adapter = adapter
}

func (r *Redis) GetAdapter() Adapter {
	if r == nil {
		return nil
	}
	return r.adapter
}

// Conn returns a raw underlying connection object,
// which expose more methods to communicate with server.
// **You should call Close function manually if you do not use this connection any further.**
func (r *Redis) Conn(ctx context.Context) (*RedisConn, error) {
	if r == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, errorNilRedis)
	}
	conn, err := r.adapter.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &RedisConn{conn: conn}, nil
}

// Stats returns pool's statistics.
func (r *Redis) Stats(ctx context.Context) (Stats, error) {
	if r == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, errorNilRedis)
	}
	return r.adapter.Stats(ctx)
}

// Do sends a command to the server and returns the received reply.
// Do automatically get a connection from pool, and close it when the reply received.
// It does not really "close" the connection, but drops it back to the connection pool.
func (r *Redis) Do(ctx context.Context, command string, args ...interface{}) (*gvar.Var, error) {
	if r == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, errorNilRedis)
	}
	conn, err := r.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			intlog.Error(ctx, err)
		}
	}()
	return conn.Do(ctx, command, args...)
}

func (r *Redis) Close(ctx context.Context) error {
	if r == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, errorNilRedis)
	}
	return r.adapter.Close(ctx)
}

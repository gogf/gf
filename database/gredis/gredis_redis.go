// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
)

// Redis client.
type Redis struct {
	adapter Adapter
}

const (
	errorNilRedis = `the Redis object is nil`
)

// SetAdapter sets custom adapter for current redis client.
func (r *Redis) SetAdapter(adapter Adapter) {
	if r == nil {
		return
	}
	r.adapter = adapter
}

// GetAdapter returns the adapter that is set in current redis client.
func (r *Redis) GetAdapter() Adapter {
	if r == nil {
		return nil
	}
	return r.adapter
}

// Conn retrieves and returns a connection object for continuous operations.
// Note that you should call Close function manually if you do not use this connection any further.
func (r *Redis) Conn(ctx context.Context) (*RedisConn, error) {
	if r == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, errorNilRedis)
	}
	conn, err := r.adapter.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &RedisConn{
		conn:  conn,
		redis: r,
	}, nil
}

// Do sends a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
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

// Close closes current redis client, closes its connection pool and releases all its related resources.
func (r *Redis) Close(ctx context.Context) error {
	if r == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, errorNilRedis)
	}
	return r.adapter.Close(ctx)
}

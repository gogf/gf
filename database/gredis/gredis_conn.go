// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"time"
)

// Do sends a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
// The timeout overrides the read timeout set when dialing the connection.
func (c *Conn) do(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error) {
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
	if timeout > 0 {
		conn, ok := c.Conn.(redis.ConnWithTimeout)
		if !ok {
			return gvar.New(nil), gerror.NewCode(gerror.CodeNotSupported, `current connection does not support "ConnWithTimeout"`)
		}
		return conn.DoWithTimeout(timeout, commandName, args...)
	}
	timestampMilli1 := gtime.TimestampMilli()
	reply, err = c.Conn.Do(commandName, args...)
	timestampMilli2 := gtime.TimestampMilli()

	// Tracing.
	c.addTracingItem(&tracingItem{
		err:         err,
		commandName: commandName,
		arguments:   args,
		costMilli:   timestampMilli2 - timestampMilli1,
	})
	return
}

// Ctx is a channing function which sets the context for next operation.
func (c *Conn) Ctx(ctx context.Context) *Conn {
	c.ctx = ctx
	return c
}

// Do sends a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (c *Conn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return c.do(0, commandName, args...)
}

// DoWithTimeout sends a command to the server and returns the received reply.
// The timeout overrides the read timeout set when dialing the connection.
func (c *Conn) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error) {
	return c.do(timeout, commandName, args...)
}

// DoVar retrieves and returns the result from command as gvar.Var.
func (c *Conn) DoVar(commandName string, args ...interface{}) (*gvar.Var, error) {
	return resultToVar(c.Do(commandName, args...))
}

// DoVarWithTimeout retrieves and returns the result from command as gvar.Var.
// The timeout overrides the read timeout set when dialing the connection.
func (c *Conn) DoVarWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (*gvar.Var, error) {
	return resultToVar(c.DoWithTimeout(timeout, commandName, args...))
}

// ReceiveVar receives a single reply as gvar.Var from the Redis server.
func (c *Conn) ReceiveVar() (*gvar.Var, error) {
	return resultToVar(c.Receive())
}

// ReceiveVarWithTimeout receives a single reply as gvar.Var from the Redis server.
// The timeout overrides the read timeout set when dialing the connection.
func (c *Conn) ReceiveVarWithTimeout(timeout time.Duration) (*gvar.Var, error) {
	conn, ok := c.Conn.(redis.ConnWithTimeout)
	if !ok {
		return gvar.New(nil), gerror.NewCode(gerror.CodeNotSupported, `current connection does not support "ConnWithTimeout"`)
	}
	return resultToVar(conn.ReceiveWithTimeout(timeout))
}

// resultToVar converts redis operation result to gvar.Var.
func resultToVar(result interface{}, err error) (*gvar.Var, error) {
	if err == nil {
		if result, ok := result.([]byte); ok {
			return gvar.New(gconv.UnsafeBytesToStr(result)), err
		}
		// It treats all returned slice as string slice.
		if result, ok := result.([]interface{}); ok {
			return gvar.New(gconv.Strings(result)), err
		}
	}
	return gvar.New(result), err
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/util/gconv"
	"github.com/gomodule/redigo/redis"
	"reflect"
)

type localAdapterRedisConn struct {
	redis.ConnWithTimeout
}

// Do sends a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
// The timeout overrides the read timeout set when dialing the connection.
func (c *localAdapterRedisConn) Do(ctx context.Context, command string, args []interface{}, option *Option) (reply *gvar.Var, err error) {
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
	if option == nil {
		option = &Option{}
	}
	//timestampMilli1 := gtime.TimestampMilli()
	reply, err = resultToVar(c.ConnWithTimeout.DoWithTimeout(option.ReadTimeout, command, args...))
	//timestampMilli2 := gtime.TimestampMilli()

	// Tracing.
	//c.addTracingItem(&tracingItem{
	//	err:         err,
	//	commandName: commandName,
	//	arguments:   args,
	//	costMilli:   timestampMilli2 - timestampMilli1,
	//})
	return
}

// Receive receives a single reply as gvar.Var from the Redis server.
func (c *localAdapterRedisConn) Receive(ctx context.Context, option *Option) (*gvar.Var, error) {
	if option == nil {
		option = defaultOption()
	}
	return resultToVar(c.ReceiveWithTimeout(option.ReadTimeout))
}

func (c *localAdapterRedisConn) Close(ctx context.Context) error {
	return c.ConnWithTimeout.Close()
}

// resultToVar converts redis operation result to gvar.Var.
func resultToVar(result interface{}, err error) (*gvar.Var, error) {
	if err == nil {
		if result, ok := result.([]byte); ok {
			return gvar.New(string(result)), err
		}
		// It treats all returned slice as string slice.
		if result, ok := result.([]interface{}); ok {
			return gvar.New(gconv.Strings(result)), err
		}
	}
	return gvar.New(result), err
}

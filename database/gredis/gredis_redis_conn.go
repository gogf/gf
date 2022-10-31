// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// RedisConn is a connection of redis client.
type RedisConn struct {
	conn  Conn
	redis *Redis
}

// traceItem holds the information for redis trace.
type traceItem struct {
	err       error
	command   string
	args      []interface{}
	costMilli int64
}

const (
	traceInstrumentName               = "github.com/gogf/gf/v2/database/gredis"
	traceAttrRedisAddress             = "redis.address"
	traceAttrRedisDb                  = "redis.db"
	traceEventRedisExecution          = "redis.execution"
	traceEventRedisExecutionCommand   = "redis.execution.command"
	traceEventRedisExecutionCost      = "redis.execution.cost"
	traceEventRedisExecutionArguments = "redis.execution.arguments"
)

// Do send a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (c *RedisConn) Do(ctx context.Context, command string, args ...interface{}) (reply *gvar.Var, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for k, v := range args {
		var (
			reflectInfo = reflection.OriginTypeAndKind(v)
		)
		switch reflectInfo.OriginKind {
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

	// Trace span start.
	tr := otel.GetTracerProvider().Tracer(traceInstrumentName, trace.WithInstrumentationVersion(gf.VERSION))
	_, span := tr.Start(ctx, "Redis."+command, trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	timestampMilli1 := gtime.TimestampMilli()
	reply, err = c.conn.Do(ctx, command, args...)
	timestampMilli2 := gtime.TimestampMilli()

	// Trace span end.
	c.traceSpanEnd(ctx, span, &traceItem{
		err:       err,
		command:   command,
		args:      args,
		costMilli: timestampMilli2 - timestampMilli1,
	})
	return
}

// Subscribe subscribes the client to the specified channels.
//
// https://redis.io/commands/subscribe/
func (c *RedisConn) Subscribe(ctx context.Context, channel string, channels ...string) ([]*Subscription, error) {
	args := append([]interface{}{channel}, gconv.Interfaces(channels)...)
	_, err := c.Do(ctx, "Subscribe", args...)
	if err != nil {
		return nil, err
	}
	subs := make([]*Subscription, len(args))
	for i := 0; i < len(subs); i++ {
		v, err := c.Receive(ctx)
		if err != nil {
			return nil, err
		}
		subs[i] = v.Val().(*Subscription)
	}
	return subs, err
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
func (c *RedisConn) PSubscribe(ctx context.Context, pattern string, patterns ...string) ([]*Subscription, error) {
	args := append([]interface{}{pattern}, gconv.Interfaces(patterns)...)
	_, err := c.Do(ctx, "PSubscribe", args...)
	if err != nil {
		return nil, err
	}
	subs := make([]*Subscription, len(args))
	for i := 0; i < len(subs); i++ {
		v, err := c.Receive(ctx)
		if err != nil {
			return nil, err
		}
		subs[i] = v.Val().(*Subscription)
	}
	return subs, err
}

// ReceiveMessage receives a single message of subscription from the Redis server.
func (c *RedisConn) ReceiveMessage(ctx context.Context) (*Message, error) {
	v, err := c.conn.Receive(ctx)
	if err != nil {
		return nil, err
	}
	return v.Val().(*Message), nil
}

// Receive receives a single reply as gvar.Var from the Redis server.
func (c *RedisConn) Receive(ctx context.Context) (*gvar.Var, error) {
	return c.conn.Receive(ctx)
}

// Close puts the connection back to connection pool.
func (c *RedisConn) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

// traceSpanEnd checks and adds redis trace information to OpenTelemetry.
func (c *RedisConn) traceSpanEnd(ctx context.Context, span trace.Span, item *traceItem) {
	if gtrace.IsUsingDefaultProvider() || !gtrace.IsTracingInternal() {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if item.err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, item.err))
	}

	span.SetAttributes(gtrace.CommonLabels()...)

	span.SetAttributes(
		attribute.String(traceAttrRedisAddress, c.redis.config.Address),
		attribute.Int(traceAttrRedisDb, c.redis.config.Db),
	)

	jsonBytes, _ := json.Marshal(item.args)
	span.AddEvent(traceEventRedisExecution, trace.WithAttributes(
		attribute.String(traceEventRedisExecutionCommand, item.command),
		attribute.String(traceEventRedisExecutionCost, fmt.Sprintf(`%d ms`, item.costMilli)),
		attribute.String(traceEventRedisExecutionArguments, string(jsonBytes)),
	))
}

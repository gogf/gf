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
func (c *RedisConn) Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	v, err := c.Do(ctx, "PUBLISH", channel, message)
	return v.Int64(), err
}

// Subscribe subscribes the client to the specified channels.
//
// https://redis.io/commands/subscribe/
func (c *RedisConn) Subscribe(ctx context.Context, channels ...string) error {
	_, err := c.Do(ctx, "SUBSCRIBE", channels)
	return err
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
func (c *RedisConn) PSubscribe(ctx context.Context, pattern string, patterns ...string) error {
	var s = []interface{}{pattern}
	s = append(s, gconv.Interfaces(patterns)...)
	_, err := c.Do(ctx, "PSUBSCRIBE", s...)
	return err
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

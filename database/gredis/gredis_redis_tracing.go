// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"fmt"
	"github.com/gogf/gf"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/net/gtrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// tracingItem holds the information for redis tracing.
type tracingItem struct {
	err       error
	command   string
	args      []interface{}
	costMilli int64
}

const (
	tracingInstrumentName               = "github.com/gogf/gf/database/gredis"
	tracingAttrRedisAddress             = "redis.address"
	tracingAttrRedisDb                  = "redis.db"
	tracingEventRedisExecution          = "redis.execution"
	tracingEventRedisExecutionCommand   = "redis.execution.command"
	tracingEventRedisExecutionCost      = "redis.execution.cost"
	tracingEventRedisExecutionArguments = "redis.execution.arguments"
)

// addTracingItem checks and adds redis tracing information to OpenTelemetry.
func (c *RedisConn) addTracingItem(ctx context.Context, item *tracingItem) {
	if !gtrace.IsTracingInternal() || !gtrace.IsActivated(ctx) {
		return
	}
	tr := otel.GetTracerProvider().Tracer(
		tracingInstrumentName,
		trace.WithInstrumentationVersion(gf.VERSION),
	)
	if ctx == nil {
		ctx = context.Background()
	}
	_, span := tr.Start(ctx, "Redis."+item.command, trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()
	if item.err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, item.err))
	}

	span.SetAttributes(gtrace.CommonLabels()...)

	if adapter, ok := c.redis.GetAdapter().(*AdapterGoRedis); ok {
		span.SetAttributes(
			attribute.String(tracingAttrRedisAddress, adapter.config.Address),
			attribute.Int(tracingAttrRedisDb, adapter.config.Db),
		)
	}

	jsonBytes, _ := json.Marshal(item.args)
	span.AddEvent(tracingEventRedisExecution, trace.WithAttributes(
		attribute.String(tracingEventRedisExecutionCommand, item.command),
		attribute.String(tracingEventRedisExecutionCost, fmt.Sprintf(`%d ms`, item.costMilli)),
		attribute.String(tracingEventRedisExecutionArguments, string(jsonBytes)),
	))
}

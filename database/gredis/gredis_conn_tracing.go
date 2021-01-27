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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

// tracingItem holds the information for redis tracing.
type tracingItem struct {
	err         error
	commandName string
	arguments   []interface{}
	costMilli   int64
}

const (
	tracingAttrRedisHost                = "redis.host"
	tracingAttrRedisPort                = "redis.port"
	tracingAttrRedisDb                  = "redis.db"
	tracingEventRedisExecution          = "redis.execution"
	tracingEventRedisExecutionCommand   = "redis.execution.command"
	tracingEventRedisExecutionCost      = "redis.execution.cost"
	tracingEventRedisExecutionArguments = "redis.execution.arguments"
)

// addTracingItem checks and adds redis tracing information to OpenTelemetry.
func (c *Conn) addTracingItem(item *tracingItem) {
	tr := otel.GetTracerProvider().Tracer(
		"github.com/gogf/gf/database/gredis",
		trace.WithInstrumentationVersion(fmt.Sprintf(`%s`, gf.VERSION)),
	)
	ctx := c.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	_, span := tr.Start(ctx, "Redis."+item.commandName, trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()
	if item.err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, item.err))
	}
	span.SetAttributes(gtrace.CommonLabels()...)
	span.SetAttributes(
		label.String(tracingAttrRedisHost, c.redis.config.Host),
		label.Int(tracingAttrRedisPort, c.redis.config.Port),
		label.Int(tracingAttrRedisDb, c.redis.config.Db),
	)
	jsonBytes, _ := json.Marshal(item.arguments)
	span.AddEvent(tracingEventRedisExecution, trace.WithAttributes(
		label.String(tracingEventRedisExecutionCommand, item.commandName),
		label.String(tracingEventRedisExecutionCost, fmt.Sprintf(`%d ms`, item.costMilli)),
		label.String(tracingEventRedisExecutionArguments, string(jsonBytes)),
	))
}

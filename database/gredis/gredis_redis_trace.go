// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/net/gtrace"
)

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

// traceSpanEnd checks and adds redis trace information to OpenTelemetry.
func (c *RedisConn) traceSpanEnd(ctx context.Context, span trace.Span, item *traceItem) {
	if gtrace.IsUsingDefaultProvider() || !gtrace.IsTracingInternal() {
		return
	}

	if item.err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, item.err))
	}

	span.SetAttributes(gtrace.CommonLabels()...)

	if adapter, ok := c.redis.GetAdapter().(*AdapterGoRedis); ok {
		span.SetAttributes(
			attribute.String(traceAttrRedisAddress, adapter.config.Address),
			attribute.Int(traceAttrRedisDb, adapter.config.Db),
		)
	}

	jsonBytes, _ := json.Marshal(item.args)
	span.AddEvent(traceEventRedisExecution, trace.WithAttributes(
		attribute.String(traceEventRedisExecutionCommand, item.command),
		attribute.String(traceEventRedisExecutionCost, fmt.Sprintf(`%d ms`, item.costMilli)),
		attribute.String(traceEventRedisExecutionArguments, string(jsonBytes)),
	))
}

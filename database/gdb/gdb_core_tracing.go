// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gdb

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/net/gtrace"
)

const (
	tracingInstrumentName       = "github.com/gogf/gf/v2/database/gdb"
	tracingAttrDbType           = "db.type"
	tracingAttrDbHost           = "db.host"
	tracingAttrDbPort           = "db.port"
	tracingAttrDbName           = "db.name"
	tracingAttrDbUser           = "db.user"
	tracingAttrDbLink           = "db.link"
	tracingAttrDbGroup          = "db.group"
	tracingEventDbExecution     = "db.execution"
	tracingEventDbExecutionSql  = "db.execution.sql"
	tracingEventDbExecutionCost = "db.execution.cost"
	tracingEventDbExecutionRows = "db.execution.rows"
	tracingEventDbExecutionTxID = "db.execution.txid"
	tracingEventDbExecutionType = "db.execution.type"
)

// addSqlToTracing adds sql information to tracer if it's enabled.
func (c *Core) addSqlToTracing(ctx context.Context, sql *Sql) {
	if gtrace.IsUsingDefaultProvider() || !gtrace.IsTracingInternal() {
		return
	}
	tr := otel.GetTracerProvider().Tracer(
		tracingInstrumentName,
		trace.WithInstrumentationVersion(gf.VERSION),
	)
	ctx, span := tr.Start(ctx, sql.Type, trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	if sql.Error != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, sql.Error))
	}
	labels := make([]attribute.KeyValue, 0)
	labels = append(labels, gtrace.CommonLabels()...)
	labels = append(labels,
		attribute.String(tracingAttrDbType, c.db.GetConfig().Type),
		semconv.DBStatementKey.String(sql.Format),
	)
	if c.db.GetConfig().Host != "" {
		labels = append(labels, attribute.String(tracingAttrDbHost, c.db.GetConfig().Host))
	}
	if c.db.GetConfig().Port != "" {
		labels = append(labels, attribute.String(tracingAttrDbPort, c.db.GetConfig().Port))
	}
	if c.db.GetConfig().Name != "" {
		labels = append(labels, attribute.String(tracingAttrDbName, c.db.GetConfig().Name))
	}
	if c.db.GetConfig().User != "" {
		labels = append(labels, attribute.String(tracingAttrDbUser, c.db.GetConfig().User))
	}
	if filteredLink := c.db.FilteredLink(); filteredLink != "" {
		labels = append(labels, attribute.String(tracingAttrDbLink, c.db.FilteredLink()))
	}
	if group := c.db.GetGroup(); group != "" {
		labels = append(labels, attribute.String(tracingAttrDbGroup, group))
	}
	span.SetAttributes(labels...)
	events := []attribute.KeyValue{
		attribute.String(tracingEventDbExecutionSql, sql.Format),
		attribute.String(tracingEventDbExecutionCost, fmt.Sprintf(`%d ms`, sql.End-sql.Start)),
		attribute.String(tracingEventDbExecutionRows, fmt.Sprintf(`%d`, sql.RowsAffected)),
	}
	if sql.IsTransaction {
		if v := ctx.Value(transactionIdForLoggerCtx); v != nil {
			events = append(events, attribute.String(
				tracingEventDbExecutionTxID, fmt.Sprintf(`%d`, v.(uint64)),
			))
		}
	}
	events = append(events, attribute.String(tracingEventDbExecutionType, sql.Type))
	span.AddEvent(tracingEventDbExecution, trace.WithAttributes(events...))
}

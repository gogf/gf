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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2/net/gtrace"
)

const (
	traceInstrumentName       = "github.com/gogf/gf/v2/database/gdb"
	traceAttrDbType           = "db.type"
	traceAttrDbHost           = "db.host"
	traceAttrDbPort           = "db.port"
	traceAttrDbName           = "db.name"
	traceAttrDbUser           = "db.user"
	traceAttrDbLink           = "db.link"
	traceAttrDbGroup          = "db.group"
	traceEventDbExecution     = "db.execution"
	traceEventDbExecutionSql  = "db.execution.sql"
	traceEventDbExecutionCost = "db.execution.cost"
	traceEventDbExecutionRows = "db.execution.rows"
	traceEventDbExecutionTxID = "db.execution.txid"
	traceEventDbExecutionType = "db.execution.type"
)

// addSqlToTracing adds sql information to tracer if it's enabled.
func (c *Core) traceSpanEnd(ctx context.Context, span trace.Span, sql *Sql) {
	if gtrace.IsUsingDefaultProvider() || !gtrace.IsTracingInternal() {
		return
	}
	if sql.Error != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, sql.Error))
	}
	labels := make([]attribute.KeyValue, 0)
	labels = append(labels, gtrace.CommonLabels()...)
	labels = append(labels,
		attribute.String(traceAttrDbType, c.db.GetConfig().Type),
		semconv.DBStatementKey.String(sql.Format),
	)
	if c.db.GetConfig().Host != "" {
		labels = append(labels, attribute.String(traceAttrDbHost, c.db.GetConfig().Host))
	}
	if c.db.GetConfig().Port != "" {
		labels = append(labels, attribute.String(traceAttrDbPort, c.db.GetConfig().Port))
	}
	if c.db.GetConfig().Name != "" {
		labels = append(labels, attribute.String(traceAttrDbName, c.db.GetConfig().Name))
	}
	if c.db.GetConfig().User != "" {
		labels = append(labels, attribute.String(traceAttrDbUser, c.db.GetConfig().User))
	}
	if filteredLink := c.db.FilteredLink(); filteredLink != "" {
		labels = append(labels, attribute.String(traceAttrDbLink, c.db.FilteredLink()))
	}
	if group := c.db.GetGroup(); group != "" {
		labels = append(labels, attribute.String(traceAttrDbGroup, group))
	}
	span.SetAttributes(labels...)
	events := []attribute.KeyValue{
		attribute.String(traceEventDbExecutionSql, sql.Format),
		attribute.String(traceEventDbExecutionCost, fmt.Sprintf(`%d ms`, sql.End-sql.Start)),
		attribute.String(traceEventDbExecutionRows, fmt.Sprintf(`%d`, sql.RowsAffected)),
	}
	if sql.IsTransaction {
		if v := ctx.Value(transactionIdForLoggerCtx); v != nil {
			events = append(events, attribute.String(
				traceEventDbExecutionTxID, fmt.Sprintf(`%d`, v.(uint64)),
			))
		}
	}
	events = append(events, attribute.String(traceEventDbExecutionType, sql.Type))
	span.AddEvent(traceEventDbExecution, trace.WithAttributes(events...))
}

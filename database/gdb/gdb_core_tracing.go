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
	"github.com/gogf/gf"
	"github.com/gogf/gf/net/gtrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracingInstrumentName       = "github.com/gogf/gf/database/gdb"
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
	tracingEventDbExecutionType = "db.execution.type"
)

// addSqlToTracing adds sql information to tracer if it's enabled.
func (c *Core) addSqlToTracing(ctx context.Context, sql *Sql) {
	if !gtrace.IsActivated(ctx) {
		return
	}
	tr := otel.GetTracerProvider().Tracer(tracingInstrumentName, trace.WithInstrumentationVersion(gf.VERSION))
	ctx, span := tr.Start(ctx, sql.Type, trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	if sql.Error != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, sql.Error))
	}
	labels := make([]label.KeyValue, 0)
	labels = append(labels, gtrace.CommonLabels()...)
	labels = append(labels,
		label.String(tracingAttrDbType, c.DB.GetConfig().Type),
	)
	if c.DB.GetConfig().Host != "" {
		labels = append(labels, label.String(tracingAttrDbHost, c.DB.GetConfig().Host))
	}
	if c.DB.GetConfig().Port != "" {
		labels = append(labels, label.String(tracingAttrDbPort, c.DB.GetConfig().Port))
	}
	if c.DB.GetConfig().Name != "" {
		labels = append(labels, label.String(tracingAttrDbName, c.DB.GetConfig().Name))
	}
	if c.DB.GetConfig().User != "" {
		labels = append(labels, label.String(tracingAttrDbUser, c.DB.GetConfig().User))
	}
	if filteredLinkInfo := c.DB.FilteredLinkInfo(); filteredLinkInfo != "" {
		labels = append(labels, label.String(tracingAttrDbLink, c.DB.FilteredLinkInfo()))
	}
	if group := c.DB.GetGroup(); group != "" {
		labels = append(labels, label.String(tracingAttrDbGroup, group))
	}
	span.SetAttributes(labels...)
	span.AddEvent(tracingEventDbExecution, trace.WithAttributes(
		label.String(tracingEventDbExecutionSql, sql.Format),
		label.String(tracingEventDbExecutionCost, fmt.Sprintf(`%d ms`, sql.End-sql.Start)),
		label.String(tracingEventDbExecutionType, sql.Type),
	))
}

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
	"github.com/gogf/gf/os/gcmd"
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

var (
	// tracingInternal enables tracing for internal type spans.
	// It's true in default.
	tracingInternal = true
)

func init() {
	tracingInternal = gcmd.GetOptWithEnv("gf.tracing.internal", true).Bool()
}

// addSqlToTracing adds sql information to tracer if it's enabled.
func (c *Core) addSqlToTracing(ctx context.Context, sql *Sql) {
	if !tracingInternal || !gtrace.IsActivated(ctx) {
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
	labels := make([]label.KeyValue, 0)
	labels = append(labels, gtrace.CommonLabels()...)
	labels = append(labels,
		label.String(tracingAttrDbType, c.db.GetConfig().Type),
	)
	if c.db.GetConfig().Host != "" {
		labels = append(labels, label.String(tracingAttrDbHost, c.db.GetConfig().Host))
	}
	if c.db.GetConfig().Port != "" {
		labels = append(labels, label.String(tracingAttrDbPort, c.db.GetConfig().Port))
	}
	if c.db.GetConfig().Name != "" {
		labels = append(labels, label.String(tracingAttrDbName, c.db.GetConfig().Name))
	}
	if c.db.GetConfig().User != "" {
		labels = append(labels, label.String(tracingAttrDbUser, c.db.GetConfig().User))
	}
	if filteredLinkInfo := c.db.FilteredLinkInfo(); filteredLinkInfo != "" {
		labels = append(labels, label.String(tracingAttrDbLink, c.db.FilteredLinkInfo()))
	}
	if group := c.db.GetGroup(); group != "" {
		labels = append(labels, label.String(tracingAttrDbGroup, group))
	}
	span.SetAttributes(labels...)
	span.AddEvent(tracingEventDbExecution, trace.WithAttributes(
		label.String(tracingEventDbExecutionSql, sql.Format),
		label.String(tracingEventDbExecutionCost, fmt.Sprintf(`%d ms`, sql.End-sql.Start)),
		label.String(tracingEventDbExecutionType, sql.Type),
	))
}

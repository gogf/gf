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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

// addSqlToTracing adds sql information to tracer if it's enabled.
func (c *Core) addSqlToTracing(ctx context.Context, sql *Sql) {
	if ctx == nil {
		return
	}
	spanCtx := trace.SpanContextFromContext(ctx)
	if traceId := spanCtx.TraceID; !traceId.IsValid() {
		return
	}

	tr := otel.GetTracerProvider().Tracer(
		"github.com/gogf/gf/database/gdb",
		trace.WithInstrumentationVersion(fmt.Sprintf(`%s`, gf.VERSION)),
	)
	ctx, span := tr.Start(ctx, sql.Type)
	defer span.End()
	if sql.Error != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, sql.Error))
	}
	labels := make([]label.KeyValue, 0)
	labels = append(labels, label.String("db.type", c.DB.GetConfig().Type))
	if c.DB.GetConfig().Host != "" {
		labels = append(labels, label.String("db.host", c.DB.GetConfig().Host))
	}
	if c.DB.GetConfig().Port != "" {
		labels = append(labels, label.String("db.port", c.DB.GetConfig().Port))
	}
	if c.DB.GetConfig().Name != "" {
		labels = append(labels, label.String("db.name", c.DB.GetConfig().Name))
	}
	if c.DB.GetConfig().User != "" {
		labels = append(labels, label.String("db.user", c.DB.GetConfig().User))
	}
	if filteredLinkInfo := c.DB.FilteredLinkInfo(); filteredLinkInfo != "" {
		labels = append(labels, label.String("db.link", c.DB.FilteredLinkInfo()))
	}
	if group := c.DB.GetGroup(); group != "" {
		labels = append(labels, label.String("db.group", group))
	}
	span.SetAttributes(labels...)
	span.AddEvent("db.execution", trace.WithAttributes(
		label.String(`db.execution.sql`, sql.Format),
		label.String(`db.execution.cost`, fmt.Sprintf(`%d ms`, sql.End-sql.Start)),
		label.String(`db.execution.type`, sql.Type),
	))
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtrace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// Span 。
type Span struct {
	trace.Span
}

// NewSpan creates a span using default tracer.
func NewSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, *Span) {
	ctx, span := NewTracer().Start(ctx, spanName, opts...)
	return ctx, &Span{
		Span: span,
	}
}

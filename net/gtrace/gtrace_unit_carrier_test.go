// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtrace_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/net/gtrace"
	"github.com/gogf/gf/test/gtest"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	traceIDStr = "4bf92f3577b34da6a3ce929d0e0e4736"
	spanIDStr  = "00f067aa0ba902b7"
)

var (
	traceID = mustTraceIDFromHex(traceIDStr)
	spanID  = mustSpanIDFromHex(spanIDStr)
)

func mustTraceIDFromHex(s string) (t trace.TraceID) {
	var err error
	t, err = trace.TraceIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustSpanIDFromHex(s string) (t trace.SpanID) {
	var err error
	t, err = trace.SpanIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func TestNewCarrier(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		}))

		ctx, _ = otel.Tracer("").Start(ctx, "inject")
		carrier1 := gtrace.NewCarrier()
		otel.GetTextMapPropagator().Inject(ctx, carrier1)
		t.Assert(carrier1.String(), `{"traceparent":"00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"}`)

		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier1)
		gotSc := trace.SpanContextFromContext(ctx)
		t.Assert(gotSc.TraceID().String(), traceID.String())
		t.Assert(gotSc.SpanID().String(), "00f067aa0ba902b7")
	})
}

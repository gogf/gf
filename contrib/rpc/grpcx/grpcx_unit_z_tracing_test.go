// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"

	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2/internal/tracing"
)

// setupPropagator installs a TraceContext propagator so Inject/Extract
// actually write and read the traceparent header.
func setupPropagator() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

// Test_Tracing_Inject_Extract verifies that a trace context injected into
// gRPC metadata can be extracted back as a valid SpanContext.
func Test_Tracing_Inject_Extract(t *testing.T) {
	setupPropagator()
	gtest.C(t, func(t *gtest.T) {
		// Create a context that carries a valid span.
		tracer := otel.Tracer("test")
		ctx, span := tracer.Start(context.Background(), "test-span")
		defer span.End()

		md := metadata.New(nil)
		tracing.Inject(ctx, md)

		// traceparent must be present after Inject.
		vals := md.Get("traceparent")
		t.Assert(len(vals) > 0, true)

		// Extract must recover a valid SpanContext.
		_, sc := tracing.Extract(context.Background(), md)
		t.Assert(sc.IsValid(), true)
		t.Assert(sc.TraceID(), span.SpanContext().TraceID())
	})
}

// Test_Tracing_Extract_Empty verifies that Extract on empty metadata returns
// an invalid SpanContext without panicking.
func Test_Tracing_Extract_Empty(t *testing.T) {
	setupPropagator()
	gtest.C(t, func(t *gtest.T) {
		md := metadata.New(nil)
		_, sc := tracing.Extract(context.Background(), md)
		t.Assert(sc.IsValid(), false)
	})
}

// Test_Tracing_Inject_NoPropagator verifies that Inject on a context with no
// active span does not write a traceparent key.
func Test_Tracing_Inject_NoPropagator(t *testing.T) {
	setupPropagator()
	gtest.C(t, func(t *gtest.T) {
		md := metadata.New(nil)
		// context.Background() carries no span.
		tracing.Inject(context.Background(), md)
		vals := md.Get("traceparent")
		// No active span → no traceparent propagated.
		t.Assert(len(vals) == 0, true)
	})
}

// Test_Tracing_Inject_Idempotent verifies that calling Inject twice on the
// same span context yields the same traceparent value.
func Test_Tracing_Inject_Idempotent(t *testing.T) {
	setupPropagator()
	gtest.C(t, func(t *gtest.T) {
		tracer := otel.Tracer("test")
		ctx, span := tracer.Start(context.Background(), "idempotent")
		defer span.End()

		md1 := metadata.New(nil)
		md2 := metadata.New(nil)
		tracing.Inject(ctx, md1)
		tracing.Inject(ctx, md2)
		t.Assert(md1.Get("traceparent"), md2.Get("traceparent"))
	})
}

// Test_Tracing_Extract_PreservesBaggage verifies that baggage encoded by
// Inject survives an Extract round-trip.
func Test_Tracing_Extract_PreservesBaggage(t *testing.T) {
	setupPropagator()
	gtest.C(t, func(t *gtest.T) {
		tracer := otel.Tracer("test")
		ctx, span := tracer.Start(context.Background(), "baggage-test")
		defer span.End()

		member, err := baggage.NewMember("test-key", "test-value")
		t.AssertNil(err)
		bag, err := baggage.New(member)
		t.AssertNil(err)
		ctx = baggage.ContextWithBaggage(ctx, bag)

		md := metadata.New(nil)
		tracing.Inject(ctx, md)

		extracted, sc := tracing.Extract(context.Background(), md)
		t.Assert(sc.IsValid(), true)
		t.Assert(sc.TraceID(), span.SpanContext().TraceID())
		t.Assert(extracted.Member("test-key").Value(), "test-value")
	})
}

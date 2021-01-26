// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gtrace provides convenience wrapping functionality for tracing feature using OpenTelemetry.
package gtrace

import (
	"context"
	"github.com/gogf/gf/container/gvar"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

// Tracer is a short function for retrieve Tracer.
func Tracer(name ...string) trace.Tracer {
	tracerName := ""
	if len(name) > 0 {
		tracerName = name[0]
	}
	return otel.Tracer(tracerName)
}

// GetTraceId retrieves and returns TraceId from context.
func GetTraceId(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).TraceID.String()
}

// GetSpanId retrieves and returns SpanId from context.
func GetSpanId(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).SpanID.String()
}

// SetBaggageValue is a convenient function for adding one key-value pair to baggage.
// Note that it uses label.Any to set the key-value pair.
func SetBaggageValue(ctx context.Context, key string, value interface{}) context.Context {
	return baggage.ContextWithValues(ctx, label.Any(key, value))
}

// SetBaggageMap is a convenient function for adding map key-value pairs to baggage.
// Note that it uses label.Any to set the key-value pair.
func SetBaggageMap(ctx context.Context, data map[string]interface{}) context.Context {
	pairs := make([]label.KeyValue, 0)
	for k, v := range data {
		pairs = append(pairs, label.Any(k, v))
	}
	return baggage.ContextWithValues(ctx, pairs...)
}

// GetBaggageVar retrieves value and returns a *gvar.Var for specified key from baggage.
func GetBaggageVar(ctx context.Context, key string) *gvar.Var {
	value := baggage.Value(ctx, label.Key(key))
	return gvar.New(value.AsInterface())
}

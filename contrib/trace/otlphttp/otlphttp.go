// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package otlphttp provides gtrace.Tracer implementation using OpenTelemetry protocol.
package otlphttp

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gipv4"
)

const (
	tracerHostnameTagKey = "hostname"
)

// Init initializes and registers `otlphttp` to global TracerProvider.
//
// The output parameter `Shutdown` is used for waiting exported trace spans to be uploaded,
// which is useful if your program is ending, and you do not want to lose recent spans.
func Init(serviceName, endpoint, path string) (func(ctx context.Context), error) {
	// Try retrieving host ip for tracing info.
	var (
		intranetIPArray, err = gipv4.GetIntranetIpArray()
		hostIP               = "NoHostIpFound"
	)

	if err != nil {
		return nil, err
	}

	if len(intranetIPArray) == 0 {
		if intranetIPArray, err = gipv4.GetIpArray(); err != nil {
			return nil, err
		}
	}
	if len(intranetIPArray) > 0 {
		hostIP = intranetIPArray[0]
	}

	ctx := context.Background()
	traceExp, err := otlptrace.New(ctx, otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithURLPath(path),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithCompression(1),
	))
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// The name of the service displayed on the traceback end。
			semconv.ServiceNameKey.String(serviceName),
			semconv.HostNameKey.String(hostIP),
			attribute.String(tracerHostnameTagKey, hostIP),
		),
	)

	tracerProvider := trace.NewTracerProvider(
		// AlwaysSample is a sampler that samples every trace.
		// see: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#AlwaysSample
		// example see: [example/trace/provider/http/main.go](../../../../../example/trace/provider/http/main.go#L84)
		trace.WithSampler(trace.AlwaysSample()),
		// WithResource returns a trace option that sets the resource to be associated with spans.
		// see: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#WithResource
		// example see: [example/trace/provider/http/main.go](../../../../../example/trace/provider/http/main.go#L33)
		trace.WithResource(res),
		// WithSpanProcessor returns a trace option that sets the span processor to be used by the trace provider.
		// see: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#WithSpanProcessor
		// example see: [example/trace/provider/http/main.go](../../../../../example/trace/provider/http/main.go#L96)
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(traceExp)),
	)

	// Set the global propagator to traceContext (not set by default).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err = tracerProvider.Shutdown(ctx); err != nil {
			g.Log().Errorf(ctx, "Shutdown tracerProvider failed err:%+v", err)
		} else {
			g.Log().Debug(ctx, "Shutdown tracerProvider success")
		}
	}, nil
}

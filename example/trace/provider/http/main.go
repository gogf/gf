// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/gogf/gf/example/trace/provider/internal"
)

func main() {
	var (
		serverIP, err = internal.GetLocalIP()
		ctx           = gctx.New()
	)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	var res *resource.Resource
	if res, err = resource.New(ctx,
		// WithFromEnv returns a resource option that sets the resource attributes from the environment.
		resource.WithFromEnv(),
		// WithProcess returns a resource option that sets the process attributes.
		resource.WithProcess(),
		// WithTelemetrySDK returns a resource option that sets the telemetry SDK attributes.
		resource.WithTelemetrySDK(),
		// WithHost returns a resource option that sets the host attributes.
		resource.WithHost(),
		// WithAttributes returns a resource option that sets the resource attributes.
		resource.WithAttributes(
			// The name of the service displayed on the traceback end。
			semconv.ServiceNameKey.String(internal.HTTPServiceName),
			// The IP address of the server.
			semconv.HostNameKey.String(serverIP),
			// The IP address of the server.
			attribute.String(internal.TracerHostnameTagKey, serverIP),
		),
		// WithOS returns a resource option that sets the OS attributes.
		resource.WithOS(),
		// WithProcessPID returns a resource option that sets the process PID attribute.
		resource.WithProcessPID(),
		// For more parameters, please customize the selection
	); err != nil {
		g.Log().Fatal(ctx, err)
	}

	var exporter *otlptrace.Exporter
	if exporter, err = otlptrace.New(ctx, otlptracehttp.NewClient(
		// WithEndpoint returns an otlptracehttp.Option that sets the endpoint to which the exporter is going to send the spans.
		otlptracehttp.WithEndpoint(internal.HTTPEndpoint),
		// WithHeaders returns an otlptracehttp.Option that sets the headers to be sent with HTTP requests.
		otlptracehttp.WithURLPath(internal.HTTPPath),
		// WithInsecure returns an otlptracehttp.Option that disables secure connection to the collector.
		otlptracehttp.WithInsecure(),
		// WithCompression returns an otlptracehttp.Option that sets the compression level for the exporter.
		otlptracehttp.WithCompression(1))); err != nil {
		g.Log().Fatal(ctx, err)
	}
	var shutdown func(ctx context.Context)
	// WithSampler sets the sampler for the trace provider.
	//  1. AlwaysSample: AlwaysSample is a sampler that samples every trace.
	//  2. NeverSample: NeverSample is a sampler that samples no traces.
	//  3. ParentBased: ParentBased is a sampler that samples a trace based on the parent span.
	//  4. TraceIDRatioBased: TraceIDRatioBased is a sampler that samples a trace based on the TraceID.
	// WithResource sets the resource for the trace provider.
	// WithSpanProcessor sets the span processor for the trace provider.
	//  1. NewSimpleSpanProcessor: NewSimpleSpanProcessor returns a new SimpleSpanProcessor.
	//  2. NewBatchSpanProcessor: NewBatchSpanProcessor returns a new BatchSpanProcessor.
	// WithRawSpanLimits sets the raw span limits for the trace provider.
	if shutdown, err = internal.InitTracer(
		// WithSampler returns a trace option that sets the sampler for the trace provider.
		// trace.WithSampler(trace.AlwaysSample()),
		// trace.WithSampler(trace.NeverSample()),
		// trace.WithSampler(trace.ParentBased(trace.AlwaysSample())),
		// WithSampler returns a trace option that sets the sampler for the trace provider.
		//  1. AlwaysSample: AlwaysSample is a sampler that samples every trace.
		//  2. NeverSample: NeverSample is a sampler that samples no traces.
		//  3. ParentBased: ParentBased is a sampler that samples a trace based on the parent span.
		//  4. TraceIDRatioBased: TraceIDRatioBased is a sampler that samples a trace based on the TraceID.
		trace.WithSampler(trace.TraceIDRatioBased(0.1)),
		// WithResource returns a trace option that sets the resource for the trace provider.
		trace.WithResource(res),
		// WithSpanProcessor returns a trace option that sets the span processor for the trace provider.
		// trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exporter)),
		// trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
		//  1. NewSimpleSpanProcessor: NewSimpleSpanProcessor returns a new SimpleSpanProcessor.
		//  2. NewBatchSpanProcessor: NewBatchSpanProcessor returns a new BatchSpanProcessor.
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
		// WithRawSpanLimits returns a trace option that sets the raw span limits for the trace provider.
		trace.WithRawSpanLimits(trace.NewSpanLimits()),
	); err != nil {
		g.Log().Fatal(ctx, err)
	}
	defer shutdown(ctx)

	internal.StartRequests()
}

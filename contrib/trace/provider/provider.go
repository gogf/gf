// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package provider

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc/encoding/gzip"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gipv4"
)

const (
	tracerHostnameTagKey = "hostname"
)

// InitTracer initializes and registers `otlpgrpc` to global TracerProvider.
func InitTracer(opts ...trace.TracerProviderOption) (func(ctx context.Context), error) {
	tracerProvider := trace.NewTracerProvider(opts...)
	// Set the global propagator to traceContext (not set by default).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		// Shutdown waits for exported trace spans to be uploaded.
		if err := tracerProvider.Shutdown(ctx); err != nil {
			g.Log().Errorf(ctx, "Shutdown tracerProvider failed err:%+v", err)
		} else {
			g.Log().Debug(ctx, "Shutdown tracerProvider success")
		}
	}, nil
}

// GetLocalIP returns the IP address of the server.
func GetLocalIP() (string, error) {
	var intranetIPArray, err = gipv4.GetIntranetIpArray()
	if err != nil {
		return "", err
	}

	if len(intranetIPArray) == 0 {
		if intranetIPArray, err = gipv4.GetIpArray(); err != nil {
			return "", err
		}
	}
	var hostIP = "NoHostIpFound"
	if len(intranetIPArray) > 0 {
		hostIP = intranetIPArray[0]
	}
	return hostIP, nil
}

// NewGRPCDefaultTraceClient creates and returns a new OTLP trace client.
// endpoint is the endpoint to which the exporter is going to send the spans.
// traceToken is the token used for authentication.
func NewGRPCDefaultTraceClient(endpoint, traceToken string) otlptrace.Client {
	return otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint), // Replace the otel Agent Addr with the access point obtained in the prerequisite。
		otlptracegrpc.WithHeaders(map[string]string{"Authentication": traceToken}),
		otlptracegrpc.WithCompressor(gzip.Name))
}

// NewGRPCCustomTraceClient creates and returns a new OTLP trace client.
func NewGRPCCustomTraceClient(opts ...otlptracegrpc.Option) otlptrace.Client {
	return otlptracegrpc.NewClient(opts...)
}

// NewHTTPDefaultTraceClient creates and returns a new OTLP trace client.
func NewHTTPDefaultTraceClient(endpoint, path string) otlptrace.Client {
	return otlptracehttp.NewClient(otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithURLPath(path),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithCompression(1))
}

// NewHTTPCustomTraceClient creates and returns a new OTLP trace client.
func NewHTTPCustomTraceClient(opts ...otlptracehttp.Option) otlptrace.Client {
	return otlptracehttp.NewClient(opts...)
}

// NewExporter creates and returns a new OTLP trace exporter.
func NewExporter(ctx context.Context, client otlptrace.Client) (*otlptrace.Exporter, error) {
	return otlptrace.New(ctx, client)
}

// NewDefaultResource creates and returns a new resource.
// serviceName is the name of the service displayed on the traceback end.
// serverIP is the IP address of the server.
func NewDefaultResource(ctx context.Context, serviceName, serverIP string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// The name of the service displayed on the traceback end。
			semconv.ServiceNameKey.String(serviceName),
			// The IP address of the server.
			semconv.HostNameKey.String(serverIP),
			// The IP address of the server.
			attribute.String(tracerHostnameTagKey, serverIP),
		),
	)
}

// NewResource creates and returns a new resource.
func NewResource(ctx context.Context, opts ...resource.Option) (*resource.Resource, error) {
	return resource.New(ctx, opts...)
}

// NewBatchSpanProcessor returns a new SpanProcessor that will batch up completed spans and send them to the exporter in batches.
// This is the recommended SpanProcessor for production use.
// The exporter will be called in a separate goroutine.
func NewBatchSpanProcessor(exporter *otlptrace.Exporter) trace.SpanProcessor {
	return trace.NewBatchSpanProcessor(exporter)
}

// NewSimpleSpanProcessor returns a new SpanProcessor that will synchronously
// send completed spans to the exporter immediately.
// This SpanProcessor is not recommended for production use.
func NewSimpleSpanProcessor(exporter *otlptrace.Exporter) trace.SpanProcessor {
	return trace.NewSimpleSpanProcessor(exporter)
}

// NewTraceSampler creates and returns a new trace sampler.
// The sampler is used to determine if a span should be recorded.
// If the sampler is nil, it will use the default sampler AlwaysSample().
func NewTraceSampler(sampler trace.Sampler) trace.Sampler {
	if sampler == nil {
		return trace.AlwaysSample()
	}
	return sampler
}

// NewAlwaysOnSampler creates and returns a new always on sampler.
func NewAlwaysOnSampler() trace.Sampler {
	return trace.AlwaysSample()
}

// NewTraceIDRatioBasedSampler creates and returns a new trace ID ratio based sampler.
// The sampler is used to determine if a span should be recorded.
// If the sampler is nil, it will use the default sampler AlwaysSample().
// traceIDRatio is the ratio of traces to sample.
func NewTraceIDRatioBasedSampler(traceIDRatio float64) trace.Sampler {
	return trace.TraceIDRatioBased(traceIDRatio)
}

// NewDefaultTracerProviderOptions creates and returns a new slice of TracerProviderOption.
// sampler is the sampler to use.
// If the sampler is nil, it will use the default sampler AlwaysSample().
// resource is the resource to use.
// If the resource is nil, it will use the default resource.
// If the resource is not nil, it will use the resource.

func NewDefaultTracerProviderOptions(sampler trace.Sampler, resource *resource.Resource, spanProcessor trace.SpanProcessor) []trace.TracerProviderOption {
	options := make([]trace.TracerProviderOption, 0)
	if sampler != nil {
		options = append(options, trace.WithSampler(sampler))
	}
	if resource != nil {
		options = append(options, trace.WithResource(resource))
	}
	if spanProcessor != nil {
		options = append(options, trace.WithSpanProcessor(spanProcessor))
	}
	return options
}

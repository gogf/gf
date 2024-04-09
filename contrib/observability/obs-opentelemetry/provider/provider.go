// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package provider provides an OpenTelemetry provider
package provider

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	runtimemetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// OtelProvider is an interface for OpenTelemetry provider
type OtelProvider interface {
	Shutdown(ctx context.Context) error
}

type otelProvider struct {
	traceProvider *sdktrace.TracerProvider
	metricsPusher *metric.MeterProvider
}

// Shutdown stops the OpenTelemetry provider
func (p *otelProvider) Shutdown(ctx context.Context) (err error) {
	if p.traceProvider != nil {
		if err = p.traceProvider.Shutdown(ctx); err != nil {
			g.Log().Errorf(ctx, "failed to shutdown trace provider: %v", err)
		}
	}

	if p.metricsPusher != nil {
		if err = p.metricsPusher.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}

	return err
}

// NewOpenTelemetryProvider Initializes an otlp trace and metrics provider
func NewOpenTelemetryProvider(ctx context.Context, opts ...Option) OtelProvider {
	var (
		err            error
		meterProvider  *metric.MeterProvider
		tracerProvider *sdktrace.TracerProvider
		cfg            = newConfig(opts)
	)

	if !cfg.enableTracing && !cfg.enableMetrics {
		return nil
	}

	// resource
	res := newResource(ctx, cfg)

	// propagator
	otel.SetTextMapPropagator(cfg.textMapPropagator)

	// Tracing
	if cfg.enableTracing {
		// trace client
		var traceClientOpts []otlptracegrpc.Option
		if cfg.exportEndpoint != "" {
			traceClientOpts = append(traceClientOpts, otlptracegrpc.WithEndpoint(cfg.exportEndpoint))
		}
		if len(cfg.exportHeaders) > 0 {
			traceClientOpts = append(traceClientOpts, otlptracegrpc.WithHeaders(cfg.exportHeaders))
		}
		if cfg.exportInsecure {
			traceClientOpts = append(traceClientOpts, otlptracegrpc.WithInsecure())
		}
		if cfg.exportEnableCompression {
			traceClientOpts = append(traceClientOpts, otlptracegrpc.WithCompressor("gzip"))
		}

		traceClient := otlptracegrpc.NewClient(traceClientOpts...)

		// trace exporter
		var traceExp *otlptrace.Exporter
		if traceExp, err = otlptrace.New(ctx, traceClient); err != nil {
			g.Log().Fatalf(ctx, "failed to create otlp trace exporter: %s", err)
			return nil
		}

		// trace processor
		bsp := sdktrace.NewBatchSpanProcessor(traceExp)

		// trace provider
		tracerProvider = cfg.sdkTracerProvider
		if tracerProvider == nil {
			tracerProvider = sdktrace.NewTracerProvider(
				sdktrace.WithSampler(cfg.sampler),
				sdktrace.WithResource(res),
				sdktrace.WithSpanProcessor(bsp),
			)
		}

		otel.SetTracerProvider(tracerProvider)
	}

	// Metrics
	if cfg.enableMetrics {
		// prometheus only supports CumulativeTemporalitySelector

		var metricsClientOpts []otlpmetricgrpc.Option
		if cfg.exportEndpoint != "" {
			metricsClientOpts = append(metricsClientOpts, otlpmetricgrpc.WithEndpoint(cfg.exportEndpoint))
		}
		if len(cfg.exportHeaders) > 0 {
			metricsClientOpts = append(metricsClientOpts, otlpmetricgrpc.WithHeaders(cfg.exportHeaders))
		}
		if cfg.exportInsecure {
			metricsClientOpts = append(metricsClientOpts, otlpmetricgrpc.WithInsecure())
		}
		if cfg.exportEnableCompression {
			metricsClientOpts = append(metricsClientOpts, otlpmetricgrpc.WithCompressor("gzip"))
		}

		meterProvider = cfg.meterProvider
		if meterProvider == nil {
			// metrics exporter
			metricExp, err := otlpmetricgrpc.New(context.Background(), metricsClientOpts...)

			handleInitErr(err, "Failed to create the metric exporter")

			// reader := metric.NewPeriodicReader(exporter)
			reader := metric.WithReader(metric.NewPeriodicReader(metricExp, metric.WithInterval(15*time.Second)))

			meterProvider = metric.NewMeterProvider(reader, metric.WithResource(res))
		}

		// metrics pusher
		otel.SetMeterProvider(meterProvider)

		err = runtimemetrics.Start()
		handleInitErr(err, "Failed to start runtime metrics collector")
	}

	return &otelProvider{
		traceProvider: tracerProvider,
		metricsPusher: meterProvider,
	}
}

func newResource(ctx context.Context, cfg *config) *resource.Resource {
	if cfg.resource != nil {
		return cfg.resource
	}

	res, err := resource.New(
		ctx,
		resource.WithHost(),
		resource.WithFromEnv(),
		resource.WithProcessPID(),
		resource.WithTelemetrySDK(),
		resource.WithDetectors(cfg.resourceDetectors...),
		resource.WithAttributes(cfg.resourceAttributes...),
	)
	if err != nil {
		return resource.Default()
	}
	return res
}

func handleInitErr(err error, message string) {
	if err != nil {
		g.Log().Fatalf(context.Background(), "%s: %v", message, err)
	}
}

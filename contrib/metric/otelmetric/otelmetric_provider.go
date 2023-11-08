// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/gogf/gf/v2/os/gmetric"
)

type localProvider struct {
	provider *metric.MeterProvider
}

func newLocalProvider(options ...metric.Option) gmetric.Provider {
	var (
		metrics = gmetric.GetAllMetrics()
		views   = createViewsByMetrics(metrics)
	)
	options = append(options, metric.WithView(views...))
	provider := &localProvider{
		provider: metric.NewMeterProvider(options...),
	}
	initializeAllMetrics(metrics, provider)
	return provider
}

func (l *localProvider) SetAsGlobal() {
	otel.SetMeterProvider(l.provider)
}

func (l *localProvider) Meter(instrument string) gmetric.Meter {
	return newMeter(l.provider)
}

func (l *localProvider) ForceFlush(ctx context.Context) error {
	return l.provider.ForceFlush(ctx)
}

func (l *localProvider) Shutdown(ctx context.Context) error {
	return l.provider.Shutdown(ctx)
}

func createViewsByMetrics(metrics []gmetric.Metric) []metric.View {
	var views = make([]metric.View, 0)
	for _, m := range metrics {
		switch m.MetricInfo().Type() {
		case gmetric.MetricTypeCounter:
		case gmetric.MetricTypeGauge:
		case gmetric.MetricTypeHistogram:
			// Custom buckets for each Histogram.
			views = append(views, metric.NewView(
				metric.Instrument{
					Name:  m.MetricInfo().Name(),
					Scope: instrumentation.Scope{Name: m.MetricInfo().Inst()},
				},
				metric.Stream{
					Aggregation: metric.AggregationExplicitBucketHistogram{
						Boundaries: m.(gmetric.Histogram).Buckets(),
					},
				},
			))
		}
	}
	return views
}

func initializeAllMetrics(metrics []gmetric.Metric, provider gmetric.Provider) {
	for _, m := range metrics {
		if initializer, ok := m.(gmetric.Initializer); ok {
			initializer.Init(provider)
		}
	}
}

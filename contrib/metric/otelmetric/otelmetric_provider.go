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

// localProvider implements interface gmetric.Provider.
type localProvider struct {
	provider *metric.MeterProvider
}

// newLocalProvider creates and returns an object that implements gmetric.Provider.
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

// SetAsGlobal sets current provider as global meter provider for current process.
func (l *localProvider) SetAsGlobal() {
	otel.SetMeterProvider(l.provider)
}

// Meter creates and returns a Meter.
// A Meter can produce types of Metric performer.
func (l *localProvider) Meter() gmetric.Meter {
	return newMeter(l.provider)
}

// ForceFlush flushes all pending metrics.
//
// This method honors the deadline or cancellation of ctx. An appropriate
// error will be returned in these situations. There is no guaranteed that all
// metrics be flushed or all resources have been released in these situations.
func (l *localProvider) ForceFlush(ctx context.Context) error {
	return l.provider.ForceFlush(ctx)
}

// Shutdown shuts down the Provider flushing all pending metrics and
// releasing any held computational resources.
func (l *localProvider) Shutdown(ctx context.Context) error {
	return l.provider.Shutdown(ctx)
}

// createViewsByMetrics creates and returns OpenTelemetry metric.View according metric type for all metrics,
// especially the Histogram which needs a metric.View to custom buckets.
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
					Name: m.MetricInfo().Name(),
					Scope: instrumentation.Scope{
						Name:    m.MetricInfo().Instrument(),
						Version: m.MetricInfo().InstrumentVersion(),
					},
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

// initializeAllMetrics initializes all metrics in provider creating.
// The initialization replaces the underlying metric performer using noop-performer with truly performer
// that implements operations for types of metric.
func initializeAllMetrics(metrics []gmetric.Metric, provider gmetric.Provider) {
	for _, m := range metrics {
		if initializer, ok := m.(gmetric.Initializer); ok {
			initializer.Init(provider)
		}
	}
}

// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"

	"go.opentelemetry.io/otel/exporters/prometheus"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmetric"
)

var (
	meter = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
		Instrument:        "github.com/gogf/gf/example/metric/basic",
		InstrumentVersion: "v1.0",
	})
	counter = meter.MustCounter(
		"goframe.metric.demo.counter",
		gmetric.MetricOption{
			Help: "This is a simple demo for Counter usage",
			Unit: "bytes",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_attr_1", 1),
			},
		},
	)
	upDownCounter = meter.MustUpDownCounter(
		"goframe.metric.demo.updown_counter",
		gmetric.MetricOption{
			Help: "This is a simple demo for UpDownCounter usage",
			Unit: "%",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_attr_2", 2),
			},
		},
	)
	histogram = meter.MustHistogram(
		"goframe.metric.demo.histogram",
		gmetric.MetricOption{
			Help: "This is a simple demo for histogram usage",
			Unit: "ms",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_attr_3", 3),
			},
			Buckets: []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
		},
	)
	observableCounter = meter.MustObservableCounter(
		"goframe.metric.demo.observable_counter",
		gmetric.MetricOption{
			Help: "This is a simple demo for ObservableCounter usage",
			Unit: "%",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_attr_4", 4),
			},
		},
	)
	observableUpDownCounter = meter.MustObservableUpDownCounter(
		"goframe.metric.demo.observable_updown_counter",
		gmetric.MetricOption{
			Help: "This is a simple demo for ObservableUpDownCounter usage",
			Unit: "%",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_attr_5", 5),
			},
		},
	)
	observableGauge = meter.MustObservableGauge(
		"goframe.metric.demo.observable_gauge",
		gmetric.MetricOption{
			Help: "This is a simple demo for ObservableGauge usage",
			Unit: "%",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_attr_6", 6),
			},
		},
	)
)

func main() {
	var ctx = gctx.New()

	// Callback for observable metrics.
	meter.MustRegisterCallback(func(ctx context.Context, obs gmetric.Observer) error {
		obs.Observe(observableCounter, 10)
		obs.Observe(observableUpDownCounter, 20)
		obs.Observe(observableGauge, 30)
		return nil
	}, observableCounter, observableUpDownCounter, observableGauge)

	// Prometheus exporter to export metrics as Prometheus format.
	exporter, err := prometheus.New(
		prometheus.WithoutCounterSuffixes(),
		prometheus.WithoutUnits(),
	)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}

	// OpenTelemetry provider.
	provider := otelmetric.MustProvider(
		otelmetric.WithReader(exporter),
		otelmetric.WithBuiltInMetrics(),
	)
	provider.SetAsGlobal()
	defer provider.Shutdown(ctx)

	// Counter.
	counter.Inc(ctx)
	counter.Add(ctx, 10)

	// UpDownCounter.
	upDownCounter.Inc(ctx)
	upDownCounter.Add(ctx, 10)
	upDownCounter.Dec(ctx)

	// Record values for histogram.
	histogram.Record(1)
	histogram.Record(20)
	histogram.Record(30)
	histogram.Record(101)
	histogram.Record(2000)
	histogram.Record(9000)
	histogram.Record(20000)

	// HTTP Server for metrics exporting.
	otelmetric.StartPrometheusMetricsServer(8000, "/metrics")
}

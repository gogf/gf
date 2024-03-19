// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmetric"
)

const (
	instrument        = "github.com/gogf/gf/example/metric/basic"
	instrumentVersion = "v1.0"
)

var (
	counter = gmetric.MustNewCounter(gmetric.MetricConfig{
		Name: "goframe.metric.demo.counter",
		Help: "This is a simple demo for Counter usage",
		Unit: "bytes",
		Attributes: gmetric.Attributes{
			gmetric.NewAttribute("const_label_1", 1),
		},
		Instrument:        instrument,
		InstrumentVersion: instrumentVersion,
	})
	upDownCounter = gmetric.MustNewUpDownCounter(gmetric.MetricConfig{
		Name: "goframe.metric.demo.updown_counter",
		Help: "This is a simple demo for UpDownCounter usage",
		Unit: "%",
		Attributes: gmetric.Attributes{
			gmetric.NewAttribute("const_label_2", 2),
		},
		Instrument:        instrument,
		InstrumentVersion: instrumentVersion,
	})
	histogram = gmetric.MustNewHistogram(gmetric.MetricConfig{
		Name: "goframe.metric.demo.histogram",
		Help: "This is a simple demo for histogram usage",
		Unit: "ms",
		Attributes: gmetric.Attributes{
			gmetric.NewAttribute("const_label_3", 3),
		},
		Instrument:        instrument,
		InstrumentVersion: instrumentVersion,
		Buckets:           []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
	})
	observableCounter = gmetric.MustNewObservableCounter(gmetric.MetricConfig{
		Name: "goframe.metric.demo.observable_counter",
		Help: "This is a simple demo for ObservableCounter usage",
		Unit: "%",
		Attributes: gmetric.Attributes{
			gmetric.NewAttribute("const_label_4", 4),
		},
		Instrument:        instrument,
		InstrumentVersion: instrumentVersion,
	})
	observableUpDownCounter = gmetric.MustNewObservableUpDownCounter(gmetric.MetricConfig{
		Name: "goframe.metric.demo.observable_updown_counter",
		Help: "This is a simple demo for ObservableUpDownCounter usage",
		Unit: "%",
		Attributes: gmetric.Attributes{
			gmetric.NewAttribute("const_label_5", 5),
		},
		Instrument:        instrument,
		InstrumentVersion: instrumentVersion,
	})
	observableGauge = gmetric.MustNewObservableGauge(gmetric.MetricConfig{
		Name: "goframe.metric.demo.observable_gauge",
		Help: "This is a simple demo for ObservableGauge usage",
		Unit: "%",
		Attributes: gmetric.Attributes{
			gmetric.NewAttribute("const_label_6", 6),
		},
		Instrument:        instrument,
		InstrumentVersion: instrumentVersion,
	})
)

func main() {
	var ctx = gctx.New()

	// Callback for observable metrics.
	gmetric.MustRegisterCallback(func(ctx context.Context, obs gmetric.Observer) error {
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
	provider := otelmetric.MustProvider(metric.WithReader(exporter))
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
	s := g.Server()
	s.BindHandler("/metrics", ghttp.WrapH(promhttp.Handler()))
	s.SetPort(8000)
	s.Run()
}

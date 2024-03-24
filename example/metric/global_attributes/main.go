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
	meter = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
		Instrument:        instrument,
		InstrumentVersion: instrumentVersion,
	})
	counter = meter.MustCounter(
		"goframe.metric.demo.counter",
		gmetric.MetricOption{
			Help: "This is a simple demo for Counter usage",
			Unit: "bytes",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_1", 1),
			},
		},
	)
	observableCounter = meter.MustObservableCounter(
		"goframe.metric.demo.observable_counter",
		gmetric.MetricOption{
			Help: "This is a simple demo for ObservableCounter usage",
			Unit: "%",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_2", 2),
			},
		},
	)
)

func main() {
	var ctx = gctx.New()

	gmetric.SetGlobalAttributes(gmetric.Attributes{
		gmetric.NewAttribute("g1", 1),
	}, gmetric.SetGlobalAttributesOption{
		Instrument:        instrument,
		InstrumentVersion: instrumentVersion,
		InstrumentPattern: "",
	})

	// Callback for observable metrics.
	meter.MustRegisterCallback(func(ctx context.Context, obs gmetric.Observer) error {
		obs.Observe(observableCounter, 10)
		return nil
	}, observableCounter)

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

	// HTTP Server for metrics exporting.
	s := g.Server()
	s.BindHandler("/metrics", ghttp.WrapH(promhttp.Handler()))
	s.SetPort(8000)
	s.Run()
}

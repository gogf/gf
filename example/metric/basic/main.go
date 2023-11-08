package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmetric"
)

var (
	counter = gmetric.NewCounter(gmetric.CounterConfig{
		MetricConfig: gmetric.MetricConfig{
			Name: "goframe.metric.demo.counter",
			Help: "This is a simple demo for Counter usage",
			Unit: "%",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_a", 1),
			},
			Instrument:        "github.com/gogf/gf/example/metric/basic",
			InstrumentVersion: "v1.0",
		},
	})
	gauge = gmetric.NewGauge(gmetric.GaugeConfig{
		MetricConfig: gmetric.MetricConfig{
			Name: "goframe.metric.demo.gauge",
			Help: "This is a simple demo for Gauge usage",
			Unit: "bytes",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_b", 2),
			},
			Instrument:        "github.com/gogf/gf/example/metric/basic",
			InstrumentVersion: "v1.0",
		},
	})
	histogram = gmetric.NewHistogram(gmetric.HistogramConfig{
		MetricConfig: gmetric.MetricConfig{
			Name: "goframe.metric.demo.histogram",
			Help: "This is a simple demo for histogram usage",
			Unit: "ms",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_c", 3),
			},
			Instrument:        "github.com/gogf/gf/example/metric/basic",
			InstrumentVersion: "v1.0",
		},
		Buckets: []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
	})
)

func main() {
	var ctx = gctx.New()

	// Prometheus exporter to export metrics as Prometheus format.
	exporter, err := prometheus.New(
		prometheus.WithoutCounterSuffixes(),
		prometheus.WithoutUnits(),
	)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}

	// OpenTelemetry provider.
	provider := otelmetric.NewProvider(metric.WithReader(exporter))
	defer provider.Shutdown(ctx)

	// Add value for counter.
	counter.Inc()
	counter.Add(10)

	// Set value for gauge.
	gauge.Set(100)
	gauge.Inc()
	gauge.Sub(1)

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
	s.SetPort(8199)
	s.Run()
}

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
	histogram1 = gmetric.NewHistogram(gmetric.HistogramConfig{
		MetricConfig: gmetric.MetricConfig{
			Name: "goframe.metric.demo.histogram1",
			Help: "This is a simple demo for histogram usage",
			Unit: "ms",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_a", 1),
			},
			Instrument:        "github.com/gogf/gf/example/metric/histogram_buckets",
			InstrumentVersion: "v1.0",
		},
		Buckets: []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
	})
	histogram2 = gmetric.NewHistogram(gmetric.HistogramConfig{
		MetricConfig: gmetric.MetricConfig{
			Name: "goframe.metric.demo.histogram2",
			Help: "This demos we can specify custom buckets in Histogram creating",
			Unit: "",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_b", 2),
			},
			Instrument:        "github.com/gogf/gf/example/metric/histogram_buckets",
			InstrumentVersion: "v1.0",
		},
		Buckets: []float64{100, 200, 300, 400, 500},
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

	// Record values for histogram1.
	histogram1.Record(1)
	histogram1.Record(20)
	histogram1.Record(30)
	histogram1.Record(101)
	histogram1.Record(2000)
	histogram1.Record(9000)
	histogram1.Record(20000)

	// Record values for histogram2.
	histogram2.Record(1)
	histogram2.Record(10)
	histogram2.Record(199)
	histogram2.Record(299)
	histogram2.Record(399)
	histogram2.Record(499)
	histogram2.Record(501)

	// HTTP Server for metrics exporting.
	s := g.Server()
	s.BindHandler("/metrics", ghttp.WrapH(promhttp.Handler()))
	s.SetPort(8199)
	s.Run()
}

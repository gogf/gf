package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"log"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmetric"
)

var (
	counter = gmetric.NewCounter(gmetric.CounterConfig{
		MetricConfig: gmetric.MetricConfig{
			Instrument: "github.com/gogf/gf/example/metric/basic",
			Name:       "goframe.metric.demo.counter",
			Help:       "This is a simple demo for Counter usage",
			Unit:       "%",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_a", 1),
				gmetric.NewAttribute("const_label_b", "value for const label b"),
			},
		},
	})
	gauge = gmetric.NewGauge(gmetric.GaugeConfig{
		MetricConfig: gmetric.MetricConfig{
			Instrument: "github.com/gogf/gf/example/metric/basic",
			Name:       "goframe.metric.demo.gauge",
			Help:       "This is a simple demo for Gauge usage",
			Unit:       "bytes",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_c", 2),
				gmetric.NewAttribute("const_label_d", "value for const label d"),
			},
		},
	})
	histogram1 = gmetric.NewHistogram(gmetric.HistogramConfig{
		MetricConfig: gmetric.MetricConfig{
			Instrument: "github.com/gogf/gf/example/metric/basic",
			Name:       "goframe.metric.demo.histogram1",
			Help:       "This is a simple demo for histogram usage",
			Unit:       "ms",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_e", 3),
				gmetric.NewAttribute("const_label_f", "value for const label f"),
			},
		},
		Buckets: []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
	})
	histogram2 = gmetric.NewHistogram(gmetric.HistogramConfig{
		MetricConfig: gmetric.MetricConfig{
			Instrument: "github.com/gogf/gf/example/metric/basic",
			Name:       "goframe.metric.demo.histogram2",
			Help:       "This demos we can specify custom buckets in Histogram creating",
			Unit:       "",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_g", 4),
				gmetric.NewAttribute("const_label_h", "value for const label h"),
			},
		},
		Buckets: []float64{100, 200, 300, 400, 500},
	})
)

func main() {
	var (
		ctx               = gctx.New()
		dynamicAttributes = gmetric.Attributes{
			gmetric.NewAttribute("dynamic_a", 1),
			gmetric.NewAttribute("dynamic_b", 0.1),
		}
	)

	// Prometheus exporter to export metrics as Prometheus format.
	exporter, err := prometheus.New(
		prometheus.WithoutCounterSuffixes(),
		prometheus.WithoutUnits(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// OpenTelemetry provider.
	provider := otelmetric.NewProvider(metric.WithReader(exporter))
	defer provider.Shutdown(ctx)

	// Add value for counter.
	counter.Inc(gmetric.Option{Attributes: dynamicAttributes})
	counter.Add(10)

	// Set value for gauge.
	gauge.Set(100)
	gauge.Inc()
	gauge.Sub(1)

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

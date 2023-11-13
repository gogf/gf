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
			Help: "This is a simple demo for dynamic attributes",
			Unit: "%",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label", 1),
			},
			Instrument:        "github.com/gogf/gf/example/metric/dynamic_attributes",
			InstrumentVersion: "v1.0",
		},
	})
)

func main() {
	var (
		ctx               = gctx.New()
		dynamicAttributes = gmetric.Option{
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("dynamic_label", 2),
			},
		}
	)

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
	defer provider.Shutdown(ctx)

	// Add value for counter.
	counter.Inc(dynamicAttributes)
	counter.Add(10, dynamicAttributes)

	// HTTP Server for metrics exporting.
	s := g.Server()
	s.BindHandler("/metrics", ghttp.WrapH(promhttp.Handler()))
	s.SetPort(8199)
	s.Run()
}

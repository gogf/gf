# GoFrame Metric In OpenTelemetry

## Installation

```
go get -u -v github.com/gogf/gf/contrib/metric/otelmetric/v2
```

suggested using `go.mod`:

```
require github.com/gogf/gf/contrib/metric/otelmetric/v2 latest
```

## Example

### [basic](../../../example/metric/basic/main.go)

```go
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
				gmetric.NewAttribute("const_label_1", 1),
			},
		},
	)
	upDownCounter = meter.MustUpDownCounter(
		"goframe.metric.demo.updown_counter",
		gmetric.MetricOption{
			Help: "This is a simple demo for UpDownCounter usage",
			Unit: "%",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_2", 2),
			},
		},
	)
	histogram = meter.MustHistogram(
		"goframe.metric.demo.histogram",
		gmetric.MetricOption{
			Help: "This is a simple demo for histogram usage",
			Unit: "ms",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_3", 3),
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
				gmetric.NewAttribute("const_label_4", 4),
			},
		},
	)
	observableUpDownCounter = meter.MustObservableUpDownCounter(
		"goframe.metric.demo.observable_updown_counter",
		gmetric.MetricOption{
			Help: "This is a simple demo for ObservableUpDownCounter usage",
			Unit: "%",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_5", 5),
			},
		},
	)
	observableGauge = meter.MustObservableGauge(
		"goframe.metric.demo.observable_gauge",
		gmetric.MetricOption{
			Help: "This is a simple demo for ObservableGauge usage",
			Unit: "%",
			Attributes: gmetric.Attributes{
				gmetric.NewAttribute("const_label_6", 6),
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
```

### [more examples](../../../example/metric/)

## License

`GoFrame Polaris` is licensed under the [MIT License](../../../LICENSE), 100% free and open-source, forever.


package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"log"
	"net/http"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/os/gmetric"
)

var counter = gmetric.NewCounter(gmetric.CounterConfig{
	MetricConfig: gmetric.MetricConfig{
		Instrument: "github.com/gogf/gf/example/metric/basic",
		Name:       "demo_counter",
		Help:       "This is a counter for demo",
		Unit:       "%",
		Attributes: gmetric.Attributes{
			gmetric.NewAttribute("A", 1),
			gmetric.NewAttribute("B", 2),
		},
	},
})

func main() {
	ctx := context.Background()

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := otelmetric.NewProvider(metric.WithReader(exporter))
	defer provider.Shutdown(ctx)

	counter.Inc()
	counter.Add(100)

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics()

	select {}
}

func serveMetrics() {
	log.Printf("serving metrics at localhost:2223/metrics")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2223", nil) //nolint:gosec // Ignoring G114: Use of net/http serve function that has no support for setting timeouts.
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}

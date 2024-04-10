// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"go.opentelemetry.io/otel/exporters/prometheus"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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
	provider := otelmetric.MustProvider(
		otelmetric.WithReader(exporter),
		otelmetric.WithBuiltInMetrics(),
	)
	provider.SetAsGlobal()
	defer provider.Shutdown(ctx)

	// A simple http client request for demonstration purpose only.
	url := `https://github.com/gogf/gf`
	content := g.Client().GetContent(ctx, url)
	g.Log().Infof(ctx, `content length from "%s": %d`, url, len(content))

	// A simple http server for metrics export.
	otelmetric.StartPrometheusMetricsServer(8000, "/metrics")
}

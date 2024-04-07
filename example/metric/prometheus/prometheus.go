// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/grand"
)

// Demo metric type Counter
var metricCounter = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "demo_counter",
		Help: "A demo counter.",
	},
)

// Demo metric type Gauge.
var metricGauge = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "demo_gauge",
		Help: "A demo gauge.",
	},
)

func main() {
	// Create prometheus metric registry.
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		metricCounter,
		metricGauge,
	)

	// Start metric http server.
	s := g.Server()
	// Fake metric values.
	// http://127.0.0.1:8000/
	s.BindHandler("/", func(r *ghttp.Request) {
		metricCounter.Add(1)
		metricGauge.Set(float64(grand.N(1, 100)))
		r.Response.Write("fake ok")
	})
	// Export metric values.
	// You can view http://127.0.0.1:8000/metrics to see all metric values.
	s.BindHandler("/metrics", ghttp.WrapH(promhttp.Handler()))
	s.SetPort(8000)
	s.Run()
}

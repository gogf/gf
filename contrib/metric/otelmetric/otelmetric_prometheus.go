// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// PrometheusHandler returns the http handler for prometheus metrics exporting.
func PrometheusHandler(r *ghttp.Request) {
	// Remove all builtin metrics that are produced by prometheus client.
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	prometheus.Unregister(collectors.NewGoCollector())

	handler := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})
	handler.ServeHTTP(r.Response.Writer, r.Request)
}

// StartPrometheusMetricsServer starts running a http server for metrics exporting.
func StartPrometheusMetricsServer(port int, path string) {
	s := g.Server()
	s.BindHandler(path, PrometheusHandler)
	s.SetPort(port)
	s.Run()
}

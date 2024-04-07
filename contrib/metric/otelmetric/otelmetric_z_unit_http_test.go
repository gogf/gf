// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmetric"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_HTTP_Server(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.BindHandler("/user/:id", func(r *ghttp.Request) {
			r.Response.Write("user")
		})
		s.BindHandler("/order/:id", func(r *ghttp.Request) {
			r.Response.Write("order")
		})
		s.BindHandler("/metrics", ghttp.WrapH(promhttp.Handler()))
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

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
		provider := otelmetric.MustProvider(otelmetric.WithReader(exporter))
		defer provider.Shutdown(ctx)

		gmetric.SetGlobalProvider(provider)
		defer gmetric.SetGlobalProvider(nil)

		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		c.GetContent(ctx, "/user/1")
		c.PutContent(ctx, "/user/1", "123")
		c.PostContent(ctx, "/user/2", "123")
		c.DeleteContent(ctx, "/user/3")
		c.GetContent(ctx, "/order/1")
		c.PutContent(ctx, "/order/1", "1234")
		c.PostContent(ctx, "/order/2", "1234")
		c.DeleteContent(ctx, "/order/3")

		var (
			metricsContent = c.GetContent(ctx, "/metrics")
			expectContent  = gtest.DataContent("http.prometheus.metrics.txt")
		)
		expectContent, _ = gregex.ReplaceString(
			`otel_scope_version=".+?"`,
			fmt.Sprintf(`otel_scope_version="%s"`, gf.VERSION),
			expectContent,
		)
		expectContent, _ = gregex.ReplaceString(
			`server_port=".+?"`,
			fmt.Sprintf(`server_port="%d"`, s.GetListenedPort()),
			expectContent,
		)
		//fmt.Println(metricsContent)
		for _, line := range gstr.SplitAndTrim(expectContent, "\n") {
			//fmt.Println(line)
			t.Assert(gstr.Contains(metricsContent, line), true)
		}
	})
}

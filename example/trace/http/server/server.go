// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"github.com/gogf/gf/contrib/trace/otlphttp/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	serviceName = "otlp-http-server"
	endpoint    = "tracing-analysis-dc-hz.aliyuncs.com"
	path        = "adapt_******_******/api/otlp/traces"
)

func main() {
	var ctx = gctx.New()
	shutdown, err := otlphttp.Init(serviceName, endpoint, path)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	defer shutdown()

	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/hello", HelloHandler)
	})
	s.SetPort(8199)
	s.Run()
}

// HelloHandler is a demo handler for tracing.
func HelloHandler(r *ghttp.Request) {
	ctx, span := gtrace.NewSpan(r.Context(), "HelloHandler")
	defer span.End()

	value := gtrace.GetBaggageVar(ctx, "name").String()

	r.Response.Write("hello:", value)
}

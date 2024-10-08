// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/gogf/gf/contrib/trace/provider/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	serviceName = "otlp-http-client"
	endpoint    = "tracing-analysis-dc-hz.aliyuncs.com"
	path        = "adapt_******_******/api/otlp/traces"
)

func main() {
	var (
		serverIP, err = provider.GetLocalIP()
		ctx           = gctx.New()
	)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	var res *resource.Resource
	if res, err = provider.NewDefaultResource(ctx, serviceName, serverIP); err != nil {
		g.Log().Fatal(ctx, err)
	}

	var (
		client   = provider.NewHTTPDefaultTraceClient(endpoint, path)
		exporter *otlptrace.Exporter
	)
	if exporter, err = provider.NewExporter(ctx, client); err != nil {
		g.Log().Fatal(ctx, err)
	}
	var (
		// BatchSpanProcessor batches spans before exporting them.
		// This is a useful way to reduce the number of calls made to the exporter.
		// The batch processor will automatically flush the spans if the batch size is reached.
		// bsp = provider.NewSimpleSpanProcessor(exporter)
		bsp = provider.NewBatchSpanProcessor(exporter)
		// AlwaysOnSampler is a sampler that samples every trace.
		// This is useful for debugging and testing.
		sampler  = provider.NewAlwaysOnSampler()
		shutdown func(ctx context.Context)
	)

	if shutdown, err = provider.InitTracer(provider.NewDefaultTracerProviderOptions(sampler, res, bsp)...); err != nil {
		g.Log().Fatal(ctx, err)
	}
	defer shutdown(ctx)

	StartRequests()
}

// StartRequests starts requests.
func StartRequests() {
	ctx, span := gtrace.NewSpan(gctx.New(), "StartRequests")
	defer span.End()

	ctx = gtrace.SetBaggageValue(ctx, "name", "john")

	content := g.Client().GetContent(ctx, "http://127.0.0.1:8199/hello")
	g.Log().Print(ctx, content)
}

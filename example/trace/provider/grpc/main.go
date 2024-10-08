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
	serviceName = "otlp-grpc-client"
	endpoint    = "tracing-analysis-dc-bj.aliyuncs.com:8090"
	traceToken  = "******_******"
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
	if res, err = provider.NewResource(ctx, serviceName, serverIP); err != nil {
		g.Log().Fatal(ctx, err)
	}

	var (
		client   = provider.NewGRPCDefaultTraceClient(endpoint, traceToken)
		exporter *otlptrace.Exporter
	)
	if exporter, err = provider.NewExporter(ctx, client); err != nil {
		g.Log().Fatal(ctx, err)
	}
	var (
		bsp = provider.NewBatchSpanProcessor(exporter)
		// AlwaysOnSampler is a sampler that samples every trace.
		// sampler  = provider.NewTraceIDRatioBasedSampler(0.1)
		sampler  = provider.NewAlwaysOnSampler()
		shutdown func(ctx context.Context)
	)

	if shutdown, err = provider.InitTracer(provider.NewTracerProviderOptions(sampler, res, bsp)...); err != nil {
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

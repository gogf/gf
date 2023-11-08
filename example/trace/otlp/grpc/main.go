// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"github.com/gogf/gf/contrib/trace/otlpgrpc/v2"

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
	var ctx = gctx.New()
	shutdown, err := otlpgrpc.Init(serviceName, endpoint, traceToken)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	defer shutdown()

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

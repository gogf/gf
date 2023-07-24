package main

import (
	"github.com/gogf/gf/contrib/trace/jaeger/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	ServiceName       = "http-client"
	JaegerUdpEndpoint = "localhost:6831"
)

func main() {
	var ctx = gctx.New()
	tp, err := jaeger.Init(ServiceName, JaegerUdpEndpoint)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	defer tp.Shutdown(ctx)

	StartRequests()
}

func StartRequests() {
	ctx, span := gtrace.NewSpan(gctx.New(), "StartRequests")
	defer span.End()

	ctx = gtrace.SetBaggageValue(ctx, "name", "john")

	content := g.Client().GetContent(ctx, "http://127.0.0.1:8199/hello")
	g.Log().Print(ctx, content)
}

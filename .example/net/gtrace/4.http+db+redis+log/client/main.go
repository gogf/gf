package main

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/net/gtrace"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	JaegerEndpoint = "http://localhost:14268/api/traces"
	ServiceName    = "tracing-http-client"
)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer() func() {
	// Create and install Jaeger export pipeline.
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(JaegerEndpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: ServiceName,
		}),
		jaeger.WithSDK(&sdkTrace.Config{DefaultSampler: sdkTrace.AlwaysSample()}),
	)
	if err != nil {
		g.Log().Fatal(err)
	}
	return flush
}

func main() {
	flush := initTracer()
	defer flush()

	ctx, span := gtrace.Tracer().Start(context.Background(), "root")
	defer span.End()

	client := g.Client().Use(ghttp.MiddlewareClientTracing)
	//client := g.Client()
	// Add user info.
	idStr := client.Ctx(ctx).PostContent(
		"http://127.0.0.1:8199/user/insert",
		g.Map{
			"name": "john",
		},
	)
	if idStr == "" {
		g.Log().Ctx(ctx).Fatal("retrieve empty id string")
	}
	g.Log().Ctx(ctx).Print("insert:", idStr)

	// Query user info.
	userJson := client.Ctx(ctx).GetContent(
		"http://127.0.0.1:8199/user/query",
		g.Map{
			"id": idStr,
		},
	)
	g.Log().Ctx(ctx).Print("query:", idStr, userJson)

	// Delete user info.
	deleteResult := client.Ctx(ctx).PostContent(
		"http://127.0.0.1:8199/user/delete",
		g.Map{
			"id": idStr,
		},
	)
	g.Log().Ctx(ctx).Print("delete:", idStr, deleteResult)
}

package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/net/gtrace"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	JaegerEndpoint = "http://localhost:14268/api/traces"
	ServiceName    = "tracing-http-server"
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

	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareServerTracing)
		group.GET("/hello", helloHandler)
	})
	s.SetPort(8199)
	s.Run()
}

func helloHandler(r *ghttp.Request) {
	ctx, span := gtrace.Tracer().Start(r.Context(), "helloHandler")
	defer span.End()

	value := baggage.Value(ctx, "name")
	r.Response.Write("hello:", value.AsString())
}

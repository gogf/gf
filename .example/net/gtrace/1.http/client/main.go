package main

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/net/gtrace"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	JaegerEndpoint = "http://localhost:14268/api/traces"
	ServiceName    = "TracingHttpClient"
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

	ctx, span := gtrace.Tracer().Start(context.Background(), "test")
	defer span.End()

	ctx = baggage.ContextWithValues(ctx, label.String("name", "john"))
	client := g.Client().Use(ghttp.MiddlewareClientTracing)
	content := client.Ctx(ctx).GetContent("http://127.0.0.1:8199/hello")
	g.Log().Print(content)
}

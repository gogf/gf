package main

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

type tracingApi struct{}

const (
	JaegerEndpoint = "http://localhost:14268/api/traces"
	ServiceName    = "tracing-demo-redis"
)

func (api *tracingApi) Set(r *ghttp.Request) {
	_, err := g.Redis().Ctx(r.Context()).Do("SET", r.GetString("key"), r.GetString("value"))
	if err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	r.Response.Write("ok")
}

func (api *tracingApi) Get(r *ghttp.Request) {
	value, err := g.Redis().Ctx(r.Context()).DoVar(
		"GET", r.GetString("key"),
	)
	if err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	r.Response.Write(value.String())
}

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
		group.ALL("/redis", new(tracingApi))
	})
	s.SetPort(8199)
	s.Run()
}

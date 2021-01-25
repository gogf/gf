package main

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	JaegerEndpoint = "http://localhost:14268/api/traces"
	ServiceName    = "TracingHttpServerWithDatabase"
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
		group.ALL("/user", new(dbTracingApi))
	})
	s.SetPort(8199)
	s.Run()
}

type dbTracingApi struct{}

func (api *dbTracingApi) Insert(r *ghttp.Request) {
	result, err := g.Table("user").Ctx(r.Context()).Insert(g.Map{
		"name": r.GetString("name"),
	})
	if err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	id, _ := result.LastInsertId()
	r.Response.Write("id:", id)
}

func (api *dbTracingApi) Query(r *ghttp.Request) {
	one, err := g.Table("user").Ctx(r.Context()).FindOne(r.GetInt("id"))
	if err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	r.Response.Write("user:", one)
}

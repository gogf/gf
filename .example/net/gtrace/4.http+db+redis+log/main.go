package main

import (
	"github.com/gogf/gcache-adapter/adapter"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"time"
)

type tracingApi struct{}

const (
	JaegerEndpoint = "http://localhost:14268/api/traces"
	ServiceName    = "TracingHttpServerWithDBRedisLog"
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

func (api *tracingApi) Insert(r *ghttp.Request) {
	result, err := g.Table("user").Ctx(r.Context()).Insert(g.Map{
		"name": r.GetString("name"),
	})
	if err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	id, _ := result.LastInsertId()
	r.Response.Write("id:", id)
}

func (api *tracingApi) Query(r *ghttp.Request) {
	one, err := g.Table("user").
		Ctx(r.Context()).
		Cache(5 * time.Second).
		FindOne(r.GetInt("id"))
	if err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	r.Response.Write("user:", one)
}

func main() {
	flush := initTracer()
	defer flush()

	g.DB().GetCache().SetAdapter(adapter.NewRedis(g.Redis()))

	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareServerTracing)
		group.ALL("/user", new(tracingApi))
	})
	s.SetPort(8199)
	s.Run()
}

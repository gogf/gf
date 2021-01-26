package main

import (
	"fmt"
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

type userApiInsert struct {
	Name string `v:"required#Please input user name."`
}

// Insert is a route handler for inserting user info into dtabase.
func (api *tracingApi) Insert(r *ghttp.Request) {
	var (
		dataReq *userApiInsert
	)
	if err := r.Parse(&dataReq); err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	result, err := g.Table("user").Ctx(r.Context()).Insert(g.Map{
		"name": dataReq.Name,
	})
	if err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	id, _ := result.LastInsertId()
	r.Response.Write(id)
}

type userApiQuery struct {
	Id int `v:"min:1#User id is required for querying."`
}

// Query is a route handler for querying user info. It firstly retrieves the info from redis,
// if there's nothing in the redis, it then does db select.
func (api *tracingApi) Query(r *ghttp.Request) {
	var (
		dataReq *userApiQuery
	)
	if err := r.Parse(&dataReq); err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	one, err := g.Table("user").
		Ctx(r.Context()).
		Cache(5*time.Second, api.userCacheKey(dataReq.Id)).
		FindOne(dataReq.Id)
	if err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	r.Response.WriteJson(one)
}

type userApiDelete struct {
	Id int `v:"min:1#User id is required for deleting."`
}

// Delete is a route handler for deleting specified user info.
func (api *tracingApi) Delete(r *ghttp.Request) {
	var (
		dataReq *userApiDelete
	)
	if err := r.Parse(&dataReq); err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	_, err := g.Table("user").
		Ctx(r.Context()).
		Cache(-1, api.userCacheKey(dataReq.Id)).
		WherePri(dataReq.Id).
		Delete()
	if err != nil {
		r.Response.WriteExit(gerror.Current(err))
	}
	r.Response.Write("ok")
}

func (api *tracingApi) userCacheKey(id int) string {
	return fmt.Sprintf(`userInfo:%d`, id)
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

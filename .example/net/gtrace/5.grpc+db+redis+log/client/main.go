package main

import (
	"context"
	"github.com/gogf/gf/.example/net/gtrace/5.grpc+db+redis+log/protobuf/user"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtrace"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

const (
	JaegerEndpoint = "http://localhost:14268/api/traces"
	ServiceName    = "tracing-grpc-client"
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

func StartRequests() {
	ctx, span := gtrace.Tracer().Start(context.Background(), "StartRequests")
	defer span.End()

	conn, err := grpc.Dial(":8000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		g.Log().Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := user.NewUserClient(conn)

	// Insert.
	insertRes, err := client.Insert(ctx, &user.InsertReq{
		Name: "john",
	})
	if err != nil {
		g.Log().Ctx(ctx).Fatalf(`%+v`, err)
	}
	g.Log().Ctx(ctx).Println("insert:", insertRes.Id)

	// Query.
	queryRes, err := client.Query(ctx, &user.QueryReq{
		Id: insertRes.Id,
	})
	if err != nil {
		g.Log().Ctx(ctx).Fatalf(`%+v`, err)
	}
	g.Log().Ctx(ctx).Println("query:", queryRes)

	// Delete.
	_, err = client.Delete(ctx, &user.DeleteReq{
		Id: insertRes.Id,
	})
	if err != nil {
		g.Log().Ctx(ctx).Fatalf(`%+v`, err)
	}
	g.Log().Ctx(ctx).Println("delete:", "ok")

}

func main() {
	flush := initTracer()
	defer flush()

	StartRequests()
}

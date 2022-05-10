package main

import (
	"github.com/gogf/gf/contrib/trace/jaeger/v2"
	"github.com/gogf/gf/example/trace/grpc_with_db/protobuf/user"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	ServiceName       = "grpc-client-with-db"
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

	var client, err = user.NewClient()
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}

	// Baggage.
	ctx = gtrace.SetBaggageValue(ctx, "uid", 100)

	// Insert.
	insertRes, err := client.User().Insert(ctx, &user.InsertReq{
		Name: "john",
	})
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}
	g.Log().Info(ctx, "insert id:", insertRes.Id)

	// Query.
	queryRes, err := client.User().Query(ctx, &user.QueryReq{
		Id: insertRes.Id,
	})
	if err != nil {
		g.Log().Errorf(ctx, `%+v`, err)
		return
	}
	g.Log().Info(ctx, "query result:", queryRes)

	// Delete.
	_, err = client.User().Delete(ctx, &user.DeleteReq{
		Id: insertRes.Id,
	})
	if err != nil {
		g.Log().Errorf(ctx, `%+v`, err)
		return
	}
	g.Log().Info(ctx, "delete id:", insertRes.Id)

	// Delete with error.
	_, err = client.User().Delete(ctx, &user.DeleteReq{
		Id: -1,
	})
	if err != nil {
		g.Log().Errorf(ctx, `%+v`, err)
		return
	}
	g.Log().Info(ctx, "delete id:", -1)
}

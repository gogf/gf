package main

import (
	"github.com/gogf/gf/contrib/trace/jaeger/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	ServiceName       = "http-client-with-db"
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

	var (
		err    error
		client = g.Client()
	)
	// Add user info.
	var insertRes = struct {
		ghttp.DefaultHandlerResponse
		Data struct{ Id int64 } `json:"data"`
	}{}
	err = client.PostVar(ctx, "http://127.0.0.1:8199/user/insert", g.Map{
		"name": "john",
	}).Scan(&insertRes)
	if err != nil {
		panic(err)
	}
	g.Log().Info(ctx, "insert result:", insertRes)
	if insertRes.Data.Id == 0 {
		g.Log().Error(ctx, "retrieve empty id string")
		return
	}

	// Query user info.
	var queryRes = struct {
		ghttp.DefaultHandlerResponse
		Data struct{ User gdb.Record } `json:"data"`
	}{}
	err = client.GetVar(ctx, "http://127.0.0.1:8199/user/query", g.Map{
		"id": insertRes.Data.Id,
	}).Scan(&queryRes)
	if err != nil {
		panic(err)
	}
	g.Log().Info(ctx, "query result:", queryRes)

	// Delete user info.
	var deleteRes = struct {
		ghttp.DefaultHandlerResponse
	}{}
	err = client.PostVar(ctx, "http://127.0.0.1:8199/user/delete", g.Map{
		"id": insertRes.Data.Id,
	}).Scan(&deleteRes)
	if err != nil {
		panic(err)
	}
	g.Log().Info(ctx, "delete result:", deleteRes)
}

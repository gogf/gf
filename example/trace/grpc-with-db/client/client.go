// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/contrib/trace/otlpgrpc/v2"
	"github.com/gogf/gf/example/trace/grpc-with-db/protobuf/user"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	serviceName = "otlp-grpc-client"
	endpoint    = "tracing-analysis-dc-bj.aliyuncs.com:8090"
	traceToken  = "******_******"
)

func main() {
	grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))

	var ctx = gctx.New()
	shutdown, err := otlpgrpc.Init(serviceName, endpoint, traceToken)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	defer shutdown()

	StartRequests()
}

// StartRequests is a demo for tracing.
func StartRequests() {
	ctx, span := gtrace.NewSpan(gctx.New(), "StartRequests")
	defer span.End()

	client := user.NewUserClient(grpcx.Client.MustNewGrpcClientConn("demo"))

	// Baggage.
	ctx = gtrace.SetBaggageValue(ctx, "uid", 100)

	// Insert.
	insertRes, err := client.Insert(ctx, &user.InsertReq{
		Name: "john",
	})
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}
	g.Log().Info(ctx, "insert id:", insertRes.Id)

	// Query.
	queryRes, err := client.Query(ctx, &user.QueryReq{
		Id: insertRes.Id,
	})
	if err != nil {
		g.Log().Errorf(ctx, `%+v`, err)
		return
	}
	g.Log().Info(ctx, "query result:", queryRes)

	// Delete.
	if _, err = client.Delete(ctx, &user.DeleteReq{
		Id: insertRes.Id,
	}); err != nil {
		g.Log().Errorf(ctx, `%+v`, err)
		return
	}
	g.Log().Info(ctx, "delete id:", insertRes.Id)

	// Delete with error.
	if _, err = client.Delete(ctx, &user.DeleteReq{
		Id: -1,
	}); err != nil {
		g.Log().Errorf(ctx, `%+v`, err)
		return
	}
	g.Log().Info(ctx, "delete id:", -1)
}

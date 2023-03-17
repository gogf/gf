// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/example/rpc/grpcx/ctx/protobuf"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var (
		conn   = grpcx.Client.MustNewGrpcClientConn("demo")
		client = protobuf.NewGreeterClient(conn)
		ctx    = grpcx.Ctx.NewOutgoing(gctx.New(), g.Map{
			"UserId":   "1000",
			"UserName": "john",
		})
	)
	g.Log().Infof(ctx, `outgoing data: %v`, grpcx.Ctx.OutgoingMap(ctx).Map())
	res, err := client.SayHello(ctx, &protobuf.HelloRequest{Name: "World"})
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Log().Debug(ctx, "Response:", res.Message)
}

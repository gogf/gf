// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package main

import (
	"context"

	"github.com/wangyougui/gf/contrib/rpc/grpcx/v2"
	"github.com/wangyougui/gf/example/rpc/grpcx/balancer/protobuf"
	"github.com/wangyougui/gf/v2/frame/g"
	"github.com/wangyougui/gf/v2/os/gctx"
)

func main() {
	var (
		ctx    context.Context
		conn   = grpcx.Client.MustNewGrpcClientConn("demo", grpcx.Balancer.WithRandom())
		client = protobuf.NewGreeterClient(conn)
	)
	for i := 0; i < 10; i++ {
		ctx = gctx.New()
		res, err := client.SayHello(ctx, &protobuf.HelloRequest{Name: "World"})
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		g.Log().Debug(ctx, "Response:", res.Message)
	}
}

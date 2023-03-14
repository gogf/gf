// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/example/rpc/grpcx/basic/protobuf"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var (
		ctx    = gctx.GetInitCtx()
		client = protobuf.NewGreeterClient(grpcx.Client.MustNewGrpcClientConn("demo"))
	)
	for i := 0; i < 100; i++ {
		res, err := client.SayHello(ctx, &protobuf.HelloRequest{Name: "gfer"})
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		g.Log().Print(ctx, "Response:", res.Message)
		time.Sleep(time.Second)
	}
}

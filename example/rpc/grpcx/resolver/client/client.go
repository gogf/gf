// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package main

import (
	"github.com/wangyougui/gf/contrib/registry/etcd/v2"
	"github.com/wangyougui/gf/contrib/rpc/grpcx/v2"
	"github.com/wangyougui/gf/example/rpc/grpcx/resolver/protobuf"
	"github.com/wangyougui/gf/v2/frame/g"
	"github.com/wangyougui/gf/v2/os/gctx"
)

func main() {
	grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))

	var (
		ctx    = gctx.New()
		conn   = grpcx.Client.MustNewGrpcClientConn("demo")
		client = protobuf.NewGreeterClient(conn)
	)
	res, err := client.SayHello(ctx, &protobuf.HelloRequest{Name: "World"})
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Log().Debug(ctx, "Response:", res.Message)
}

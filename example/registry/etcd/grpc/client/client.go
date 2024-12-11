// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"
	"time"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/example/registry/etcd/grpc/protobuf"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))

	var (
		conn   = grpcx.Client.MustNewGrpcClientConn("demo")
		client = protobuf.NewGreeterClient(conn)
	)

	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			res, err := client.SayHello(ctx, &protobuf.HelloRequest{Name: "World"})
			if err != nil {
				g.Log().Errorf(ctx, `%+v`, err)
			} else {
				g.Log().Debug(ctx, "Response:", res.Message)
			}
		}()

		time.Sleep(time.Second)
	}

}

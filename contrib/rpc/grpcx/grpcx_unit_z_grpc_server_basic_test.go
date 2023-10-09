// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2/testdata/controller"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2/testdata/protobuf"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Grpcx_Grpc_Server_Basic(t *testing.T) {
	c := grpcx.Server.NewConfig()
	c.Name = guid.S()
	s := grpcx.Server.New(c)
	controller.Register(s)
	s.Start()
	time.Sleep(time.Millisecond * 100)
	defer s.Stop()

	// use service discovery.
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx    = gctx.New()
			conn   = grpcx.Client.MustNewGrpcClientConn(c.Name)
			client = protobuf.NewGreeterClient(conn)
		)
		res, err := client.SayHello(ctx, &protobuf.HelloRequest{Name: "World"})
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		t.Assert(res.Message, `Hello World`)
	})

	// use direct address.
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = gctx.New()
			address = fmt.Sprintf(`%s:%d`, gipv4.MustGetIntranetIp(), s.GetListenedPort())
			conn    = grpcx.Client.MustNewGrpcClientConn(address)
			client  = protobuf.NewGreeterClient(conn)
		)
		res, err := client.SayHello(ctx, &protobuf.HelloRequest{Name: "World"})
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		t.Assert(res.Message, `Hello World`)
	})
}

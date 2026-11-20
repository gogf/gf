// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2/testdata/protobuf"
)

// slowGreeter blocks for the duration given by the request name and ignores the context on purpose.
type slowGreeter struct {
	protobuf.UnimplementedGreeterServer
}

func (slowGreeter) SayHello(_ context.Context, in *protobuf.HelloRequest) (*protobuf.HelloReply, error) {
	d, err := time.ParseDuration(in.GetName())
	if err != nil {
		return nil, err
	}
	time.Sleep(d)
	return &protobuf.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func startTimeoutServer(c *grpcx.GrpcServerConfig) (protobuf.GreeterClient, func()) {
	c.Name = guid.S()
	c.Address = "127.0.0.1:0"
	s := grpcx.Server.New(c)
	protobuf.RegisterGreeterServer(s.Server, slowGreeter{})
	s.Start()
	time.Sleep(time.Millisecond * 100)
	conn, err := grpc.NewClient(
		fmt.Sprintf("127.0.0.1:%d", s.GetListenedPort()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	return protobuf.NewGreeterClient(conn), func() {
		_ = conn.Close()
		s.Stop()
	}
}

func Test_Grpcx_Grpc_Server_Timeout(t *testing.T) {
	c := grpcx.Server.NewConfig()
	c.Timeout = 200 * time.Millisecond
	client, stop := startTimeoutServer(c)
	defer stop()

	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		// Finishes in time.
		res, err := client.SayHello(ctx, &protobuf.HelloRequest{Name: "10ms"})
		t.AssertNil(err)
		t.Assert(res.Message, "Hello 10ms")

		// The client has no deadline, the server limits the request.
		start := time.Now()
		_, err = client.SayHello(ctx, &protobuf.HelloRequest{Name: "3s"})
		t.Assert(status.Code(err), codes.DeadlineExceeded)
		t.Assert(time.Since(start) < 2*time.Second, true)
	})
}

func Test_Grpcx_Grpc_Server_Timeout_Disabled(t *testing.T) {
	client, stop := startTimeoutServer(grpcx.Server.NewConfig())
	defer stop()

	gtest.C(t, func(t *gtest.T) {
		res, err := client.SayHello(gctx.New(), &protobuf.HelloRequest{Name: "300ms"})
		t.AssertNil(err)
		t.Assert(res.Message, "Hello 300ms")
	})
}

func Test_Grpcx_Grpc_Server_Timeout_Method(t *testing.T) {
	c := grpcx.Server.NewConfig()
	c.Timeout = 100 * time.Millisecond
	c.MethodTimeouts = map[string]time.Duration{
		"/protobuf.Greeter/SayHello": 0,
	}
	client, stop := startTimeoutServer(c)
	defer stop()

	// The method specific value overrides the global one, 0 means unlimited.
	gtest.C(t, func(t *gtest.T) {
		res, err := client.SayHello(gctx.New(), &protobuf.HelloRequest{Name: "300ms"})
		t.AssertNil(err)
		t.Assert(res.Message, "Hello 300ms")
	})

	c = grpcx.Server.NewConfig()
	c.MethodTimeouts = map[string]time.Duration{
		"/protobuf.Greeter/SayHello": 100 * time.Millisecond,
	}
	client2, stop2 := startTimeoutServer(c)
	defer stop2()
	gtest.C(t, func(t *gtest.T) {
		_, err := client2.SayHello(gctx.New(), &protobuf.HelloRequest{Name: "2s"})
		t.Assert(status.Code(err), codes.DeadlineExceeded)
	})
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package controller

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/example/rpc/grpcx/basic/protobuf"
)

type Controller struct {
	protobuf.UnimplementedGreeterServer
}

func Register(s *grpcx.GrpcServer) {
	protobuf.RegisterGreeterServer(s.Server, &Controller{})
}

// SayHello implements helloworld.GreeterServer
func (s *Controller) SayHello(ctx context.Context, in *protobuf.HelloRequest) (*protobuf.HelloReply, error) {
	return &protobuf.HelloReply{Message: "Hello " + in.GetName()}, nil
}

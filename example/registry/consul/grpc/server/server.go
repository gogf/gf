// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/contrib/registry/consul/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"

	"github.com/gogf/gf/example/registry/consul/grpc/controller"
)

func main() {
	registry, err := consul.New(consul.WithAddress("127.0.0.1:8500"))
	if err != nil {
		g.Log().Fatal(context.Background(), err)
	}
	grpcx.Resolver.Register(registry)

	s := grpcx.Server.New()
	controller.Register(s)
	s.Run()
}

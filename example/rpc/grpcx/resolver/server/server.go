// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/example/rpc/grpcx/resolver/controller"
)

func main() {
	grpcx.Resolver.Registry(etcd.New("127.0.0.1:2379"))

	s := grpcx.Server.New()
	controller.Register(s)
	s.Run()
}

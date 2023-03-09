// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/example/rpc/grpcx/basic/protocol"
	"github.com/gogf/gf/example/rpc/grpcx/basic/service"
)

func main() {
	s := grpcx.Server.New()
	protocol.RegisterEchoServer(s.Server, new(service.Echo))
	s.Run()
}

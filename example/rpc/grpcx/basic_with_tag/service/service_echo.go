// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package service

import (
	"fmt"

	"context"

	"github.com/gogf/gf/example/rpc/grpcx/basic_with_tag/protocol"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

// Echo is the service for echo.
type Echo struct{}

// Say implements the protobuf.EchoServer interface.
func (s *Echo) Say(ctx context.Context, r *protocol.SayReq) (*protocol.SayRes, error) {
	g.Log().Print(ctx, "Received:", r.Content)
	text := fmt.Sprintf(`%s: > %s`, gcmd.GetOpt("node", "default"), r.Content)
	return &protocol.SayRes{Content: text}, nil
}

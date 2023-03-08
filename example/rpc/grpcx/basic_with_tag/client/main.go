// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"time"

	"github.com/gogf/gf/example/rpc/grpcx/basic_with_tag/protocol"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var (
		ctx         = gctx.GetInitCtx()
		client, err = protocol.NewClient()
	)
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}
	for i := 0; i < 100; i++ {
		res, err := client.Echo().Say(ctx, &protocol.SayReq{Content: "Hello"})
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		g.Log().Print(ctx, "Response:", res.Content)
		time.Sleep(time.Second)
	}
}

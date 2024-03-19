// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package main

import (
	"fmt"

	_ "github.com/wangyougui/gf/contrib/nosql/redis/v2"

	"github.com/wangyougui/gf/v2/frame/g"
	"github.com/wangyougui/gf/v2/os/gctx"
)

func main() {
	var (
		ctx     = gctx.New()
		channel = "channel"
	)
	conn, _ := g.Redis().Conn(ctx)
	defer conn.Close(ctx)
	_, err := conn.Subscribe(ctx, channel)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	for {
		msg, err := conn.ReceiveMessage(ctx)
		if err != nil {
			g.Log().Fatal(ctx, err)
		}
		fmt.Println(msg.Payload)
	}
}

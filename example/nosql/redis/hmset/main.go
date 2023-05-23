package main

import (
	"fmt"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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

package main

import (
	"fmt"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	config = gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      1,
	}
	group = "cache"
	ctx   = gctx.New()
)

func main() {
	gredis.SetConfig(&config, group)

	_, err := g.Redis(group).Set(ctx, "key", "value")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	value, err := g.Redis(group).Get(ctx, "key")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	fmt.Println(value.String())
}

package main

import (
	"fmt"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var ctx = gctx.New()
	_, err := g.Redis().Set(ctx, "key", "value")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	value, err := g.Redis().Get(ctx, "key")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	fmt.Println(value.String())
}

package main

import (
	"fmt"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/gctx"
)

func main() {
	var (
		ctx           = gctx.New()
		key1  int32   = 1
		key2  float64 = 1
		value         = `value`
	)
	_ = gcache.Ctx(ctx).Set(key1, value, 0)
	fmt.Println(gcache.Ctx(ctx).Get(key1))
	fmt.Println(gcache.Ctx(ctx).Get(key2))
}

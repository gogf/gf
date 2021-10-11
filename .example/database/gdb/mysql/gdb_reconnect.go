package main

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"time"
)

func main() {
	var (
		db  = g.DB()
		ctx = gctx.New()
	)
	db.SetDebug(true)
	for {
		r, err := db.Ctx(ctx).Model("user").All()
		fmt.Println(err)
		fmt.Println(r)
		time.Sleep(time.Second * 10)
	}
}

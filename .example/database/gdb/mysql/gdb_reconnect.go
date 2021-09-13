package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gctx"
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

package main

import (
	"github.com/gogf/gf/os/gctx"
	"time"

	"github.com/gogf/gf/frame/g"
)

func main() {
	var (
		db  = g.DB()
		ctx = gctx.New()
	)

	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)

	for {
		for i := 0; i < 10; i++ {
			go db.Ctx(ctx).Model("user").All()
		}
		time.Sleep(time.Millisecond * 100)
	}

}

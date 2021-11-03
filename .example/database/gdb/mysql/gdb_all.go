package main

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var (
		db  = g.DB()
		ctx = gctx.New()
	)
	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)

	r, e := db.Ctx(ctx).GetAll("SELECT * from `user` where id in(?)", g.Slice{})
	if e != nil {
		fmt.Println(e)
	}
	if r != nil {
		fmt.Println(r)
	}
}

package main

import (
	"fmt"
	"github.com/gogf/gf/os/gctx"

	"github.com/gogf/gf/frame/g"
)

func main() {
	var (
		db  = g.DB()
		ctx = gctx.New()
	)
	db.SetDebug(true)

	r, e := db.Ctx(ctx).Model("test").All()
	if e != nil {
		panic(e)
	}
	if r != nil {
		fmt.Println(r.List())
	}
}

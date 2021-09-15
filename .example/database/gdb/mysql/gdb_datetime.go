package main

import (
	"fmt"
	"github.com/gogf/gf/os/gctx"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

func main() {
	var (
		db  = g.DB()
		ctx = gctx.New()
	)
	db.SetDebug(true)

	r, err := db.Ctx(ctx).Model("user").Data(g.Map{
		"name":        "john",
		"create_time": gtime.Now().String(),
	}).Insert()
	if err == nil {
		fmt.Println(r.LastInsertId())
	} else {
		panic(err)
	}
}

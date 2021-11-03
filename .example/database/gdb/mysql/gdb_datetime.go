package main

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
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

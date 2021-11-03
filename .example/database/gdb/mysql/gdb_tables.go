package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var (
		db  = g.DB()
		ctx = gctx.New()
	)
	db.SetDebug(true)

	tables, err := db.Tables(ctx)
	if err != nil {
		panic(err)
	}
	if tables != nil {
		g.Dump(tables)
	}
}

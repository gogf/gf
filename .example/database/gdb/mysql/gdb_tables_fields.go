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

	tables, e := db.Tables(ctx)
	if e != nil {
		panic(e)
	}
	if tables != nil {
		g.Dump(tables)
		for _, table := range tables {
			fields, err := db.TableFields(ctx, table)
			if err != nil {
				panic(err)
			}
			g.Dump(fields)
		}
	}
}

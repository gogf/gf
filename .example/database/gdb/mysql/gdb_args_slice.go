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

	db.Ctx(ctx).Model("user").
		Where("nickname like ? and passport like ?", g.Slice{"T3", "t3"}).
		OrderAsc("id").All()

	conditions := g.Map{
		"nickname like ?":    "%T%",
		"id between ? and ?": g.Slice{1, 3},
		"id >= ?":            1,
		"create_time > ?":    0,
		"id in(?)":           g.Slice{1, 2, 3},
	}
	db.Ctx(ctx).Model("user").Where(conditions).OrderAsc("id").All()

	var params []interface{}
	db.Ctx(ctx).Model("user").Where("1=1", params).OrderAsc("id").All()
}

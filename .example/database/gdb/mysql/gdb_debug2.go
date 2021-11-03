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

	// 执行3条SQL查询
	for i := 1; i <= 3; i++ {
		db.Ctx(ctx).Model("user").Where("id=?", i).One()
	}
	// 构造一条错误查询
	db.Ctx(ctx).Model("user").Where("no_such_field=?", "just_test").One()

	db.Ctx(ctx).Model("user").Data(g.Map{"name": "smith"}).Where("uid=?", 1).Save()
}

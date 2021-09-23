package main

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gctx"
	"github.com/gogf/gf/os/glog"
)

func main() {
	gdb.AddDefaultConfigNode(gdb.ConfigNode{
		Host:    "127.0.0.1",
		Port:    "3306",
		User:    "root",
		Pass:    "12345678",
		Name:    "test",
		Type:    "mysql",
		Role:    "master",
		Charset: "utf8",
	})
	var (
		ctx = gctx.New()
	)
	db, err := gdb.New()
	if err != nil {
		panic(err)
	}
	//db.SetDebug(false)

	glog.SetPath("/tmp")

	// 执行3条SQL查询
	for i := 1; i <= 3; i++ {
		db.Ctx(ctx).Model("user").Where("uid=?", i).One()
	}
	// 构造一条错误查询
	db.Model("user").Where("no_such_field=?", "just_test").One()

	db.Ctx(ctx).Model("user").Data(g.Map{"name": "smith"}).Where("uid=?", 1).Save()

}

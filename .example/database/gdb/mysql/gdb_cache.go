package main

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/os/gctx"
	"github.com/gogf/gf/util/gutil"
	"time"
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
	//db.GetCache().SetAdapter(adapter.NewRedis(g.Redis()))
	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)

	// 执行2次查询并将查询结果缓存3秒，并可执行缓存名称(可选)
	for i := 0; i < 3; i++ {
		r, _ := db.Ctx(ctx).Model("user").Cache(3000*time.Second).Where("id=?", 1).One()
		gutil.Dump(r.Map())
	}

	// 执行更新操作，并清理指定名称的查询缓存
	//db.Table("user").Cache(-1, "vip-user").Data(gdb.Map{"name": "smith"}).Where("id=?", 1).Update()

	// 再次执行查询，启用查询缓存特性
	//r, _ := db.Table("user").Cache(300000*time.Second, "vip-user").Where("id=?", 1).One()
	//gutil.Dump(r.Map())
}

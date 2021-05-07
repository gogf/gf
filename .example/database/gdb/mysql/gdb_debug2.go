package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	db := g.DB()

	// 执行3条SQL查询
	for i := 1; i <= 3; i++ {
		db.Table("user").Where("id=?", i).One()
	}
	// 构造一条错误查询
	db.Table("user").Where("no_such_field=?", "just_test").One()

	db.Table("user").Data(g.Map{"name": "smith"}).Where("uid=?", 1).Save()
}

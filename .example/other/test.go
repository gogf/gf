package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	db := g.DB()
	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)

	db.Table("user").Delete("score < ", 60)
}

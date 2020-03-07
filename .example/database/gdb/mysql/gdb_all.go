package main

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
)

func main() {
	db := g.DB()
	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)

	r, e := db.Table("test").Order("id asc").All()
	if e != nil {
		fmt.Println(e)
	}
	if r != nil {
		fmt.Println(r.List())
	}
}

package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/frame/g"
)

func main() {
	db := g.DB()
	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)

	fmt.Println(time.Now())
	r, e := db.Table("user").Data(g.Map{
		"create_time": time.Now().Local(),
	}).Insert()
	if e != nil {
		panic(e)
	}
	if r != nil {
		fmt.Println(r.LastInsertId())
	}
}

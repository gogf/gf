package main

import (
	"fmt"
	"github.com/gogf/gf/g"
	"time"
)

func main() {
	db := g.DB()
	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)

	r, e := db.Table("user").Data(g.Map{
		"passport"    : "1",
		"password"    : "1",
		"nickname"    : "1",
		"create_time" : time.Now(),
	}).Insert()
	if e != nil {
		panic(e)
	}
	if r != nil {
		fmt.Println(r.LastInsertId())
	}
}

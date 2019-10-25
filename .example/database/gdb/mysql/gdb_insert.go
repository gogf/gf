package main

import (
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"time"
)

func main() {
	//db := g.DB()

	gdb.AddDefaultConfigNode(gdb.ConfigNode{
		LinkInfo: "root:12345678@tcp(127.0.0.1:3306)/test?parseTime=true&loc=Local",
		Type:     "mysql",
		Charset:  "utf8",
	})
	db, _ := gdb.New()

	db.SetDebug(true)

	type User struct {
		CreateTime time.Time `orm:"create_time"`
	}
	r, e := db.Table("user").Data(User{CreateTime: time.Now()}).Insert()
	if e != nil {
		panic(e)
	}
	if r != nil {
		fmt.Println(r.LastInsertId())
	}
}

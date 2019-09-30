package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

func main() {
	db := g.DB()
	db.SetDebug(true)

	type User struct {
		Id   int
		Name *gtime.Time
	}

	user := new(User)
	e := db.Table("test").Where("id", 10000).Struct(user)
	if e != nil {
		panic(e)
	}
	g.Dump(user)

}

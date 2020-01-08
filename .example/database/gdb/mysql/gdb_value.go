package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	db := g.DB()
	db.SetDebug(true)

	one, e := db.Table("order.order o").LeftJoin("user.user u", "o.uid=u.id").Where("u.id", 1).One()
	if e != nil {
		panic(e)
	}
	g.Dump(one)

}

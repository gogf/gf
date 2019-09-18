package main

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
)

func main() {
	db := g.DB()
	db.SetDebug(true)

	r, e := db.Table("test").Where("id IN (?)", []interface{}{1, 2}).All()
	if e != nil {
		panic(e)
	}
	if r != nil {
		fmt.Println(r.ToList())
	}
}

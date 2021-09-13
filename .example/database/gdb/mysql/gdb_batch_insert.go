package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
)

func main() {
	db := g.DB()
	db.SetDebug(true)
	list := make(g.List, 0)
	for i := 0; i < 100; i++ {
		list = append(list, g.Map{
			"name": fmt.Sprintf(`name_%d`, i),
		})
	}
	r, e := db.Table("user").Data(list).Batch(2).Insert()
	if e != nil {
		panic(e)
	}
	if r != nil {
		fmt.Println(r.LastInsertId())
	}
}

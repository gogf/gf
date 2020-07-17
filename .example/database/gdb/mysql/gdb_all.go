package main

import (
	"fmt"
	"github.com/jin502437344/gf/frame/g"
)

func main() {
	db := g.DB()
	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)

	r, e := db.GetAll("SELECT * from `user` where id in(?)", g.Slice{})
	if e != nil {
		fmt.Println(e)
	}
	if r != nil {
		fmt.Println(r)
	}
	return
	//r, e := db.Table("user").Where("id in(?)", g.Slice{}).All()
	//if e != nil {
	//	fmt.Println(e)
	//}
	//if r != nil {
	//	fmt.Println(r.List())
	//}
}

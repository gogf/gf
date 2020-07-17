package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/os/gtime"
)

func main() {
	db := g.Database()
	db.SetDebug(true)

	//r, err := db.Table("user").Data("create_time", gtime.Now().String()).Insert()
	//if err == nil {
	//    fmt.Println(r.LastInsertId())
	//} else {
	//    panic(err)
	//}

	r, err := db.Table("user").Data(g.Map{
		"name":        "john",
		"create_time": gtime.Now().String(),
	}).Insert()
	if err == nil {
		fmt.Println(r.LastInsertId())
	} else {
		panic(err)
	}
}

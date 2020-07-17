package main

import (
	"github.com/jin502437344/gf/frame/g"
)

func main() {
	db := g.DB()
	db.SetDebug(true)

	db.Table("user").Data("num=num+1").Where("id", 8).Update()
}

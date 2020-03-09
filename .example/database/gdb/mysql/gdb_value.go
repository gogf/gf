package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	db := g.DB()
	db.SetDebug(true)

	db.Table("user").Fields("DISTINCT id,nickname").Filter().All()
}

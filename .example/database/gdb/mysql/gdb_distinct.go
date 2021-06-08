package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.DB().Model("user").Distinct().CountColumn("uid,name")
}

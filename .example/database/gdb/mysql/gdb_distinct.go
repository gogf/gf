package main

import (
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	g.DB().Model("user").Distinct().CountColumn("uid,name")
}

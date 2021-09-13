package main

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Config().SetFileName("config2.toml")
	if r, err := g.DB().Model("user").Where("uid=?", 1).One(); err == nil {
		fmt.Println(r["uid"].Int())
		fmt.Println(r["name"].String())
	} else {
		fmt.Println(err)
	}
}

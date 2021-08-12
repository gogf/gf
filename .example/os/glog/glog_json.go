package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Log().Debug(g.Map{"uid": 100, "name": "john"})

	type User struct {
		Uid  int    `json:"uid"`
		Name string `json:"name"`
	}
	g.Log().Debug(User{100, "john"})
}

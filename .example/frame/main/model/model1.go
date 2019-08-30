package main

import (
	"github.com/gogf/gf/.example/frame/mvc/model/test"
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.DB().SetDebug(true)
	user, err := test.ModelUser().One()
	g.Dump(err)
	g.Dump(user)
	user.Password = "1"
	g.Dump(user.Update())
}

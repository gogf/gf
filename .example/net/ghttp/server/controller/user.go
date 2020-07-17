package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/frame/gmvc"
)

type User struct {
	gmvc.Controller
}

func (c *User) Index() {
	c.View.Display("index.html")
}

// 不符合规范，不会被自动注册
func (c *User) Test(value interface{}) {
	c.View.Display("index.html")
}

func main() {
	//g.View().SetPath("C:/www/static")
	s := g.Server()
	s.BindController("/user", new(User))
	s.SetPort(8199)
	s.Run()
}

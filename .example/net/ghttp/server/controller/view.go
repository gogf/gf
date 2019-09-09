package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

type Controller struct {
	gmvc.Controller
}

func (c *Controller) Index() {
	c.View.Display("layout.html")
}

// 不符合规范，不会被自动注册
func (c *Controller) Test(value interface{}) {
	c.View.Display("layout.html")
}

func main() {
	s := g.Server()
	s.BindController("/view", new(Controller))
	s.SetPort(8199)
	s.Run()
}

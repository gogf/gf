package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

type Controller struct {
	gmvc.Controller
}

func (c *Controller) Test() {
	c.View.Display("layout.html")
}
func main() {
	s := g.Server()
	s.BindControllerMethod("/", new(Controller), "Test")
	s.SetPort(8199)
	s.Run()
}

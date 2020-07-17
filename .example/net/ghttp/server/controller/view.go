package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/frame/gmvc"
)

type Controller struct {
	gmvc.Controller
}

func (c *Controller) Index() {
	c.View.Assign("name", "john")
	c.View.Assign("mainTpl", "main/main2.html")
	c.View.Display("layout.html")
}

func main() {
	s := g.Server()
	s.BindController("/view", new(Controller))
	s.SetPort(8199)
	s.Run()
}

package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
	"github.com/gogf/gf/net/ghttp"
)

type ControllerIndex struct {
	gmvc.Controller
}

func (c *ControllerIndex) Info() {
	c.View.Assign("title", "Go Frame 第一个网站")
	c.View.Assigns(g.Map{
		"name":  "很开心1",
		"score": 100,
	})
	c.View.Display("index.html")
}
func main() {
	s := ghttp.GetServer()
	s.BindController("/", new(ControllerIndex))
	s.SetPort(8199)
	s.Run()
}

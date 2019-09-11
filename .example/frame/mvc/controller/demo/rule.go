package demo

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

type ControllerRule struct {
	gmvc.Controller
}

func init() {
	g.Server().BindController("/rule/{method}/:name", &ControllerRule{})
}

func (c *ControllerRule) Show() {
	c.Response.Write(c.Request.Get("name"))
}

package demo

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/frame/gmvc"
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

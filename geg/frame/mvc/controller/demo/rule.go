package demo

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/frame/gmvc"
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


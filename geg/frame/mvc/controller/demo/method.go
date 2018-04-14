package demo

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gmvc"
)

type ControllerMethod struct {
    gmvc.Controller
}

func init() {
    ghttp.GetServer().BindControllerMethod("/method", &ControllerMethod{}, "Name, Age")
}

func (c *ControllerMethod) Name() {
    c.Response.Write("John")
}

func (c *ControllerMethod) Age() {
    c.Response.Write("18")
}

func (c *ControllerMethod) Info() {
    c.Response.Write("Info")
}




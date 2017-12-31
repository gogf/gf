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
    c.Response.WriteString("John")
}

func (c *ControllerMethod) Age() {
    c.Response.WriteString("18")
}

func (c *ControllerMethod) Info() {
    c.Response.WriteString("Info")
}




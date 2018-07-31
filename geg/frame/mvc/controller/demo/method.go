package demo

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/frame/gmvc"
)

type Method struct {
    gmvc.Controller
}

func init() {
    g.Server().BindControllerMethod("/method", &Method{}, "Name, Age")
}

func (c *Method) Name() {
    c.Response.Write("John")
}

func (c *Method) Age() {
    c.Response.Write("18")
}

func (c *Method) Info() {
    c.Response.Write("Info")
}




package demo

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gmvc"
    "gitee.com/johng/gf/g/frame/gins"
)

type ControllerTemplate struct {
    gmvc.Controller
}

func init() {
    gins.View().
    ghttp.GetServer().BindController("/template", &ControllerTemplate{})
}

func (c *ControllerTemplate) Info() {
    c.View.Assign("name", "john")
    c.View.Display("user/index")
}




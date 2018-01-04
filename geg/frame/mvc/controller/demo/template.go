package demo

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gmvc"
    "gitee.com/johng/gf/g/frame/gins"
    "gitee.com/johng/gf/g/os/gview"
)

type ControllerTemplate struct {
    gmvc.Controller
}

func init() {
    ghttp.GetServer().BindHandler("/template/handler-info", Info)
    ghttp.GetServer().BindControllerMethod("/template/ctl-info", &ControllerTemplate{}, "Info")
}

func Info(r *ghttp.Request) {
    gins.View().SetPath("")
}

func (c *ControllerTemplate) Info() {
    c.View.Assign("name", "john")
    c.View.Display("user/index")
}




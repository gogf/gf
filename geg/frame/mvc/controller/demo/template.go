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
    ghttp.GetServer().BindControllerMethod("/template", &ControllerTemplate{}, "Info")
}

func (c *ControllerTemplate) Info() {
    c.View.Assign("name", "john")
    c.View.Assigns(map[string]interface{}{
        "age"   : 18,
        "score" : 100,
    })
    c.View.Display("user/index.tpl")
}




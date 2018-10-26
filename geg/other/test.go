package main

import (
    "gitee.com/johng/gf/g/frame/gmvc"
    "gitee.com/johng/gf/g/net/ghttp"
)
type ControllerIndex struct {
    gmvc.Controller
}
func (c *ControllerIndex) Info() {
    c.View.Assign("title", "Go Frame 第一个网站")
    c.View.Assigns(map[string]interface{}{
        "name"   : "很开心",
        "score" : 100,
    })
    c.View.Display("web/index.html")
}
func main() {
    s := ghttp.GetServer()
    s.SetServerRoot("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/other/web/")
    s.BindController("/", new(ControllerIndex))
    s.SetPort(8199)
    s.Run()
}
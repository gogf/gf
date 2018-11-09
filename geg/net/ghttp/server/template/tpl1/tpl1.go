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
        "name"   : "很开心1",
        "score" : 100,
    })
    c.View.Display("index.html")
}
func main() {
    s := ghttp.GetServer()
    s.BindController("/", new(ControllerIndex))
    s.SetPort(8199)
    s.Run()
}

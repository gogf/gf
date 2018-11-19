package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

// 允许对同一个路由同一个事件注册多个回调函数，按照注册顺序进行优先级调用
func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Writeln(r.GetParam("name").String())
        r.Response.Writeln(r.GetParam("site").String())
    })
    s.BindHookHandler("/", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
        r.SetParam("name", "GoFrame")
        r.Response.Writeln("set name")
    })
    s.BindHookHandler("/", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
        r.SetParam("site", "https://gfer.me")
        r.Response.Writeln("set site")
    })
    s.SetPort(8199)
    s.Run()
}
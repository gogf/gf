package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/frame/gmvc"
    "gitee.com/johng/gf/g/net/ghttp"
)

type Object struct {}

func (o *Object) Show(r *ghttp.Request) {
    r.Response.Writeln("Object Show")
}

func (o *Object) Delete(r *ghttp.Request) {
    r.Response.Writeln("Object REST Delete")
}

func (o *Object) Shut(r *ghttp.Request) {
    r.Response.Writeln("Object Shut")
}

type Controller struct {
    gmvc.Controller
}

func (c *Controller) Show() {
    c.Response.Writeln("Controller Show")
}

func (c *Controller) Post() {
    c.Response.Writeln("Controller REST Post")
}

func (c *Controller) Shut() {
    c.Response.Writeln("Controller Shut")
}

func Handler(r *ghttp.Request) {
    r.Response.Writeln("Handler")
}

func HookHandler(r *ghttp.Request) {
    r.Response.Writeln("Hook Handler")
}

func main() {
    s   := g.Server()
    obj := new(Object)
    ctl := new(Controller)

    // 分组路由方法注册
    //g := s.Group("/api")
    //g.ALL ("*",            HookHandler, ghttp.HOOK_BEFORE_SERVE)
    //g.ALL ("/handler",     Handler)
    //g.ALL ("/ctl",         ctl)
    //g.GET ("/ctl/my-show", ctl, "Show")
    //g.REST("/ctl/rest",    ctl)
    //g.ALL ("/obj",         obj)
    //g.GET ("/obj/my-show", obj, "Show")
    //g.REST("/obj/rest",    obj)

    // 分组路由批量注册
    s.Group("/api").Bind("/api", []ghttp.GroupItem{

        {"ALL",  "/handler",     Handler},
        {"ALL",  "/ctl",         ctl},
        {"GET",  "/ctl/my-show", ctl, "Show"},
        {"REST", "/ctl/rest",    ctl},
        {"ALL",  "/obj",         obj},
        {"GET",  "/obj/my-show", obj, "Show"},
        {"REST", "/obj/rest",    obj},
        {"ALL",  "*",            HookHandler, ghttp.HOOK_BEFORE_SERVE},
    })

    s.SetPort(8199)
    s.Run()
}

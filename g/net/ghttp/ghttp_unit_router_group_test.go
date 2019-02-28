// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// 分组路由测试
package ghttp_test

import (
    "fmt"
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/frame/gmvc"
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
    "time"
)

// 执行对象
type Object struct {}

func (o *Object) Init(r *ghttp.Request) {
    r.Response.Write("1")
}

func (o *Object) Shut(r *ghttp.Request) {
    r.Response.Write("2")
}

func (o *Object) Index(r *ghttp.Request) {
    r.Response.Write("Object Index")
}

func (o *Object) Show(r *ghttp.Request) {
    r.Response.Write("Object Show")
}

func (o *Object) Delete(r *ghttp.Request) {
    r.Response.Write("Object REST Delete")
}

// 控制器
type Controller struct {
    gmvc.Controller
}

func (c *Controller) Init(r *ghttp.Request) {
    c.Controller.Init(r)
    c.Response.Write("1")
}

func (c *Controller) Shut() {
    c.Response.Write("2")
}

func (c *Controller) Index() {
    c.Response.Write("Controller Index")
}

func (c *Controller) Show() {
    c.Response.Write("Controller Show")
}

func (c *Controller) Post() {
    c.Response.Write("Controller REST Post")
}

func Handler(r *ghttp.Request) {
    r.Response.Write("Handler")
}

func Test_Router_Group1(t *testing.T) {
    p   := ports.PopRand()
    s   := g.Server(p)
    obj := new(Object)
    ctl := new(Controller)
    // 分组路由方法注册
    g := s.Group("/api")
    g.ALL ("/handler",     Handler)
    g.ALL ("/ctl",         ctl)
    g.GET ("/ctl/my-show", ctl, "Show")
    g.REST("/ctl/rest",    ctl)
    g.ALL ("/obj",         obj)
    g.GET ("/obj/my-show", obj, "Show")
    g.REST("/obj/rest",    obj)
    s.SetPort(p)
    s.SetDumpRouteMap(false)
    s.Start()
    defer s.Shutdown()

    time.Sleep(time.Second)
    gtest.Case(t, func() {
        client := ghttp.NewClient()
        client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

        gtest.Assert(client.GetContent ("/api/handler"),     "Handler")

        gtest.Assert(client.GetContent ("/api/ctl"),         "1Controller Index2")
        gtest.Assert(client.GetContent ("/api/ctl/"),        "1Controller Index2")
        gtest.Assert(client.GetContent ("/api/ctl/index"),   "1Controller Index2")
        gtest.Assert(client.GetContent ("/api/ctl/my-show"), "1Controller Show2")
        gtest.Assert(client.GetContent ("/api/ctl/post"),    "1Controller REST Post2")
        gtest.Assert(client.GetContent ("/api/ctl/show"),    "1Controller Show2")
        gtest.Assert(client.PostContent("/api/ctl/rest"),    "1Controller REST Post2")

        gtest.Assert(client.GetContent ("/api/obj"),         "1Object Index2")
        gtest.Assert(client.GetContent ("/api/obj/"),        "1Object Index2")
        gtest.Assert(client.GetContent ("/api/obj/index"),   "1Object Index2")
        gtest.Assert(client.GetContent ("/api/obj/delete"),  "1Object REST Delete2")
        gtest.Assert(client.GetContent ("/api/obj/my-show"), "1Object Show2")
        gtest.Assert(client.GetContent ("/api/obj/show"),    "1Object Show2")
        gtest.Assert(client.DeleteContent("/api/obj/rest"),  "1Object REST Delete2")

        // 测试404
        resp, err := client.Get("/ThisDoesNotExist")
        defer resp.Close()
        gtest.Assert(err,             nil)
        gtest.Assert(resp.StatusCode, 404)
    })
}

func Test_Router_Group2(t *testing.T) {
    p   := ports.PopRand()
    s   := g.Server(p)
    obj := new(Object)
    ctl := new(Controller)
    // 分组路由批量注册
    s.Group("/api").Bind("/api", []ghttp.GroupItem{
        {"ALL",  "/handler",     Handler},
        {"ALL",  "/ctl",         ctl},
        {"GET",  "/ctl/my-show", ctl, "Show"},
        {"REST", "/ctl/rest",    ctl},
        {"ALL",  "/obj",         obj},
        {"GET",  "/obj/my-show", obj, "Show"},
        {"REST", "/obj/rest",    obj},
    })
    s.SetPort(p)
    s.SetDumpRouteMap(false)
    s.Start()
    defer s.Shutdown()

    time.Sleep(time.Second)
    gtest.Case(t, func() {
        client := ghttp.NewClient()
        client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

        gtest.Assert(client.GetContent ("/api/handler"),     "Handler")

        gtest.Assert(client.GetContent ("/api/ctl/my-show"), "1Controller Show2")
        gtest.Assert(client.GetContent ("/api/ctl/post"),    "1Controller REST Post2")
        gtest.Assert(client.GetContent ("/api/ctl/show"),    "1Controller Show2")
        gtest.Assert(client.PostContent("/api/ctl/rest"),    "1Controller REST Post2")

        gtest.Assert(client.GetContent ("/api/obj/delete"),  "1Object REST Delete2")
        gtest.Assert(client.GetContent ("/api/obj/my-show"), "1Object Show2")
        gtest.Assert(client.GetContent ("/api/obj/show"),    "1Object Show2")
        gtest.Assert(client.DeleteContent("/api/obj/rest"),  "1Object REST Delete2")

        // 测试404
        resp, err := client.Get("/ThisDoesNotExist")
        defer resp.Close()
        gtest.Assert(err,             nil)
        gtest.Assert(resp.StatusCode, 404)
    })
}

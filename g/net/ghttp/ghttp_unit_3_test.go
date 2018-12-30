// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 分组路由测试
package ghttp_test

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/frame/gmvc"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

// 执行对象
type Object struct {}

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
    s   := g.Server(gtime.Nanosecond())
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
    s.SetPort(8199)
    s.SetDumpRouteMap(false)
    go s.Run()
    defer s.Shutdown()
    time.Sleep(time.Second)
    gtest.Case(func() {
        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/handler"),     "Handler")

        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/ctl/my-show"), "Controller Show")
        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/ctl/post"),    "Controller REST Post")
        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/ctl/show"),    "Controller Show")
        gtest.Assert(ghttp.PostContent("http://127.0.0.1:8199/api/ctl/rest"),    "Controller REST Post")

        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/obj/delete"),  "Object REST Delete")
        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/obj/my-show"), "Object Show")
        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/obj/show"),    "Object Show")
        gtest.Assert(ghttp.DeleteContent("http://127.0.0.1:8199/api/obj/rest"),  "Object REST Delete")

    })
}

func Test_Router_Group2(t *testing.T) {
    s   := g.Server(gtime.Nanosecond())
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
    s.SetPort(8199)
    s.SetDumpRouteMap(false)
    go s.Run()
    defer s.Shutdown()
    time.Sleep(time.Second)
    gtest.Case(func() {
        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/handler"),     "Handler")

        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/ctl/my-show"), "Controller Show")
        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/ctl/post"),    "Controller REST Post")
        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/ctl/show"),    "Controller Show")
        gtest.Assert(ghttp.PostContent("http://127.0.0.1:8199/api/ctl/rest"),    "Controller REST Post")

        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/obj/delete"),  "Object REST Delete")
        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/obj/my-show"), "Object Show")
        gtest.Assert(ghttp.GetContent ("http://127.0.0.1:8199/api/obj/show"),    "Object Show")
        gtest.Assert(ghttp.DeleteContent("http://127.0.0.1:8199/api/obj/rest"),  "Object REST Delete")
    })
}

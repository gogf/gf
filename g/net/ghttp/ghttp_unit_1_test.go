// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package ghttp_test

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

// 基本路由功能以及优先级测试
func Test_Router_Basic(t *testing.T) {
    s := g.Server(gtime.Nanosecond())
    s.BindHandler("/:name", func(r *ghttp.Request){
        r.Response.Write("/:name")
    })
    s.BindHandler("/:name/update", func(r *ghttp.Request){
        r.Response.Write(r.Get("name"))
    })
    s.BindHandler("/:name/:action", func(r *ghttp.Request){
        r.Response.Write(r.Get("action"))
    })
    s.BindHandler("/:name/*any", func(r *ghttp.Request){
        r.Response.Write(r.Get("any"))
    })
    s.BindHandler("/user/list/{field}.html", func(r *ghttp.Request){
        r.Response.Write(r.Get("field"))
    })
    s.SetPort(8199)
    s.SetDumpRouteMap(false)
    go s.Run()
    defer s.Shutdown()
    // 等待启动完成
    time.Sleep(time.Second)
    gtest.Case(func() {
        gtest.Assert(ghttp.GetContent("http://127.0.0.1:8199/john"),               "")
        gtest.Assert(ghttp.GetContent("http://127.0.0.1:8199/john/update"),        "john")
        gtest.Assert(ghttp.GetContent("http://127.0.0.1:8199/john/edit"),          "edit")
        gtest.Assert(ghttp.GetContent("http://127.0.0.1:8199/user/list/100.html"), "100")
    })
}

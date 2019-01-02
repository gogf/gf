// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 基本路由功能以及优先级测试
package ghttp_test

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)


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
    s.SetPort(8100)
    s.SetDumpRouteMap(false)
    go s.Run()
    defer func() {
        s.Shutdown()
        time.Sleep(time.Second)
    }()
    // 等待启动完成
    time.Sleep(time.Second)
    gtest.Case(t, func() {
        client := ghttp.NewClient()
        client.SetPrefix("http://127.0.0.1:8100")
        
        gtest.Assert(client.GetContent("/john"),               "")
        gtest.Assert(client.GetContent("/john/update"),        "john")
        gtest.Assert(client.GetContent("/john/edit"),          "edit")
        gtest.Assert(client.GetContent("/user/list/100.html"), "100")
    })
}

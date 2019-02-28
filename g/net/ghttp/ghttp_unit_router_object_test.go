// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
    "fmt"
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
    "time"
)

// 执行对象
type ObjectBasic struct {}

func (o *ObjectBasic) Init(r *ghttp.Request) {
    r.Response.Write("1")
}

func (o *ObjectBasic) Shut(r *ghttp.Request) {
    r.Response.Write("2")
}

func (o *ObjectBasic) Index(r *ghttp.Request) {
    r.Response.Write("Object Index")
}

func (o *ObjectBasic) Show(r *ghttp.Request) {
    r.Response.Write("Object Show")
}

// 执行对象注册
func Test_Router_Object(t *testing.T) {
    p := ports.PopRand()
    s := g.Server(p)
    s.BindObject("/", new(ObjectBasic))
    s.SetPort(p)
    s.SetDumpRouteMap(false)
    s.Start()
    defer s.Shutdown()

    // 等待启动完成
    time.Sleep(time.Second)
    gtest.Case(t, func() {
        client := ghttp.NewClient()
        client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

        gtest.Assert(client.GetContent("/"),            "1Object Index2")
        gtest.Assert(client.GetContent("/init"),        "Not Found")
        gtest.Assert(client.GetContent("/shut"),        "Not Found")
        gtest.Assert(client.GetContent("/index"),       "1Object Index2")
        gtest.Assert(client.GetContent("/show"),        "1Object Show2")
        gtest.Assert(client.GetContent("/none-exist"),  "Not Found")
    })
}

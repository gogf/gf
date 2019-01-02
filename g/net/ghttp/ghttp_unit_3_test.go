// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 请求参数测试
package ghttp_test

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

func Test_Params(t *testing.T) {
    type User struct {
        Id    int
        Name  string
        Pass1 string `params:"password1"`
        Pass2 string `params:"password2"`
    }
    s := g.Server(gtime.Nanosecond())
    s.BindHandler("/get", func(r *ghttp.Request){
        if r.GetQuery("slice") != nil {
            r.Response.Write(r.GetQuery("slice"))
        }
        if r.GetQuery("bool") != nil {
            r.Response.Write(r.GetQueryBool("bool"))
        }
        if r.GetQuery("float32") != nil {
            r.Response.Write(r.GetQueryFloat32("float32"))
        }
        if r.GetQuery("float64") != nil {
            r.Response.Write(r.GetQueryFloat64("float64"))
        }
        if r.GetQuery("int") != nil {
            r.Response.Write(r.GetQueryInt("int"))
        }
        if r.GetQuery("uint") != nil {
            r.Response.Write(r.GetQueryUint("uint"))
        }
        if r.GetQuery("string") != nil {
            r.Response.Write(r.GetQueryString("string"))
        }
    })
    s.BindHandler("/post", func(r *ghttp.Request){
        if r.GetPost("slice") != nil {
            r.Response.Write(r.GetPost("slice"))
        }
        if r.GetPost("bool") != nil {
            r.Response.Write(r.GetPostBool("bool"))
        }
        if r.GetPost("float32") != nil {
            r.Response.Write(r.GetPostFloat32("float32"))
        }
        if r.GetPost("float64") != nil {
            r.Response.Write(r.GetPostFloat64("float64"))
        }
        if r.GetPost("int") != nil {
            r.Response.Write(r.GetPostInt("int"))
        }
        if r.GetPost("uint") != nil {
            r.Response.Write(r.GetPostUint("uint"))
        }
        if r.GetPost("string") != nil {
            r.Response.Write(r.GetPostString("string"))
        }
    })
    s.BindHandler("/map", func(r *ghttp.Request){
        if m := r.GetQueryMap(); len(m) > 0 {
            r.Response.Write(m["name"])
        }
        if m := r.GetPostMap(); len(m) > 0 {
            r.Response.Write(m["name"])
        }
    })
    s.BindHandler("/raw", func(r *ghttp.Request){
        r.Response.Write(r.GetRaw())
    })
    s.BindHandler("/json", func(r *ghttp.Request){
        r.Response.Write(r.GetJson().Get("name"))
    })
    s.BindHandler("/struct", func(r *ghttp.Request){
        if m := r.GetQueryMap(); len(m) > 0 {
            user := new(User)
            r.GetQueryToStruct(user)
            r.Response.Write(user.Id, user.Name, user.Pass1, user.Pass2)
        }
        if m := r.GetPostMap(); len(m) > 0 {
            user := new(User)
            r.GetPostToStruct(user)
            r.Response.Write(user.Id, user.Name, user.Pass1, user.Pass2)
        }
    })
    s.SetPort(8400)
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
        client.SetPrefix("http://127.0.0.1:8400")
        // GET
        gtest.Assert(client.GetContent("/get", "slice=1&slice=2"), `["1","2"]`)
        gtest.Assert(client.GetContent("/get", "bool=1"),          `true`)
        gtest.Assert(client.GetContent("/get", "bool=0"),          `false`)
        gtest.Assert(client.GetContent("/get", "float32=0.11"),    `0.11`)
        gtest.Assert(client.GetContent("/get", "float64=0.22"),    `0.22`)
        gtest.Assert(client.GetContent("/get", "int=-10000"),      `-10000`)
        gtest.Assert(client.GetContent("/get", "int=10000"),       `10000`)
        gtest.Assert(client.GetContent("/get", "uint=-10000"),     `10000`)
        gtest.Assert(client.GetContent("/get", "uint=9"),          `9`)
        gtest.Assert(client.GetContent("/get", "string=key"),      `key`)

        // POST
        gtest.Assert(client.PostContent("/post", "slice=1&slice=2"), `["1","2"]`)
        gtest.Assert(client.PostContent("/post", "bool=1"),          `true`)
        gtest.Assert(client.PostContent("/post", "bool=0"),          `false`)
        gtest.Assert(client.PostContent("/post", "float32=0.11"),    `0.11`)
        gtest.Assert(client.PostContent("/post", "float64=0.22"),    `0.22`)
        gtest.Assert(client.PostContent("/post", "int=-10000"),      `-10000`)
        gtest.Assert(client.PostContent("/post", "int=10000"),       `10000`)
        gtest.Assert(client.PostContent("/post", "uint=-10000"),     `10000`)
        gtest.Assert(client.PostContent("/post", "uint=9"),          `9`)
        gtest.Assert(client.PostContent("/post", "string=key"),      `key`)

        // Map
        gtest.Assert(client.GetContent ("/map",  "id=1&name=john"), `john`)
        gtest.Assert(client.PostContent("/map",  "id=1&name=john"), `john`)

        // Raw
        gtest.Assert(client.PutContent("/raw",   "id=1&name=john"), `id=1&name=john`)

        // Json
        gtest.Assert(client.PostContent("/json", `{"id":1, "name":"john"}`), `john`)

        // Struct
        gtest.Assert(client.GetContent("/struct",  `id=1&name=john&password1=123&password2=456`), `1john123456`)
        gtest.Assert(client.PostContent("/struct", `id=1&name=john&password1=123&password2=456`), `1john123456`)
    })
}

package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    // 基本事件回调使用
    p := "/:name/info/{uid}"
    s := ghttp.GetServer()
    s.BindHookHandlerByMap(p, map[string]ghttp.HandlerFunc{
        "BeforeServe"  : func(r *ghttp.Request){ fmt.Println("BeforeServe") },
        "AfterServe"   : func(r *ghttp.Request){ fmt.Println("AfterServe") },
        "BeforeOutput" : func(r *ghttp.Request){ fmt.Println("BeforeOutput") },
        "AfterOutput"  : func(r *ghttp.Request){ fmt.Println("AfterOutput") },
        "BeforeClose"  : func(r *ghttp.Request){ fmt.Println("BeforeClose") },
        "AfterClose"   : func(r *ghttp.Request){ fmt.Println("AfterClose") },
    })
    s.BindHandler(p, func(r *ghttp.Request) {
       r.Response.Write("用户:", r.Get("name"), ", uid:", r.Get("uid"))
    })
    s.SetPort(8199)
    s.Run()
}
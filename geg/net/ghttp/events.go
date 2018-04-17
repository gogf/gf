package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    p := "/"
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
        r.Response.Write("哈喽世界！")
    })
    s.SetPort(8199)
    s.Run()
}
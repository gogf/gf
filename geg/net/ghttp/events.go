package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "fmt"
)

func main() {
    pattern := "/"
    ghttp.GetServer().BindHookHandlerByMap(pattern, map[string]ghttp.HandlerFunc{
        "BeforeServe"        : func(r *ghttp.Request){ fmt.Println("BeforeServe") },
        "AfterServe"         : func(r *ghttp.Request){ fmt.Println("AfterServe") },
        "BeforeRouterPatch"  : func(r *ghttp.Request){ fmt.Println("BeforeRouterPatch") },
        "AfterRouterPatch"   : func(r *ghttp.Request){ fmt.Println("AfterRouterPatch") },
        "BeforeCookieOutput" : func(r *ghttp.Request){ fmt.Println("BeforeCookieOutput") },
        "AfterCookieOutput"  : func(r *ghttp.Request){ fmt.Println("AfterCookieOutput") },
        "BeforeBufferOutput" : func(r *ghttp.Request){ fmt.Println("BeforeBufferOutput") },
        "AfterBufferOutput"  : func(r *ghttp.Request){ fmt.Println("AfterBufferOutput") },
        "BeforeRequestClose" : func(r *ghttp.Request){ fmt.Println("BeforeRequestClose") },
        "AfterRequestClose"  : func(r *ghttp.Request){ fmt.Println("AfterRequestClose") },
    })
    ghttp.GetServer().BindHandler(pattern, func(r *ghttp.Request) {
        r.Response.WriteString("Hello World!")
    })
    ghttp.GetServer().SetPort(10000)
    ghttp.GetServer().Run()

    select { }
}
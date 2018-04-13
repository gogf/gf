package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "fmt"
)

func main() {
    pattern := "/:name/action"
    ghttp.GetServer().BindHookHandlerByMap(pattern, map[string]ghttp.HandlerFunc{
        "BeforeServe"  : func(r *ghttp.Request){ fmt.Println("BeforeServe") },
        "AfterServe"   : func(r *ghttp.Request){ fmt.Println("AfterServe") },
        "BeforePatch"  : func(r *ghttp.Request){ fmt.Println("BeforePatch") },
        "AfterPatch"   : func(r *ghttp.Request){ fmt.Println("AfterPatch") },
        "BeforeOutput" : func(r *ghttp.Request){ fmt.Println("BeforeOutput") },
        "AfterOutput"  : func(r *ghttp.Request){ fmt.Println("AfterOutput") },
        "BeforeClose"  : func(r *ghttp.Request){ fmt.Println("BeforeClose") },
        "AfterClose"   : func(r *ghttp.Request){ fmt.Println("AfterClose") },
    })
    ghttp.GetServer().BindHandler(pattern, func(r *ghttp.Request) {
        r.Response.WriteString("Hello World!")
    })
    ghttp.GetServer().SetPort(10000)
    ghttp.GetServer().Run()

}
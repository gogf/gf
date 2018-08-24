package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g"
    "fmt"
)

func main() {
    s := g.Server()
    s.BindHookHandler("/*any", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
        fmt.Println(r.Router)
        fmt.Println(r.Get("customer_id"))
    })
    s.BindHandler("/admin/customer/{customer_id}/edit", func(r *ghttp.Request){
        r.Response.Writeln(r.Get("customer_id"))
        r.Response.Writeln(r.Router.Uri)
    })
    s.BindHandler("/admin/customer/{customer_id}/disable", func(r *ghttp.Request){
        r.Response.Writeln(r.Get("customer_id"))
        r.Response.Writeln(r.Router.Uri)
    })
    s.SetPort(8199)
    s.Run()
}
package main

import "gitee.com/johng/gf/g/net/ghttp"

func main() {
    ghttp.GetServer().BindHandler("/:name/*any", func(r *ghttp.Request) {
        r.Response.WriteString("any")
        r.Response.WriteString(r.GetQueryString("name"))
        r.Response.WriteString(r.GetQueryString("any"))
    })
    //ghttp.GetServer().BindHandler("/:name/action", func(r *ghttp.Request) {
    //    r.Response.WriteString(r.GetQueryString("name"))
    //})
    ghttp.GetServer().BindHandler("/:name/:action/:aaa", func(r *ghttp.Request) {
        r.Response.WriteString("name")
        r.Response.WriteString(r.GetQueryString("name"))
        r.Response.WriteString(r.GetQueryString("action"))
    })
    ghttp.GetServer().SetPort(10000)
    ghttp.GetServer().Run()
}
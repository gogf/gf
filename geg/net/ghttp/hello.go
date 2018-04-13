package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    ghttp.GetServer().BindHandler("get:/h", func(r *ghttp.Request) {
        //r.Response.RedirectTo("http://www.baidu.com/")
        r.Response.Write("hello world")
        //r.Response.WriteStatus(302)
    })

    ghttp.GetServer().BindHandler("/:name/*any", func(r *ghttp.Request) {
        r.Response.Write("any")
        r.Response.Write(r.GetQueryString("name"))
        r.Response.Write(r.GetQueryString("any"))
    })
    //ghttp.GetServer().BindHandler("/:name/action", func(r *ghttp.Request) {
    //    r.Response.WriteString(r.GetQueryString("name"))
    //})
    ghttp.GetServer().BindHandler("/:name/:action/:aaa", func(r *ghttp.Request) {
        r.Response.Write("name")
        r.Response.Write(r.GetQueryString("name"))
        r.Response.Write(r.GetQueryString("action"))
    })
    ghttp.GetServer().SetPort(10000)
    ghttp.GetServer().Run()
}
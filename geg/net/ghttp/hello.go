package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    ghttp.GetServer().BindHandler("/", func(r *ghttp.Request) {
        //r.Response.RedirectTo("http://www.baidu.com/")
        r.Response.Write("哈喽世界！")
        //r.Response.WriteStatus(302)
    })

    //ghttp.GetServer().BindHandler("/:name/*any", func(r *ghttp.Request) {
    //    r.Response.Write("any")
    //    r.Response.Write(r.Get("name"))
    //    r.Response.Write(r.Get("any"))
    //})
    ////ghttp.GetServer().BindHandler("/:name/action", func(r *ghttp.Request) {
    ////    r.Response.Write(r.Get("name"))
    ////})
    //ghttp.GetServer().BindHandler("/:name/:action/:aaa", func(r *ghttp.Request) {
    //    r.Response.Write("name")
    //    r.Response.Write(r.Get("name"))
    //    r.Response.Write(r.Get("action"))
    //})
    ghttp.GetServer().SetPort(10000)
    ghttp.GetServer().Run()
}
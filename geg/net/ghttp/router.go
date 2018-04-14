package main

import "gitee.com/johng/gf/g/net/ghttp"

func main () {
    ghttp.GetServer().BindHandler("/router/*name", func(r *ghttp.Request) {
        r.Response.Write(r.Get("name"))
    })
    ghttp.GetServer().SetPort(10000)
    ghttp.GetServer().Run()
}
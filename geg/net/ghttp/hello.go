package main

import "gitee.com/johng/gf/g/net/ghttp"

func init() {
    ghttp.GetServer().BindHandler("/", func(r *ghttp.Request) {
        r.Response.WriteString("Hello World!")
    })
}
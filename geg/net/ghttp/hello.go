package main

import "gitee.com/johng/gf/g/net/ghttp"

func main() {
    ghttp.GetServer().BindHandler("/:name", func(r *ghttp.Request) {
        r.Response.WriteString("Hello World!")
    })
    ghttp.GetServer().SetPort(10000)
    ghttp.GetServer().Run()
}
package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)


func main() {
    s := g.Server()
    s.Domain("www.a.com").BindHandler("/*", func(r *ghttp.Request) {
        r.Response.ServeFile("/home/john/www1" + r.URL.Path)
    })
    s.Domain("www.b.com").BindHandler("/*", func(r *ghttp.Request) {
        r.Response.ServeFile("/home/john/www2" + r.URL.Path)
    })
    s.SetIndexFolder(true)
    s.SetPort(8080)
    s.Run()
}
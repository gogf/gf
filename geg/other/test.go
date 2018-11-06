package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    g.Server().BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write(r.GetInt("amount"))
    })
    g.Server().SetPort(8199)
    g.Server().Run()
}
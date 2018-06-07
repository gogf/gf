package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/:name", func(r *ghttp.Request){
        r.Response.Write(r.Get("name"))
    })
    s.SetPort(8199)
    s.Run()
}
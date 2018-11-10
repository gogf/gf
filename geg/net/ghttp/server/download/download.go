package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/download", func(r *ghttp.Request){
        r.Response.ServeFile("text.txt")
    })
    s.SetPort(8199)
    s.Run()
}
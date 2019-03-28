package main

import (
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/os/gproc"
)

func main() {
    s := g.Server()
    s.SetIndexFolder(true)
    s.BindHandler("/", func(r *ghttp.Request){
        r.Response.Write("pid:", gproc.Pid())
    })
    s.BindHandler("/panic", func(r *ghttp.Request){
        panic("error")
    })
    s.SetAccessLogEnabled(true)
    s.SetPort(8199)
    s.Run()
}
package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request){
        r.Response.Writeln("hello")
    })
    s.BindHandler("/restart", func(r *ghttp.Request){
        r.Response.Writeln("restart server")
        r.Server.Restart()
    })
    s.BindHandler("/shutdown", func(r *ghttp.Request){
        r.Response.Writeln("shutdown server")
        r.Server.Shutdown()
    })
    s.SetPort(8199, 8200)
    s.Run()
}
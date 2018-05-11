package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s1 := g.Server("s1")
    s1.BindHandler("/", func(r *ghttp.Request){
        r.Response.Writeln("hello s1")
    })
    s1.BindHandler("/restart", func(r *ghttp.Request){
        r.Response.Writeln("restart server")
        r.Server.Restart()
    })
    s1.BindHandler("/shutdown", func(r *ghttp.Request){
        r.Response.Writeln("shutdown server")
        r.Server.Shutdown()
    })
    s1.SetPort(8199, 8200)
    go s1.Run()

    s2 := g.Server("s2")
    s2.BindHandler("/", func(r *ghttp.Request){
        r.Response.Writeln("hello s2")
    })
    s2.SetPort(8300, 8080)
    go s2.Run()

    ghttp.Wait()
}
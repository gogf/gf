package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/os/gtime"
    "time"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request){
        r.Response.Writeln("hello")
    })
    s.BindHandler("/restart", func(r *ghttp.Request){
        r.Response.Writeln("restart server in 2 seconds")
        gtime.SetTimeout(2*time.Second, func() {
            r.Server.Restart()
        })
    })
    s.SetPort(8199)
    s.Run()
}
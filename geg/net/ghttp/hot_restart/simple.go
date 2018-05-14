package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/os/gproc"
    "time"
)

func main() {
    s := g.Server()
    s.BindHandler("/sleep", func(r *ghttp.Request){
        r.Response.Writeln(gproc.Pid())
        time.Sleep(10*time.Second)
        r.Response.Writeln(gproc.Pid())
    })
    s.BindHandler("/pid", func(r *ghttp.Request){
        r.Response.Writeln(gproc.Pid())
    })
    s.EnableAdmin()
    s.SetPort(8199)
    s.Run()
}
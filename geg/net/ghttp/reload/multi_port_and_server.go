package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "time"
    "gitee.com/johng/gf/g/os/gproc"
)

func main() {
    s1 := g.Server("s1")
    s1.EnableAdmin()
    s1.BindHandler("/", func(r *ghttp.Request) {
        pid := gproc.Pid()
        r.Response.Writeln("before restart, pid:", pid)
        time.Sleep(10*time.Second)
        r.Response.Writeln("after restart, pid:", gproc.Pid())
    })
    s1.BindHandler("/pid", func(r *ghttp.Request) {
        r.Response.Write(gproc.Pid())
    })
    s1.SetPort(8199, 8200)
    s1.Start()

    s2 := g.Server("s2")
    s2.EnableAdmin()
    s2.SetPort(8300, 8080)
    s2.Start()

    g.Wait()
}
package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/os/gproc"
)

func main() {
    s := g.Server()
    s.SetIndexFolder(true)
    s.BindHandler("/", func(r *ghttp.Request){
        r.Response.Write("pid:", gproc.Pid())
    })
    s.SetPort(8199)
    s.Run()
}
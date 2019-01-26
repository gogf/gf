package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    s1 := ghttp.GetServer("s1")
    s1.SetPort(8882)
    s1.BindHandler("/", func(r *ghttp.Request) {
        glog.Println("s1")
        r.Response.Writeln("s1")
    })
    go s1.Run()

    s2 := ghttp.GetServer("s2")
    s2.SetPort(8882)
    s2.BindHandler("/", func(r *ghttp.Request) {
        glog.Println("s2")
        r.Response.Writeln("s1")
    })
    go s2.Run()

    select{}
}
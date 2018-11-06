package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    p := "/"
    s := g.Server()
    s.BindHandler(p, func(r *ghttp.Request) {
        r.Response.Writeln("start")
        r.Exit()
        r.Response.Writeln("end")
    })
    s.BindHookHandlerByMap(p, map[string]ghttp.HandlerFunc{
        ghttp.HOOK_BEFORE_SERVE : func(r *ghttp.Request){
            glog.To(r.Response.Writer).Println("BeforeServe")
        },
        ghttp.HOOK_AFTER_SERVE  : func(r *ghttp.Request){
            glog.To(r.Response.Writer).Println("AfterServe")
        },
    })
    s.SetPort(8199)
    s.Run()
}
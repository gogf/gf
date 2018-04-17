package main

import (
    //"gitee.com/johng/gf/g/net/ghttp"
    _"net/http/pprof"
    "log"
    "net/http"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    //s := ghttp.GetServer()
    s.BindHandler("/", func(r *ghttp.Request){
        r.Response.Writeln("哈喽世界！")
    })
    //s.SetPort(8199)
    //s.Run()
    log.Println(http.ListenAndServe("127.0.0.1:6060", nil))
}
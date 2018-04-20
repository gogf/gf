package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := ghttp.GetServer()
    s.BindHandler("/", func(r *ghttp.Request){
        r.Response.Writeln("哈喽世界！")
        panic("test")
    })
    s.SetAccessLogEnabled(true)
    s.SetErrorLogEnabled(true)
    s.SetPort(8199)
    s.Run()
}
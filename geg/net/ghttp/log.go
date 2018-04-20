package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := ghttp.GetServer()
    s.BindHandler("/log/test", func(r *ghttp.Request){
        r.Response.Writeln("哈喽世界！")
    })
    s.SetLogPath("/tmp/gf.log")
    s.SetAccessLogEnabled(true)
    s.SetErrorLogEnabled(true)
    s.SetPort(8199)
    s.Run()
}
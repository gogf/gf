package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := ghttp.GetServer()
    s.BindHandler("/", func(r *ghttp.Request){
        r.Response.Writeln("来自于HTTPS的：哈喽世界！")
    })
    s.EnableHTTPS("/home/john/temp/server.crt", "/home/john/temp/server.key")
    s.SetPort(8199)
    s.Run()
}
package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request){
        r.Response.Write("哈喽世界！")
    })
    s.EnablePprof()
    s.SetPort(8199)
    s.Run()
}
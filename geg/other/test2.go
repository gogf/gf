package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/{class}-{course}/:name/*act", func(r *ghttp.Request){
        r.Response.Writeln(r.Get("class"))
        r.Response.Writeln(r.Get("course"))
        r.Response.Writeln(r.Get("name"))
        r.Response.Writeln(r.Get("act"))
    })
    s.SetPort(8199)
    s.Run()
}
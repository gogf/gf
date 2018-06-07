package main

import "gitee.com/johng/gf/g/net/ghttp"

func main() {
    s := ghttp.GetServer()
    s.BindHandler("/:name", func(r *ghttp.Request){
        r.Response.Writeln("pattern: /:name match")
        r.Response.Writeln(r.Get("name"))
    })
    s.BindHandler("/:name/:action", func(r *ghttp.Request){
        r.Response.Writeln("pattern: /:name/:action match")
        r.Response.Writeln(r.Get("name"))
        r.Response.Writeln(r.Get("action"))
    })
    s.BindHandler("/:name/*any", func(r *ghttp.Request){
        r.Response.Writeln("pattern: /:name/*any match")
        r.Response.Writeln(r.Get("name"))
        r.Response.Writeln(r.Get("any"))
    })
    s.SetPort(8199)
    s.Run()
}
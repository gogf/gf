package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "fmt"
)

func main() {
    s := g.Server()
    s.BindHandler("/admin", func(r *ghttp.Request) {
        fmt.Println("admin")
    })
    s.BindHandler("/admin/", func(r *ghttp.Request) {
        fmt.Println("admin/")
    })
    s.SetPort(8199)
    s.Run()
}
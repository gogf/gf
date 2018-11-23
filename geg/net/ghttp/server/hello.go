package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.SetIndexFolder(true)
    s.BindHandler("/", func(r *ghttp.Request){
        r.Response.ServeFileDownload("/Users/john/Temp/edrawmax-cn-9.3.dmg")
    })
    s.SetPort(8199)
    s.Run()
}
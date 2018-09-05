package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request){
        r.Cookie.Set("theme", "default")
        r.Session.Set("name", "john")
        content :=`Config:{{.Config.redis.cache}}, Cookie:{{.Cookie.theme}}, Session:{{.Session.name}}`
        r.Response.WriteTplContent(content, nil)
    })
    s.SetPort(8199)
    s.Run()
}
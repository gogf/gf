package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request){
        content :=
            `
            {{if (get "name")}}
                {{get "name"}}
            {{else}}
                NoName
            {{end}}
            `
        r.Response.WriteTplContent(content, nil)
    })
    s.SetPort(8199)
    s.Run()
}
package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gins"
    "gitee.com/johng/gf/g"
)

func main() {
    s := ghttp.GetServer()
    s.BindHandler("/template2", func(r *ghttp.Request){
        tplcontent := `id:{{.id}}, name:{{.name}}`
        content, _ := gins.View().ParseContent(tplcontent, g.Map{
            "id"   : 123,
            "name" : "john",
        })
        r.Response.Write(content)
    })
    //s.SetPort(80)
    s.Run()
}
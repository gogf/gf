package demo

import (
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g"
)

func init() {
    ghttp.GetServer().BindHandler("/template2", func(r *ghttp.Request){
        content, _ := g.View().Parse("index.tpl", map[string]interface{}{
            "id"   : 123,
            "name" : "john",
        })
        r.Response.Write(content)
    })
}
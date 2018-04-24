package demo

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g"
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
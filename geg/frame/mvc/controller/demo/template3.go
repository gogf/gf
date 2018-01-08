package demo

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gins"
)

func init() {
    gins.View().SetPath("/home/www/template/")
    ghttp.GetServer().BindHandler("/template3", func(r *ghttp.Request){
        content, _ := gins.View().Parse("index.tpl", map[string]interface{}{
            "id"   : 123,
            "name" : "john",
        })
        r.Response.Write(content)
    })
}
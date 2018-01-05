package demo

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gins"
)

func init() {
    ghttp.GetServer().BindHandler("/template2", func(r *ghttp.Request){
        view := gins.View()
        view.SetPath("/home/www/template/")
        content, _ := view.Parse("index", map[string]interface{}{
            "id"   : 123,
            "name" : "john",
        })
        r.Response.Write(content)
    })
}
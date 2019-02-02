package demo

import (
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/frame/gins"
)

func init() {
    ghttp.GetServer().BindHandler("/config", func (r *ghttp.Request) {
        r.Response.Write(gins.Config().GetString("database.default.0.host"))
    })
}

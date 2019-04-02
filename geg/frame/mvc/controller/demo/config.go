package demo

import (
	"github.com/gogf/gf/g/frame/gins"
	"github.com/gogf/gf/g/net/ghttp"
)

func init() {
	ghttp.GetServer().BindHandler("/config", func(r *ghttp.Request) {
		r.Response.Write(gins.Config().GetString("database.default.0.host"))
	})
}

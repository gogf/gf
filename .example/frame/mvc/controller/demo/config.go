package demo

import (
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	ghttp.GetServer().BindHandler("/config", func(r *ghttp.Request) {
		r.Response.Write(gins.Config().GetString("database.default.0.host"))
	})
}

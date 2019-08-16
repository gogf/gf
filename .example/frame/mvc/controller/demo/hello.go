package demo

import "github.com/gogf/gf/net/ghttp"

func init() {
	ghttp.GetServer().BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("Hello World!")
	})
}

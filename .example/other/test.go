package main

import (
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := ghttp.GetServer()
	s.BindHandler("/admin", func(r *ghttp.Request) {
		r.Response.Write("admin")
	})
	s.BindHandler("/admin-{page}", func(r *ghttp.Request) {
		r.Response.Write("admin-{page}", r.GetInt("page"))
	})
	s.BindHandler("/admin-goods", func(r *ghttp.Request) {
		r.Response.Write("admin-goods")
	})
	s.BindHandler("/admin-goods-{page}", func(r *ghttp.Request) {
		r.Response.Write("admin-goods-{page}", r.GetInt("page"))
	})
	s.SetPort(8199)
	s.Run()
}

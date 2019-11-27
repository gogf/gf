package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.SetIndexFolder(true)
	s.BindHandler("/admin.html", func(r *ghttp.Request) {
		r.Response.Write("admin")
	})
	s.BindHandler("/admin-do-{page}.html", func(r *ghttp.Request) {
		r.Response.Write("admin-do-" + r.GetString("page"))
	})
	s.SetPort(8999)
	s.Run()
}

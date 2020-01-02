package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/user/:name", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.BindHandler("/user/member/:name/*any", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.BindHandler("/user/member/:name/edit/*any", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.BindHandler("/user/member/:name/edit/sex", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.BindHandler("/user/member/:name/edit/info/*any", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.BindHandler("/user/community/female/:name", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.BindHandler("/admin/stats/today/:hour", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.SetPort(8199)
	s.Run()
}

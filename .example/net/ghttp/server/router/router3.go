package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	s := g.Server()
	// 一个简单的分页路由示例
	s.BindHandler("/user/list/{page}.html", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("page"))
	})
	// {xxx} 规则与 :xxx 规则混合使用
	s.BindHandler("/{object}/:attr/{act}.php", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("object"))
		r.Response.Writeln(r.Get("attr"))
		r.Response.Writeln(r.Get("act"))
	})
	// 多种模糊匹配规则混合使用
	s.BindHandler("/{class}-{course}/:name/*act", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("class"))
		r.Response.Writeln(r.Get("course"))
		r.Response.Writeln(r.Get("name"))
		r.Response.Writeln(r.Get("act"))
	})
	s.SetPort(8199)
	s.Run()
}

package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.RedirectTo("/login")
	})
	s.BindHandler("/login", func(r *ghttp.Request) {
		r.Response.Writeln("Login First")
	})
	s.SetPort(8199)
	s.Run()
}

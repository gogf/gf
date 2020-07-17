package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("amount"))
		r.Response.Writeln(r.GetInt("amount"))
		r.Response.Writeln(r.GetFloat32("amount"))
	})
	s.SetPort(8199)
	s.Run()
}

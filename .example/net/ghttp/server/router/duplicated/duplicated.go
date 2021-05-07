package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/test", func(r *ghttp.Request) {
			r.Response.Writeln(1)
		})
		group.ALL("/test", func(r *ghttp.Request) {
			r.Response.Writeln(2)
		})
	})
	s.SetPort(8199)
	s.Run()
}

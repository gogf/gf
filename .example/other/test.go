package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("POST:/login", func(r *ghttp.Request) {
		r.Response.Write("login handler")
	})
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("for authenticated handler testing")
		})
	})
	s.SetPort(8199)
	s.Run()
}

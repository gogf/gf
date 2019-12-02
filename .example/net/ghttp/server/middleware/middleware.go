package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			r.Response.Write("start")
			r.Middleware.Next()
			r.Response.Write("end")
		})
		g.Group("/order", func(group *ghttp.RouterGroup) {
			g.GET("/list", func(r *ghttp.Request) {
				r.Response.Write("list")
			})
		})
		g.Group("/user", func(group *ghttp.RouterGroup) {
			g.GET("/info", func(r *ghttp.Request) {
				r.Response.Write("info")
			})
			g.POST("/edit", func(r *ghttp.Request) {
				r.Response.Write("edit")
			})
		})
		g.Group("/hook", func(group *ghttp.RouterGroup) {
			group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
				r.Response.Write("hook any")
			})
			group.Hook("/:name", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
				r.Response.Write("hook name")
			})
		})
	})
	s.SetPort(8199)
	s.Run()
}

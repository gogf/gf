package main

import (
	"net/http"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

func MiddlewareAuth(r *ghttp.Request) {
	token := r.Get("token")
	if token == "123456" {
		r.Middleware.Next()
	} else {
		r.Response.WriteStatus(http.StatusForbidden)
	}
}

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func MiddlewareLog(r *ghttp.Request) {
	r.Middleware.Next()
	glog.Println(r.Response.Status, r.URL.Path)
}

func main() {
	s := g.Server()
	s.Group("/", func(g *ghttp.RouterGroup) {
		g.Middleware(MiddlewareLog)
	})
	s.Group("/api.v2", func(g *ghttp.RouterGroup) {
		g.Middleware(MiddlewareAuth, MiddlewareCORS)
		g.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})
		g.Group("/order", func(g *ghttp.RouterGroup) {
			g.GET("/list", func(r *ghttp.Request) {
				r.Response.Write("list")
			})
			g.PUT("/update", func(r *ghttp.Request) {
				r.Response.Write("update")
			})
		})
		g.Group("/user", func(g *ghttp.RouterGroup) {
			g.GET("/info", func(r *ghttp.Request) {
				r.Response.Write("info")
			})
			g.POST("/edit", func(r *ghttp.Request) {
				r.Response.Write("edit")
			})
			g.DELETE("/drop", func(r *ghttp.Request) {
				r.Response.Write("drop")
			})
		})
		g.Group("/hook", func(g *ghttp.RouterGroup) {
			g.Hook("/*", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
				r.Response.Write("hook any")
			})
			g.Hook("/:name", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
				r.Response.Write("hook name")
			})
		})
	})
	s.SetPort(8199)
	s.Run()
}

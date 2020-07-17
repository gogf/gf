package main

import (
	"net/http"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
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
	g.Log().Println(r.Response.Status, r.URL.Path)
}

func main() {
	s := g.Server()
	s.Use(MiddlewareLog)
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareAuth, MiddlewareCORS)
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})
		group.Group("/order", func(group *ghttp.RouterGroup) {
			group.GET("/list", func(r *ghttp.Request) {
				r.Response.Write("list")
			})
			group.PUT("/update", func(r *ghttp.Request) {
				r.Response.Write("update")
			})
		})
		group.Group("/user", func(group *ghttp.RouterGroup) {
			group.GET("/info", func(r *ghttp.Request) {
				r.Response.Write("info")
			})
			group.POST("/edit", func(r *ghttp.Request) {
				r.Response.Write("edit")
			})
			group.DELETE("/drop", func(r *ghttp.Request) {
				r.Response.Write("drop")
			})
		})
		group.Group("/hook", func(group *ghttp.RouterGroup) {
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

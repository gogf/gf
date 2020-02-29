package main

import (
	"net/http"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func MiddlewareAuth(r *ghttp.Request) {
	token := r.Get("token")
	if token == "123456" {
		r.Response.Writeln("auth")
		r.Middleware.Next()
	} else {
		r.Response.WriteStatus(http.StatusForbidden)
	}
}

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.Writeln("cors")
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func main() {
	s := g.Server()
	s.Use(MiddlewareCORS)
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareAuth)
		group.ALL("/user/list", func(r *ghttp.Request) {
			r.Response.Writeln("list")
		})
	})
	s.SetPort(8199)
	s.Run()
}

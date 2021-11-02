package main

import (
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func MiddlewareAuth(r *ghttp.Request) {
	token := r.Get("token")
	if token == "123456" {
		r.Middleware.Next()
	} else {
		r.Response.WriteStatus(http.StatusForbidden)
	}
}

func main() {
	s := g.Server()
	s.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			if action := r.GetRouterString("action"); action != "" {
				switch action {
				case "login":
					r.Middleware.Next()
					return
				}
			}
			MiddlewareAuth(r)
		})
		group.ALL("/login", func(r *ghttp.Request) {
			r.Response.Write("login")
		})
		group.ALL("/dashboard", func(r *ghttp.Request) {
			r.Response.Write("dashboard")
		})
	})
	s.SetPort(8199)
	s.Run()
}

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

func MiddlewareError(r *ghttp.Request) {
	r.Middleware.Next()
	if r.Response.Status >= http.StatusInternalServerError {
		r.Response.ClearBuffer()
		r.Response.Write("Internal error occurred, please try again later.")
	}
}

func main() {
	s := g.Server()
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareAuth, MiddlewareCORS, MiddlewareError)
		group.ALL("/user/list", func(r *ghttp.Request) {
			panic("db error: sql is xxxxxxx")
		})
	})
	s.SetPort(8199)
	s.Run()
}

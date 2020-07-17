package main

import (
	"net/http"

	"github.com/jin502437344/gf/os/glog"

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
		group.ALL("/user/list", func(r *ghttp.Request) {
			panic("custom error")
		})
	})
	s.SetPort(8199)
	s.Run()
}

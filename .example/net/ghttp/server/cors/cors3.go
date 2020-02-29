package main

import (
	"net/http"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func MiddlewareCORS(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	corsOptions.AllowDomain = []string{"goframe.org"}
	if !r.Response.CORSAllowedOrigin(corsOptions) {
		r.Response.WriteStatus(http.StatusForbidden)
		return
	}
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

func Order(r *ghttp.Request) {
	r.Response.Write("GET")
}

func main() {
	s := g.Server()
	s.Group("/api.v1", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS)
		group.GET("/order", Order)
	})
	s.SetPort(8199)
	s.Run()
}

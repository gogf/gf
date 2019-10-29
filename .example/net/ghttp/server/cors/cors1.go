package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func Order(r *ghttp.Request) {
	r.Response.Write("GET")
}

func main() {
	s := g.Server()
	s.Group("/api.v1", func(g *ghttp.RouterGroup) {
		g.Middleware(MiddlewareCORS)
		g.GET("/order", Order)
	})
	s.SetPort(8199)
	s.Run()
}

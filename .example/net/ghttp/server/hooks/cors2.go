package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func Order(r *ghttp.Request) {
	r.Response.Write("GET")
}

func main() {
	s := g.Server()
	s.Group("/api.v1", func(group *ghttp.RouterGroup) {
		group.Hook("/*any", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
			r.Response.CORSDefault()
		})
		g.GET("/order", Order)
	})
	s.SetPort(8199)
	s.Run()
}

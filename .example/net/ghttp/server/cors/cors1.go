package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func Order(r *ghttp.Request) {
	glog.Print("order")
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

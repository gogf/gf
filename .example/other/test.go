package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.ALL("/test", func(r *ghttp.Request) {
			r.Response.Write(r.GetRequest("nickname"))
		})
	})
	s.SetPort(8199)
	s.Run()

}

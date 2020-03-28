package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	_ "github.com/gogf/gf/os/gres/testdata/example/boot"
)

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/template", func(r *ghttp.Request) {
			r.Response.WriteTplDefault(g.Map{
				"name": "GoFrame",
			})
		})
	})
	s.Run()
}

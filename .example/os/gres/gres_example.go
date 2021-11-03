package main

import (
	_ "github.com/gogf/gf/v2/os/gres/testdata/example/boot"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
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

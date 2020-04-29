package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/", func(r *ghttp.Request) {
			fmt.Println(r.GetBodyString())
			fmt.Println(r.Header)
			r.Response.Write(r.GetBodyString())
			r.Response.Write(r.Header)
		})
	})
	s.SetPort(8199)
	s.Run()
}

package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"net/http"
)

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/", func(r *ghttp.Request) {
			r.Response.Write("halo world!")
		})
		group.ALL("/log/handler", func(r *ghttp.Request) {
			r.Response.WriteStatus(http.StatusNotFound, "File Not Found!")
		})
	})
	s.SetPort(8199)
	s.Run()
}

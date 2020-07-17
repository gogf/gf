package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("halo 世界！")
	})
	s.BindStatusHandler(404, func(r *ghttp.Request) {
		r.Response.Writeln("This is customized 404 page")
	})
	s.SetPort(8199)
	s.Run()
}

package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/input", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("amount"))
	})
	s.BindHandler("/query", func(r *ghttp.Request) {
		r.Response.Writeln(r.GetQuery("amount"))
	})
	s.SetPort(8199)
	s.Run()
}

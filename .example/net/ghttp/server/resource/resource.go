package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gres"
	_ "github.com/gogf/gf/os/gres/testdata"
)

func main() {
	gres.Dump()

	v := g.View()
	v.SetPath("/template/layout1")

	s := g.Server()
	s.SetIndexFolder(true)
	s.SetServerRoot("/root")
	s.BindHandler("/template", func(r *ghttp.Request) {
		r.Response.WriteTpl("layout.html")
	})
	s.SetPort(8199)
	s.Run()
}

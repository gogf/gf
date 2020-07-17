package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/main1", func(r *ghttp.Request) {
		r.Response.WriteTpl("layout.html", g.Map{
			"mainTpl": "main/main1.html",
		})
	})
	s.BindHandler("/main2", func(r *ghttp.Request) {
		r.Response.WriteTpl("layout.html", g.Map{
			"mainTpl": "main/main2.html",
		})
	})
	g.View().SetPath("template")
	s.SetPort(8199)
	s.Run()
}

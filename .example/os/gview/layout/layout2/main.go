package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/main1", func(r *ghttp.Request) {
		r.Response.WriteTpl("layout.html", g.Map{
			"name":    "smith",
			"mainTpl": "main/main1.html",
		})
	})
	s.BindHandler("/main2", func(r *ghttp.Request) {
		r.Response.WriteTpl("layout.html", g.Map{
			"name":    "john",
			"mainTpl": "main/main2.html",
		})
	})
	s.SetPort(8199)
	s.Run()
}

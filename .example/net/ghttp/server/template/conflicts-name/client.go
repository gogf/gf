package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
)

// https://github.com/jin502437344/gf/issues/437
func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteTpl("client/layout.html")
	})
	s.SetPort(8199)
	s.Run()
}

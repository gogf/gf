package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteTplContent(`${.name}`, g.Map{
			"name": "john",
		})
	})
	s.SetPort(8199)
	s.Run()
}

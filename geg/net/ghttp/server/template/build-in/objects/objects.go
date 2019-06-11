package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		content := `{{.Request.Get "name"}}`
		r.Response.WriteTplContent(content)
	})
	s.SetPort(8199)
	s.Run()
}

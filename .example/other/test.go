package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	//g.Server().Run()
	s := g.Server()
	s.BindHandler("/aaa", func(r *ghttp.Request) {
		r.Cookie.Set("theme", "default")
		r.Session.Set("name", "john")
		content := `Config:{{.Config.redis.cache}}, Cookie:{{.Cookie.theme}}, Session:{{.Session.name}}, Query:{{.Query.name}}`
		r.Response.WriteTplContent(content, nil)
	})
	s.SetPort(8199)
	s.Run()
}

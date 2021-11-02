package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Cookie.Set("theme", "default")
		r.Session.Set("name", "john")
		content := `
Get: {{.Get.name}}
Post: {{.Post.name}}
Config: {{.Config.redis}}
Cookie: {{.Cookie.theme}}, 
Session: {{.Session.name}}`
		r.Response.WriteTplContent(content)
	})
	s.SetPort(8199)
	s.Run()
}

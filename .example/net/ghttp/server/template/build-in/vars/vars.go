package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
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

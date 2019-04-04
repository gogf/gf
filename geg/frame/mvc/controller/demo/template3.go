package demo

import (
	"github.com/gogf/gf/g/frame/gins"
	"github.com/gogf/gf/g/net/ghttp"
)

func init() {
	gins.View().SetPath("/home/www/template/")
	ghttp.GetServer().BindHandler("/template3", func(r *ghttp.Request) {
		content, _ := gins.View().Parse("index.tpl", map[string]interface{}{
			"id":   123,
			"name": "john",
		})
		r.Response.Write(content)
	})
}

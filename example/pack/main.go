package main

import (
	_ "github.com/gogf/gf/example/pack/packed"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gres"
)

func main() {
	gres.Dump()

	s := g.Server()
	s.SetPort(8199)
	s.SetServerRoot("resource/public")
	s.BindHandler("/i18n", func(r *ghttp.Request) {
		var (
			lang    = r.Get("lang", "zh-CN").String()
			ctx     = gi18n.WithLanguage(r.Context(), lang)
			content string
		)
		content = g.I18n().T(ctx, `{#hello} {#world}!`)
		r.Response.Write(content)
	})
	s.Run()
}

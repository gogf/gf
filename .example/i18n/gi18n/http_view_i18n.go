package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/i18n/gi18n"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			r.SetCtx(gi18n.WithLanguage(r.Context(), "zh-CN"))
			r.Middleware.Next()
		})
		group.ALL("/", func(r *ghttp.Request) {
			r.Response.WriteTplContent(`{#hello}{#world}!`)
		})
	})
	s.SetPort(8199)
	s.Run()
}

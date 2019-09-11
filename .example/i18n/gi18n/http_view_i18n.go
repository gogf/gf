package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
	"github.com/gogf/gf/net/ghttp"
)

type Controller struct {
	gmvc.Controller
}

func (c *Controller) Init(r *ghttp.Request) {
	c.Controller.Init(r)
	c.View.BindFunc("T", c.translate)
}

func (c *Controller) Index() {
	c.View.Display("index.html")
}

func (c *Controller) translate(langKey, msg string) string {
	lang := c.Request.Get("lang", "en")
	g.I18n().SetLanguage(lang)
	t := g.I18n().T(langKey, lang)
	fmt.Println(t, lang, c.Request.Request.Header.Get("Accept-Language"))
	return t
}

func main() {
	g.I18n().SetPath("i18n-file")
	g.View().SetPath(`D:\Workspace\Go\GOPATH\src\github.com\gogf\gf\.example\i18n\gi18n\template`)
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		lang := r.Get("lang", "en")
		g.I18n().SetLanguage(lang)
		r.Response.WriteTplContent(`{#hello}{#world}`)
	})
	s.BindHandler("/template", func(r *ghttp.Request) {
		lang := r.Get("lang", "en")
		g.I18n().SetLanguage(lang)
		r.Response.WriteTplContent(`{#hello}{#world}`)
	})
	s.BindController("/controller", new(Controller))
	s.SetPort(8199)
	s.Run()
}

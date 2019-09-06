package main

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.SetSessionMaxAge(60)
	s.BindHandler("/set", func(r *ghttp.Request) {
		r.Session.Set("captcha", map[string]string{
			"key": "value",
		})
		r.Response.Write("ok")
	})
	s.BindHandler("/get", func(r *ghttp.Request) {
		fmt.Println(r.Session.Get("captcha"))
		r.Response.Write(r.Session.Get("captcha"))
	})
	s.SetPort(8199)
	s.Run()
}

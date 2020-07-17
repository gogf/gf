package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		fmt.Println("cors")
		r.Response.CORSDefault()
		r.Middleware.Next()
	})
	s.BindHandler("/api/captcha", func(r *ghttp.Request) {
		r.Response.Write("captcha")
	})
	s.SetPort(8010)
	s.Run()
}

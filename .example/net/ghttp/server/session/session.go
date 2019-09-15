package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"time"
)

func main() {
	s := g.Server()
	s.SetSessionMaxAge(2 * time.Second)
	s.BindHandler("/set", func(r *ghttp.Request) {
		r.Session.Set("time", gtime.Second())
		r.Response.Write("ok")
	})
	s.BindHandler("/get", func(r *ghttp.Request) {
		r.Response.WriteJson(r.Session.Map())
	})
	s.BindHandler("/clear", func(r *ghttp.Request) {
		r.Session.Clear()
	})
	s.SetPort(8199)
	s.Run()
}

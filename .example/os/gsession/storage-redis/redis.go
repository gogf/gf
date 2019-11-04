package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gsession"
	"github.com/gogf/gf/os/gtime"
	"time"
)

func main() {
	s := g.Server()
	s.SetConfigWithMap(g.Map{
		"SessionMaxAge":  time.Minute,
		"SessionStorage": gsession.NewStorageRedis(g.Redis()),
	})
	s.BindHandler("/set", func(r *ghttp.Request) {
		r.Session.Set("time", gtime.Second())
		r.Response.Write("ok")
	})
	s.BindHandler("/get", func(r *ghttp.Request) {
		r.Response.Write(r.Session.Map())
	})
	s.BindHandler("/del", func(r *ghttp.Request) {
		r.Session.Clear()
		r.Response.Write("ok")
	})
	s.SetPort(8199)
	s.Run()
}

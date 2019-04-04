package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	s := g.Server()
	s.BindHandler("/priority/show", func(r *ghttp.Request) {
		r.Response.Write("priority test")
	})

	s.BindHookHandlerByMap("/priority/:name", map[string]ghttp.HandlerFunc{
		"BeforeServe": func(r *ghttp.Request) {
			glog.Println(r.Router.Uri)
		},
	})
	s.BindHookHandlerByMap("/priority/*any", map[string]ghttp.HandlerFunc{
		"BeforeServe": func(r *ghttp.Request) {
			glog.Println(r.Router.Uri)
		},
	})
	s.BindHookHandlerByMap("/priority/show", map[string]ghttp.HandlerFunc{
		"BeforeServe": func(r *ghttp.Request) {
			glog.Println(r.Router.Uri)
		},
	})
	s.SetPort(8199)
	s.Run()
}

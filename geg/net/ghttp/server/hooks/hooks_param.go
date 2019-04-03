package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln(r.GetParam("name").String())
	})
	s.BindHookHandlerByMap("/", map[string]ghttp.HandlerFunc{
		ghttp.HOOK_BEFORE_SERVE: func(r *ghttp.Request) {
			r.SetParam("name", "john")
		},
	})
	s.SetPort(8199)
	s.Run()
}

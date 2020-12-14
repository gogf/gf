package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHookHandler("/*any", ghttp.HookBeforeServe, func(r *ghttp.Request) {
		r.Response.Writeln("/*any")
	})
	s.BindHookHandler("/v1/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
		r.Response.Writeln("/v1/*")
		r.ExitHook()
	})
	s.SetPort(8199)
	s.Run()
}

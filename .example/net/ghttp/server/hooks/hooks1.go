package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

func main() {
	// 基本事件回调使用
	p := "/:name/info/{uid}"
	s := g.Server()
	s.BindHookHandlerByMap(p, map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe:  func(r *ghttp.Request) { glog.Println(ghttp.HookBeforeServe) },
		ghttp.HookAfterServe:   func(r *ghttp.Request) { glog.Println(ghttp.HookAfterServe) },
		ghttp.HookBeforeOutput: func(r *ghttp.Request) { glog.Println(ghttp.HookBeforeOutput) },
		ghttp.HookAfterOutput:  func(r *ghttp.Request) { glog.Println(ghttp.HookAfterOutput) },
	})
	s.BindHandler(p, func(r *ghttp.Request) {
		r.Response.Write("用户:", r.Get("name"), ", uid:", r.Get("uid"))
	})
	s.SetPort(8199)
	s.Run()
}

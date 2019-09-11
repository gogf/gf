package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

func main() {
	// 基本事件回调使用
	p := "/:name/info/{uid}"
	s := g.Server()
	s.BindHookHandlerByMap(p, map[string]ghttp.HandlerFunc{
		"BeforeServe":  func(r *ghttp.Request) { glog.Println("BeforeServe") },
		"AfterServe":   func(r *ghttp.Request) { glog.Println("AfterServe") },
		"BeforeOutput": func(r *ghttp.Request) { glog.Println("BeforeOutput") },
		"AfterOutput":  func(r *ghttp.Request) { glog.Println("AfterOutput") },
		"BeforeClose":  func(r *ghttp.Request) { glog.Println("BeforeClose") },
		"AfterClose":   func(r *ghttp.Request) { glog.Println("AfterClose") },
	})
	s.BindHandler(p, func(r *ghttp.Request) {
		r.Response.Write("用户:", r.Get("name"), ", uid:", r.Get("uid"))
	})
	s.SetPort(8199)
	s.Run()
}

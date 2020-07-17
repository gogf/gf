package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"github.com/jin502437344/gf/os/glog"
)

func main() {
	p := "/"
	s := g.Server()
	s.BindHandler(p, func(r *ghttp.Request) {
		r.Response.Writeln("start")
		r.Exit()
		r.Response.Writeln("end")
	})
	s.BindHookHandlerByMap(p, map[string]ghttp.HandlerFunc{
		ghttp.HOOK_BEFORE_SERVE: func(r *ghttp.Request) {
			glog.To(r.Response.Writer).Println("BeforeServe")
		},
		ghttp.HOOK_AFTER_SERVE: func(r *ghttp.Request) {
			glog.To(r.Response.Writer).Println("AfterServe")
		},
	})
	s.SetPort(8199)
	s.Run()
}

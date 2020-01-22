package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

func main() {
	s := g.Server()
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.ALL("/user/list", func(r *ghttp.Request) {
			glog.Debug(r.Method, r.RequestURI)

			paramKey := "X-CSRF-Token"

			// // www-form or query
			// glog.Debug("go:", r.Request.FormValue(paramKey))

			// // post form-data
			// glog.Debug("go form:", r.Request.PostFormValue(paramKey))

			glog.Debug("gf GetString:", r.GetString(paramKey))
			glog.Debug("gf GetFormString:", r.GetFormString(paramKey))
			r.Response.Writeln("list")
		})
	})
	s.SetPort(8199)
	s.Run()
}

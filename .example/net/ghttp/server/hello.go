package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

func main() {
	s := g.Server()
	s.SetIndexFolder(true)
	s.BindHandler("/", func(r *ghttp.Request) {
		glog.Println(r.Header)
		r.Response.Write("hello world")
	})
	s.SetPort(8999)
	s.Run()
}

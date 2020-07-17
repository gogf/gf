package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"github.com/jin502437344/gf/os/glog"
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

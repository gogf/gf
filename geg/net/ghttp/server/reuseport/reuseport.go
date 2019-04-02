package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	s1 := ghttp.GetServer("s1")
	s1.SetPort(8882)
	s1.BindHandler("/", func(r *ghttp.Request) {
		glog.Println("s1")
		r.Response.Writeln("s1")
	})
	s1.Start()

	s2 := ghttp.GetServer("s2")
	s2.SetPort(8882)
	s2.BindHandler("/", func(r *ghttp.Request) {
		glog.Println("s2")
		r.Response.Writeln("s2")
	})
	s2.Start()

	g.Wait()
}

package main

import (
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := ghttp.GetServer()
	s.BindHandler("/log/path", func(r *ghttp.Request) {
		r.Response.Writeln("请到/tmp/gf.log目录查看日志")
	})
	s.SetLogPath("/tmp/gf.log")
	s.SetAccessLogEnabled(true)
	s.SetErrorLogEnabled(true)
	s.SetPort(8199)
	s.Run()
}

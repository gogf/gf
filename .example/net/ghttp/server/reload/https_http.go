package main

import (
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	s := ghttp.GetServer()
	s.EnableAdmin()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("您可以同时通过HTTP和HTTPS方式看到该内容！")
	})
	s.EnableHTTPS("/home/john/temp/server.crt", "/home/john/temp/server.key")
	s.SetHTTPSPort(8198, 8199)
	s.SetPort(8200, 8300)
	s.EnableAdmin()
	s.Run()
}

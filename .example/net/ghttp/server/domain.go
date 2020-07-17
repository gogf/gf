package main

import "github.com/jin502437344/gf/net/ghttp"

func Hello1(r *ghttp.Request) {
	r.Response.Write("127.0.0.1: Hello World1!")
}

func Hello2(r *ghttp.Request) {
	r.Response.Write("localhost: Hello World2!")
}

func main() {
	s := ghttp.GetServer()
	s.Domain("127.0.0.1").BindHandler("/", Hello1)
	s.Domain("localhost, local").BindHandler("/", Hello2)
	s.SetPort(8199)
	s.Run()
}

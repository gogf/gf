package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {

	})
	s.BindHandler("/user", func(r *ghttp.Request) {

	})
	s.BindHandler("/user/:id", func(r *ghttp.Request) {
		r.Response.Write(r.GetRouterString("id"))
	})
	s.EnablePprof()
	s.SetPort(3000)
	s.Run()
}

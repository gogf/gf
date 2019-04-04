package main

import (
	"github.com/gogf/gf/g/net/ghttp"
)

func main() {
	s := ghttp.GetServer()
	s.EnablePprof()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("哈喽世界！")
	})
	s.SetPort(8199)
	s.Run()
}

package main

import (
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	s := ghttp.GetServer()
	s.EnablePProf()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("哈喽世界！")
	})
	s.SetPort(8199)
	s.Run()
}

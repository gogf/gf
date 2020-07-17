package main

import (
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := ghttp.GetServer()
	s.BindHandler("/log/error", func(r *ghttp.Request) {
		panic("OMG")
	})
	s.SetErrorLogEnabled(true)
	s.SetPort(8199)
	s.Run()
}

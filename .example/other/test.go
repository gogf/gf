package main

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := ghttp.GetServer()
	s.BindHandler("/*", func(r *ghttp.Request) {
		fmt.Println(r.URL.RawPath)
		r.Response.Write(r.GetUrl())
	})
	s.SetPort(8199)
	s.Run()
}

package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		fmt.Println(r.GetPostMap())
		r.Response.Write("ok")
	})
	s.SetPort(8999)
	s.Run()
}

package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/test", func(r *ghttp.Request) {
		fmt.Println(r.GetBody())
		r.Response.Write(r.GetBody())
	})
	s.SetPort(8199)
	s.Run()
}

package main

import (
	"fmt"

	"github.com/gogf/gf"
	"github.com/gogf/gf/g/net/ghttp"

	"github.com/gogf/gf/g"
)

func main() {
	// fmt.Print(g.)
	fmt.Println(gf.VERSION)
	s := g.Server()

	s.BindHandler("/status/:status", func(r *ghttp.Request) {
		r.Response.Write("woops, status ", r.Get("status"), " found")
	})
	s.BindStatusHandler(404, func(r *ghttp.Request) {
		r.Response.RedirectTo("/status/404")
	})

	s.SetErrorLogEnabled(true)
	s.SetAccessLogEnabled(true)
	s.SetPort(8890)
	s.Run()
}
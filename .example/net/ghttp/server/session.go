package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

func main() {
	s := g.Server()
	s.BindHandler("/session", func(r *ghttp.Request) {
		id := r.Session.GetInt("id")
		r.Session.Set("id", id+1)
		r.Response.Write("id:" + gconv.String(id))
	})
	s.SetPort(8199)
	s.Run()
}

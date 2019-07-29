package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
	"github.com/gogf/gf/net/ghttp"
)

type Order struct {
	gmvc.Controller
}

func (o *Order) Get() {
	o.Response.Write("GET")
}

func main() {
	s := g.Server()
	s.BindHookHandlerByMap("/api.v1/*any", map[string]ghttp.HandlerFunc{
		"BeforeServe": func(r *ghttp.Request) {
			r.Response.SetAllowCrossDomainRequest("*", "PUT,GET,POST,DELETE,OPTIONS")
		},
	})
	s.BindControllerRest("/api.v1/{.struct}", new(Order))
	s.SetPort(8199)
	s.Run()
}

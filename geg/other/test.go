package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

type Order struct{}

func (order *Order) Get(r *ghttp.Request) {
	r.Response.Write("GET")
}

func main() {
	s := g.Server()
	s.BindHookHandlerByMap("/api.v1/*any", map[string]ghttp.HandlerFunc{
		"BeforeServe": func(r *ghttp.Request) {
			r.Response.CORSDefault()
		},
	})
	s.BindObjectRest("/api.v1/{.struct}", new(Order))
	s.SetPort(8199)
	s.Run()
}

package demo

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

type Order struct{}

func init() {
	g.Server().BindObject("/{.struct}-{.method}", new(Order))
}

func (o *Order) List(r *ghttp.Request) {
	r.Response.Write("List")
}

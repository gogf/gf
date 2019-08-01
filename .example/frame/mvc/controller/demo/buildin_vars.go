package demo

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type Order struct{}

func init() {
	g.Server().BindObject("/{.struct}-{.method}", new(Order))
}

func (o *Order) List(r *ghttp.Request) {
	r.Response.Write("List")
}

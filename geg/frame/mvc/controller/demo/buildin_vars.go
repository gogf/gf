package demo

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

type Order struct { }

func init() {
    g.Server().BindObject("/{.struct}-{.method}", new(Order))
}

func (o *Order) List(r *ghttp.Request) {
    r.Response.Write("List")
}

package demo

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

type Product struct {
    total int
}

func init() {
    p := &Product{}
    g.Server().BindHandler("/product/total", p.Total)
    g.Server().BindHandler("/product/list/{page}.html", p.List)
}

func (p *Product) Total(r *ghttp.Request) {
    p.total++
    r.Response.Write("total: ", p.total)
}

func (p *Product) List(r *ghttp.Request) {
    r.Response.Write("page: ", r.Get("page"))
}

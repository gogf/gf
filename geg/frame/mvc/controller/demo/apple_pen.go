package demo

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g"
)

func init() {
    s := g.Server()
    s.BindHandler("/apple",     Apple)
    s.BindHandler("/pen",       Pen)
    s.BindHandler("/apple-pen", ApplePen)
}

func Apple(r *ghttp.Request) {
    r.Response.Write("Apple")
}

func Pen(r *ghttp.Request) {
    r.Response.Write("Pen")
}

func ApplePen(r *ghttp.Request) {
    r.Response.Write("Apple-Pen")
}
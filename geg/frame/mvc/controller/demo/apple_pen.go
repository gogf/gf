package demo

import "gitee.com/johng/gf/g/net/ghttp"

func init() {
    ghttp.GetServer().BindHandler("/apple",     Apple)
    ghttp.GetServer().BindHandler("/pen",       Pen)
    ghttp.GetServer().BindHandler("/apple-pen", ApplePen)
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
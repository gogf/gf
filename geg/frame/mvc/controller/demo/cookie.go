package demo

import (
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/gtime"
)

func init() {
	ghttp.GetServer().BindHandler("/cookie", Cookie)
}

func Cookie(r *ghttp.Request) {
	datetime := r.Cookie.Get("datetime")
	r.Cookie.Set("datetime", gtime.Datetime())
	r.Response.Write("datetime:" + datetime)
}

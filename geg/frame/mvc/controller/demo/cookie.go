package demo

import (
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/net/ghttp"
)


func init() {
    ghttp.GetServer().BindHandler("/cookie", Cookie)
}

// 用于函数映射
func Cookie(s *ghttp.Server, r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    datetime := r.Cookie.Get("datetime")
    r.Cookie.Set("datetime", gtime.Datetime())
    
    w.WriteString("datetime:" + datetime)
}
package demo

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "strconv"
)


func init() {
    ghttp.GetServer().BindHandler("/session", Session)
}

// 用于函数映射
func Session(s *ghttp.Server, r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    id := r.Session.GetInt("id")
    r.Session.Set("id", id + 1)
    
    w.WriteString("id:" + strconv.Itoa(id))
}
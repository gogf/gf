package demo

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/net/gsession"
    "strconv"
)


// 初始化控制器对象，并绑定操作到Web Server
func init() {
    // 将URI映射到指定的方法中执行
    ghttp.GetServer().BindHandler("/hello", Hello)
}

// 用于函数映射
func Hello(s *ghttp.Server, r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    cookie  := ghttp.GetCookie(r.Id())
    session := gsession.Get(cookie.SessionId())

    id := 0
    for i := 0; i < 1; i++ {
        if r := session.Get("id"); r != nil {
            id = r.(int)
        }
        id++
        session.Set("id", id)
    }

    w.WriteString("Hello World!" + strconv.Itoa(id))
}
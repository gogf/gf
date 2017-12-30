package demo

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

// 测试绑定对象
type T struct {}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
    // 只能通过RESTFul方式访问接口，这里为测试方便，使用的是Get
    //g.HTTPServer().BindObject("/test", &T{})
    // 只有localhost域名下才能访问该对象，
    // 对应URL为：http://localhost:8199/test/show
    // 通过该地址将无法访问到内容：http://127.0.0.1:8199/test/show
    g.HTTPServer().Domain("localhost").BindObject("/test", &T{})
}

// 用于对象映射
func (t *T) Show(s *ghttp.Server, r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    w.WriteString("It's show time bibi!")
}


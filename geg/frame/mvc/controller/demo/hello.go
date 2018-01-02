package demo

import "gitee.com/johng/gf/g/net/ghttp"

// 初始化控制器对象，并绑定操作到Web Server
func init() {
    // 将URI映射到指定的方法中执行
    ghttp.GetServer().BindHandler("/hello", Hello)
}

// 用于函数映射
func Hello(r *ghttp.Request) {
    r.Response.WriteString("Hello World!")
}
package user

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/mvc"
)

// 定义业务相关的控制器对象
type Controller_User struct {
    mvc.Controller
}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
    u := &Controller_User{}
    ghttp.GetServer("johng.cn").BindHandle("/user/info", u.Info)
}

// 定义操作逻辑
func (cu *Controller_User) Info(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    uid := r.GetQueryString("uid")
    if uid != "" {
        w.Write([]byte("uid: " + uid + "\n"))
    }
    w.Write([]byte("name: John\n"))
    w.Write([]byte("..."))
}




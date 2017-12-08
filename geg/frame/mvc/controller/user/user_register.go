package user

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/mvc"
)

// 定义业务相关的控制器对象
type Controller_User_Register struct {
    mvc.Controller
}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
    ur := &Controller_User_Register{}
    ghttp.GetServer("johng.cn").BindHandle("/user/register", ur.Show)
}

// 定义操作逻辑
func (cu *Controller_User_Register) Show(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    w.Write([]byte("user register page"))
}




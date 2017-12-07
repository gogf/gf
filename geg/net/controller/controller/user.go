package controller

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

type Controller_User struct {
    ghttp.Controller
}

func (cu *Controller_User) Hello(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    w.Write([]byte("Hello"))
}

func init() {
    user := &Controller_User{}
    ghttp.GetServer("johng.cn").BindHandle("/hello", user.Hello)
}


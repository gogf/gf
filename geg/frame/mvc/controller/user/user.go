package user

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gmvc"
    "fmt"
)

// 定义业务相关的控制器对象
type ControllerUser struct {
    gmvc.Controller
}

func Add(i1, i2 int) int {
    return i1 + i2
}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
    //ghttp.GetServer("johng").Domain("localhost").BindHandler("/user", u.Info)
    ghttp.GetServer("johng").BindController("/user", &ControllerUser{})
}

// 定义操作逻辑
func (c *ControllerUser) Info() {
    fmt.Println(c.Db)
    c.View.Assign("name", "john")
    c.View.Display("user/index")
}




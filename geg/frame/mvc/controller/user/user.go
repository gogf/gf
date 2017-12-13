package user

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gmvc"
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
    c.Cookie.Set("name", "john2")
    c.Response.Write([]byte("user information page"))
    //t, err := template.New("test").Funcs(template.FuncMap{"add": Add}).Parse(gfile.GetContents("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame/mvc/view/user/info.tpl"))
    //if err != nil {
    //    fmt.Println(err)
    //}
    //
    //t.Execute(w.ResponseWriter, map[string]string{
    //    "name" : "john",
    //})
}




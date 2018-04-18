package demo

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gmvc"
)

// 测试控制器
type ControllerRest struct {
    gmvc.Controller
}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
    // 控制器公开方法中与HTTP Method方法同名的方法将会自动绑定映射
    ghttp.GetServer().BindControllerRest("/john", &ControllerRest{})
}

// RESTFul - GET
func (c *ControllerRest) Get() {
    c.Response.Write("RESTFul HTTP Method GET")
}

// RESTFul - POST
func (c *ControllerRest) Post() {
    c.Response.Write("RESTFul HTTP Method POST")
}

// RESTFul - DELETE
func (c *ControllerRest) Delete() {
    c.Response.Write("RESTFul HTTP Method DELETE")
}

// 该方法无法映射，将会无法访问到
func (c *ControllerRest) Hello() {
    c.Response.Write("Hello")
}




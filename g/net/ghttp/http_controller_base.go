package ghttp

import (
    "gitee.com/johng/gf/g/net/gsession"
)

// 控制器基类
type ControllerBase struct {
    Request  *ClientRequest
    Response *ServerResponse
    Cookie   *Cookie
    Session  *gsession.Session
}

// 控制器初始化
func (c *ControllerBase) Init() {
    c.Cookie = NewCookie(c.Request, c.Response)
    if r := c.Cookie.Get("gfsessionid"); r != "" {
        c.Session = gsession.Get(r)
    } else {
        c.Session = gsession.Get(gsession.Id())
    }
}

// 控制器结束请求
func (c *ControllerBase) Shut() {
    c.Cookie.Output()
}



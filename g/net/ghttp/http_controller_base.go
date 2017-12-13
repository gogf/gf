package ghttp

import (
    "gitee.com/johng/gf/g/net/gsession"
)

const (
    gDEFAULT_SESSION_ID_NAME = "gfsessionid"
)

// 控制器基类
type ControllerBase struct {
    Request  *ClientRequest
    Response *ServerResponse
    Cookie   *Cookie
    Session  *gsession.Session
}

// 控制器初始化
func (c *ControllerBase) Init(r *ClientRequest, w *ServerResponse) {
    c.Request  = r
    c.Response = w
    c.Cookie   = NewCookie(c.Request, c.Response)
    if r := c.Cookie.Get(gDEFAULT_SESSION_ID_NAME); r != "" {
        c.Session = gsession.Get(r)
    } else {
        c.Session = gsession.Get(gsession.Id())
    }
}

// 控制器结束请求
func (c *ControllerBase) Shut() {
    if c.Cookie.Get(gDEFAULT_SESSION_ID_NAME) == "" {
        c.Cookie.Set(gDEFAULT_SESSION_ID_NAME, c.Session.Id())
    }
    c.Cookie.Output()
    c.Response.Output()
}



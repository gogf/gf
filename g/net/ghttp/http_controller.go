package ghttp

import "gitee.com/johng/gf/g/net/gsession"

// 控制器基类
type Controller struct {

}

func (c *Controller) Get(*ClientRequest, *ServerResponse)     {}
func (c *Controller) Put(*ClientRequest, *ServerResponse)     {}
func (c *Controller) Post(*ClientRequest, *ServerResponse)    {}
func (c *Controller) Delete(*ClientRequest, *ServerResponse)  {}
func (c *Controller) Head(*ClientRequest, *ServerResponse)    {}
func (c *Controller) Patch(*ClientRequest, *ServerResponse)   {}
func (c *Controller) Connect(*ClientRequest, *ServerResponse) {}
func (c *Controller) Options(*ClientRequest, *ServerResponse) {}
func (c *Controller) Trace(*ClientRequest, *ServerResponse)   {}

// 获取当前请求的session对象
func (c *Controller) Session(r *ClientRequest, w *ServerResponse) *gsession.Session {
    sessionid := ""
    if r, err := r.Cookie("gfsessionid"); err == nil {
        sessionid = r.Value
    } else {
        sessionid = gsession.Id()
    }
    return gsession.Get(sessionid)
}

// 请求初始化时的回调函数
func (c *Controller) __init(r *ClientRequest, w *ServerResponse) {

}

// 请求结束时的回调函数
func (c *Controller) __shut(r *ClientRequest, w *ServerResponse) {
    
}
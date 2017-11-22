package ghttp

// 控制器基类
type Controller struct {
    Server *Server
}

// 控制器接口
type ControllerApi interface {
    Get(r *ClientRequest, w *ServerResponse)
    Put(r *ClientRequest, w *ServerResponse)
    Post(r *ClientRequest, w *ServerResponse)
    Delete(r *ClientRequest, w *ServerResponse)
    Head(r *ClientRequest, w *ServerResponse)
    Patch(r *ClientRequest, w *ServerResponse)
    Connect(r *ClientRequest, w *ServerResponse)
    Options(r *ClientRequest, w *ServerResponse)
    Trace(r *ClientRequest, w *ServerResponse)
}

func (c *Controller) Get(r *ClientRequest, w *ServerResponse)     {}
func (c *Controller) Put(r *ClientRequest, w *ServerResponse)     {}
func (c *Controller) Post(r *ClientRequest, w *ServerResponse)    {}
func (c *Controller) Delete(r *ClientRequest, w *ServerResponse)  {}
func (c *Controller) Head(r *ClientRequest, w *ServerResponse)    {}
func (c *Controller) Patch(r *ClientRequest, w *ServerResponse)   {}
func (c *Controller) Connect(r *ClientRequest, w *ServerResponse) {}
func (c *Controller) Options(r *ClientRequest, w *ServerResponse) {}
func (c *Controller) Trace(r *ClientRequest, w *ServerResponse)   {}
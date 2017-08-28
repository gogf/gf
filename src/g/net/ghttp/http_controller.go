package ghttp

// 控制器基类
type Controller struct {
    Server *Server
}

// 控制器接口
type ControllerApi interface {
    GET(r *ClientRequest, w *ServerResponse)
    PUT(r *ClientRequest, w *ServerResponse)
    POST(r *ClientRequest, w *ServerResponse)
    DELETE(r *ClientRequest, w *ServerResponse)
    HEAD(r *ClientRequest, w *ServerResponse)
    PATCH(r *ClientRequest, w *ServerResponse)
    CONNECT(r *ClientRequest, w *ServerResponse)
    OPTIONS(r *ClientRequest, w *ServerResponse)
    TRACE(r *ClientRequest, w *ServerResponse)
}

func (c *Controller) GET(r *ClientRequest, w *ServerResponse)     {}
func (c *Controller) PUT(r *ClientRequest, w *ServerResponse)     {}
func (c *Controller) POST(r *ClientRequest, w *ServerResponse)    {}
func (c *Controller) DELETE(r *ClientRequest, w *ServerResponse)  {}
func (c *Controller) HEAD(r *ClientRequest, w *ServerResponse)    {}
func (c *Controller) PATCH(r *ClientRequest, w *ServerResponse)   {}
func (c *Controller) CONNECT(r *ClientRequest, w *ServerResponse) {}
func (c *Controller) OPTIONS(r *ClientRequest, w *ServerResponse) {}
func (c *Controller) TRACE(r *ClientRequest, w *ServerResponse)   {}
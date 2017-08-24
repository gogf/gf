package ghttp

// 控制器基类
type ControllerBase struct {
    Server *Server
}

// 控制器接口
type Controller interface {
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

func (c *ControllerBase) GET(r *ClientRequest, w *ServerResponse)     {}
func (c *ControllerBase) PUT(r *ClientRequest, w *ServerResponse)     {}
func (c *ControllerBase) POST(r *ClientRequest, w *ServerResponse)    {}
func (c *ControllerBase) DELETE(r *ClientRequest, w *ServerResponse)  {}
func (c *ControllerBase) HEAD(r *ClientRequest, w *ServerResponse)    {}
func (c *ControllerBase) PATCH(r *ClientRequest, w *ServerResponse)   {}
func (c *ControllerBase) CONNECT(r *ClientRequest, w *ServerResponse) {}
func (c *ControllerBase) OPTIONS(r *ClientRequest, w *ServerResponse) {}
func (c *ControllerBase) TRACE(r *ClientRequest, w *ServerResponse)   {}
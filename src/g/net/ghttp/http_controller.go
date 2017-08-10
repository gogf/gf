package ghttp

// 控制器基类
type ControllerBase struct {

}

// 控制器接口
type Controller interface {
    GET(r *Request, w *ServerResponse)
    PUT(r *Request, w *ServerResponse)
    POST(r *Request, w *ServerResponse)
    DELETE(r *Request, w *ServerResponse)
    HEAD(r *Request, w *ServerResponse)
    PATCH(r *Request, w *ServerResponse)
    CONNECT(r *Request, w *ServerResponse)
    OPTIONS(r *Request, w *ServerResponse)
    TRACE(r *Request, w *ServerResponse)
}

func (c *ControllerBase) GET(r *Request, w *ServerResponse)     {}
func (c *ControllerBase) PUT(r *Request, w *ServerResponse)     {}
func (c *ControllerBase) POST(r *Request, w *ServerResponse)    {}
func (c *ControllerBase) DELETE(r *Request, w *ServerResponse)  {}
func (c *ControllerBase) HEAD(r *Request, w *ServerResponse)    {}
func (c *ControllerBase) PATCH(r *Request, w *ServerResponse)   {}
func (c *ControllerBase) CONNECT(r *Request, w *ServerResponse) {}
func (c *ControllerBase) OPTIONS(r *Request, w *ServerResponse) {}
func (c *ControllerBase) TRACE(r *Request, w *ServerResponse)   {}
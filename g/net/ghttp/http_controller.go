package ghttp

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
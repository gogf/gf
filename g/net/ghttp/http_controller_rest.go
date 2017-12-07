package ghttp

// RESTful控制器接口
type ControllerRest interface {
    Get(*ClientRequest, *ServerResponse)
    Put(*ClientRequest, *ServerResponse)
    Post(*ClientRequest, *ServerResponse)
    Delete(*ClientRequest, *ServerResponse)
    Head(*ClientRequest, *ServerResponse)
    Patch(*ClientRequest, *ServerResponse)
    Connect(*ClientRequest, *ServerResponse)
    Options(*ClientRequest, *ServerResponse)
    Trace(*ClientRequest, *ServerResponse)
}
package ghttp


// 控制器接口
type Controller interface {
    Init(*Server, *ClientRequest, *ServerResponse)
    Shut()
}

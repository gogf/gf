package ghttp


// 控制器接口
type Controller interface {
    Init(*ClientRequest, *ServerResponse)
    Shut()
}

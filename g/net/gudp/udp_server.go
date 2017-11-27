package gudp

import "gitee.com/johng/gf/g/os/glog"

// 执行监听
func (s *gUdpServer) Run() {
    if s == nil || s.listener == nil {
        glog.Println("start running failed: socket address bind failed")
        return
    }
    if s.handler == nil {
        glog.Println("start running failed: socket handler not defined")
        return
    }
    for {
        s.handler(s.listener)
    }

}

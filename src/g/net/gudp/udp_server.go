package gudp

import "log"

// 执行监听
func (s *gUdpServer) Run() {
    if s == nil || s.listener == nil {
        log.Println("start running failed: socket address bind failed")
        return
    }
    if s.handler == nil {
        log.Println("start running failed: socket handler not defined")
        return
    }
    for {
        s.handler(s.listener)
    }

}

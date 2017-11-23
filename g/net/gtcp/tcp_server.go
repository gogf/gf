package gtcp

import (
    "gf/g/os/glog"
)

// 执行监听
func (s *gTcpServer) Run() {
    if s == nil || s.listener == nil {
        glog.Println("start running failed: socket address bind failed")
        return
    }
    for  {
        conn, err := s.listener.Accept()
        if err != nil {
            glog.Error(err)
        } else if conn != nil {
            go s.handler(conn)
        }
    }
}

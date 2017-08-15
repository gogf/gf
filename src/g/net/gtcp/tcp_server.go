package gtcp

import (
    "g/os/glog"
)

// 执行监听
func (s *gTcpServer) Run() {
    if s == nil || s.listener == nil {
        glog.Println("start running failed: socket address bind failed")
        return
    }
    //if s.handler == nil {
    //    glog.Println("start running failed: socket handler not defined")
    //    return
    //}

    //fmt.Println("listening on address", s.address)
    for  {
        conn, err := s.listener.Accept()
        if err != nil {
            glog.Println(err)
        } else if conn != nil {
            go s.handler(conn)
        }
    }
    //fmt.Println("tcp server closed on address", s.address)
}

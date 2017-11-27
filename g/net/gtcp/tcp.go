package gtcp

import (
    "net"
    "gitee.com/johng/gf/g/os/glog"
)

// tcp server结构体
type gTcpServer struct {
    address   string
    listener *net.TCPListener
    handler   func (net.Conn)
}

// 创建一个tcp server对象
func NewServer (address string, handler func (net.Conn)) *gTcpServer {
    tcpaddr, err := net.ResolveTCPAddr("tcp4", address)
    if err != nil {
        glog.Fatalln(err)
        return nil
    }
    listen, err := net.ListenTCP("tcp", tcpaddr)
    if err != nil {
        glog.Fatalln(err)
        return nil
    }
    return &gTcpServer{ address, listen, handler}
}


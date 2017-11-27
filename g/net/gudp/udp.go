package gudp

import (
    "net"
    "gitee.com/johng/gf/g/os/glog"
)

// tcp server结构体
type gUdpServer struct {
    address   string
    listener *net.UDPConn
    handler   func (*net.UDPConn)
}

// 创建一个tcp server对象
func NewServer (address string, handler func (*net.UDPConn)) *gUdpServer {
    tcpaddr, err := net.ResolveUDPAddr("udp4", address)
    if err != nil {
        glog.Println(err)
        return nil
    }
    listen, err := net.ListenUDP("udp", tcpaddr)
    if err != nil {
        glog.Println(err)
        return nil
    }
    return &gUdpServer{ address, listen, handler}
}


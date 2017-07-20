package gtcp

import (
    "net"
    "fmt"
    "sync"
)

// 用于tcp server的gorutine监听管理
var ServerWaitGroup sync.WaitGroup

// tcp server结构体
type gTcpServer struct {
    address   string
    listener *net.TCPListener
    handler   func (net.Conn)
}

// 创建一个tcp server对象
func NewTCPServer (address string, handler func (net.Conn)) *gTcpServer {
    tcpaddr, err := net.ResolveTCPAddr("tcp4", address)
    if err != nil {
        return nil
    }
    listen, err := net.ListenTCP("tcp", tcpaddr)
    if err != nil {
        return nil
    }
    return &gTcpServer{ address, listen, handler}
}

// 执行监听
func (s *gTcpServer) Run() {
    ServerWaitGroup.Add(1)
    go func() {
        fmt.Println("listening on address", s.address)
        for  {
            conn, err := s.listener.Accept()
            if err != nil {
                conn.Close()
            }
            go s.handler(conn)
        }
        fmt.Println("tcp server closed on address", s.address)
        ServerWaitGroup.Add(-1)
    }()
}
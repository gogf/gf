// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtcp

import (
    "net"
    "gitee.com/johng/gf/g/os/glog"
)

// tcp server结构体
type Server struct {
    address   string
    listener *net.TCPListener
    handler   func (net.Conn)
}

// 创建一个tcp server对象
func NewServer(address string, handler func (net.Conn)) *Server {
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
    return &Server{ address, listen, handler}
}


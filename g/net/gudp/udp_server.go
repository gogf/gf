// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package gudp

import (
    "net"
    "errors"
)

// 设置参数 - address
func (s *Server) SetAddress (address string) {
    s.address = address
}

// 设置参数 - handler
func (s *Server) SetHandler (handler func (*net.UDPConn)) {
    s.handler = handler
}

// 执行监听
func (s *Server) Run() error {
    if s.handler == nil {
        return errors.New("start running failed: socket handler not defined")
    }
    tcpaddr, err := net.ResolveUDPAddr("udp4", s.address)
    if err != nil {
        return err
    }
    listen, err := net.ListenUDP("udp", tcpaddr)
    if err != nil {
        return err
    }
    for {
        s.handler(listen)
    }
    return nil
}

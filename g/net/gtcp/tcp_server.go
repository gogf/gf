// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package gtcp

import (
    "errors"
    "gitee.com/johng/gf/g/os/glog"
    "net"
)

// 设置参数 - address
func (s *Server) SetAddress (address string) {
    s.address = address
}

// 设置参数 - handler
func (s *Server) SetHandler (handler func (net.Conn)) {
    s.handler = handler
}

// 执行监听
func (s *Server) Run() error {
    if s.handler == nil {
        return errors.New("start running failed: socket handler not defined")
    }
    tcpaddr, err := net.ResolveTCPAddr("tcp4", s.address)
    if err != nil {
        return err
    }
    listen, err := net.ListenTCP("tcp", tcpaddr)
    if err != nil {
        return err
    }
    for  {
        if conn, err := listen.Accept(); err != nil {
            glog.Error(err)
        } else if conn != nil {
            go s.handler(conn)
        }
    }
    return nil
}

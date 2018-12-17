// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package gudp

import (
    "net"
    "errors"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/util/gconv"
)

const (
    gDEFAULT_SERVER = "default"
)

// tcp server结构体
type Server struct {
    address   string
    handler   func (*Conn)
}

// Server表，用以存储和检索名称与Server对象之间的关联关系
var serverMapping = gmap.NewStringInterfaceMap()

// 获取/创建一个空配置的UDP Server
// 单例模式，请保证name的唯一性
func GetServer(name...interface{}) (*Server) {
    sname := gDEFAULT_SERVER
    if len(name) > 0 {
        sname = gconv.String(name[0])
    }
    if s := serverMapping.Get(sname); s != nil {
        return s.(*Server)
    }
    s := NewServer("", nil)
    serverMapping.Set(sname, s)
    return s
}

// 创建一个tcp server对象，并且可以选择指定一个单例名字
func NewServer (address string, handler func (*Conn), names...string) *Server {
    s := &Server{address, handler}
    if len(names) > 0 {
        serverMapping.Set(names[0], s)
    }
    return s
}

// 设置参数 - address
func (s *Server) SetAddress (address string) {
    s.address = address
}

// 设置参数 - handler
func (s *Server) SetHandler (handler func (*Conn)) {
    s.handler = handler
}

// 执行监听
func (s *Server) Run() error {
    if s.handler == nil {
        return errors.New("start running failed: socket handler not defined")
    }
    addr, err := net.ResolveUDPAddr("udp", s.address)
    if err != nil {
        return err
    }
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        return err
    }
    for {
        s.handler(NewConnByNetConn(conn))
    }
}

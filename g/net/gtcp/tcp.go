// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package gtcp

import (
    "net"
    "gitee.com/johng/gf/g/container/gmap"
)

const (
    gDEFAULT_SERVER = "default"
)

// tcp server结构体
type Server struct {
    address   string
    handler   func (net.Conn)
}

// Server表，用以存储和检索名称与Server对象之间的关联关系
var serverMapping = gmap.NewStringInterfaceMap()

// 获取/创建一个空配置的TCP Server
// 单例模式，请保证name的唯一性
func GetServer(names...string) (*Server) {
    name := gDEFAULT_SERVER
    if len(names) > 0 {
        name = names[0]
    }
    if s := serverMapping.Get(name); s != nil {
        return s.(*Server)
    }
    s := NewServer("", nil)
    serverMapping.Set(name, s)
    return s
}

// 创建一个tcp server对象，并且可以选择指定一个单例名字
func NewServer(address string, handler func (net.Conn), names...string) *Server {
    s := &Server{address, handler}
    if len(names) > 0 {
        serverMapping.Set(names[0], s)
    }
    return s
}


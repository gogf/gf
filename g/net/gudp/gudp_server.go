// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp

import (
	"errors"
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/util/gconv"
	"net"
)

const (
	gDEFAULT_SERVER = "default"
)

// tcp server结构体
type Server struct {
	conn    *Conn  // UDP server connection object.
	address string // Listening address.
	handler func(*Conn)
}

// Server表，用以存储和检索名称与Server对象之间的关联关系
var serverMapping = gmap.NewStrAnyMap()

// 获取/创建一个空配置的UDP Server
// 单例模式，请保证name的唯一性
func GetServer(name ...interface{}) *Server {
	serverName := gDEFAULT_SERVER
	if len(name) > 0 {
		serverName = gconv.String(name[0])
	}
	if s := serverMapping.Get(serverName); s != nil {
		return s.(*Server)
	}
	s := NewServer("", nil)
	serverMapping.Set(serverName, s)
	return s
}

// 创建一个tcp server对象，并且可以选择指定一个单例名字
func NewServer(address string, handler func(*Conn), names ...string) *Server {
	s := &Server{
		address: address,
		handler: handler,
	}
	if len(names) > 0 {
		serverMapping.Set(names[0], s)
	}
	return s
}

// 设置参数 - address
func (s *Server) SetAddress(address string) {
	s.address = address
}

// 设置参数 - handler
func (s *Server) SetHandler(handler func(*Conn)) {
	s.handler = handler
}

// Close closes the connection.
// It will make server shutdowns immediately.
func (s *Server) Close() error {
	return s.conn.Close()
}

// 执行监听
func (s *Server) Run() error {
	if s.handler == nil {
		err := errors.New("start running failed: socket handler not defined")
		glog.Error(err)
		return err
	}
	addr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		glog.Error(err)
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		glog.Error(err)
		return err
	}
	s.conn = NewConnByNetConn(conn)
	s.handler(s.conn)
	return nil
}

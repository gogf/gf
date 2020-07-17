// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp

import (
	"errors"
	"net"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
)

const (
	gDEFAULT_SERVER = "default"
)

// Server is the UDP server.
type Server struct {
	conn    *Conn       // UDP server connection object.
	address string      // UDP server listening address.
	handler func(*Conn) // Handler for UDP connection.
}

var (
	// serverMapping is used for instance name to its UDP server mappings.
	serverMapping = gmap.NewStrAnyMap(true)
)

// GetServer creates and returns a UDP server instance with given name.
func GetServer(name ...interface{}) *Server {
	serverName := gDEFAULT_SERVER
	if len(name) > 0 && name[0] != "" {
		serverName = gconv.String(name[0])
	}
	if s := serverMapping.Get(serverName); s != nil {
		return s.(*Server)
	}
	s := NewServer("", nil)
	serverMapping.Set(serverName, s)
	return s
}

// NewServer creates and returns a UDP server.
// The optional parameter <name> is used to specify its name, which can be used for
// GetServer function to retrieve its instance.
func NewServer(address string, handler func(*Conn), name ...string) *Server {
	s := &Server{
		address: address,
		handler: handler,
	}
	if len(name) > 0 && name[0] != "" {
		serverMapping.Set(name[0], s)
	}
	return s
}

// SetAddress sets the server address for UDP server.
func (s *Server) SetAddress(address string) {
	s.address = address
}

// SetHandler sets the connection handler for UDP server.
func (s *Server) SetHandler(handler func(*Conn)) {
	s.handler = handler
}

// Close closes the connection.
// It will make server shutdowns immediately.
func (s *Server) Close() error {
	return s.conn.Close()
}

// Run starts listening UDP connection.
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

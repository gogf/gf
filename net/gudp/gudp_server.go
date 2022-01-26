// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp

import (
	"net"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	defaultServer = "default"
)

// Server is the UDP server.
type Server struct {
	conn    *Conn       // UDP server connection object.
	address string      // UDP server listening address.
	handler func(*Conn) // Handler for UDP connection.
}

// serverMapping is used for instance name to its UDP server mappings.
var serverMapping = gmap.NewStrAnyMap(true)

// GetServer creates and returns a UDP server instance with given name.
func GetServer(name ...interface{}) *Server {
	serverName := defaultServer
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
// The optional parameter `name` is used to specify its name, which can be used for
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
func (s *Server) Close() (err error) {
	err = s.conn.Close()
	if err != nil {
		err = gerror.Wrap(err, "connection failed")
	}
	return
}

// Run starts listening UDP connection.
func (s *Server) Run() error {
	if s.handler == nil {
		err := gerror.NewCode(gcode.CodeMissingConfiguration, "start running failed: socket handler not defined")
		return err
	}
	addr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		err = gerror.Wrapf(err, `net.ResolveUDPAddr failed for address "%s"`, s.address)
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		err = gerror.Wrapf(err, `net.ListenUDP failed for address "%s"`, s.address)
		return err
	}
	s.conn = NewConnByNetConn(conn)
	s.handler(s.conn)
	return nil
}

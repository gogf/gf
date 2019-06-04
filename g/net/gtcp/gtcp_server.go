// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.


package gtcp

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/util/gconv"
	"net"
	"time"
)

const (
    gDEFAULT_SERVER = "default"
)

// TCP Server.
type Server struct {
    address   string
    handler   func (*Conn)
	tlsConfig *tls.Config
}

// Map for name to server, for singleton purpose.
var serverMapping = gmap.NewStrAnyMap()

// GetServer returns the TCP server with specified <name>,
// or it returns a new normal TCP server named <name> if it does not exist.
// The parameter <name> is used to specify the TCP server
func GetServer(name...interface{}) *Server {
    serverName := gDEFAULT_SERVER
    if len(name) > 0 {
        serverName = gconv.String(name[0])
    }
    return serverMapping.GetOrSetFuncLock(serverName, func() interface{} {
	    return NewServer("", nil)
    }).(*Server)
}

// NewServer creates and returns a new normal TCP server.
// The param <name> is optional, which is used to specify the instance name of the server.
func NewServer(address string, handler func (*Conn), name...string) *Server {
    s := &Server{
    	address : address,
    	handler : handler,
    }
    if len(name) > 0 {
        serverMapping.Set(name[0], s)
    }
    return s
}

// NewTlsServer creates and returns a new TCP server with TLS support.
// The param <name> is optional, which is used to specify the instance name of the server.
func NewTLSServer(address, crtFile, keyFile string, handler func (*Conn), name...string) *Server {
	s := NewServer(address, handler, name...)
	s.SetTLSKeyCrt(crtFile, keyFile)
	return s
}

// SetAddress sets the listening address for server.
func (s *Server) SetAddress (address string) {
    s.address = address
}

// SetHandler sets the connection handler for server.
func (s *Server) SetHandler (handler func (*Conn)) {
    s.handler = handler
}

// SetTlsKeyCrt sets the certificate and key file for TLS configuration of server.
func (s *Server) SetTLSKeyCrt (crtFile, keyFile string) error {
	crt, err := tls.LoadX509KeyPair(crtFile,keyFile)
	if err != nil {
		return err
	}
	s.tlsConfig              = &tls.Config{}
	s.tlsConfig.Certificates = []tls.Certificate{crt}
	s.tlsConfig.Time         = time.Now
	s.tlsConfig.Rand         = rand.Reader
	return nil
}

// SetTlsConfig sets the TLS configuration of server.
func (s *Server) SetTLSConfig(tlsConfig *tls.Config) {
	s.tlsConfig = tlsConfig
}

// Run starts running the TCP Server.
func (s *Server) Run() (err error) {
    if s.handler == nil {
        err = errors.New("start running failed: socket handler not defined")
        glog.Error(err)
        return
    }
    listen := net.Listener(nil)
    if s.tlsConfig != nil {
    	// TLS Server
	    listen, err = tls.Listen("tcp", s.address, s.tlsConfig)
	    if err != nil {
	    	glog.Error(err)
	    	return
	    }
    } else {
    	// Normal Server
	    addr, err := net.ResolveTCPAddr("tcp", s.address)
	    if err != nil {
		    glog.Error(err)
		    return err
	    }
	    listen, err = net.ListenTCP("tcp", addr)
	    if err != nil {
		    glog.Error(err)
		    return err
	    }
    }
    for {
        if conn, err := listen.Accept(); err != nil {
            glog.Error(err)
            return err
        } else if conn != nil {
            go s.handler(NewConnByNetConn(conn))
        }
    }
}

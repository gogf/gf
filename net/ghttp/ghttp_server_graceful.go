// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/text/gstr"
)

// gracefulServer wraps the net/http.Server with graceful reload/restart feature.
type gracefulServer struct {
	server      *Server      // Belonged server.
	fd          uintptr      // File descriptor for passing to the child process when graceful reload.
	address     string       // Listening address like:":80", ":8080".
	httpServer  *http.Server // Underlying http.Server.
	rawListener net.Listener // Underlying net.Listener.
	rawLnMu     sync.RWMutex // Concurrent safety mutex for `rawListener`.
	listener    net.Listener // Wrapped net.Listener.
	isHttps     bool         // Is HTTPS.
	status      int          // Status of current server.
}

// newGracefulServer creates and returns a graceful http server with a given address.
// The optional parameter `fd` specifies the file descriptor which is passed from parent server.
func (s *Server) newGracefulServer(address string, fd ...int) *gracefulServer {
	// Change port to address like: 80 -> :80
	if gstr.IsNumeric(address) {
		address = ":" + address
	}
	gs := &gracefulServer{
		server:     s,
		address:    address,
		httpServer: s.newHttpServer(address),
	}
	if len(fd) > 0 && fd[0] > 0 {
		gs.fd = uintptr(fd[0])
	}
	if s.config.Listeners != nil {
		addrArray := gstr.SplitAndTrim(address, ":")
		addrPort, err := strconv.Atoi(addrArray[len(addrArray)-1])
		if err == nil {
			for _, v := range s.config.Listeners {
				if listenerPort := (v.Addr().(*net.TCPAddr)).Port; listenerPort == addrPort {
					gs.rawListener = v
					break
				}
			}
		}
	}
	return gs
}

// newHttpServer creates and returns an underlying http.Server with a given address.
func (s *Server) newHttpServer(address string) *http.Server {
	server := &http.Server{
		Addr:           address,
		Handler:        http.HandlerFunc(s.config.Handler),
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		IdleTimeout:    s.config.IdleTimeout,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
		ErrorLog:       log.New(&errorLogger{logger: s.config.Logger}, "", 0),
	}
	server.SetKeepAlivesEnabled(s.config.KeepAlive)
	return server
}

// Fd retrieves and returns the file descriptor of the current server.
// It is available ony in *nix like operating systems like linux, unix, darwin.
func (s *gracefulServer) Fd() uintptr {
	if ln := s.getRawListener(); ln != nil {
		file, err := ln.(*net.TCPListener).File()
		if err == nil {
			return file.Fd()
		}
	}
	return 0
}

// setFd sets the file descriptor for current server.
func (s *gracefulServer) setFd(fd int) {
	s.fd = uintptr(fd)
}

// CreateListener creates listener on configured address.
func (s *gracefulServer) CreateListener() error {
	ln, err := s.getNetListener()
	if err != nil {
		return err
	}
	s.listener = ln
	s.setRawListener(ln)
	return nil
}

// CreateListenerTLS creates listener on configured address with HTTPS.
// The parameter `certFile` and `keyFile` specify the necessary certification and key files for HTTPS.
// The optional parameter `tlsConfig` specifies the custom TLS configuration.
func (s *gracefulServer) CreateListenerTLS(certFile, keyFile string, tlsConfig ...*tls.Config) error {
	var config *tls.Config
	if len(tlsConfig) > 0 && tlsConfig[0] != nil {
		config = tlsConfig[0]
	} else if s.httpServer.TLSConfig != nil {
		config = s.httpServer.TLSConfig
	} else {
		config = &tls.Config{}
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}
	var err error
	if len(config.Certificates) == 0 {
		config.Certificates = make([]tls.Certificate, 1)
		if gres.Contains(certFile) {
			config.Certificates[0], err = tls.X509KeyPair(
				gres.GetContent(certFile),
				gres.GetContent(keyFile),
			)
		} else {
			config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		}
	}
	if err != nil {
		return gerror.Wrapf(err, `open certFile "%s" and keyFile "%s" failed`, certFile, keyFile)
	}
	ln, err := s.getNetListener()
	if err != nil {
		return err
	}

	s.listener = tls.NewListener(ln, config)
	s.setRawListener(ln)
	return nil
}

// Serve starts the serving with blocking way.
func (s *gracefulServer) Serve(ctx context.Context) error {
	if s.rawListener == nil {
		return gerror.NewCode(gcode.CodeInvalidOperation, `call CreateListener/CreateListenerTLS before Serve`)
	}

	action := "started"
	if s.fd != 0 {
		action = "reloaded"
	}
	s.server.Logger().Infof(
		ctx,
		`pid[%d]: %s server %s listening on [%s]`,
		gproc.Pid(), s.getProto(), action, s.GetListenedAddress(),
	)
	s.status = ServerStatusRunning
	err := s.httpServer.Serve(s.listener)
	s.status = ServerStatusStopped
	return err
}

// GetListenedAddress retrieves and returns the address string which are listened by current server.
func (s *gracefulServer) GetListenedAddress() string {
	if !gstr.Contains(s.address, freePortAddress) {
		return s.address
	}
	var (
		address      = s.address
		listenedPort = s.GetListenedPort()
	)
	address = gstr.Replace(address, freePortAddress, fmt.Sprintf(`:%d`, listenedPort))
	return address
}

// GetListenedPort retrieves and returns one port which is listened to by current server.
func (s *gracefulServer) GetListenedPort() int {
	if ln := s.getRawListener(); ln != nil {
		return ln.Addr().(*net.TCPAddr).Port
	}
	return 0
}

// getProto retrieves and returns the proto string of current server.
func (s *gracefulServer) getProto() string {
	proto := "http"
	if s.isHttps {
		proto = "https"
	}
	return proto
}

// getNetListener retrieves and returns the wrapped net.Listener.
func (s *gracefulServer) getNetListener() (net.Listener, error) {
	if s.rawListener != nil {
		return s.rawListener, nil
	}
	var (
		ln  net.Listener
		err error
	)
	if s.fd > 0 {
		f := os.NewFile(s.fd, "")
		ln, err = net.FileListener(f)
		if err != nil {
			err = gerror.Wrap(err, "net.FileListener failed")
			return nil, err
		}
	} else {
		ln, err = net.Listen("tcp", s.httpServer.Addr)
		if err != nil {
			err = gerror.Wrapf(err, `net.Listen address "%s" failed`, s.httpServer.Addr)
		}
	}
	return ln, err
}

// shutdown shuts down the server gracefully.
func (s *gracefulServer) shutdown(ctx context.Context) {
	if s.status == ServerStatusStopped {
		return
	}
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, gracefulShutdownTimeout)
	defer cancelFunc()
	if err := s.httpServer.Shutdown(timeoutCtx); err != nil {
		s.server.Logger().Errorf(
			ctx,
			"%d: %s server [%s] shutdown error: %v",
			gproc.Pid(), s.getProto(), s.address, err,
		)
	}
}

// setRawListener sets `rawListener` with given net.Listener.
func (s *gracefulServer) setRawListener(ln net.Listener) {
	s.rawLnMu.Lock()
	defer s.rawLnMu.Unlock()
	s.rawListener = ln
}

// setRawListener returns the `rawListener` of current server.
func (s *gracefulServer) getRawListener() net.Listener {
	s.rawLnMu.RLock()
	defer s.rawLnMu.RUnlock()
	return s.rawListener
}

// close shuts down the server forcibly.
func (s *gracefulServer) close(ctx context.Context) {
	if s.status == ServerStatusStopped {
		return
	}
	if err := s.httpServer.Close(); err != nil {
		s.server.Logger().Errorf(
			ctx,
			"%d: %s server [%s] closed error: %v",
			gproc.Pid(), s.getProto(), s.address, err,
		)
	}
}

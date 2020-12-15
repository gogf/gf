// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gres"
	"github.com/gogf/gf/text/gstr"
	"log"
	"net"
	"net/http"
	"os"
)

// gracefulServer wraps the net/http.Server with graceful reload/restart feature.
type gracefulServer struct {
	server      *Server      // Belonged server.
	fd          uintptr      // File descriptor for passing to child process when graceful reload.
	address     string       // Listening address like:":80", ":8080".
	httpServer  *http.Server // Underlying http.Server.
	rawListener net.Listener // Underlying net.Listener.
	listener    net.Listener // Wrapped net.Listener.
	isHttps     bool         // Is HTTPS.
	status      int          // Status of current server.
}

// newGracefulServer creates and returns a graceful http server with given address.
// The optional parameter <fd> specifies the file descriptor which is passed from parent server.
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
	return gs
}

// newGracefulServer creates and returns a underlying http.Server with given address.
func (s *Server) newHttpServer(address string) *http.Server {
	server := &http.Server{
		Addr:           address,
		Handler:        s.config.Handler,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		IdleTimeout:    s.config.IdleTimeout,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
		ErrorLog:       log.New(&errorLogger{logger: s.config.Logger}, "", 0),
	}
	server.SetKeepAlivesEnabled(s.config.KeepAlive)
	return server
}

// ListenAndServe starts listening on configured address.
func (s *gracefulServer) ListenAndServe() error {
	ln, err := s.getNetListener()
	if err != nil {
		return err
	}
	s.listener = ln
	s.rawListener = ln
	return s.doServe()
}

// Fd retrieves and returns the file descriptor of current server.
// It is available ony in *nix like operation systems like: linux, unix, darwin.
func (s *gracefulServer) Fd() uintptr {
	if s.rawListener != nil {
		file, err := s.rawListener.(*net.TCPListener).File()
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

// ListenAndServeTLS starts listening on configured address with HTTPS.
// The parameter <certFile> and <keyFile> specify the necessary certification and key files for HTTPS.
// The optional parameter <tlsConfig> specifies the custom TLS configuration.
func (s *gracefulServer) ListenAndServeTLS(certFile, keyFile string, tlsConfig ...*tls.Config) error {
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
	err := error(nil)
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
		return errors.New(fmt.Sprintf(`open cert file "%s","%s" failed: %s`, certFile, keyFile, err.Error()))
	}
	ln, err := s.getNetListener()
	if err != nil {
		return err
	}

	s.listener = tls.NewListener(ln, config)
	s.rawListener = ln
	return s.doServe()
}

// getProto retrieves and returns the proto string of current server.
func (s *gracefulServer) getProto() string {
	proto := "http"
	if s.isHttps {
		proto = "https"
	}
	return proto
}

// doServe does staring the serving.
func (s *gracefulServer) doServe() error {
	action := "started"
	if s.fd != 0 {
		action = "reloaded"
	}
	s.server.Logger().Printf(
		"%d: %s server %s listening on [%s]",
		gproc.Pid(), s.getProto(), action, s.address,
	)
	s.status = ServerStatusRunning
	err := s.httpServer.Serve(s.listener)
	s.status = ServerStatusStopped
	return err
}

// getNetListener retrieves and returns the wrapped net.Listener.
func (s *gracefulServer) getNetListener() (net.Listener, error) {
	var ln net.Listener
	var err error
	if s.fd > 0 {
		f := os.NewFile(s.fd, "")
		ln, err = net.FileListener(f)
		if err != nil {
			err = fmt.Errorf("%d: net.FileListener error: %v", gproc.Pid(), err)
			return nil, err
		}
	} else {
		ln, err = net.Listen("tcp", s.httpServer.Addr)
		if err != nil {
			err = fmt.Errorf("%d: net.Listen error: %v", gproc.Pid(), err)
		}
	}
	return ln, err
}

// shutdown shuts down the server gracefully.
func (s *gracefulServer) shutdown() {
	if s.status == ServerStatusStopped {
		return
	}
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		s.server.Logger().Errorf(
			"%d: %s server [%s] shutdown error: %v",
			gproc.Pid(), s.getProto(), s.address, err,
		)
	}
}

// close shuts down the server forcibly.
func (s *gracefulServer) close() {
	if s.status == ServerStatusStopped {
		return
	}
	if err := s.httpServer.Close(); err != nil {
		s.server.Logger().Errorf(
			"%d: %s server [%s] closed error: %v",
			gproc.Pid(), s.getProto(), s.address, err,
		)
	}
}

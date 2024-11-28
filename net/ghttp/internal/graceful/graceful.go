// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package graceful implements graceful reload/restart features for HTTP servers.
// It provides the ability to gracefully shutdown or restart HTTP servers without
// interrupting existing connections. This is particularly useful for zero-downtime
// deployments and maintenance operations.
//
// The package wraps the standard net/http.Server and provides additional functionality
// for graceful server management, including:
// - Graceful server shutdown with timeout
// - Support for both HTTP and HTTPS servers
// - File descriptor inheritance for server reload/restart
// - Connection management during shutdown
package graceful

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/text/gstr"
)

// ServerStatus is the server status enum type.
type ServerStatus = int

const (
	// FreePortAddress marks the server listens using random free port.
	FreePortAddress = ":0"
	// ServerStatusStopped indicates the server is stopped.
	ServerStatusStopped ServerStatus = 0
	// ServerStatusRunning indicates the server is running.
	ServerStatusRunning ServerStatus = 1
)

// Server wraps the net/http.Server with graceful reload/restart feature.
type Server struct {
	fd          uintptr      // File descriptor for passing to the child process when graceful reload.
	address     string       // Listening address like ":80", ":8080".
	httpServer  *http.Server // Underlying http.Server.
	rawListener net.Listener // Underlying net.Listener.
	rawLnMu     sync.RWMutex // Concurrent safety mutex for rawListener.
	listener    net.Listener // Wrapped net.Listener with TLS support if necessary.
	isHttps     bool         // Whether server is running in HTTPS mode.
	status      *gtype.Int   // Server status using gtype for concurrent safety.
	config      ServerConfig // Server configuration.
}

// ServerConfig is the graceful Server configuration manager.
type ServerConfig struct {
	// Listeners specifies the custom listeners.
	Listeners []net.Listener `json:"listeners"`

	// Handler the handler for HTTP request.
	Handler func(w http.ResponseWriter, r *http.Request) `json:"-"`

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration `json:"readTimeout"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeout time.Duration `json:"writeTimeout"`

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alive are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	IdleTimeout time.Duration `json:"idleTimeout"`

	// GracefulShutdownTimeout set the maximum survival time (seconds) before stopping the server.
	GracefulShutdownTimeout int `json:"gracefulShutdownTimeout"`

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	//
	// It can be configured in configuration file using string like: 1m, 10m, 500kb etc.
	// It's 10240 bytes in default.
	MaxHeaderBytes int `json:"maxHeaderBytes"`

	// KeepAlive enables HTTP keep-alive.
	KeepAlive bool `json:"keepAlive"`

	// Logger specifies the logger for server.
	Logger *glog.Logger `json:"logger"`
}

// New creates and returns a graceful http server with a given address.
// The optional parameter `fd` specifies the file descriptor which is passed from parent server.
func New(
	address string,
	fd int,
	loggerWriter io.Writer,
	config ServerConfig,
) *Server {
	// Change port to address like: 80 -> :80
	if gstr.IsNumeric(address) {
		address = ":" + address
	}
	gs := &Server{
		address:    address,
		httpServer: newHttpServer(address, loggerWriter, config),
		status:     gtype.NewInt(),
		config:     config,
	}
	if fd != 0 {
		gs.fd = uintptr(fd)
	}
	if len(config.Listeners) > 0 {
		addrArray := gstr.SplitAndTrim(address, ":")
		addrPort, err := strconv.Atoi(addrArray[len(addrArray)-1])
		if err == nil {
			for _, v := range config.Listeners {
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
func newHttpServer(
	address string,
	loggerWriter io.Writer,
	config ServerConfig,
) *http.Server {
	server := &http.Server{
		Addr:           address,
		Handler:        http.HandlerFunc(config.Handler),
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		IdleTimeout:    config.IdleTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
		ErrorLog:       log.New(loggerWriter, "", 0),
	}
	server.SetKeepAlivesEnabled(config.KeepAlive)
	return server
}

// Fd retrieves and returns the file descriptor of the current server.
// It is available ony in *nix like operating systems like linux, unix, darwin.
func (s *Server) Fd() uintptr {
	if ln := s.getRawListener(); ln != nil {
		file, err := ln.(*net.TCPListener).File()
		if err == nil {
			return file.Fd()
		}
	}
	return 0
}

// CreateListener creates listener on configured address.
func (s *Server) CreateListener() error {
	ln, err := s.getNetListener()
	if err != nil {
		return err
	}
	s.listener = ln
	s.setRawListener(ln)
	return nil
}

// IsHttps returns whether the server is running in HTTPS mode.
func (s *Server) IsHttps() bool {
	return s.isHttps
}

// GetAddress returns the server's configured address.
func (s *Server) GetAddress() string {
	return s.address
}

// SetIsHttps sets the HTTPS mode for the server.
// The parameter isHttps determines whether to enable HTTPS mode.
func (s *Server) SetIsHttps(isHttps bool) {
	s.isHttps = isHttps
}

// CreateListenerTLS creates listener on configured address with HTTPS.
// The parameter `certFile` and `keyFile` specify the necessary certification and key files for HTTPS.
// The optional parameter `tlsConfig` specifies the custom TLS configuration.
func (s *Server) CreateListenerTLS(certFile, keyFile string, tlsConfig ...*tls.Config) error {
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
func (s *Server) Serve(ctx context.Context) error {
	if s.rawListener == nil {
		return gerror.NewCode(gcode.CodeInvalidOperation, `call CreateListener/CreateListenerTLS before Serve`)
	}

	var action = "started"
	if s.fd != 0 {
		action = "reloaded"
	}
	s.config.Logger.Infof(
		ctx,
		`pid[%d]: %s server %s listening on [%s]`,
		gproc.Pid(), s.getProto(), action, s.GetListenedAddress(),
	)
	s.status.Set(ServerStatusRunning)
	err := s.httpServer.Serve(s.listener)
	s.status.Set(ServerStatusStopped)
	return err
}

// GetListenedAddress retrieves and returns the address string which are listened by current server.
func (s *Server) GetListenedAddress() string {
	if !gstr.Contains(s.address, FreePortAddress) {
		return s.address
	}
	var (
		address      = s.address
		listenedPort = s.GetListenedPort()
	)
	address = gstr.Replace(address, FreePortAddress, fmt.Sprintf(`:%d`, listenedPort))
	return address
}

// GetListenedPort retrieves and returns one port which is listened to by current server.
// Note that this method is only available if the server is listening on one port.
func (s *Server) GetListenedPort() int {
	if ln := s.getRawListener(); ln != nil {
		return ln.Addr().(*net.TCPAddr).Port
	}
	return -1
}

// Status returns the current status of the server.
// It returns either ServerStatusStopped or ServerStatusRunning.
func (s *Server) Status() ServerStatus {
	return s.status.Val()
}

// getProto retrieves and returns the proto string of current server.
func (s *Server) getProto() string {
	proto := "http"
	if s.isHttps {
		proto = "https"
	}
	return proto
}

// getNetListener retrieves and returns the wrapped net.Listener.
func (s *Server) getNetListener() (net.Listener, error) {
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

// Shutdown shuts down the server gracefully.
func (s *Server) Shutdown(ctx context.Context) {
	if s.status.Val() == ServerStatusStopped {
		return
	}
	timeoutCtx, cancelFunc := context.WithTimeout(
		ctx,
		time.Duration(s.config.GracefulShutdownTimeout)*time.Second,
	)
	defer cancelFunc()
	if err := s.httpServer.Shutdown(timeoutCtx); err != nil {
		s.config.Logger.Errorf(
			ctx,
			"%d: %s server [%s] shutdown error: %v",
			gproc.Pid(), s.getProto(), s.address, err,
		)
	}
}

// setRawListener sets `rawListener` with given net.Listener.
func (s *Server) setRawListener(ln net.Listener) {
	s.rawLnMu.Lock()
	defer s.rawLnMu.Unlock()
	s.rawListener = ln
}

// getRawListener returns the `rawListener` of current server.
func (s *Server) getRawListener() net.Listener {
	s.rawLnMu.RLock()
	defer s.rawLnMu.RUnlock()
	return s.rawListener
}

// Close shuts down the server forcibly.
// for graceful shutdown, please use Server.shutdown.
func (s *Server) Close(ctx context.Context) {
	if s.status.Val() == ServerStatusStopped {
		return
	}
	if err := s.httpServer.Close(); err != nil {
		s.config.Logger.Errorf(
			ctx,
			"%d: %s server [%s] closed error: %v",
			gproc.Pid(), s.getProto(), s.address, err,
		)
	}
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Adapter implementations for built-in server types.

package gapp

import (
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/net/gudp"
)

// adapterStartTimeout is the maximum duration to wait for a server
// to indicate it is listening after Start() is called.
const adapterStartTimeout = time.Second * 2

// adapterStartPollInterval is the interval between readiness polls.
const adapterStartPollInterval = time.Millisecond * 10

// httpServerAdapter wraps ghttp.Server to implement the Server interface.
type httpServerAdapter struct {
	server *ghttp.Server
}

// NewHTTPServerAdapter creates and returns a Server adapter for ghttp.Server.
func NewHTTPServerAdapter(server *ghttp.Server) Server {
	return &httpServerAdapter{server: server}
}

// Start starts the HTTP server in non-blocking way.
func (a *httpServerAdapter) Start() error {
	return a.server.Start()
}

// Stop stops the HTTP server.
// When graceful is true, it waits for in-flight requests to complete.
// When graceful is false, it forcibly closes all connections.
func (a *httpServerAdapter) Stop(graceful bool) error {
	if graceful {
		return a.server.Shutdown()
	}
	return a.server.Close()
}

// tcpServerAdapter wraps gtcp.Server to implement the Server interface.
type tcpServerAdapter struct {
	server *gtcp.Server
}

// NewTCPServerAdapter creates and returns a Server adapter for gtcp.Server.
func NewTCPServerAdapter(server *gtcp.Server) Server {
	return &tcpServerAdapter{server: server}
}

// Start starts the TCP server in non-blocking way.
// Since gtcp.Server.Run() is blocking, it is launched in a goroutine
// and the method polls until the listener is ready.
func (a *tcpServerAdapter) Start() error {
	var (
		errCh = make(chan error, 1)
	)
	go func() {
		if err := a.server.Run(); err != nil {
			select {
			case errCh <- err:
			default:
			}
		}
	}()

	// Poll until the listener is established or timeout.
	deadline := time.Now().Add(adapterStartTimeout)
	for time.Now().Before(deadline) {
		if a.server.GetListenedPort() > 0 {
			return nil
		}
		select {
		case err := <-errCh:
			return gerror.WrapCode(gcode.CodeInternalError, err, "tcp server start failed")
		default:
		}
		time.Sleep(adapterStartPollInterval)
	}

	select {
	case err := <-errCh:
		return gerror.WrapCode(gcode.CodeInternalError, err, "tcp server start failed")
	default:
	}
	// Best-effort cleanup of the leaked goroutine and port on timeout.
	a.server.Close()
	return gerror.NewCode(gcode.CodeOperationFailed, "tcp server failed to start within timeout")
}

// Post-start errors from the underlying Run() goroutine are sent to errCh
// but are not read after Start() returns successfully. A server that fails
// after startup will appear healthy while being non-functional. Callers
// needing post-start error monitoring should wrap the adapter with their
// own error notification mechanism.

// Stop stops the TCP server.
// Since TCP has no built-in graceful shutdown concept, both graceful and forceful
// stop close the listener immediately.
func (a *tcpServerAdapter) Stop(_ bool) error {
	return a.server.Close()
}

// udpServerAdapter wraps gudp.Server to implement the Server interface.
type udpServerAdapter struct {
	server *gudp.Server
}

// NewUDPServerAdapter creates and returns a Server adapter for gudp.Server.
func NewUDPServerAdapter(server *gudp.Server) Server {
	return &udpServerAdapter{server: server}
}

// Start starts the UDP server in non-blocking way.
// Since gudp.Server.Run() is blocking, it is launched in a goroutine
// and the method polls until the connection is ready.
func (a *udpServerAdapter) Start() error {
	var (
		errCh = make(chan error, 1)
	)
	go func() {
		if err := a.server.Run(); err != nil {
			select {
			case errCh <- err:
			default:
			}
		}
	}()

	// Poll until the connection is established or timeout.
	deadline := time.Now().Add(adapterStartTimeout)
	for time.Now().Before(deadline) {
		if a.server.GetListenedPort() > 0 {
			return nil
		}
		select {
		case err := <-errCh:
			return gerror.WrapCode(gcode.CodeInternalError, err, "udp server start failed")
		default:
		}
		time.Sleep(adapterStartPollInterval)
	}

	select {
	case err := <-errCh:
		return gerror.WrapCode(gcode.CodeInternalError, err, "udp server start failed")
	default:
	}
	// Best-effort cleanup of the leaked goroutine and connection on timeout.
	a.server.Close()
	return gerror.NewCode(gcode.CodeOperationFailed, "udp server failed to start within timeout")
}

// Stop stops the UDP server.
// Since UDP has no built-in graceful shutdown concept, both graceful and forceful
// stop close the connection immediately.
func (a *udpServerAdapter) Stop(_ bool) error {
	return a.server.Close()
}

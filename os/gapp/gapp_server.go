// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Adapter implementations for built-in server types.

package gapp

import (
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/net/gudp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
)

// adapterStartTimeout is the maximum duration to wait for a server
// to indicate it is listening after Start() is called.
// Tests shorten it to reach the timeout path without waiting.
var adapterStartTimeout = time.Second * 2

// adapterStartPollInterval is the interval between readiness polls.
const adapterStartPollInterval = time.Millisecond * 10

// adapterServeGate, when non-nil, runs at the start of an adapter serve
// goroutine before the underlying server blocks in Run. Tests hold it so the
// startup timeout path is observed before a listener exists.
var adapterServeGate func()

// serveResult delivers a background Run error to Start, or logs it once Start
// has already reported success. An abandoned start ignores the error caused by
// the follow-up Close.
type serveResult struct {
	mu        sync.Mutex
	errCh     chan error
	accepted  bool
	abandoned bool
}

// newServeResult creates an empty serve result for one Start call.
func newServeResult() *serveResult {
	return &serveResult{errCh: make(chan error, 1)}
}

// watch records Run errors and closes the result channel even when Run exits cleanly.
func (r *serveResult) watch(run func() error) {
	defer close(r.errCh)
	if adapterServeGate != nil {
		adapterServeGate()
	}
	err := run()
	if err == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.abandoned {
		return
	}
	if r.accepted {
		glog.Errorf(gctx.GetInitCtx(), "server stopped after start: %+v", err)
		return
	}
	select {
	case r.errCh <- err:
	default:
		glog.Errorf(gctx.GetInitCtx(), "server stopped after start: %+v", err)
	}
}

// accept marks startup successful and returns a Run error that won the race.
func (r *serveResult) accept() error {
	r.mu.Lock()
	r.accepted = true
	r.mu.Unlock()
	err, _ := r.takeResult()
	return err
}

// markAbandoned ignores Run errors caused by closing an unready server.
func (r *serveResult) markAbandoned() {
	r.mu.Lock()
	r.abandoned = true
	r.mu.Unlock()
}

// takeResult distinguishes a completed Run from one that has not returned yet.
func (r *serveResult) takeResult() (error, bool) {
	select {
	case err := <-r.errCh:
		return err, true
	default:
		return nil, false
	}
}

// abandonUnreadyServer waits until the server is listening, then closes it.
// Waiting for the port prevents a late bind from escaping timeout cleanup.
// A close error is returned to the caller.
// serveErr also closes on a clean Run exit, so cleanup cannot wait for a socket that will never appear.
func abandonUnreadyServer(port func() int, closeFn func() error, serveErr <-chan error, message string) error {
	var serveDone bool
	for {
		if serveErr != nil {
			select {
			case err, open := <-serveErr:
				serveDone = !open
				if err != nil {
					return gerror.WrapCode(gcode.CodeInternalError, err, message)
				}
			default:
			}
		}
		if port() > 0 || serveDone {
			if closeErr := closeFn(); closeErr != nil {
				return gerror.WrapCode(gcode.CodeOperationFailed, closeErr, message)
			}
			return gerror.NewCode(gcode.CodeOperationFailed, message)
		}
		select {
		case err, open := <-serveErr:
			serveDone = !open
			if err != nil {
				return gerror.WrapCode(gcode.CodeInternalError, err, message)
			}
		case <-time.After(adapterStartPollInterval):
		}
	}
}

// startBlockingServer polls until a blocking server is listening.
// On timeout it keeps waiting only long enough to close the listener once the
// port appears, and returns that close error instead of leaving the server running.
func startBlockingServer(run func() error, port func() int, closeFn func() error, failPrefix, timeoutMessage string) error {
	// Reject another Run before it can replace the socket that Stop must close.
	if port() > 0 {
		return gerror.NewCode(gcode.CodeInvalidOperation, "server is already listening")
	}
	result := newServeResult()
	go result.watch(run)

	deadline := time.Now().Add(adapterStartTimeout)
	for time.Now().Before(deadline) {
		if port() > 0 {
			if err := result.accept(); err != nil {
				return gerror.WrapCode(gcode.CodeInternalError, err, failPrefix)
			}
			return nil
		}
		if err, completed := result.takeResult(); err != nil {
			return gerror.WrapCode(gcode.CodeInternalError, err, failPrefix)
		} else if completed {
			timeoutMessage = "server stopped before listening"
			break
		}
		time.Sleep(adapterStartPollInterval)
	}
	if err, _ := result.takeResult(); err != nil {
		return gerror.WrapCode(gcode.CodeInternalError, err, failPrefix)
	}
	return abandonUnreadyServer(port, func() error {
		result.markAbandoned()
		return closeFn()
	}, result.errCh, timeoutMessage)
}

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
	return a.server.StartManaged()
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
	return startBlockingServer(
		a.server.Run,
		a.server.GetListenedPort,
		a.server.Close,
		"tcp server start failed",
		"tcp server failed to start within timeout",
	)
}

// Stop stops the TCP server.
// Since TCP has no built-in graceful shutdown concept, both graceful and forceful
// stop close the listener immediately. Stop before the listener exists is a no-op
// so a nil listener is not treated as a successful close of a later bind.
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
	return startBlockingServer(
		a.server.Run,
		a.server.GetListenedPort,
		a.server.Close,
		"udp server start failed",
		"udp server failed to start within timeout",
	)
}

// Stop stops the UDP server.
// Since UDP has no built-in graceful shutdown concept, both graceful and forceful
// stop close the connection immediately. Stop before the connection exists is a no-op.
func (a *udpServerAdapter) Stop(_ bool) error {
	return a.server.Close()
}

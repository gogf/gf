// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests for boot failure, cleanup-hook detachment, and adapter startup timeout.

package gapp

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/net/gudp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestRunHooksReverseDetachesBackingArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			calledA int32
			calledC int32
		)
		started := make(chan struct{})
		appended := make(chan struct{})
		app := New()
		app.hooks = []func(context.Context){
			func(context.Context) {
				atomic.AddInt32(&calledA, 1)
			},
			func(context.Context) {
				close(started)
				<-appended
			},
		}

		done := make(chan struct{})
		go func() {
			app.runHooksReverse(context.Background())
			close(done)
		}()

		<-started
		app.mu.Lock()
		app.hooks = append(app.hooks, func(context.Context) {
			atomic.AddInt32(&calledC, 1)
		})
		app.mu.Unlock()
		close(appended)
		<-done

		t.Assert(atomic.LoadInt32(&calledA), int32(1))
		t.Assert(atomic.LoadInt32(&calledC), int32(0))

		app.runHooksReverse(context.Background())
		t.Assert(atomic.LoadInt32(&calledC), int32(1))
	})
}

func TestBootPanicDoesNotMarkBooted(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			applies  int32
			cleanups int32
			waitOnce sync.Once
		)
		inApply := make(chan struct{})
		releasePanic := make(chan struct{})
		waiterWaiting := make(chan struct{})
		bootWaitHook = func() {
			waitOnce.Do(func() {
				close(waiterWaiting)
			})
		}
		defer func() {
			bootWaitHook = nil
		}()

		app := New()
		app.Option(NewOptionWithHook(func(ctx context.Context, a *App) (func(ctx context.Context), error) {
			return func(ctx context.Context) {
				atomic.AddInt32(&cleanups, 1)
			}, nil
		}))
		app.Option(NewOption(func(ctx context.Context, a *App) {
			if atomic.AddInt32(&applies, 1) == 1 {
				close(inApply)
				<-releasePanic
				panic("boot boom")
			}
		}))

		leaderErr := make(chan error, 1)
		waiterErr := make(chan error, 1)
		go func() {
			leaderErr <- app.Boot(context.Background())
		}()
		<-inApply
		go func() {
			waiterErr <- app.Boot(context.Background())
		}()
		<-waiterWaiting
		close(releasePanic)

		err := <-leaderErr
		t.AssertNE(err, nil)
		t.Assert(strings.Contains(err.Error(), "app boot failed"), true)
		t.Assert(strings.Contains(err.Error(), "boot boom"), true)
		waited := <-waiterErr
		t.AssertNE(waited, nil)
		t.Assert(strings.Contains(waited.Error(), "app boot failed"), true)
		t.Assert(app.Booted(), false)
		t.Assert(atomic.LoadInt32(&cleanups), int32(1))

		err = app.Boot(context.Background())
		t.AssertNil(err)
		t.Assert(app.Booted(), true)
		t.Assert(atomic.LoadInt32(&applies), int32(2))
	})
}

func TestAbandonUnreadyServerClosesWhenPortAppears(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			port   int32 = -1
			called int32
		)
		done := make(chan error, 1)
		go func() {
			done <- abandonUnreadyServer(func() int {
				return int(atomic.LoadInt32(&port))
			}, func() error {
				atomic.StoreInt32(&called, 1)
				return nil
			}, nil, "udp server failed to start within timeout")
		}()

		time.Sleep(30 * time.Millisecond)
		t.Assert(atomic.LoadInt32(&called), int32(0))
		atomic.StoreInt32(&port, 1)

		err := <-done
		t.Assert(atomic.LoadInt32(&called), int32(1))
		t.AssertNE(err, nil)
		t.Assert(strings.Contains(err.Error(), "udp server failed to start within timeout"), true)
	})
}

func TestAbandonUnreadyServerClosesWhenListening(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		called := false
		err := abandonUnreadyServer(func() int { return 1 }, func() error {
			called = true
			return nil
		}, nil, "tcp server failed to start within timeout")
		t.Assert(called, true)
		t.AssertNE(err, nil)
		t.Assert(strings.Contains(err.Error(), "tcp server failed to start within timeout"), true)
		t.Assert(strings.Contains(err.Error(), "context canceled"), false)
	})
}

func TestAbandonUnreadyServerReturnsCloseError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := abandonUnreadyServer(func() int { return 1 }, func() error {
			return context.Canceled
		}, nil, "tcp server failed to start within timeout")
		t.AssertNE(err, nil)
		t.Assert(strings.Contains(err.Error(), "tcp server failed to start within timeout"), true)
		t.Assert(strings.Contains(err.Error(), context.Canceled.Error()), true)
	})
}

func TestAbandonUnreadyServerReturnsServeError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		serveErr := make(chan error, 1)
		serveErr <- context.Canceled
		called := false
		err := abandonUnreadyServer(func() int { return -1 }, func() error {
			called = true
			return nil
		}, serveErr, "tcp server failed to start within timeout")
		t.Assert(called, false)
		t.AssertNE(err, nil)
		t.Assert(strings.Contains(err.Error(), context.Canceled.Error()), true)
	})
}

func TestStopBeforeListenDoesNotPanic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tcpServer := gtcp.NewServer("127.0.0.1:0", func(*gtcp.Conn) {})
		err := (&tcpServerAdapter{server: tcpServer}).Stop(false)
		t.AssertNil(err)

		udpServer := gudp.NewServer("127.0.0.1:0", func(*gudp.ServerConn) {})
		err = (&udpServerAdapter{server: udpServer}).Stop(false)
		t.AssertNil(err)
	})
}

func TestTCPAndUDPStartTimeoutClosesWhenPortAppears(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		previousTimeout := adapterStartTimeout
		adapterStartTimeout = 0
		defer func() {
			adapterStartTimeout = previousTimeout
			adapterServeGate = nil
		}()

		release := make(chan struct{})
		adapterServeGate = func() {
			<-release
		}
		go func() {
			time.Sleep(20 * time.Millisecond)
			close(release)
		}()

		tcpServer := gtcp.NewServer("127.0.0.1:0", func(*gtcp.Conn) {})
		tcpAdapter := &tcpServerAdapter{server: tcpServer}
		tcpErr := tcpAdapter.Start()
		t.AssertNE(tcpErr, nil)
		t.Assert(strings.Contains(tcpErr.Error(), "tcp server failed to start within timeout"), true)
		t.Assert(tcpServer.GetListenedPort(), -1)
		t.AssertNil(tcpAdapter.Stop(false))

		udpRelease := make(chan struct{})
		udpReceived := make(chan struct{}, 1)
		adapterServeGate = func() {
			<-udpRelease
		}
		go func() {
			time.Sleep(20 * time.Millisecond)
			close(udpRelease)
		}()
		udpServer := gudp.NewServer("127.0.0.1:0", func(conn *gudp.ServerConn) {
			data, _, err := conn.Recv(-1)
			if err == nil && len(data) > 0 {
				udpReceived <- struct{}{}
			}
		})
		udpAdapter := &udpServerAdapter{server: udpServer}
		udpErr := udpAdapter.Start()
		t.AssertNE(udpErr, nil)
		t.Assert(strings.Contains(udpErr.Error(), "udp server failed to start within timeout"), true)
		t.Assert(udpServer.GetListenedPort(), -1)
		select {
		case <-udpReceived:
			t.Fatal("udp server still accepted data after startup was abandoned")
		case <-time.After(200 * time.Millisecond):
		}
		t.AssertNil(udpAdapter.Stop(false))
	})
}

func TestTCPServerLogsServeErrorAfterStart(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var buf lockedBuffer
		logger := glog.NewWithWriter(&buf)
		logger.SetStdoutPrint(false)
		previous := glog.DefaultLogger()
		glog.SetDefaultLogger(logger)
		defer glog.SetDefaultLogger(previous)

		tcpServer := gtcp.NewServer("127.0.0.1:0", func(*gtcp.Conn) {})
		adapter := &tcpServerAdapter{server: tcpServer}
		t.AssertNil(adapter.Start())
		t.AssertNil(tcpServer.Close())

		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) && !strings.Contains(buf.String(), "server stopped after start") {
			time.Sleep(10 * time.Millisecond)
		}
		t.Assert(strings.Contains(buf.String(), "server stopped after start"), true)
	})
}

// lockedBuffer is a concurrent-safe buffer for log capture.
type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// Write implements io.Writer.
func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// String returns the captured text.
func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

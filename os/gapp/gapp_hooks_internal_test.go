// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests for boot failure, cleanup-hook detachment, and adapter startup timeout.

package gapp

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/net/gudp"
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

func TestAbandonUnreadyServerSkipsCloseBeforeListen(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		called := false
		err := abandonUnreadyServer(-1, func() error {
			called = true
			return context.Canceled
		}, "udp server failed to start within timeout")
		t.Assert(called, false)
		t.Assert(strings.Contains(err.Error(), "udp server failed to start within timeout"), true)
		t.Assert(strings.Contains(err.Error(), context.Canceled.Error()), false)
	})
}

func TestAbandonUnreadyServerClosesWhenListening(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		called := false
		err := abandonUnreadyServer(1, func() error {
			called = true
			return nil
		}, "tcp server failed to start within timeout")
		t.Assert(called, true)
		t.AssertNE(err, nil)
		t.Assert(strings.Contains(err.Error(), "tcp server failed to start within timeout"), true)
		t.Assert(strings.Contains(err.Error(), "context canceled"), false)
	})
}

func TestAbandonUnreadyServerReturnsCloseError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := abandonUnreadyServer(1, func() error {
			return context.Canceled
		}, "tcp server failed to start within timeout")
		t.AssertNE(err, nil)
		t.Assert(strings.Contains(err.Error(), "tcp server failed to start within timeout"), true)
		t.Assert(strings.Contains(err.Error(), context.Canceled.Error()), true)
	})
}

func TestTCPAndUDPStartTimeoutDoNotCloseBeforeListen(t *testing.T) {
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

		tcpServer := gtcp.NewServer("127.0.0.1:0", func(*gtcp.Conn) {})
		tcpErr := (&tcpServerAdapter{server: tcpServer}).Start()
		t.AssertNE(tcpErr, nil)
		t.Assert(strings.Contains(tcpErr.Error(), "tcp server failed to start within timeout"), true)

		udpServer := gudp.NewServer("127.0.0.1:0", func(*gudp.ServerConn) {})
		udpErr := (&udpServerAdapter{server: udpServer}).Start()
		t.AssertNE(udpErr, nil)
		t.Assert(strings.Contains(udpErr.Error(), "udp server failed to start within timeout"), true)

		close(release)
		waitUntilListening(t, tcpServer.GetListenedPort, udpServer.GetListenedPort)
		t.AssertNil(tcpServer.Close())
		t.AssertNil(udpServer.Close())
	})
}

// waitUntilListening waits until both port readers report a listened port.
func waitUntilListening(t *gtest.T, ports ...func() int) {
	deadline := time.Now().Add(2 * time.Second)
	for _, port := range ports {
		for time.Now().Before(deadline) && port() <= 0 {
			time.Sleep(time.Millisecond)
		}
		t.Assert(port() > 0, true)
	}
}

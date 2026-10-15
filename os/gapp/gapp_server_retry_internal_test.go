// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests for TCP/UDP adapter retries after timeout and application rollback.

package gapp

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/net/gudp"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestBlockingAdapterRejectsDuplicateStart preserves the original socket until Stop.
func TestBlockingAdapterRejectsDuplicateStart(t *testing.T) {
	for _, protocol := range []string{"tcp", "udp"} {
		t.Run(protocol, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					adapter Server
					port    func() int
					address func() string
				)
				if protocol == "tcp" {
					server := gtcp.NewServer("127.0.0.1:0", func(conn *gtcp.Conn) {
						if err := conn.Close(); err != nil {
							t.Error(err)
						}
					})
					adapter, port, address = NewTCPServerAdapter(server), server.GetListenedPort, server.GetListenedAddress
				} else {
					server := gudp.NewServer("127.0.0.1:0", func(conn *gudp.ServerConn) {
						if _, _, err := conn.Recv(-1); err != nil && !gerror.Is(err, net.ErrClosed) {
							t.Error(err)
						}
					})
					adapter, port, address = NewUDPServerAdapter(server), server.GetListenedPort, server.GetListenedAddress
				}
				defer func() { t.AssertNil(adapter.Stop(false)) }()
				app := New(adapter)
				t.AssertNil(app.Start(context.Background()))
				originalAddress := address()
				err := app.Start(context.Background())
				t.AssertNE(err, nil)
				t.Assert(gerror.Code(err), gcode.CodeInternalError)
				t.Assert(gerror.Code(adapter.Start()), gcode.CodeInvalidOperation)
				t.Assert(address(), originalAddress)
				t.AssertNil(app.Stop(context.Background(), false))
				t.Assert(port(), -1)
				if protocol == "tcp" {
					listener, err := net.Listen(protocol, originalAddress)
					t.AssertNil(err)
					t.AssertNil(listener.Close())
				} else {
					listener, err := net.ListenPacket(protocol, originalAddress)
					t.AssertNil(err)
					t.AssertNil(listener.Close())
				}
				t.AssertNil(adapter.Start())
				t.AssertGT(port(), 0)
				t.AssertNil(adapter.Stop(false))
				t.Assert(port(), -1)
			})
		})
	}
}

// TestBlockingAdapterRetryUsesNewListener proves retry waits for its own bind and Stop closes it.
func TestBlockingAdapterRetryUsesNewListener(t *testing.T) {
	for _, protocol := range []string{"tcp", "udp"} {
		for _, timeout := range []bool{true, false} {
			name := protocol + "/rollback"
			if timeout {
				name = protocol + "/timeout"
			}
			t.Run(name, func(t *testing.T) {
				gtest.C(t, func(t *gtest.T) {
					var (
						adapter Server
						port    func() int
						address func() string
					)
					accepted := make(chan struct{}, 1)
					if protocol == "tcp" {
						server := gtcp.NewServer("127.0.0.1:0", func(conn *gtcp.Conn) {
							if err := conn.Close(); err != nil {
								t.Error(err)
							}
							accepted <- struct{}{}
						})
						adapter, port, address = NewTCPServerAdapter(server), server.GetListenedPort, server.GetListenedAddress
					} else {
						server := gudp.NewServer("127.0.0.1:0", func(conn *gudp.ServerConn) {
							if _, _, err := conn.Recv(-1); err == nil {
								accepted <- struct{}{}
							}
						})
						adapter, port, address = NewUDPServerAdapter(server), server.GetListenedPort, server.GetListenedAddress
					}
					app := New(adapter)
					previousTimeout := adapterStartTimeout
					defer func() {
						adapterStartTimeout = previousTimeout
						adapterServeGate = nil
						t.AssertNil(adapter.Stop(false))
					}()
					if timeout {
						adapterStartTimeout = 0
						t.AssertNE(adapter.Start(), nil)
						adapterStartTimeout = previousTimeout
					} else {
						app.Add(&retryFailureServer{})
						t.AssertNE(app.Start(context.Background()), nil)
					}

					var releaseOnce sync.Once
					entered := make(chan struct{})
					release := make(chan struct{})
					adapterServeGate = func() {
						close(entered)
						<-release
					}
					defer releaseOnce.Do(func() { close(release) })
					done := make(chan error, 1)
					go func() { done <- app.Start(context.Background()) }()
					<-entered
					select {
					case err := <-done:
						t.Fatalf("retry returned before its Run could listen: %v", err)
					case <-time.After(30 * time.Millisecond):
					}
					releaseOnce.Do(func() { close(release) })
					select {
					case err := <-done:
						t.AssertNil(err)
					case <-time.After(3 * time.Second):
						t.Fatal("retry did not finish after Run was released")
					}
					t.AssertGT(port(), 0)
					listenedAddress := address()
					conn, err := net.DialTimeout(protocol, listenedAddress, time.Second)
					t.AssertNil(err)
					if protocol == "udp" {
						_, err = conn.Write([]byte("retry"))
						t.AssertNil(err)
					}
					t.AssertNil(conn.Close())
					select {
					case <-accepted:
					case <-time.After(time.Second):
						t.Fatal("retry listener did not serve traffic")
					}
					t.AssertNil(app.Stop(context.Background(), false))
					t.Assert(port(), -1)
					if protocol == "tcp" {
						listener, err := net.Listen(protocol, listenedAddress)
						t.AssertNil(err)
						t.AssertNil(listener.Close())
					} else {
						listener, err := net.ListenPacket(protocol, listenedAddress)
						t.AssertNil(err)
						t.AssertNil(listener.Close())
					}
				})
			})
		}
	}
}

// TestUDPAdapterRetryAfterHandlerClose waits for a new bind after handler-initiated shutdown.
func TestUDPAdapterRetryAfterHandlerClose(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := gudp.NewServer("127.0.0.1:0", func(conn *gudp.ServerConn) {
			t.AssertNil(conn.Close())
		})
		t.AssertNil(server.Run())
		adapter := NewUDPServerAdapter(server)
		t.AssertNil(adapter.Stop(false))
		t.Assert(server.GetListenedPort(), -1)

		var releaseOnce sync.Once
		entered := make(chan struct{})
		release := make(chan struct{})
		received := make(chan error, 1)
		server.SetHandler(func(conn *gudp.ServerConn) {
			_, _, err := conn.Recv(-1)
			received <- err
		})
		adapterServeGate = func() {
			close(entered)
			<-release
		}
		defer func() {
			releaseOnce.Do(func() { close(release) })
			adapterServeGate = nil
			t.AssertNil(adapter.Stop(false))
		}()
		done := make(chan error, 1)
		go func() { done <- adapter.Start() }()
		select {
		case <-entered:
		case <-time.After(time.Second):
			t.Fatal("retry did not enter the serve gate")
		}
		select {
		case err := <-done:
			t.Fatalf("retry returned before its Run could listen: %v", err)
		case <-time.After(30 * time.Millisecond):
		}
		releaseOnce.Do(func() { close(release) })
		select {
		case err := <-done:
			t.AssertNil(err)
		case <-time.After(3 * time.Second):
			t.Fatal("retry did not finish after the serve gate was released")
		}
		client, err := net.DialTimeout("udp", server.GetListenedAddress(), time.Second)
		t.AssertNil(err)
		_, err = client.Write([]byte("retry"))
		t.AssertNil(err)
		t.AssertNil(client.Close())
		select {
		case err := <-received:
			t.AssertNil(err)
		case <-time.After(time.Second):
			t.Fatal("retry listener did not receive the datagram")
		}
		t.AssertNil(adapter.Stop(false))
		t.Assert(server.GetListenedPort(), -1)
	})
}

// TestUDPAdapterRejectsEarlyClose returns promptly when its handler closes the socket before readiness.
func TestUDPAdapterRejectsEarlyClose(t *testing.T) {
	for _, closeServer := range []bool{true, false} {
		name := "connection"
		if closeServer {
			name = "server"
		}
		t.Run(name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				server := gudp.NewServer("127.0.0.1:0", nil)
				server.SetHandler(func(conn *gudp.ServerConn) {
					if closeServer {
						t.AssertNil(server.Close())
					} else {
						t.AssertNil(conn.Close())
					}
				})
				adapter := NewUDPServerAdapter(server)
				defer func() { t.AssertNil(adapter.Stop(false)) }()
				result := make(chan error, 1)
				go func() { result <- adapter.Start() }()
				select {
				case err := <-result:
					t.AssertNE(err, nil)
				case <-time.After(time.Second):
					t.Fatal("startup did not observe the handler's clean exit")
				}
			})
		})
	}
}

// TestUDPAdapterSupportsAsyncHandler preserves handlers that transfer their live socket to a goroutine.
func TestUDPAdapterSupportsAsyncHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			received = make(chan string, 1)
			exited   = make(chan struct{})
			returned = make(chan struct{})
		)
		server := gudp.NewServer("127.0.0.1:0", func(conn *gudp.ServerConn) {
			go func() {
				defer close(exited)
				data, _, err := conn.Recv(-1)
				if err != nil {
					if !gerror.Is(err, net.ErrClosed) {
						t.Error(err)
					}
					return
				}
				received <- string(data)
			}()
			close(returned)
		})
		adapter := NewUDPServerAdapter(server)
		defer func() {
			t.AssertNil(adapter.Stop(false))
			<-exited
		}()
		t.AssertNil(adapter.Start())
		<-returned
		client, err := net.DialTimeout("udp", server.GetListenedAddress(), time.Second)
		t.AssertNil(err)
		defer func() { t.AssertNil(client.Close()) }()
		_, err = client.Write([]byte("async"))
		t.AssertNil(err)
		select {
		case data := <-received:
			t.Assert(data, "async")
		case <-time.After(time.Second):
			t.Fatal("asynchronous handler did not receive the datagram")
		}
	})
}

// TestServeResultCleanExit keeps completion observable before and after startup acceptance.
func TestServeResultCleanExit(t *testing.T) {
	for _, accepted := range []bool{false, true} {
		name := "before-start"
		if accepted {
			name = "after-start"
		}
		t.Run(name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				result := newServeResult()
				if accepted {
					t.AssertNil(result.accept())
				}
				result.watch(func() error { return nil })
				err, completed := result.takeResult()
				t.AssertNil(err)
				t.Assert(completed, true)
				t.AssertNil(result.accept())
			})
		})
	}
}

// TestAbandonUnreadyServerCleanExit cleans a completed Run without waiting for nonexistent readiness.
func TestAbandonUnreadyServerCleanExit(t *testing.T) {
	for _, closeBeforeWait := range []bool{true, false} {
		name := "during-wait"
		if closeBeforeWait {
			name = "before-wait"
		}
		t.Run(name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					serveErr = make(chan error)
					calls    int
					closed   int
				)
				if closeBeforeWait {
					close(serveErr)
					closed++
				}
				err := abandonUnreadyServer(func() int {
					if closed == 0 {
						close(serveErr)
						closed++
					}
					return -1
				}, func() error {
					calls++
					return nil
				}, serveErr, "startup stopped")
				t.AssertNE(err, nil)
				t.Assert(calls, 1)
			})
		})
	}
}

// retryFailureServer fails once to trigger rollback of another server.
type retryFailureServer struct {
	failed bool
}

// Start fails its first call and succeeds on retry.
func (s *retryFailureServer) Start() error {
	if !s.failed {
		s.failed = true
		return context.Canceled
	}
	return nil
}

// Stop has no resources to release.
func (*retryFailureServer) Stop(bool) error {
	return nil
}

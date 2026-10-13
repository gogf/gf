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

	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/net/gudp"
	"github.com/gogf/gf/v2/test/gtest"
)

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

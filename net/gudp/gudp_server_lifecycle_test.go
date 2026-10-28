// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests for connection cleanup and restart after closing a UDP server.

package gudp_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/net/gudp"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestServerCloseClearsConnection verifies idempotent close and readiness on retry.
func TestServerCloseClearsConnection(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := gudp.NewServer("127.0.0.1:0", func(conn *gudp.ServerConn) {
			if _, _, err := conn.Recv(-1); err == nil {
				t.Error("receive succeeded without a datagram")
			}
		})
		t.AssertNil(server.Close())
		for range 2 {
			done := make(chan error, 1)
			go func() { done <- server.Run() }()
			deadline := time.Now().Add(time.Second)
			for server.GetListenedPort() <= 0 && time.Now().Before(deadline) {
				time.Sleep(time.Millisecond)
			}
			t.AssertGT(server.GetListenedPort(), 0)
			t.AssertNil(server.Close())
			t.Assert(server.GetListenedPort(), -1)
			t.AssertNil(server.Close())
			select {
			case err := <-done:
				t.AssertNil(err)
			case <-time.After(time.Second):
				t.Fatal("UDP handler did not exit after close")
			}
		}
	})
}

// TestServerCloseAfterHandlerClose clears readiness when the handler closed the socket first.
func TestServerCloseAfterHandlerClose(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := gudp.NewServer("127.0.0.1:0", func(conn *gudp.ServerConn) {
			t.AssertNil(conn.Close())
		})
		t.AssertNil(server.Run())
		t.Assert(server.GetListenedPort(), -1)
		t.AssertNil(server.Close())
		t.Assert(server.GetListenedPort(), -1)
		t.AssertNil(server.Close())
	})
}

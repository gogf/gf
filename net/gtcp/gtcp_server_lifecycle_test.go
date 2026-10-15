// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests for listener cleanup and restart after closing a TCP server.

package gtcp_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestServerCloseClearsListener verifies idempotent close and listener readiness on retry.
func TestServerCloseClearsListener(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := gtcp.NewServer("127.0.0.1:0", func(conn *gtcp.Conn) {
			if err := conn.Close(); err != nil {
				t.Error(err)
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
				t.AssertNE(err, nil)
			case <-time.After(time.Second):
				t.Fatal("TCP serve loop did not exit after close")
			}
		}
	})
}

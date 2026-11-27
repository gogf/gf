// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

// Test_Grpcx_Grpc_Server_GetListenedPort_Concurrent asserts that reading the listened port
// while the server starts does not race with the listener assignment.
//
// Start runs the server asynchronously, and the listener is assigned by that goroutine. Reading
// the port right after Start, as the tests and the users do, must be concurrency safe, or the
// race detector reports a data race and the server may return a wrong port.
func Test_Grpcx_Grpc_Server_GetListenedPort_Concurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			c = grpcx.Server.NewConfig()
		)
		c.Name = guid.S()
		c.Address = "127.0.0.1:0"
		s := grpcx.Server.New(c)
		s.Start()
		defer s.Stop()

		// Read the port concurrently with the asynchronous listener assignment. The port may
		// be -1 before the server starts listening, which is expected and documented.
		var done = make(chan struct{})
		go func() {
			defer close(done)
			for i := 0; i < 1000; i++ {
				s.GetListenedPort()
			}
		}()
		<-done

		// Wait for the asynchronous start to complete before asserting the port.
		for i := 0; i < 100 && s.GetListenedPort() <= 0; i++ {
			time.Sleep(time.Millisecond * 10)
		}
		t.Assert(s.GetListenedPort() > 0, true)
	})
}

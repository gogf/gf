// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx_test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

const grpcServerStopBeforeServeEnv = "GRPCX_TEST_STOP_BEFORE_SERVE"

// Test_Grpcx_Grpc_Server_Stop_Before_Serve asserts that stopping the server before the
// serving goroutine enters Serve does not exit the process. Serve returns
// grpc.ErrServerStopped in that case, which is a normal stop instead of a server error.
// The server is stopped before Start, so Serve deterministically observes a stopped server
// instead of depending on the goroutine scheduling. The scenario runs in a child process,
// as the process used to be killed by the fatal logging on grpc.ErrServerStopped.
func Test_Grpcx_Grpc_Server_Stop_Before_Serve(t *testing.T) {
	if os.Getenv(grpcServerStopBeforeServeEnv) == "1" {
		c := grpcx.Server.NewConfig()
		c.Name = guid.S()
		c.Address = "127.0.0.1:0"
		s := grpcx.Server.New(c)
		s.Stop()
		s.Start()
		// Run() continues and registers the service after the stop, wait for it and
		// deregister it, so no stale service is left in the registry.
		time.Sleep(time.Millisecond * 500)
		s.Stop()
		return
	}
	gtest.C(t, func(t *gtest.T) {
		cmd := exec.Command(os.Args[0], "-test.run=^Test_Grpcx_Grpc_Server_Stop_Before_Serve$")
		cmd.Env = append(os.Environ(), grpcServerStopBeforeServeEnv+"=1")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal("server exited from a normal stop: ", err, "\n", string(output))
		}
	})
}

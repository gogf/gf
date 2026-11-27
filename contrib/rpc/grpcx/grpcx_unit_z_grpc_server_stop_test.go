// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

// serverStopHelperEnv marks the child process that stops the server before it serves.
const serverStopHelperEnv = "GF_TEST_GRPCX_STOP_BEFORE_SERVE_HELPER"

// runServerStopBeforeServeHelper stops the server before it starts serving. grpc.Server.Serve
// then returns grpc.ErrServerStopped, which must not be handled as a fatal error.
func runServerStopBeforeServeHelper() {
	c := grpcx.Server.NewConfig()
	c.Name = guid.S()
	c.Address = "127.0.0.1:0"
	s := grpcx.Server.New(c)
	// Stopping before Start makes the ordering deterministic: Serve always observes an
	// already stopped server, whatever the goroutine scheduling is.
	s.Stop()
	s.Start()
	// Give the serving goroutine enough time to hit grpc.ErrServerStopped and, on the
	// faulty implementation, to call Fatalf which exits the process.
	time.Sleep(500 * time.Millisecond)
}

// Test_Grpcx_Grpc_Server_StopBeforeServeDoesNotExitProcess asserts that stopping the server
// right after starting it does not terminate the process.
//
// The helper runs in a child process, because glog.Logger.Fatalf calls os.Exit and would
// terminate the test binary itself.
func Test_Grpcx_Grpc_Server_StopBeforeServeDoesNotExitProcess(t *testing.T) {
	if os.Getenv(serverStopHelperEnv) != "" {
		runServerStopBeforeServeHelper()
		return
	}

	var (
		testName = t.Name()
	)
	gtest.C(t, func(t *gtest.T) {
		cmd := exec.Command(os.Args[0], "-test.run=^"+testName+"$")
		cmd.Env = append(os.Environ(), serverStopHelperEnv+"=1")
		output, err := cmd.CombinedOutput()
		// A fatal log exits the child process with a non-zero status code.
		if err != nil {
			t.Fatal(fmt.Sprintf("child process exited unexpectedly: %v\n%s", err, output))
		}
	})
}

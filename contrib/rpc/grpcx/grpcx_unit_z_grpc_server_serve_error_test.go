// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests for Server.Serve failure handling on Run versus StartManaged.

package grpcx

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// testRunServeFatalEnv enables the subprocess path that expects Run() to Fatalf.
const testRunServeFatalEnv = "GF_GRPCX_TEST_RUN_SERVE_FATAL"

// errAcceptFailed is returned by acceptFailListener to force Server.Serve to fail.
var errAcceptFailed = errors.New("accept failed")

// safeBuffer is a concurrent-safe bytes.Buffer for capturing logs from Serve goroutines.
type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// Write implements io.Writer.
func (b *safeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// String returns the captured log content.
func (b *safeBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// acceptFailListener is a net.Listener whose Accept always fails.
type acceptFailListener struct{}

// Accept always returns errAcceptFailed.
func (l *acceptFailListener) Accept() (net.Conn, error) {
	return nil, errAcceptFailed
}

// Close implements net.Listener.
func (l *acceptFailListener) Close() error {
	return nil
}

// Addr implements net.Listener.
func (l *acceptFailListener) Addr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
}

// newServeErrorTestServer creates a GrpcServer that writes logs to w.
func newServeErrorTestServer(w io.Writer) *GrpcServer {
	c := Server.NewConfig()
	c.Name = guid.S()
	c.Address = "127.0.0.1:0"
	logger := glog.NewWithWriter(w)
	logger.SetStdoutPrint(false)
	c.Logger = logger
	c.LogStdout = false
	return Server.New(c)
}

// closeListenerWhenReady closes the server listener after it is bound.
func closeListenerWhenReady(s *GrpcServer) error {
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if s.GetListenedPort() > 0 {
			s.listenerMu.RLock()
			ln := s.listener
			s.listenerMu.RUnlock()
			if ln != nil {
				return ln.Close()
			}
		}
		time.Sleep(time.Millisecond)
	}
	return errors.New("listener was not ready to close")
}

// runServeErrorFatalChild starts Run() and closes the listener so Serve fails.
func runServeErrorFatalChild() {
	logger := glog.New()
	logger.SetStdoutPrint(true)
	c := Server.NewConfig()
	c.Name = guid.S()
	c.Address = "127.0.0.1:0"
	c.Logger = logger
	c.LogStdout = false
	s := Server.New(c)
	go func() {
		if err := closeListenerWhenReady(s); err != nil {
			os.Exit(2)
		}
	}()
	s.Run()
}

func Test_doServeAsynchronously_LogsServeErrorWhenManaged(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var buf safeBuffer
		s := newServeErrorTestServer(&buf)
		s.doServeAsynchronously(
			gctx.GetInitCtx(),
			&acceptFailListener{},
			serveFailureLog,
		)
		logged := buf.String()
		t.Assert(strings.Contains(logged, "grpc server serve error"), true)
		t.Assert(strings.Contains(logged, errAcceptFailed.Error()), true)
		t.Assert(strings.Contains(logged, "FATA"), false)
	})
}

func Test_doServeAsynchronously_IgnoresErrServerStopped(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var buf safeBuffer
		s := newServeErrorTestServer(&buf)
		s.Server.Stop()
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		t.AssertNil(err)

		s.doServeAsynchronously(gctx.GetInitCtx(), ln, serveFailureFatal)
		t.Assert(strings.Contains(buf.String(), "grpc server serve error"), false)
		t.Assert(strings.Contains(buf.String(), "FATA"), false)
	})
}

func Test_StartManaged_ServeErrorDoesNotExit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var buf safeBuffer
		s := newServeErrorTestServer(&buf)
		err := s.StartManaged()
		t.AssertNil(err)
		defer s.StopForceful()

		err = closeListenerWhenReady(s)
		t.AssertNil(err)

		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) {
			logged := buf.String()
			if strings.Contains(logged, "grpc server serve error") {
				t.Assert(strings.Contains(logged, "FATA"), false)
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
		t.Fatalf("expected serve error log, got %q", buf.String())
	})
}

func Test_Grpcx_Run_ServeErrorIsFatal(t *testing.T) {
	if os.Getenv(testRunServeFatalEnv) == "1" {
		runServeErrorFatalChild()
		return
	}

	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()

		cmd := exec.CommandContext(
			ctx,
			os.Args[0],
			"-test.run=^Test_Grpcx_Run_ServeErrorIsFatal$",
			"-test.count=1",
			"-test.v=false",
		)
		cmd.Env = append(os.Environ(), testRunServeFatalEnv+"=1")
		output, err := cmd.CombinedOutput()
		if ctx.Err() != nil {
			t.Fatalf("subprocess timed out, output=%s", output)
		}
		t.AssertNE(err, nil)

		var exitErr *exec.ExitError
		t.Assert(errors.As(err, &exitErr), true)
		t.Assert(exitErr.ExitCode(), 1)
		t.Assert(strings.Contains(string(output), "FATA"), true)
	})
}

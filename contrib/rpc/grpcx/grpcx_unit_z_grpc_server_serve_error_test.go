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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// testRunServeFatalEnv enables the subprocess path that expects Run() to Fatalf.
const testRunServeFatalEnv = "GF_GRPCX_TEST_RUN_SERVE_FATAL"

// testRegisterFatalEnv enables the subprocess path that expects register failure to Fatalf.
const testRegisterFatalEnv = "GF_GRPCX_TEST_REGISTER_FATAL"

// testRegisterHealthEnv isolates the legacy path that terminates on registration failure.
const testRegisterHealthEnv = "GF_GRPCX_TEST_REGISTER_HEALTH"

// healthCheckRegistrar verifies that a legacy server serves RPCs during registration.
type healthCheckRegistrar struct{}

// Register performs a real health RPC before accepting the service.
func (healthCheckRegistrar) Register(ctx context.Context, service gsvc.Service) (gsvc.Service, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	client, err := grpc.NewClient(service.GetEndpoints()[0].String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	response, rpcErr := grpc_health_v1.NewHealthClient(client).Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err := client.Close(); err != nil {
		return nil, err
	}
	if rpcErr != nil {
		return nil, rpcErr
	}
	if response.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
		return nil, errors.New("registered service is not serving")
	}
	return service, nil
}

// Deregister releases the test service without additional registry state.
func (healthCheckRegistrar) Deregister(context.Context, gsvc.Service) error {
	return nil
}

// Test_Grpcx_Legacy_ServesDuringRegister keeps registrars with synchronous health checks compatible.
func Test_Grpcx_Legacy_ServesDuringRegister(t *testing.T) {
	if os.Getenv(testRegisterHealthEnv) == "1" {
		gtest.C(t, func(t *gtest.T) {
			s := newServeErrorTestServer(os.Stdout)
			s.registrar = healthCheckRegistrar{}
			grpc_health_v1.RegisterHealthServer(s.Server, health.NewServer())
			defer s.StopForceful()
			// Run and Start share this fatal serving path, without entering signal wait here.
			t.AssertNil(s.serve(serveFailureFatal))
		})
		return
	}
	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^Test_Grpcx_Legacy_ServesDuringRegister$", "-test.count=1")
		cmd.Env = append(os.Environ(), testRegisterHealthEnv+"=1")
		output, err := cmd.CombinedOutput()
		if ctx.Err() != nil {
			t.Fatalf("subprocess timed out, output=%s", output)
		}
		if err != nil {
			t.Fatalf("legacy registration failed: %v, output=%s", err, output)
		}
	})
}

// registerFailRegistrar fails every Register call.
type registerFailRegistrar struct{}

// Register returns a fixed registration error.
func (registerFailRegistrar) Register(context.Context, gsvc.Service) (gsvc.Service, error) {
	return nil, errors.New("register failed")
}

// Deregister implements gsvc.Registrar.
func (registerFailRegistrar) Deregister(context.Context, gsvc.Service) error {
	return nil
}

// runRegisterFatalChild registers with a failing registrar on the fatal path.
func runRegisterFatalChild() {
	logger := glog.New()
	logger.SetStdoutPrint(true)
	c := Server.NewConfig()
	c.Name = guid.S()
	c.Address = "127.0.0.1:0"
	c.Logger = logger
	c.LogStdout = false
	s := Server.New(c)
	s.registrar = registerFailRegistrar{}
	if err := s.doServiceRegister(serveFailureFatal); err != nil {
		os.Exit(2)
	}
}

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
			s.Server,
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

		s.doServeAsynchronously(gctx.GetInitCtx(), s.Server, ln, serveFailureFatal)
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

func Test_StartManaged_RegisterErrorIsReturned(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var buf safeBuffer
		s := newServeErrorTestServer(&buf)
		s.registrar = registerFailRegistrar{}
		err := s.StartManaged()
		t.AssertNE(err, nil)
		t.Assert(strings.Contains(err.Error(), "register failed"), true)
		t.Assert(strings.Contains(buf.String(), "FATA"), false)
	})
}

func Test_doServiceRegister_FatalOnBlockingPath(t *testing.T) {
	if os.Getenv(testRegisterFatalEnv) == "1" {
		runRegisterFatalChild()
		return
	}

	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()

		cmd := exec.CommandContext(
			ctx,
			os.Args[0],
			"-test.run=^Test_doServiceRegister_FatalOnBlockingPath$",
			"-test.count=1",
			"-test.v=false",
		)
		cmd.Env = append(os.Environ(), testRegisterFatalEnv+"=1")
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

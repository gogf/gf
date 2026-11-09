// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

//go:build !windows

package gproc_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/test/gtest"
)

// These tests raise real OS signals in a child process. They cannot be expressed in the
// same process, because what they assert is precisely that the process gets terminated.
//
// Note that Test_Signal in gproc_z_signal_test.go feeds signalChan directly and therefore
// never exercises signal.Notify, so it cannot observe the behavior verified here.

const (
	signalHelperEnv = "GF_TEST_GPROC_SIGNAL_HELPER"

	// helperModeBlock installs a shutdown handler that never returns.
	helperModeBlock = "block"

	// helperModeReAdd installs a shutdown handler that never returns and, before blocking,
	// registers another handler. Registering re-arms signal.Notify, which must not happen
	// once the listening loop is over: signalChan would have no reader any more.
	helperModeReAdd = "readd"
)

// runSignalHelper blocks forever inside a shutdown handler, so that the process can only
// be terminated if the default signal behavior has been restored.
func runSignalHelper(mode string) {
	gproc.AddSigHandlerShutdown(func(sig os.Signal) {
		if mode == helperModeReAdd {
			gproc.AddSigHandlerShutdown(func(os.Signal) {})
		}
		// Announced only once the handler is fully set up, so that the second signal is
		// never raced against the work done above.
		fmt.Println("HANDLING")
		select {}
	})
	fmt.Println("READY")
	gproc.Listen()
}

func Test_Signal_SecondShutdownSignalTerminatesProcess(t *testing.T) {
	assertSecondSignalTerminates(t, helperModeBlock)
}

func Test_Signal_HandlersAddedAfterListenEndedDoNotReArmNotify(t *testing.T) {
	assertSecondSignalTerminates(t, helperModeReAdd)
}

func assertSecondSignalTerminates(t *testing.T, mode string) {
	if os.Getenv(signalHelperEnv) != "" {
		runSignalHelper(os.Getenv(signalHelperEnv))
		return
	}

	gtest.C(t, func(t *gtest.T) {
		output := newOutputSink()
		cmd := exec.Command(os.Args[0], "-test.run=^"+t.Name()+"$")
		cmd.Env = append(os.Environ(), signalHelperEnv+"="+mode)
		cmd.Stdout = output
		cmd.Stderr = output
		t.AssertNil(cmd.Start())
		defer func() {
			if err := killHelperProcess(cmd); err != nil {
				t.Logf("kill helper process failed: %v", err)
			}
		}()

		output.waitLine(t, "READY", 30*time.Second)
		t.AssertNil(cmd.Process.Signal(syscall.SIGINT))
		output.waitLine(t, "HANDLING", 30*time.Second)

		// The shutdown handler is blocked and never returns. The process may only be
		// terminated by the second signal if the default behavior has been restored.
		t.AssertNil(cmd.Process.Signal(syscall.SIGINT))

		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()
		select {
		case err := <-done:
			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				t.Error(fmt.Sprintf("helper process should be terminated by signal, got: %v", err))
			}
			status, ok := exitErr.Sys().(syscall.WaitStatus)
			if !ok {
				t.Error(fmt.Sprintf("unexpected wait status type: %T", exitErr.Sys()))
			}
			t.Assert(status.Signaled(), true)
			t.Assert(status.Signal(), syscall.SIGINT)
		case <-time.After(30 * time.Second):
			killErr := killHelperProcess(cmd)
			<-done
			t.Error(fmt.Sprintf(
				"helper process ignored the second signal (kill error: %v), output:\n%s",
				killErr, output.String(),
			))
		}
	})
}

// killHelperProcess kills the helper process. A process that has already exited is not
// an error: in the passing case it has been terminated by the second signal.
func killHelperProcess(cmd *exec.Cmd) error {
	if err := cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return err
	}
	return nil
}

// outputSink collects the output of the helper process, both line by line for waiting and
// as a whole for failure reporting.
type outputSink struct {
	mu      sync.Mutex
	all     bytes.Buffer
	pending []byte
	lines   chan string
}

func newOutputSink() *outputSink {
	return &outputSink{lines: make(chan string, 128)}
}

func (s *outputSink) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.all.Write(p)
	s.pending = append(s.pending, p...)
	for {
		i := bytes.IndexByte(s.pending, '\n')
		if i < 0 {
			break
		}
		line := string(s.pending[:i])
		s.pending = s.pending[i+1:]
		select {
		case s.lines <- line:
		default:
		}
	}
	return len(p), nil
}

func (s *outputSink) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.all.String()
}

func (s *outputSink) waitLine(t *gtest.T, expect string, timeout time.Duration) {
	deadline := time.After(timeout)
	for {
		select {
		case line := <-s.lines:
			if strings.Contains(line, expect) {
				return
			}
		case <-deadline:
			t.Error(fmt.Sprintf("timeout waiting for %q, output:\n%s", expect, s.String()))
		}
	}
}

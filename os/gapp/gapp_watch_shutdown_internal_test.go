// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests verify waitForRunExit without blocking on gproc.Listen waitChan.
package gapp

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

// recorderServer captures Stop for waitForRunExit tests.
type recorderServer struct {
	stopped int32
}

// Start succeeds immediately with no listeners.
func (r *recorderServer) Start() error {
	return nil
}

// Stop records that shutdown ran.
func (r *recorderServer) Stop(_ bool) error {
	atomic.StoreInt32(&r.stopped, 1)
	return nil
}

func TestWaitForRunExitTriggersGracefulStopOnContextCancel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rec := &recorderServer{}
		appInstance := New(rec)

		root, cancel := context.WithCancel(context.Background())
		defer cancel()

		t.AssertNil(appInstance.Boot(root))
		t.AssertNil(appInstance.Start(root))

		var (
			stopOnce sync.Once
			exitCh   = make(chan struct{})
		)
		doShutdown := func() {
			stopOnce.Do(func() {
				t.AssertNil(appInstance.Stop(gctx.NeverDone(root), true))
				close(exitCh)
			})
		}

		done := make(chan struct{})
		go func() {
			appInstance.waitForRunExit(root, doShutdown, exitCh)
			close(done)
		}()

		cancel()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("waitForRunExit did not return after context cancellation")
		}
		t.Assert(atomic.LoadInt32(&rec.stopped), int32(1))
	})
}

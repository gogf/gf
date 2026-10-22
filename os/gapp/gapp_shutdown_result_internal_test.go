// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests cover canceled startup, shared shutdown results, and completed Run waiters.

package gapp

import (
	"context"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
)

// shutdownResultServer controls lifecycle calls without opening listeners.
type shutdownResultServer struct {
	start func() error
	stop  func(bool) error
}

// Start runs the configured startup behavior.
func (s shutdownResultServer) Start() error { return s.start() }

// Stop runs the configured shutdown behavior.
func (s shutdownResultServer) Stop(graceful bool) error { return s.stop(graceful) }

// lifecycleGoroutineStacks returns stacks for checking lifecycle waiters.
func lifecycleGoroutineStacks() string {
	buffer := make([]byte, 1<<20)
	return string(buffer[:runtime.Stack(buffer, true)])
}

// TestStartRejectsAlreadyCanceledContext avoids starting servers for a canceled request.
func TestStartRejectsAlreadyCanceledContext(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx, cancel = context.WithCancel(context.Background())
			starts      atomic.Int32
			stops       atomic.Int32
			app         = New(shutdownResultServer{
				start: func() error { starts.Add(1); return nil },
				stop:  func(bool) error { stops.Add(1); return nil },
			})
		)
		cancel()
		t.AssertNil(app.Boot(context.Background()))
		t.Assert(app.Start(ctx), context.Canceled)
		t.Assert(starts.Load(), int32(0))
		t.Assert(stops.Load(), int32(0))
	})
}

// TestStartRejectsCancellationDuringBoot avoids launching servers after Boot cancels its context.
func TestStartRejectsCancellationDuringBoot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx, cancel = context.WithCancel(context.Background())
			starts      atomic.Int32
			app         = New(shutdownResultServer{
				start: func() error { starts.Add(1); return nil },
				stop:  func(bool) error { return nil },
			})
		)
		defer cancel()
		app.Option(NewOption(func(context.Context, *App) { cancel() }))
		t.Assert(app.Start(ctx), context.Canceled)
		t.Assert(starts.Load(), int32(0))
	})
}

// TestStartRollsBackCancellationOnSuccessfulStartup covers both ready select cases.
func TestStartRollsBackCancellationOnSuccessfulStartup(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		previous := runtime.GOMAXPROCS(1)
		defer runtime.GOMAXPROCS(previous)
		for i := 0; i < 1000; i++ {
			var (
				ctx, cancel = context.WithCancel(context.Background())
				stops       atomic.Int32
				app         = New(shutdownResultServer{
					start: func() error { cancel(); return nil },
					stop: func(graceful bool) error {
						if !graceful {
							stops.Add(1)
						}
						return nil
					},
				})
			)
			app.logger.SetStdoutPrint(false)
			err := app.Start(ctx)
			cancel()
			t.Assert(err, context.Canceled)
			t.Assert(stops.Load(), int32(1))
		}
	})
}

// TestStopSharesFirstResult returns the first shutdown error to concurrent and repeated callers.
func TestStopSharesFirstResult(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			stopErr = gerror.New("shutdown failed")
			stops   atomic.Int32
			entered = make(chan struct{})
			release = make(chan struct{})
			results = make(chan error, 10)
			app     = New(shutdownResultServer{
				start: func() error { return nil },
				stop: func(bool) error {
					stops.Add(1)
					close(entered)
					<-release
					return stopErr
				},
			})
		)
		t.AssertNil(app.Start(context.Background()))
		go func() { results <- app.Stop(context.Background(), true) }()
		<-entered
		for i := 0; i < 9; i++ {
			go func() { results <- app.Stop(context.Background(), false) }()
		}
		close(release)
		first := <-results
		t.Assert(gerror.Is(first, stopErr), true)
		for i := 0; i < 9; i++ {
			t.Assert(<-results, first)
		}
		t.Assert(app.Stop(context.Background(), true), first)
		t.Assert(stops.Load(), int32(1))
	})
}

// TestRunPreservesExternalStopError reports failed shutdown even if Stop already ran.
func TestRunPreservesExternalStopError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx, cancel = context.WithCancel(context.Background())
			stopErr     = gerror.New("shutdown failed")
			returned    = make(chan error, 1)
			app         = New(shutdownResultServer{
				start: func() error { return nil },
				stop:  func(bool) error { return stopErr },
			})
		)
		defer cancel()
		go func() { returned <- app.Run(ctx) }()
		deadline := time.Now().Add(time.Second)
		for !strings.Contains(lifecycleGoroutineStacks(), "(*App).waitForRunExit(") {
			if time.Now().After(deadline) {
				t.Fatal("Run did not finish startup")
			}
			runtime.Gosched()
		}
		first := app.Stop(context.Background(), true)
		t.Assert(gerror.Is(first, stopErr), true)
		cancel()
		select {
		case err := <-returned:
			t.Assert(err, first)
		case <-time.After(time.Second):
			t.Fatal("Run did not finish after context cancellation")
		}
	})
}

// TestRunSignalCompletionLeavesNoContextWatcher covers shutdown without canceling the root.
func TestRunSignalCompletionLeavesNoContextWatcher(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		before := strings.Count(lifecycleGoroutineStacks(), "(*App).waitForRunExit.func1")
		for i := 0; i < 3; i++ {
			var (
				done     = make(chan struct{})
				returned = make(chan struct{})
			)
			go func() {
				New().waitForRunExit(ctx, func() {}, done)
				close(returned)
			}()
			close(done)
			<-returned
		}
		runtime.Gosched()
		t.Assert(strings.Count(lifecycleGoroutineStacks(), "(*App).waitForRunExit.func1"), before)
	})
}

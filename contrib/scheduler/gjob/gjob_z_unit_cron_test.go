// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests cron entry ownership during shutdown and startup rollback.

package gjob

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestWorkerServerNilHandler checks all tasks before starting any worker.
func TestWorkerServerNilHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		called := make(chan struct{}, 1)
		s := NewWorkerServer(nil, WorkerTask{
			Name: "valid", Handler: func(context.Context) func() { called <- struct{}{}; return nil },
		})
		t.AssertNil(s.Add(WorkerTask{Name: "nil-handler"}))
		err := s.Start()
		defer func() { t.AssertNil(s.Stop(false)) }()
		t.Assert(gerror.Code(err), gcode.CodeInvalidConfiguration)
		t.Assert(s.started, false)
		t.AssertNil(s.ctx)
		select {
		case <-called:
			t.Fatal("valid worker executed before all handlers were validated")
		default:
		}
		s.tasks = s.tasks[:1]
		t.AssertNil(s.Start())
		select {
		case <-called:
		case <-time.After(time.Second):
			t.Fatal("worker did not start after fixing the handler")
		}
	})
}

// TestCronServerNilHandler rejects missing callbacks without registering earlier tasks.
func TestCronServerNilHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := NewCronServer(nil, CronTask{Name: "valid", Spec: "* * * * * *", Handler: func(context.Context) error { return nil }})
		t.AssertNil(s.Add(CronTask{Name: "nil-handler", Spec: "* * * * * *"}))
		err := s.Start()
		defer func() { t.AssertNil(s.Stop(false)) }()
		t.Assert(gerror.Code(err), gcode.CodeInvalidConfiguration)
		t.Assert(s.started, false)
		t.AssertNil(s.ctx)
		t.AssertNil(s.cron)
		s.tasks = s.tasks[:1]
		t.AssertNil(s.Start())
		t.Assert(s.cron.Size(), 1)
	})
}

// TestCronServerStopClosesSparseEntries releases entries without waiting for their schedule.
func TestCronServerStopClosesSparseEntries(t *testing.T) {
	for _, graceful := range []bool{true, false} {
		gtest.C(t, func(t *gtest.T) {
			s := NewCronServer(nil, CronTask{
				Name: "yearly", Spec: "0 0 0 1 1 *",
				Handler: func(context.Context) error { return nil },
			})
			t.AssertNil(s.Start())
			entry := s.cron.Search("yearly")
			t.AssertNE(entry, nil)
			t.AssertNil(s.Stop(graceful))
			t.Assert(s.cron.Size(), 0)
			t.Assert(entry.Status(), gcron.StatusClosed)
		})
	}
}

// TestCronServerFailedStartClosesEntries releases entries from each failed attempt.
func TestCronServerFailedStartClosesEntries(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var called atomic.Int32
		s := NewCronServer(nil,
			CronTask{Name: "ready", Spec: "* * * * * *", Handler: func(context.Context) error { called.Add(1); return nil }},
			CronTask{Name: "invalid", Spec: "invalid-spec", Handler: func(context.Context) error { return nil }},
		)
		t.AssertNE(s.Start(), nil)
		t.Assert(called.Load(), 0)
		t.AssertNE(s.cron, nil)
		t.Assert(s.cron.Size(), 0)
		attempt := s.cron
		s.tasks = s.tasks[:1]
		t.AssertNil(s.Start())
		t.Assert(s.cron == attempt, false)
		t.Assert(s.cron.Size(), 1)
		t.Assert(attempt.Size(), 0)
		t.AssertNil(s.Stop(true))
	})
}

// TestCronServerConcurrentStop preserves the first caller's graceful shutdown.
func TestCronServerConcurrentStop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			running = make(chan context.Context, 1)
			release = make(chan struct{})
			first   = make(chan error, 1)
			second  = make(chan error, 1)
		)
		defer close(release)
		s := NewCronServer(nil, CronTask{
			Name: "blocking-cron",
			Spec: "* * * * * *",
			Handler: func(ctx context.Context) error {
				running <- ctx
				<-release
				return nil
			},
		})
		t.AssertNil(s.Start())
		var jobCtx context.Context
		select {
		case jobCtx = <-running:
		case <-time.After(2 * time.Second):
			t.Fatal("cron job did not start")
		}
		go func() { first <- s.Stop(true) }()
		// Wait until the first caller owns shutdown before starting the second.
		deadline := time.Now().Add(time.Second)
		for {
			s.mu.Lock()
			stopping := !s.started
			s.mu.Unlock()
			if stopping {
				break
			}
			if time.Now().After(deadline) {
				t.Fatal("cron Stop did not begin")
			}
			time.Sleep(time.Millisecond)
		}
		go func() { second <- s.Stop(false) }()
		select {
		case err := <-second:
			t.Fatalf("concurrent Stop returned before the in-flight job: %v", err)
		case <-time.After(50 * time.Millisecond):
		}
		t.AssertNil(jobCtx.Err())
		release <- struct{}{}
		t.AssertNil(<-first)
		t.AssertNil(<-second)
		t.Assert(jobCtx.Err(), context.Canceled)
		t.AssertNil(s.Stop(true))
	})
}

// TestCronServerStopOwnsGeneratedNameEntries closes timers displaced by automatic names.
func TestCronServerStopOwnsGeneratedNameEntries(t *testing.T) {
	for _, graceful := range []bool{true, false} {
		gtest.C(t, func(t *gtest.T) {
			s := NewCronServer(nil,
				CronTask{Name: "cron-1", Spec: "0 0 0 1 1 *", Handler: func(context.Context) error { return nil }},
				CronTask{Spec: "0 0 0 1 1 *", Handler: func(context.Context) error { return nil }},
			)
			t.AssertNil(s.Start())
			defer func() { t.AssertNil(s.Stop(false)) }()
			t.Assert(len(s.entries), 2)
			entries := s.entries
			t.AssertNil(s.Stop(graceful))
			for _, entry := range entries {
				t.Assert(entry.Status(), gcron.StatusClosed)
			}
		})
	}
}

// TestCronServerFailedStartOwnsGeneratedNameEntries closes every timer before retry.
func TestCronServerFailedStartOwnsGeneratedNameEntries(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := NewCronServer(nil,
			CronTask{Name: "cron-1", Spec: "0 0 0 1 1 *", Handler: func(context.Context) error { return nil }},
			CronTask{Spec: "0 0 0 1 1 *", Handler: func(context.Context) error { return nil }},
			CronTask{Name: "invalid", Spec: "invalid-spec", Handler: func(context.Context) error { return nil }},
		)
		t.AssertNE(s.Start(), nil)
		t.Assert(len(s.entries), 2)
		failedEntries := s.entries
		for _, entry := range failedEntries {
			t.Assert(entry.Status(), gcron.StatusClosed)
		}
		s.tasks = s.tasks[:2]
		t.AssertNil(s.Start())
		defer func() { t.AssertNil(s.Stop(false)) }()
		t.Assert(len(s.entries), 2)
		t.Assert(s.entries[0] == failedEntries[0], false)
		t.AssertNil(s.Stop(true))
		for _, entry := range s.entries {
			t.Assert(entry.Status(), gcron.StatusClosed)
		}
	})
}

// TestJobServerTaskSliceOwnership keeps independently added tasks isolated between servers.
func TestJobServerTaskSliceOwnership(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tasks := make([]WorkerTask, 1, 2)
		tasks[0] = WorkerTask{Name: "common", Handler: func(context.Context) func() { return nil }}
		var (
			first  = NewWorkerServer(nil, tasks...)
			second = NewWorkerServer(nil, tasks...)
		)
		t.AssertNil(first.Add(WorkerTask{Name: "first", Handler: func(context.Context) func() { return nil }}))
		t.AssertNil(second.Add(WorkerTask{Name: "second", Handler: func(context.Context) func() { return nil }}))
		t.Assert(first.tasks[1].Name, "first")
		t.Assert(second.tasks[1].Name, "second")
	})
	gtest.C(t, func(t *gtest.T) {
		tasks := make([]CronTask, 1, 2)
		tasks[0] = CronTask{Name: "common", Spec: "* * * * * *", Handler: func(context.Context) error { return nil }}
		var (
			first  = NewCronServer(nil, tasks...)
			second = NewCronServer(nil, tasks...)
		)
		t.AssertNil(first.Add(CronTask{Name: "first", Spec: "* * * * * *", Handler: func(context.Context) error { return nil }}))
		t.AssertNil(second.Add(CronTask{Name: "second", Spec: "* * * * * *", Handler: func(context.Context) error { return nil }}))
		t.Assert(first.tasks[1].Name, "first")
		t.Assert(second.tasks[1].Name, "second")
	})
}

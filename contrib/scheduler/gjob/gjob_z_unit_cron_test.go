// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests cron entry ownership during shutdown and startup rollback.

package gjob

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/test/gtest"
)

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
		s := NewCronServer(nil,
			CronTask{Name: "yearly", Spec: "0 0 0 1 1 *", Handler: func(context.Context) error { return nil }},
			CronTask{Name: "invalid", Spec: "invalid-spec", Handler: func(context.Context) error { return nil }},
		)
		t.AssertNE(s.Start(), nil)
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

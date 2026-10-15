// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests cover concurrent Boot results across immediate retries.

package gapp

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

// TestBootWaitersKeepAttemptResult makes the retry finish before failed-attempt waiters resume.
func TestBootWaitersKeepAttemptResult(t *testing.T) {
	previous := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(previous)
	gtest.C(t, func(t *gtest.T) {
		var (
			app      = New()
			entered  = make(chan struct{})
			release  = make(chan struct{})
			joined   = make(chan struct{}, 3)
			first    = make(chan error, 1)
			retry    = make(chan error, 1)
			waiters  = make(chan error, 3)
			failure  = errors.New("first boot failed")
			attempts int
		)
		app.Option(NewOptionWithHook(func(context.Context, *App) (func(context.Context), error) {
			attempts++
			if attempts == 1 {
				close(entered)
				<-release
				return nil, failure
			}
			return nil, nil
		}))
		bootWaitHook = func() { joined <- struct{}{} }
		defer func() { bootWaitHook = nil }()
		go func() {
			first <- app.Boot(context.Background())
			retry <- app.Boot(context.Background())
		}()
		<-entered
		for i := 0; i < cap(waiters); i++ {
			go func() { waiters <- app.Boot(context.Background()) }()
			<-joined
		}
		close(release)
		firstErr := <-first
		t.Assert(errors.Is(firstErr, failure), true)
		t.AssertNil(<-retry)
		for i := 0; i < cap(waiters); i++ {
			select {
			case err := <-waiters:
				t.Assert(err, firstErr)
			case <-time.After(time.Second):
				t.Fatal("Boot waiter did not finish")
			}
		}
		t.Assert(attempts, 2)
		t.Assert(app.Booted(), true)
	})
}

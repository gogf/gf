// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests cover rollback of registrations made during failed boot attempts.

package gapp

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

// bootRetryServer records how many times each registration starts its server.
type bootRetryServer struct {
	starts atomic.Int32
}

// TestBootFailureRollsBackConcurrentRegistrations keeps pre-attempt registrations.
func TestBootFailureRollsBackConcurrentRegistrations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			original = &bootRetryServer{}
			app      = New(original)
			entered  = make(chan struct{})
			release  = make(chan struct{})
			result   = make(chan error, 1)
			fail     = true
			applies  int
		)
		app.Option(NewOptionWithHook(func(context.Context, *App) (func(context.Context), error) {
			if fail {
				close(entered)
				<-release
				return nil, context.Canceled
			}
			return nil, nil
		}))
		go func() {
			result <- app.Boot(context.Background())
		}()
		select {
		case <-entered:
		case <-time.After(time.Second):
			close(release)
			t.Fatal("boot option did not run")
		}
		app.Add(&bootRetryServer{})
		app.Option(NewOption(func(context.Context, *App) { applies++ }))
		close(release)
		select {
		case err := <-result:
			t.AssertNE(err, nil)
		case <-time.After(time.Second):
			t.Fatal("failed boot did not return")
		}
		t.Assert(app.Servers(), []Server{original})
		t.Assert(len(app.options), 1)
		fail = false
		t.AssertNil(app.Boot(context.Background()))
		t.Assert(applies, 0)
		t.AssertNil(app.Start(context.Background()))
		t.Assert(original.starts.Load(), int32(1))
		t.AssertNil(app.Stop(context.Background(), true))
	})
}

// Start records the startup without opening a listener.
func (s *bootRetryServer) Start() error {
	s.starts.Add(1)
	return nil
}

// Stop satisfies the lifecycle contract without acquiring resources.
func (*bootRetryServer) Stop(_ bool) error {
	return nil
}

// TestBootRetryRollsBackRegistrations covers both error and panic rollback.
func TestBootRetryRollsBackRegistrations(t *testing.T) {
	for _, panicOnFailure := range []bool{false, true} {
		name := "error"
		if panicOnFailure {
			name = "panic"
		}
		t.Run(name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					original       = &bootRetryServer{}
					added          = &bootRetryServer{}
					app            = New(original)
					fail           = true
					dynamicApplies int
					cleanupOrder   []int
				)
				app.Option(NewOptionWithHook(func(_ context.Context, a *App) (func(context.Context), error) {
					a.Add(added)
					a.Option(NewOptionWithHook(func(context.Context, *App) (func(context.Context), error) {
						dynamicApplies++
						return func(context.Context) { cleanupOrder = append(cleanupOrder, 2) }, nil
					}))
					a.Option(NewOptionWithHook(func(context.Context, *App) (func(context.Context), error) {
						if fail {
							if panicOnFailure {
								panic(context.Canceled)
							}
							return nil, context.Canceled
						}
						return nil, nil
					}))
					return func(context.Context) { cleanupOrder = append(cleanupOrder, 1) }, nil
				}))

				t.AssertNE(app.Boot(context.Background()), nil)
				t.Assert(app.Booted(), false)
				t.Assert(app.Servers(), []Server{original})
				t.Assert(len(app.options), 1)
				t.Assert(cleanupOrder, []int{2, 1})

				fail = false
				t.AssertNil(app.Boot(context.Background()))
				t.Assert(app.Servers(), []Server{original, added})
				t.Assert(len(app.options), 3)
				t.Assert(dynamicApplies, 2)
				t.AssertNil(app.Start(context.Background()))
				t.Assert(original.starts.Load(), int32(1))
				t.Assert(added.starts.Load(), int32(1))
				t.AssertNil(app.Stop(context.Background(), true))
				t.Assert(cleanupOrder, []int{2, 1, 2, 1})
			})
		})
	}
}

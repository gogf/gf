// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests cover cancellation-independent automatic cleanup and explicit Stop contexts.

package gapp

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

// cleanupContextKey identifies propagation metadata in cleanup tests.
type cleanupContextKey struct{}

// cancelStartupServer cancels the lifecycle context while startup fails.
type cancelStartupServer struct {
	cancel context.CancelFunc
}

// Start simulates startup cancellation.
func (s cancelStartupServer) Start() error {
	s.cancel()
	return context.Canceled
}

// Stop satisfies the server lifecycle without owning resources.
func (cancelStartupServer) Stop(bool) error {
	return nil
}

// TestBootRollbackCleanupKeepsContextValues covers both error and panic rollback.
func TestBootRollbackCleanupKeepsContextValues(t *testing.T) {
	for _, panicOnFailure := range []bool{false, true} {
		name := "error"
		if panicOnFailure {
			name = "panic"
		}
		t.Run(name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					ctx, cancel = context.WithTimeout(context.WithValue(context.Background(), cleanupContextKey{}, "trace"), time.Hour)
					app         = New()
					cleanupCtx  context.Context
				)
				defer cancel()
				app.Option(NewOptionWithHook(func(context.Context, *App) (func(context.Context), error) {
					return func(ctx context.Context) { cleanupCtx = ctx }, nil
				}))
				app.Option(NewOptionWithHook(func(context.Context, *App) (func(context.Context), error) {
					cancel()
					if panicOnFailure {
						panic(context.Canceled)
					}
					return nil, context.Canceled
				}))
				t.AssertNE(app.Boot(ctx), nil)
				t.AssertNE(cleanupCtx, nil)
				t.AssertNil(cleanupCtx.Err())
				t.AssertNil(cleanupCtx.Done())
				_, hasDeadline := cleanupCtx.Deadline()
				t.Assert(hasDeadline, false)
				t.Assert(cleanupCtx.Value(cleanupContextKey{}), "trace")
				t.Assert(ctx.Err(), context.Canceled)
			})
		})
	}
}

// TestRunStartupFailureCleanupKeepsContextValues preserves metadata after canceled startup.
func TestRunStartupFailureCleanupKeepsContextValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx, cancel = context.WithCancel(context.WithValue(context.Background(), cleanupContextKey{}, "trace"))
			app         = New(cancelStartupServer{cancel: cancel})
			cleanupCtx  context.Context
		)
		defer cancel()
		app.Option(NewOptionWithHook(func(context.Context, *App) (func(context.Context), error) {
			return func(ctx context.Context) { cleanupCtx = ctx }, nil
		}))
		t.AssertNE(app.Run(ctx), nil)
		t.AssertNE(cleanupCtx, nil)
		t.AssertNil(cleanupCtx.Err())
		t.AssertNil(cleanupCtx.Done())
		t.Assert(cleanupCtx.Value(cleanupContextKey{}), "trace")
		t.Assert(ctx.Err(), context.Canceled)
	})
}

// TestStopPreservesExplicitCanceledContext keeps the caller's explicit cleanup contract.
func TestStopPreservesExplicitCanceledContext(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx, cancel = context.WithCancel(context.WithValue(context.Background(), cleanupContextKey{}, "trace"))
			app         = New()
			cleanupCtx  context.Context
		)
		defer cancel()
		app.Option(NewOptionWithHook(func(context.Context, *App) (func(context.Context), error) {
			return func(ctx context.Context) { cleanupCtx = ctx }, nil
		}))
		t.AssertNil(app.Boot(ctx))
		cancel()
		t.AssertNil(app.Stop(ctx, false))
		t.Assert(cleanupCtx, ctx)
		t.Assert(cleanupCtx.Err(), context.Canceled)
		t.Assert(cleanupCtx.Value(cleanupContextKey{}), "trace")
	})
}

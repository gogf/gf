// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests cover signal callback release and active application dispatch.

package gapp

import (
	"context"
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

// activeRunSignals reports callback ownership under the dispatch lock.
func activeRunSignals() int {
	runSignalMu.Lock()
	defer runSignalMu.Unlock()
	return len(runSignals)
}

// TestRunReleasesSignalCallbacks covers consecutive context-driven lifecycles and live dispatch.
func TestRunReleasesSignalCallbacks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		before := activeRunSignals()
		for i := 0; i < 4; i++ {
			var (
				ctx, cancel = context.WithCancel(context.Background())
				returned    = make(chan error, 1)
			)
			go func() { returned <- New().Run(ctx) }()
			deadline := time.Now().Add(time.Second)
			for activeRunSignals() != before+1 {
				if time.Now().After(deadline) {
					cancel()
					t.Fatal("Run did not register its signal callback")
				}
				runtime.Gosched()
			}
			if i == 3 {
				dispatchRunSignal(syscall.SIGTERM)
			} else {
				cancel()
			}
			select {
			case err := <-returned:
				t.AssertNil(err)
			case <-time.After(time.Second):
				cancel()
				t.Fatal("Run did not finish shutdown")
			}
			cancel()
			t.Assert(activeRunSignals(), before)
		}
	})
}

// TestRunSignalDispatchIsolatesPanic allows other applications to stop after one callback panics.
func TestRunSignalDispatchIsolatesPanic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		before := activeRunSignals()
		releasePanic := registerRunSignal(func(os.Signal) { panic("test panic") })
		defer releasePanic()
		var release func()
		release = registerRunSignal(func(sig os.Signal) {
			t.Assert(sig, syscall.SIGTERM)
			release()
		})
		defer release()
		dispatchRunSignal(syscall.SIGTERM)
		t.Assert(activeRunSignals(), before+1)
		releasePanic()
		t.Assert(activeRunSignals(), before)
	})
}

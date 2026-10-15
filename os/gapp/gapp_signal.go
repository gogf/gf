// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Signal dispatch retains only callbacks belonging to active Run calls.

package gapp

import (
	"context"
	"os"
	"sync"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/util/gutil"
)

var (
	// runSignalOnce installs one process handler without capturing an application.
	runSignalOnce sync.Once
	// runSignalMu protects registration and release of live callbacks.
	runSignalMu sync.Mutex
	// runSignals uses callback addresses as unique registration identities.
	runSignals = make(map[*gproc.SigHandler]struct{})
)

// registerRunSignal registers a live Run callback and returns its release function.
func registerRunSignal(handler gproc.SigHandler) func() {
	runSignalOnce.Do(func() { gproc.AddSigHandlerShutdown(dispatchRunSignal) })
	runSignalMu.Lock()
	runSignals[&handler] = struct{}{}
	runSignalMu.Unlock()
	return func() {
		runSignalMu.Lock()
		delete(runSignals, &handler)
		runSignalMu.Unlock()
	}
}

// dispatchRunSignal snapshots callbacks so shutdown can release registrations without deadlock.
func dispatchRunSignal(sig os.Signal) {
	runSignalMu.Lock()
	handlers := make([]gproc.SigHandler, 0, len(runSignals))
	for handler := range runSignals {
		handlers = append(handlers, *handler)
	}
	runSignalMu.Unlock()
	for _, handler := range handlers {
		gutil.TryCatch(gctx.GetInitCtx(), func(context.Context) {
			handler(sig)
		}, func(ctx context.Context, err error) {
			glog.Errorf(ctx, "application signal handler failed: %+v", err)
		})
	}
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer

import (
	"context"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Entry is the timing job.
type Entry struct {
	job         JobFunc         // The job function.
	ctx         context.Context // The context for the job, for READ ONLY.
	timer       *Timer          // Belonged timer.
	ticks       int64           // The job runs every tick.
	times       *gtype.Int      // Limit running times.
	status      *gtype.Int      // Job status.
	isSingleton *gtype.Bool     // Singleton mode.
	nextTicks   *gtype.Int64    // Next run ticks of the job.
	infinite    *gtype.Bool     // No times limit.
}

// JobFunc is the timing called job function in timer.
type JobFunc = func(ctx context.Context)

// Status returns the status of the job.
func (entry *Entry) Status() int {
	return entry.status.Val()
}

// Run runs the timer job asynchronously.
func (entry *Entry) Run() {
	if !entry.infinite.Val() {
		leftRunningTimes := entry.times.Add(-1)
		// It checks its running times exceeding.
		if leftRunningTimes < 0 {
			entry.status.Set(StatusClosed)
			return
		}
	}
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				if exception != panicExit {
					if v, ok := exception.(error); ok && gerror.HasStack(v) {
						panic(v)
					} else {
						panic(gerror.Newf(`exception recovered: %+v`, exception))
					}
				} else {
					entry.Close()
					return
				}
			}
			if entry.Status() == StatusRunning {
				entry.SetStatus(StatusReady)
			}
		}()
		entry.job(entry.ctx)
	}()
}

// doCheckAndRunByTicks checks the if job can run in given timer ticks,
// it runs asynchronously if the given `currentTimerTicks` meets or else
// it increments its ticks and waits for next running check.
func (entry *Entry) doCheckAndRunByTicks(currentTimerTicks int64) {
	// Ticks check.
	if currentTimerTicks < entry.nextTicks.Val() {
		return
	}
	entry.nextTicks.Set(currentTimerTicks + entry.ticks)
	// Perform job checking.
	switch entry.status.Val() {
	case StatusRunning:
		if entry.IsSingleton() {
			return
		}
	case StatusReady:
		if !entry.status.Cas(StatusReady, StatusRunning) {
			return
		}
	case StatusStopped:
		return
	case StatusClosed:
		return
	}
	// Perform job running.
	entry.Run()
}

// SetStatus custom sets the status for the job.
func (entry *Entry) SetStatus(status int) int {
	return entry.status.Set(status)
}

// Start starts the job.
func (entry *Entry) Start() {
	entry.status.Set(StatusReady)
}

// Stop stops the job.
func (entry *Entry) Stop() {
	entry.status.Set(StatusStopped)
}

// Close closes the job, and then it will be removed from the timer.
func (entry *Entry) Close() {
	entry.status.Set(StatusClosed)
}

// Reset resets the job, which resets its ticks for next running.
func (entry *Entry) Reset() {
	entry.nextTicks.Set(entry.timer.ticks.Val() + entry.ticks)
}

// IsSingleton checks and returns whether the job in singleton mode.
func (entry *Entry) IsSingleton() bool {
	return entry.isSingleton.Val()
}

// SetSingleton sets the job singleton mode.
func (entry *Entry) SetSingleton(enabled bool) {
	entry.isSingleton.Set(enabled)
}

// Job returns the job function of this job.
func (entry *Entry) Job() JobFunc {
	return entry.job
}

// Ctx returns the initialized context of this job.
func (entry *Entry) Ctx() context.Context {
	return entry.ctx
}

// SetTimes sets the limit running times for the job.
func (entry *Entry) SetTimes(times int) {
	entry.times.Set(times)
	entry.infinite.Set(false)
}

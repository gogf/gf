// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/util/gconv"
)

// JobFunc is the timing called job function in cron.
type JobFunc = gtimer.JobFunc

// Entry is timing task entry.
type Entry struct {
	cron         *Cron         // Cron object belonged to.
	timerEntry   *gtimer.Entry // Associated timer Entry.
	schedule     *cronSchedule // Timed schedule object.
	jobName      string        // Callback function name(address info).
	times        *gtype.Int    // Running times limit.
	infinite     *gtype.Bool   // No times limit.
	Name         string        // Entry name.
	RegisterTime time.Time     // Registered time.
	Job          JobFunc       `json:"-"` // Callback function.
}

type doAddEntryInput struct {
	Name        string          // Name names this entry for manual control.
	Job         JobFunc         // Job is the callback function for timed task execution.
	Ctx         context.Context // The context for the job.
	Times       int             // Times specifies the running limit times for the entry.
	Pattern     string          // Pattern is the crontab style string for scheduler.
	IsSingleton bool            // Singleton specifies whether timed task executing in singleton mode.
	Infinite    bool            // Infinite specifies whether this entry is running with no times limit.
}

// doAddEntry creates and returns a new Entry object.
func (c *Cron) doAddEntry(in doAddEntryInput) (*Entry, error) {
	if in.Name != "" {
		if c.Search(in.Name) != nil {
			return nil, gerror.NewCodef(
				gcode.CodeInvalidOperation,
				`duplicated cron job name "%s", already exists`,
				in.Name,
			)
		}
	}
	schedule, err := newSchedule(in.Pattern)
	if err != nil {
		return nil, err
	}
	// No limit for `times`, for timer checking scheduling every second.
	entry := &Entry{
		cron:         c,
		schedule:     schedule,
		jobName:      runtime.FuncForPC(reflect.ValueOf(in.Job).Pointer()).Name(),
		times:        gtype.NewInt(in.Times),
		infinite:     gtype.NewBool(in.Infinite),
		RegisterTime: time.Now(),
		Job:          in.Job,
	}
	if in.Name != "" {
		entry.Name = in.Name
	} else {
		entry.Name = "cron-" + gconv.String(c.idGen.Add(1))
	}
	// When you add a scheduled task, you cannot allow it to run.
	// It cannot start running when added to timer.
	// It should start running after the entry is added to the Cron entries map, to avoid the task
	// from running during adding where the entries do not have the entry information, which might cause panic.
	entry.timerEntry = gtimer.AddEntry(
		in.Ctx,
		time.Second,
		entry.checkAndRun,
		in.IsSingleton,
		-1,
		gtimer.StatusStopped,
	)
	c.entries.Set(entry.Name, entry)
	entry.timerEntry.Start()
	return entry, nil
}

// IsSingleton return whether this entry is a singleton timed task.
func (e *Entry) IsSingleton() bool {
	return e.timerEntry.IsSingleton()
}

// SetSingleton sets the entry running in singleton mode.
func (e *Entry) SetSingleton(enabled bool) {
	e.timerEntry.SetSingleton(enabled)
}

// SetTimes sets the times which the entry can run.
func (e *Entry) SetTimes(times int) {
	e.times.Set(times)
	e.infinite.Set(false)
}

// Status returns the status of entry.
func (e *Entry) Status() int {
	return e.timerEntry.Status()
}

// SetStatus sets the status of the entry.
func (e *Entry) SetStatus(status int) int {
	return e.timerEntry.SetStatus(status)
}

// Start starts running the entry.
func (e *Entry) Start() {
	e.timerEntry.Start()
}

// Stop stops running the entry.
func (e *Entry) Stop() {
	e.timerEntry.Stop()
}

// Close stops and removes the entry from cron.
func (e *Entry) Close() {
	e.cron.entries.Remove(e.Name)
	e.timerEntry.Close()
}

// checkAndRun is the core timing task check logic.
// This function is called every second.
func (e *Entry) checkAndRun(ctx context.Context) {
	currentTime := time.Now()
	if !e.schedule.checkMeetAndUpdateLastSeconds(ctx, currentTime) {
		return
	}
	switch e.cron.status.Val() {
	case StatusStopped:
		return

	case StatusClosed:
		e.logDebugf(ctx, `cron job "%s" is removed`, e.getJobNameWithPattern())
		e.Close()

	case StatusReady, StatusRunning:
		defer func() {
			if exception := recover(); exception != nil {
				// Exception caught, it logs the error content to logger in default behavior.
				e.logErrorf(ctx,
					`cron job "%s(%s)" end with error: %+v`,
					e.jobName, e.schedule.pattern, exception,
				)
			} else {
				e.logDebugf(ctx, `cron job "%s" ends`, e.getJobNameWithPattern())
			}
			if e.timerEntry.Status() == StatusClosed {
				e.Close()
			}
		}()

		// Running times check.
		if !e.infinite.Val() {
			times := e.times.Add(-1)
			if times <= 0 {
				if e.timerEntry.SetStatus(StatusClosed) == StatusClosed || times < 0 {
					return
				}
			}
		}
		e.logDebugf(ctx, `cron job "%s" starts`, e.getJobNameWithPattern())
		e.Job(ctx)
	}
}

func (e *Entry) getJobNameWithPattern() string {
	return fmt.Sprintf(`%s(%s)`, e.jobName, e.schedule.pattern)
}

func (e *Entry) logDebugf(ctx context.Context, format string, v ...interface{}) {
	if logger := e.cron.GetLogger(); logger != nil {
		logger.Debugf(ctx, format, v...)
	}
}

func (e *Entry) logErrorf(ctx context.Context, format string, v ...interface{}) {
	logger := e.cron.GetLogger()
	if logger == nil {
		logger = glog.DefaultLogger()
	}
	logger.Errorf(ctx, format, v...)
}

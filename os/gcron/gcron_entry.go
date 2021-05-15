// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron

import (
	"reflect"
	"runtime"
	"time"

	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/util/gconv"
)

// Timed task entry.
type Job struct {
	cron     *Cron         // Cron object belonged to.
	job      *gtimer.Job   // Associated gtimer.Job.
	schedule *cronSchedule // Timed schedule object.
	jobName  string        // Callback function name(address info).
	times    *gtype.Int    // Running times limit.
	Name     string        // Job name.
	Job      func()        `json:"-"` // Callback function.
	Time     time.Time     // Registered time.
}

// addJob creates and returns a new Job object.
// Param <job> is the callback function for timed task execution.
// Param <singleton> specifies whether timed task executing in singleton mode.
// Param <name> names this entry for manual control.
func (c *Cron) addJob(pattern string, job func(), singleton bool, name ...string) (*Job, error) {
	schedule, err := newSchedule(pattern)
	if err != nil {
		return nil, err
	}
	// No limit for <times>, for gtimer checking scheduling every second.
	entry := &Job{
		cron:     c,
		schedule: schedule,
		jobName:  runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name(),
		times:    gtype.NewInt(defaultTimes),
		Job:      job,
		Time:     time.Now(),
	}
	if len(name) > 0 {
		entry.Name = name[0]
	} else {
		entry.Name = "gcron-" + gconv.String(c.idGen.Add(1))
	}
	// When you add a scheduled task, you cannot allow it to run.
	// It cannot start running when added to gtimer.
	// It should start running after the entry is added to the entries map,
	// to avoid the task from running during adding where the entries
	// does not have the entry information, which might cause panic.
	entry.job = gtimer.AddJob(time.Second, entry.check, singleton, -1, gtimer.StatusStopped)
	c.entries.Set(entry.Name, entry)
	entry.job.Start()
	return entry, nil
}

// IsSingleton return whether this entry is a singleton timed task.
func (entry *Job) IsSingleton() bool {
	return entry.job.IsSingleton()
}

// SetSingleton sets the entry running in singleton mode.
func (entry *Job) SetSingleton(enabled bool) {
	entry.job.SetSingleton(true)
}

// SetTimes sets the times which the entry can run.
func (entry *Job) SetTimes(times int) {
	entry.times.Set(times)
}

// Status returns the status of entry.
func (entry *Job) Status() int {
	return entry.job.Status()
}

// SetStatus sets the status of the entry.
func (entry *Job) SetStatus(status int) int {
	return entry.job.SetStatus(status)
}

// Start starts running the entry.
func (entry *Job) Start() {
	entry.job.Start()
}

// Stop stops running the entry.
func (entry *Job) Stop() {
	entry.job.Stop()
}

// Close stops and removes the entry from cron.
func (entry *Job) Close() {
	entry.cron.entries.Remove(entry.Name)
	entry.job.Close()
}

// Timed task check execution.
// The running times limits feature is implemented by gcron.Job and cannot be implemented by gtimer.Job.
// gcron.Job relies on gtimer to implement a scheduled task check for gcron.Job per second.
func (entry *Job) check() {
	if entry.schedule.meet(time.Now()) {
		path := entry.cron.GetLogPath()
		level := entry.cron.GetLogLevel()
		switch entry.cron.status.Val() {
		case StatusStopped:
			return

		case StatusClosed:
			glog.Path(path).Level(level).Debugf("[gcron] %s(%s) %s removed", entry.Name, entry.schedule.pattern, entry.jobName)
			entry.Close()

		case StatusReady:
			fallthrough
		case StatusRunning:
			// Running times check.
			times := entry.times.Add(-1)
			if times <= 0 {
				if entry.job.SetStatus(StatusClosed) == StatusClosed || times < 0 {
					return
				}
			}
			if times < 2000000000 && times > 1000000000 {
				entry.times.Set(defaultTimes)
			}
			glog.Path(path).Level(level).Debugf("[gcron] %s(%s) %s start", entry.Name, entry.schedule.pattern, entry.jobName)
			defer func() {
				if err := recover(); err != nil {
					glog.Path(path).Level(level).Errorf("[gcron] %s(%s) %s end with error: %v", entry.Name, entry.schedule.pattern, entry.jobName, err)
				} else {
					glog.Path(path).Level(level).Debugf("[gcron] %s(%s) %s end", entry.Name, entry.schedule.pattern, entry.jobName)
				}
				if entry.job.Status() == StatusClosed {
					entry.Close()
				}
			}()
			entry.Job()

		}
	}
}

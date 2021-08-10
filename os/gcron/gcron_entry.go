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

// Entry is timing task entry.
type Entry struct {
	cron     *Cron         // Cron object belonged to.
	entry    *gtimer.Entry // Associated gtimer.Entry.
	schedule *cronSchedule // Timed schedule object.
	jobName  string        // Callback function name(address info).
	times    *gtype.Int    // Running times limit.
	Name     string        // Entry name.
	Job      func()        `json:"-"` // Callback function.
	TimedJob func(int64)   `json:"-"` // Callback function accepting the schedule time.
	Time     time.Time     // Registered time.
}

// addEntry creates and returns a new Entry object.
// Param <job> is the callback function for timed task execution.
// Param <singleton> specifies whether timed task executing in singleton mode.
// Param <name> names this entry for manual control.
func (c *Cron) addEntry(pattern string, job func(), singleton bool, name ...string) (*Entry, error) {
	entry, err := c.addEntryCommon(pattern, singleton, name...)
	if err != nil {
		return nil, err
	}
	// When you add a scheduled task, you cannot allow it to run.
	// It cannot start running when added to gtimer.
	// It should start running after the entry is added to the entries map,
	// to avoid the task from running during adding where the entries
	// does not have the entry information, which might cause panic.
	entry.entry = gtimer.AddEntry(time.Second, entry.check, singleton, -1, gtimer.StatusStopped)
	entry.jobName = runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name()
	entry.Job = job
	entry.entry.Start()

	return entry, nil
}


// addTimedJobEntry creates and returns a new Entry object.
// Param <timedJob> is the callback function for timed task execution which accepts the triggered tick as parameter
// Param <singleton> specifies whether timed task executing in singleton mode.
// Param <name> names this entry for manual control.
func (c *Cron) addTimedJobEntry(pattern string, timedJob func(int64), singleton bool, name ...string) (*Entry, error) {
	entry, err := c.addEntryCommon(pattern, singleton, name...)
	if err != nil {
		return nil, err
	}
	// When you add a scheduled task, you cannot allow it to run.
	// It cannot start running when added to gtimer.
	// It should start running after the entry is added to the entries map,
	// to avoid the task from running during adding where the entries
	// does not have the entry information, which might cause panic.
	entry.entry = gtimer.AddTimedJobEntry(time.Second, entry.checkTimedJob, singleton, -1, gtimer.StatusStopped)
	entry.jobName = runtime.FuncForPC(reflect.ValueOf(timedJob).Pointer()).Name()
	entry.TimedJob = timedJob
	entry.entry.Start()
	return entry, nil
}

func (c *Cron) addEntryCommon(pattern string, singleton bool, name ...string) (*Entry, error) {
	schedule, err := newSchedule(pattern)
	if err != nil {
		return nil, err
	}
	// No limit for <times>, for gtimer checking scheduling every second.
	entry := &Entry{
		cron:     c,
		schedule: schedule,
		times:    gtype.NewInt(defaultTimes),
		Time:     time.Now(),
	}
	if len(name) > 0 {
		entry.Name = name[0]
	} else {
		entry.Name = "gcron-" + gconv.String(c.idGen.Add(1))
	}
	c.entries.Set(entry.Name, entry)
	return entry, nil
}

// IsSingleton return whether this entry is a singleton timed task.
func (entry *Entry) IsSingleton() bool {
	return entry.entry.IsSingleton()
}

// SetSingleton sets the entry running in singleton mode.
func (entry *Entry) SetSingleton(enabled bool) {
	entry.entry.SetSingleton(true)
}

// SetTimes sets the times which the entry can run.
func (entry *Entry) SetTimes(times int) {
	entry.times.Set(times)
}

// Status returns the status of entry.
func (entry *Entry) Status() int {
	return entry.entry.Status()
}

// SetStatus sets the status of the entry.
func (entry *Entry) SetStatus(status int) int {
	return entry.entry.SetStatus(status)
}

// Start starts running the entry.
func (entry *Entry) Start() {
	entry.entry.Start()
}

// Stop stops running the entry.
func (entry *Entry) Stop() {
	entry.entry.Stop()
}

// Close stops and removes the entry from cron.
func (entry *Entry) Close() {
	entry.cron.entries.Remove(entry.Name)
	entry.entry.Close()
}

// check verifies if the job can be run and if yes, runs it after the checks are successful
func (entry *Entry) check() {
	if entry.canRunJob() {
		entry.Job()
	}
}

// checkTimedJob verifies if the timedJob can be run and if yes, runs it after the checks are successful
func (entry *Entry) checkTimedJob(ticks int64) {
	if entry.canRunJob() {
		entry.TimedJob(ticks)
	}
}

// canRunJob checks if the timing task can run the job function.
// The running times limits feature is implemented by gcron.Entry and cannot be implemented by gtimer.Entry.
// gcron.Entry relies on gtimer to implement a scheduled task check for gcron.Entry per second.
func (entry *Entry) canRunJob() bool {
	if !entry.schedule.meet(time.Now()) {
		return false
	}

	var (
		path  = entry.cron.GetLogPath()
		level = entry.cron.GetLogLevel()
	)
	switch entry.cron.status.Val() {
	case StatusStopped:
		return false

	case StatusClosed:
		glog.Path(path).Level(level).Debugf("[gcron] %s(%s) %s removed", entry.Name, entry.schedule.pattern, entry.jobName)
		entry.Close()

	case StatusReady:
		fallthrough
	case StatusRunning:
		defer func() {
			if err := recover(); err != nil {
				glog.Path(path).Level(level).Errorf(
					"[gcron] %s(%s) %s end with error: %+v",
					entry.Name, entry.schedule.pattern, entry.jobName, err,
				)
			} else {
				glog.Path(path).Level(level).Debugf(
					"[gcron] %s(%s) %s end",
					entry.Name, entry.schedule.pattern, entry.jobName,
				)
			}
			if entry.entry.Status() == StatusClosed {
				entry.Close()
			}
		}()

		// Running times check.
		times := entry.times.Add(-1)
		if times <= 0 {
			if entry.entry.SetStatus(StatusClosed) == StatusClosed || times < 0 {
				return false
			}
		}
		if times < 2000000000 && times > 1000000000 {
			entry.times.Set(defaultTimes)
		}
		glog.Path(path).Level(level).Debugf("[gcron] %s(%s) %s start", entry.Name, entry.schedule.pattern, entry.jobName)
		return true
	}
	return false
}


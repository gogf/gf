// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gcron

import (
	"reflect"
	"runtime"
	"time"

	"github.com/jin502437344/gf/container/gtype"
	"github.com/jin502437344/gf/os/glog"
	"github.com/jin502437344/gf/os/gtimer"
	"github.com/jin502437344/gf/util/gconv"
)

// Timed task entry.
type Entry struct {
	cron     *Cron         // Cron object belonged to.
	entry    *gtimer.Entry // Associated gtimer.Entry.
	schedule *cronSchedule // Timed schedule object.
	jobName  string        // Callback function name(address info).
	times    *gtype.Int    // Running times limit.
	Name     string        // Entry name.
	Job      func()        `json:"-"` // Callback function.
	Time     time.Time     // Registered time.
}

// addEntry creates and returns a new Entry object.
// Param <job> is the callback function for timed task execution.
// Param <singleton> specifies whether timed task executing in singleton mode.
// Param <name> names this entry for manual control.
func (c *Cron) addEntry(pattern string, job func(), singleton bool, name ...string) (*Entry, error) {
	schedule, err := newSchedule(pattern)
	if err != nil {
		return nil, err
	}
	// No limit for <times>, for gtimer checking scheduling every second.
	entry := &Entry{
		cron:     c,
		schedule: schedule,
		jobName:  runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name(),
		times:    gtype.NewInt(gDEFAULT_TIMES),
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
	entry.entry = gtimer.AddEntry(time.Second, entry.check, singleton, -1, gtimer.STATUS_STOPPED)
	c.entries.Set(entry.Name, entry)
	entry.entry.Start()
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

// Timed task check execution.
// The running times limits feature is implemented by gcron.Entry and cannot be implemented by gtimer.Entry.
// gcron.Entry relies on gtimer to implement a scheduled task check for gcron.Entry per second.
func (entry *Entry) check() {
	if entry.schedule.meet(time.Now()) {
		path := entry.cron.GetLogPath()
		level := entry.cron.GetLogLevel()
		switch entry.cron.status.Val() {
		case STATUS_STOPPED:
			return

		case STATUS_CLOSED:
			glog.Path(path).Level(level).Debugf("[gcron] %s(%s) %s removed", entry.Name, entry.schedule.pattern, entry.jobName)
			entry.Close()

		case STATUS_READY:
			fallthrough
		case STATUS_RUNNING:
			// Running times check.
			times := entry.times.Add(-1)
			if times <= 0 {
				if entry.entry.SetStatus(STATUS_CLOSED) == STATUS_CLOSED || times < 0 {
					return
				}
			}
			if times < 2000000000 && times > 1000000000 {
				entry.times.Set(gDEFAULT_TIMES)
			}
			glog.Path(path).Level(level).Debugf("[gcron] %s(%s) %s start", entry.Name, entry.schedule.pattern, entry.jobName)
			defer func() {
				if err := recover(); err != nil {
					glog.Path(path).Level(level).Errorf("[gcron] %s(%s) %s end with error: %v", entry.Name, entry.schedule.pattern, entry.jobName, err)
				} else {
					glog.Path(path).Level(level).Debugf("[gcron] %s(%s) %s end", entry.Name, entry.schedule.pattern, entry.jobName)
				}
				if entry.entry.Status() == STATUS_CLOSED {
					entry.Close()
				}
			}()
			entry.Job()

		}
	}
}

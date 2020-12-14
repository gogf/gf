// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron

import (
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtimer"
)

type Cron struct {
	idGen    *gtype.Int64    // Used for unique name generation.
	status   *gtype.Int      // Timed task status(0: Not Start; 1: Running; 2: Stopped; -1: Closed)
	entries  *gmap.StrAnyMap // All timed task entries.
	logPath  *gtype.String   // Logging path(folder).
	logLevel *gtype.Int      // Logging level.
}

// New returns a new Cron object with default settings.
func New() *Cron {
	return &Cron{
		idGen:    gtype.NewInt64(),
		status:   gtype.NewInt(StatusRunning),
		entries:  gmap.NewStrAnyMap(true),
		logPath:  gtype.NewString(),
		logLevel: gtype.NewInt(glog.LEVEL_PROD),
	}
}

// SetLogPath sets the logging folder path.
func (c *Cron) SetLogPath(path string) {
	c.logPath.Set(path)
}

// GetLogPath return the logging folder path.
func (c *Cron) GetLogPath() string {
	return c.logPath.Val()
}

// SetLogLevel sets the logging level.
func (c *Cron) SetLogLevel(level int) {
	c.logLevel.Set(level)
}

// GetLogLevel returns the logging level.
func (c *Cron) GetLogLevel() int {
	return c.logLevel.Val()
}

// Add adds a timed task.
// A unique <name> can be bound with the timed task.
// It returns and error if the <name> is already used.
func (c *Cron) Add(pattern string, job func(), name ...string) (*Entry, error) {
	if len(name) > 0 {
		if c.Search(name[0]) != nil {
			return nil, errors.New(fmt.Sprintf(`cron job "%s" already exists`, name[0]))
		}
	}
	return c.addEntry(pattern, job, false, name...)
}

// AddSingleton adds a singleton timed task.
// A singleton timed task is that can only be running one single instance at the same time.
// A unique <name> can be bound with the timed task.
// It returns and error if the <name> is already used.
func (c *Cron) AddSingleton(pattern string, job func(), name ...string) (*Entry, error) {
	if entry, err := c.Add(pattern, job, name...); err != nil {
		return nil, err
	} else {
		entry.SetSingleton(true)
		return entry, nil
	}
}

// AddOnce adds a timed task which can be run only once.
// A unique <name> can be bound with the timed task.
// It returns and error if the <name> is already used.
func (c *Cron) AddOnce(pattern string, job func(), name ...string) (*Entry, error) {
	if entry, err := c.Add(pattern, job, name...); err != nil {
		return nil, err
	} else {
		entry.SetTimes(1)
		return entry, nil
	}
}

// AddTimes adds a timed task which can be run specified times.
// A unique <name> can be bound with the timed task.
// It returns and error if the <name> is already used.
func (c *Cron) AddTimes(pattern string, times int, job func(), name ...string) (*Entry, error) {
	if entry, err := c.Add(pattern, job, name...); err != nil {
		return nil, err
	} else {
		entry.SetTimes(times)
		return entry, nil
	}
}

// DelayAdd adds a timed task after <delay> time.
func (c *Cron) DelayAdd(delay time.Duration, pattern string, job func(), name ...string) {
	gtimer.AddOnce(delay, func() {
		if _, err := c.Add(pattern, job, name...); err != nil {
			panic(err)
		}
	})
}

// DelayAddSingleton adds a singleton timed task after <delay> time.
func (c *Cron) DelayAddSingleton(delay time.Duration, pattern string, job func(), name ...string) {
	gtimer.AddOnce(delay, func() {
		if _, err := c.AddSingleton(pattern, job, name...); err != nil {
			panic(err)
		}
	})
}

// DelayAddOnce adds a timed task after <delay> time.
// This timed task can be run only once.
func (c *Cron) DelayAddOnce(delay time.Duration, pattern string, job func(), name ...string) {
	gtimer.AddOnce(delay, func() {
		if _, err := c.AddOnce(pattern, job, name...); err != nil {
			panic(err)
		}
	})
}

// DelayAddTimes adds a timed task after <delay> time.
// This timed task can be run specified times.
func (c *Cron) DelayAddTimes(delay time.Duration, pattern string, times int, job func(), name ...string) {
	gtimer.AddOnce(delay, func() {
		if _, err := c.AddTimes(pattern, times, job, name...); err != nil {
			panic(err)
		}
	})
}

// Search returns a scheduled task with the specified <name>.
// It returns nil if no found.
func (c *Cron) Search(name string) *Entry {
	if v := c.entries.Get(name); v != nil {
		return v.(*Entry)
	}
	return nil
}

// Start starts running the specified timed task named <name>.
func (c *Cron) Start(name ...string) {
	if len(name) > 0 {
		for _, v := range name {
			if entry := c.Search(v); entry != nil {
				entry.Start()
			}
		}
	} else {
		c.status.Set(StatusReady)
	}
}

// Stop stops running the specified timed task named <name>.
func (c *Cron) Stop(name ...string) {
	if len(name) > 0 {
		for _, v := range name {
			if entry := c.Search(v); entry != nil {
				entry.Stop()
			}
		}
	} else {
		c.status.Set(StatusStopped)
	}
}

// Remove deletes scheduled task which named <name>.
func (c *Cron) Remove(name string) {
	if v := c.entries.Get(name); v != nil {
		v.(*Entry).Close()
	}
}

// Close stops and closes current cron.
func (c *Cron) Close() {
	c.status.Set(StatusClosed)
}

// Size returns the size of the timed tasks.
func (c *Cron) Size() int {
	return c.entries.Size()
}

// Entries return all timed tasks as slice(order by registered time asc).
func (c *Cron) Entries() []*Entry {
	array := garray.NewSortedArraySize(c.entries.Size(), func(v1, v2 interface{}) int {
		entry1 := v1.(*Entry)
		entry2 := v2.(*Entry)
		if entry1.Time.Nanosecond() > entry2.Time.Nanosecond() {
			return 1
		}
		return -1
	}, true)
	c.entries.RLockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			array.Add(v.(*Entry))
		}
	})
	entries := make([]*Entry, array.Len())
	array.RLockFunc(func(array []interface{}) {
		for k, v := range array {
			entries[k] = v.(*Entry)
		}
	})
	return entries
}

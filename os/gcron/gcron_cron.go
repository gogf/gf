// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron

import (
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtimer"
)

type Cron struct {
	idGen   *gtype.Int64    // Used for unique name generation.
	status  *gtype.Int      // Timed task status(0: Not Start; 1: Running; 2: Stopped; -1: Closed)
	entries *gmap.StrAnyMap // All timed task entries.
	logger  *glog.Logger    // Logger, it is nil in default.
}

// New returns a new Cron object with default settings.
func New() *Cron {
	return &Cron{
		idGen:   gtype.NewInt64(),
		status:  gtype.NewInt(StatusRunning),
		entries: gmap.NewStrAnyMap(true),
	}
}

// SetLogger sets the logger for cron.
func (c *Cron) SetLogger(logger *glog.Logger) {
	c.logger = logger
}

// GetLogger returns the logger in the cron.
func (c *Cron) GetLogger() *glog.Logger {
	return c.logger
}

// AddEntry creates and returns a new Entry object.
func (c *Cron) AddEntry(pattern string, job func(), times int, singleton bool, name ...string) (*Entry, error) {
	var (
		entryName = ""
		infinite  = false
	)
	if len(name) > 0 {
		entryName = name[0]
	}
	if times <= 0 {
		infinite = true
	}
	return c.doAddEntry(addEntryInput{
		Name:      entryName,
		Job:       job,
		Times:     times,
		Pattern:   pattern,
		Singleton: singleton,
		Infinite:  infinite,
	})
}

// Add adds a timed task.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func (c *Cron) Add(pattern string, job func(), name ...string) (*Entry, error) {
	return c.AddEntry(pattern, job, -1, false, name...)
}

// AddSingleton adds a singleton timed task.
// A singleton timed task is that can only be running one single instance at the same time.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func (c *Cron) AddSingleton(pattern string, job func(), name ...string) (*Entry, error) {
	return c.AddEntry(pattern, job, -1, true, name...)
}

// AddTimes adds a timed task which can be run specified times.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func (c *Cron) AddTimes(pattern string, times int, job func(), name ...string) (*Entry, error) {
	return c.AddEntry(pattern, job, times, false, name...)
}

// AddOnce adds a timed task which can be run only once.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func (c *Cron) AddOnce(pattern string, job func(), name ...string) (*Entry, error) {
	return c.AddEntry(pattern, job, 1, false, name...)
}

// DelayAddEntry adds a timed task after `delay` time.
func (c *Cron) DelayAddEntry(delay time.Duration, pattern string, job func(), times int, singleton bool, name ...string) {
	gtimer.AddOnce(delay, func() {
		if _, err := c.AddEntry(pattern, job, times, singleton, name...); err != nil {
			panic(err)
		}
	})
}

// DelayAdd adds a timed task after `delay` time.
func (c *Cron) DelayAdd(delay time.Duration, pattern string, job func(), name ...string) {
	gtimer.AddOnce(delay, func() {
		if _, err := c.Add(pattern, job, name...); err != nil {
			panic(err)
		}
	})
}

// DelayAddSingleton adds a singleton timed task after `delay` time.
func (c *Cron) DelayAddSingleton(delay time.Duration, pattern string, job func(), name ...string) {
	gtimer.AddOnce(delay, func() {
		if _, err := c.AddSingleton(pattern, job, name...); err != nil {
			panic(err)
		}
	})
}

// DelayAddOnce adds a timed task after `delay` time.
// This timed task can be run only once.
func (c *Cron) DelayAddOnce(delay time.Duration, pattern string, job func(), name ...string) {
	gtimer.AddOnce(delay, func() {
		if _, err := c.AddOnce(pattern, job, name...); err != nil {
			panic(err)
		}
	})
}

// DelayAddTimes adds a timed task after `delay` time.
// This timed task can be run specified times.
func (c *Cron) DelayAddTimes(delay time.Duration, pattern string, times int, job func(), name ...string) {
	gtimer.AddOnce(delay, func() {
		if _, err := c.AddTimes(pattern, times, job, name...); err != nil {
			panic(err)
		}
	})
}

// Search returns a scheduled task with the specified `name`.
// It returns nil if not found.
func (c *Cron) Search(name string) *Entry {
	if v := c.entries.Get(name); v != nil {
		return v.(*Entry)
	}
	return nil
}

// Start starts running the specified timed task named `name`.
// If no`name` specified, it starts the entire cron.
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

// Stop stops running the specified timed task named `name`.
// If no`name` specified, it stops the entire cron.
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

// Remove deletes scheduled task which named `name`.
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

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcron implements a cron pattern parser and job runner.
package gcron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtimer"
)

const (
	StatusReady   = gtimer.StatusReady
	StatusRunning = gtimer.StatusRunning
	StatusStopped = gtimer.StatusStopped
	StatusClosed  = gtimer.StatusClosed
)

var (
	// Default cron object.
	defaultCron = New()
)

// SetLogger sets the logger for cron.
func SetLogger(logger glog.ILogger) {
	defaultCron.SetLogger(logger)
}

// GetLogger returns the logger in the cron.
func GetLogger() glog.ILogger {
	return defaultCron.GetLogger()
}

// Add adds a timed task to default cron object.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func Add(ctx context.Context, pattern string, job JobFunc, name ...string) (*Entry, error) {
	return defaultCron.Add(ctx, pattern, job, name...)
}

// AddSingleton adds a singleton timed task, to default cron object.
// A singleton timed task is that can only be running one single instance at the same time.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func AddSingleton(ctx context.Context, pattern string, job JobFunc, name ...string) (*Entry, error) {
	return defaultCron.AddSingleton(ctx, pattern, job, name...)
}

// AddParallelSingleton adds a singleton timed task, to default cron object.
// Multiple tasks will be executed parallel (async).
// Simplifying form of the AddSingletons that have same pattern.
func AddParallelSingleton(ctx context.Context, pattern string, jobs ...JobFunc) (*Entry, error) {
	return AddParallelSingletonWithName(ctx, pattern, jobs)
}

// AddParallelSingletonWithName adds a singleton timed task, to default cron object.
// Multiple tasks with name will be executed parallel (async).
// Simplifying form of the AddSingletons that have same pattern.
func AddParallelSingletonWithName(ctx context.Context, pattern string, jobs []JobFunc, name ...string) (*Entry, error) {
	return defaultCron.AddSingleton(ctx, pattern, func(ctx context.Context) { waitAsync(ctx, jobs...) }, name...)
}

// AddSerialSingleton adds a singleton timed task, to default cron object.
// Multiple tasks will be executed serially (sync).
func AddSerialSingleton(ctx context.Context, pattern string, jobs ...JobFunc) (*Entry, error) {
	return AddSerialSingletonWithName(ctx, pattern, jobs)
}

// AddSerialSingletonWithName adds a singleton timed task, to default cron object.
// Multiple tasks with name will be executed serially (sync).
func AddSerialSingletonWithName(ctx context.Context, pattern string, jobs []JobFunc, name ...string) (*Entry, error) {
	return defaultCron.AddSingleton(ctx, pattern, func(ctx context.Context) {
		var i int
		defer func() {
			if e := recover(); e != nil {
				panic(fmt.Errorf("jobs[%v] error: \n%v", i, e))
			}
		}()
		for i = range jobs {
			jobs[i](ctx)
		}
	}, name...)
}

// AddSerialGroupSingleton adds a singleton timed task, to default cron object.
// The groups will be executed serially (sync), and the tasks in a group will be executed parallel (async).
func AddSerialGroupSingleton(ctx context.Context, pattern string, jobGroups ...[]JobFunc) (*Entry, error) {
	return AddSerialGroupSingletonWithName(ctx, pattern, jobGroups)
}

// AddSerialGroupSingletonWithName adds a singleton timed task, to default cron object.
// The groups will be executed serially (sync), and the tasks in a group will be executed parallel (async).
func AddSerialGroupSingletonWithName(ctx context.Context, pattern string, jobGroups [][]JobFunc, name ...string) (*Entry, error) {
	return defaultCron.AddSingleton(ctx, pattern, func(ctx context.Context) {
		var i int
		defer func() {
			if e := recover(); e != nil {
				panic(fmt.Errorf("waitGroupSync Groups[%v] error: \n%v", i, e))
			}
		}()
		for i = range jobGroups {
			waitAsync(ctx, jobGroups[i]...)
		}
	}, name...)
}

// waitAsync execute jobs (async), then wait util all jobs done.
func waitAsync(ctx context.Context, jobs ...JobFunc) {
	var (
		err error
		wg  = sync.WaitGroup{}
		mu  = sync.Mutex{}
	)
	newErr := func(i int, e interface{}) {
		mu.Lock()
		if err != nil {
			err = fmt.Errorf("%v\nwaitAsync Tasks[%v] error: %v ", err, i, e)
		} else {
			err = fmt.Errorf("waitAsync Tasks[%v] error: %v ", i, e)
		}
		mu.Unlock()
	}
	wg.Add(len(jobs))
	for i := 0; i < len(jobs); i++ {
		go func(i int) {
			defer func() {
				if e := recover(); e != nil {
					newErr(i, e)
					wg.Done()
				}
			}()
			jobs[i](ctx)
			wg.Done()
		}(i)
	}
	wg.Wait()
	if err != nil {
		panic(err)
	}
}

// AddOnce adds a timed task which can be run only once, to default cron object.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func AddOnce(ctx context.Context, pattern string, job JobFunc, name ...string) (*Entry, error) {
	return defaultCron.AddOnce(ctx, pattern, job, name...)
}

// AddTimes adds a timed task which can be run specified times, to default cron object.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func AddTimes(ctx context.Context, pattern string, times int, job JobFunc, name ...string) (*Entry, error) {
	return defaultCron.AddTimes(ctx, pattern, times, job, name...)
}

// DelayAdd adds a timed task to default cron object after `delay` time.
func DelayAdd(ctx context.Context, delay time.Duration, pattern string, job JobFunc, name ...string) {
	defaultCron.DelayAdd(ctx, delay, pattern, job, name...)
}

// DelayAddSingleton adds a singleton timed task after `delay` time to default cron object.
func DelayAddSingleton(ctx context.Context, delay time.Duration, pattern string, job JobFunc, name ...string) {
	defaultCron.DelayAddSingleton(ctx, delay, pattern, job, name...)
}

// DelayAddOnce adds a timed task after `delay` time to default cron object.
// This timed task can be run only once.
func DelayAddOnce(ctx context.Context, delay time.Duration, pattern string, job JobFunc, name ...string) {
	defaultCron.DelayAddOnce(ctx, delay, pattern, job, name...)
}

// DelayAddTimes adds a timed task after `delay` time to default cron object.
// This timed task can be run specified times.
func DelayAddTimes(ctx context.Context, delay time.Duration, pattern string, times int, job JobFunc, name ...string) {
	defaultCron.DelayAddTimes(ctx, delay, pattern, times, job, name...)
}

// Search returns a scheduled task with the specified `name`.
// It returns nil if no found.
func Search(name string) *Entry {
	return defaultCron.Search(name)
}

// Remove deletes scheduled task which named `name`.
func Remove(name string) {
	defaultCron.Remove(name)
}

// Size returns the size of the timed tasks of default cron.
func Size() int {
	return defaultCron.Size()
}

// Entries return all timed tasks as slice.
func Entries() []*Entry {
	return defaultCron.Entries()
}

// Start starts running the specified timed task named `name`.
// If no`name` specified, it starts the entire cron.
func Start(name ...string) {
	defaultCron.Start(name...)
}

// Stop stops running the specified timed task named `name`.
// If no`name` specified, it stops the entire cron.
func Stop(name ...string) {
	defaultCron.Stop(name...)
}

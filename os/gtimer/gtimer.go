// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gtimer implements timer for interval/delayed jobs running and management.
//
// This package is designed for management for millions of timing jobs. The differences
// between gtimer and gcron are as follows:
//  1. package gcron is implemented based on package gtimer.
//  2. gtimer is designed for high performance and for millions of timing jobs.
//  3. gcron supports configuration pattern grammar like linux crontab, which is more manually
//     readable.
//  4. gtimer's benchmark OP is measured in nanoseconds, and gcron's benchmark OP is measured
//     in microseconds.
//
// ALSO VERY NOTE the common delay of the timer: https://github.com/golang/go/issues/14410
package gtimer

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/command"
)

// Timer is the timer manager, which uses ticks to calculate the timing interval.
type Timer struct {
	mu      sync.RWMutex
	queue   *priorityQueue // queue is a priority queue based on heap structure.
	status  *gtype.Int     // status is the current timer status.
	ticks   *gtype.Int64   // ticks is the proceeded interval number by the timer.
	options TimerOptions   // timer options is used for timer configuration.
}

// TimerOptions is the configuration object for Timer.
type TimerOptions struct {
	Interval time.Duration // Interval is the interval escaped of the timer.
}

// internalPanic is the custom panic for internal usage.
type internalPanic string

const (
	StatusReady                        = 0      // Job or Timer is ready for running.
	StatusRunning                      = 1      // Job or Timer is already running.
	StatusStopped                      = 2      // Job or Timer is stopped.
	StatusClosed                       = -1     // Job or Timer is closed and waiting to be deleted.
	panicExit            internalPanic = "exit" // panicExit is used for custom job exit with panic.
	defaultTimerInterval               = "100"  // defaultTimerInterval is the default timer interval in milliseconds.
	// commandEnvKeyForInterval is the key for command argument or environment configuring default interval duration for timer.
	commandEnvKeyForInterval = "gf.gtimer.interval"
)

var (
	defaultInterval = getDefaultInterval()
	defaultTimer    = New()
)

func getDefaultInterval() time.Duration {
	interval := command.GetOptWithEnv(commandEnvKeyForInterval, defaultTimerInterval)
	n, err := strconv.Atoi(interval)
	if err != nil {
		panic(gerror.WrapCodef(
			gcode.CodeInvalidConfiguration, err, `error converting string "%s" to int number`,
			interval,
		))
	}
	return time.Duration(n) * time.Millisecond
}

// DefaultOptions creates and returns a default options object for Timer creation.
func DefaultOptions() TimerOptions {
	return TimerOptions{
		Interval: defaultInterval,
	}
}

// SetTimeout runs the job once after duration of `delay`.
// It is like the one in javascript.
func SetTimeout(ctx context.Context, delay time.Duration, job JobFunc) {
	AddOnce(ctx, delay, job)
}

// SetInterval runs the job every duration of `delay`.
// It is like the one in javascript.
func SetInterval(ctx context.Context, interval time.Duration, job JobFunc) {
	Add(ctx, interval, job)
}

// Add adds a timing job to the default timer, which runs in interval of `interval`.
func Add(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.Add(ctx, interval, job)
}

// AddEntry adds a timing job to the default timer with detailed parameters.
//
// The parameter `interval` specifies the running interval of the job.
//
// The parameter `singleton` specifies whether the job running in singleton mode.
// There's only one of the same job is allowed running when its a singleton mode job.
//
// The parameter `times` specifies limit for the job running times, which means the job
// exits if its run times exceeds the `times`.
//
// The parameter `status` specifies the job status when it's firstly added to the timer.
func AddEntry(ctx context.Context, interval time.Duration, job JobFunc, isSingleton bool, times int, status int) *Entry {
	return defaultTimer.AddEntry(ctx, interval, job, isSingleton, times, status)
}

// AddSingleton is a convenience function for add singleton mode job.
func AddSingleton(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.AddSingleton(ctx, interval, job)
}

// AddOnce is a convenience function for adding a job which only runs once and then exits.
func AddOnce(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.AddOnce(ctx, interval, job)
}

// AddTimes is a convenience function for adding a job which is limited running times.
func AddTimes(ctx context.Context, interval time.Duration, times int, job JobFunc) *Entry {
	return defaultTimer.AddTimes(ctx, interval, times, job)
}

// DelayAdd adds a timing job after delay of `interval` duration.
// Also see Add.
func DelayAdd(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAdd(ctx, delay, interval, job)
}

// DelayAddEntry adds a timing job after delay of `interval` duration.
// Also see AddEntry.
func DelayAddEntry(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc, isSingleton bool, times int, status int) {
	defaultTimer.DelayAddEntry(ctx, delay, interval, job, isSingleton, times, status)
}

// DelayAddSingleton adds a timing job after delay of `interval` duration.
// Also see AddSingleton.
func DelayAddSingleton(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAddSingleton(ctx, delay, interval, job)
}

// DelayAddOnce adds a timing job after delay of `interval` duration.
// Also see AddOnce.
func DelayAddOnce(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAddOnce(ctx, delay, interval, job)
}

// DelayAddTimes adds a timing job after delay of `interval` duration.
// Also see AddTimes.
func DelayAddTimes(ctx context.Context, delay time.Duration, interval time.Duration, times int, job JobFunc) {
	defaultTimer.DelayAddTimes(ctx, delay, interval, times, job)
}

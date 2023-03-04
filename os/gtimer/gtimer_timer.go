// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
)

// New creates and returns a Timer.
func New(options ...TimerOptions) *Timer {
	t := &Timer{
		queue:  newPriorityQueue(),
		status: gtype.NewInt(StatusRunning),
		ticks:  gtype.NewInt64(),
	}
	if len(options) > 0 {
		t.options = options[0]
	} else {
		t.options = DefaultOptions()
	}
	go t.loop()
	return t
}

// Add adds a timing job to the timer, which runs in interval of `interval`.
func (t *Timer) Add(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Ctx:         ctx,
		Interval:    interval,
		Job:         job,
		IsSingleton: false,
		Times:       -1,
		Status:      StatusReady,
	})
}

// AddEntry adds a timing job to the timer with detailed parameters.
//
// The parameter `interval` specifies the running interval of the job.
//
// The parameter `singleton` specifies whether the job running in singleton mode.
// There's only one of the same job is allowed running when it's a singleton mode job.
//
// The parameter `times` specifies limit for the job running times, which means the job
// exits if its run times exceeds the `times`.
//
// The parameter `status` specifies the job status when it's firstly added to the timer.
func (t *Timer) AddEntry(ctx context.Context, interval time.Duration, job JobFunc, isSingleton bool, times int, status int) *Entry {
	return t.createEntry(createEntryInput{
		Ctx:         ctx,
		Interval:    interval,
		Job:         job,
		IsSingleton: isSingleton,
		Times:       times,
		Status:      status,
	})
}

// AddSingleton is a convenience function for add singleton mode job.
func (t *Timer) AddSingleton(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Ctx:         ctx,
		Interval:    interval,
		Job:         job,
		IsSingleton: true,
		Times:       -1,
		Status:      StatusReady,
	})
}

// AddOnce is a convenience function for adding a job which only runs once and then exits.
func (t *Timer) AddOnce(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Ctx:         ctx,
		Interval:    interval,
		Job:         job,
		IsSingleton: true,
		Times:       1,
		Status:      StatusReady,
	})
}

// AddTimes is a convenience function for adding a job which is limited running times.
func (t *Timer) AddTimes(ctx context.Context, interval time.Duration, times int, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Ctx:         ctx,
		Interval:    interval,
		Job:         job,
		IsSingleton: true,
		Times:       times,
		Status:      StatusReady,
	})
}

// DelayAdd adds a timing job after delay of `delay` duration.
// Also see Add.
func (t *Timer) DelayAdd(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(ctx, delay, func(ctx context.Context) {
		t.Add(ctx, interval, job)
	})
}

// DelayAddEntry adds a timing job after delay of `delay` duration.
// Also see AddEntry.
func (t *Timer) DelayAddEntry(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc, isSingleton bool, times int, status int) {
	t.AddOnce(ctx, delay, func(ctx context.Context) {
		t.AddEntry(ctx, interval, job, isSingleton, times, status)
	})
}

// DelayAddSingleton adds a timing job after delay of `delay` duration.
// Also see AddSingleton.
func (t *Timer) DelayAddSingleton(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(ctx, delay, func(ctx context.Context) {
		t.AddSingleton(ctx, interval, job)
	})
}

// DelayAddOnce adds a timing job after delay of `delay` duration.
// Also see AddOnce.
func (t *Timer) DelayAddOnce(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(ctx, delay, func(ctx context.Context) {
		t.AddOnce(ctx, interval, job)
	})
}

// DelayAddTimes adds a timing job after delay of `delay` duration.
// Also see AddTimes.
func (t *Timer) DelayAddTimes(ctx context.Context, delay time.Duration, interval time.Duration, times int, job JobFunc) {
	t.AddOnce(ctx, delay, func(ctx context.Context) {
		t.AddTimes(ctx, interval, times, job)
	})
}

// Start starts the timer.
func (t *Timer) Start() {
	t.status.Set(StatusRunning)
}

// Stop stops the timer.
func (t *Timer) Stop() {
	t.status.Set(StatusStopped)
}

// Close closes the timer.
func (t *Timer) Close() {
	t.status.Set(StatusClosed)
}

type createEntryInput struct {
	Ctx         context.Context
	Interval    time.Duration
	Job         JobFunc
	IsSingleton bool
	Times       int
	Status      int
}

// createEntry creates and adds a timing job to the timer.
func (t *Timer) createEntry(in createEntryInput) *Entry {
	var (
		infinite  = false
		nextTicks int64
	)
	if in.Times <= 0 {
		infinite = true
	}
	var (
		intervalTicksOfJob = int64(in.Interval / t.options.Interval)
	)
	if intervalTicksOfJob == 0 {
		// If the given interval is lesser than the one of the wheel,
		// then sets it to one tick, which means it will be run in one interval.
		intervalTicksOfJob = 1
	}
	if t.options.Quick {
		// If the quick mode is enabled, which means it will be run right now.
		// Don't need to wait for the first interval.
		nextTicks = t.ticks.Val()
	} else {
		nextTicks = t.ticks.Val() + intervalTicksOfJob
	}
	var (
		entry = &Entry{
			job:         in.Job,
			ctx:         in.Ctx,
			timer:       t,
			ticks:       intervalTicksOfJob,
			times:       gtype.NewInt(in.Times),
			status:      gtype.NewInt(in.Status),
			isSingleton: gtype.NewBool(in.IsSingleton),
			nextTicks:   gtype.NewInt64(nextTicks),
			infinite:    gtype.NewBool(infinite),
		}
	)
	t.queue.Push(entry, nextTicks)
	return entry
}

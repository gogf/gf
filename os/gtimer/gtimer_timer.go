// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer

import (
	"github.com/gogf/gf/container/gtype"
	"time"
)

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

// Add adds a timing job to the timer, which runs in interval of <interval>.
func (t *Timer) Add(interval time.Duration, job JobFunc) *Job {
	return t.createJob(interval, job, false, defaultTimes, StatusReady)
}

// AddJob adds a timing job to the timer with detailed parameters.
//
// The parameter <interval> specifies the running interval of the job.
//
// The parameter <singleton> specifies whether the job running in singleton mode.
// There's only one of the same job is allowed running when its a singleton mode job.
//
// The parameter <times> specifies limit for the job running times, which means the job
// exits if its run times exceeds the <times>.
//
// The parameter <status> specifies the job status when it's firstly added to the timer.
func (t *Timer) AddJob(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Job {
	return t.createJob(interval, job, singleton, times, status)
}

// AddSingleton is a convenience function for add singleton mode job.
func (t *Timer) AddSingleton(interval time.Duration, job JobFunc) *Job {
	return t.createJob(interval, job, true, defaultTimes, StatusReady)
}

// AddOnce is a convenience function for adding a job which only runs once and then exits.
func (t *Timer) AddOnce(interval time.Duration, job JobFunc) *Job {
	return t.createJob(interval, job, true, 1, StatusReady)
}

// AddTimes is a convenience function for adding a job which is limited running times.
func (t *Timer) AddTimes(interval time.Duration, times int, job JobFunc) *Job {
	return t.createJob(interval, job, true, times, StatusReady)
}

// DelayAdd adds a timing job after delay of <interval> duration.
// Also see Add.
func (t *Timer) DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.Add(interval, job)
	})
}

// DelayAddJob adds a timing job after delay of <interval> duration.
// Also see AddJob.
func (t *Timer) DelayAddJob(delay time.Duration, interval time.Duration, job JobFunc, singleton bool, times int, status int) {
	t.AddOnce(delay, func() {
		t.AddJob(interval, job, singleton, times, status)
	})
}

// DelayAddSingleton adds a timing job after delay of <interval> duration.
// Also see AddSingleton.
func (t *Timer) DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddSingleton(interval, job)
	})
}

// DelayAddOnce adds a timing job after delay of <interval> duration.
// Also see AddOnce.
func (t *Timer) DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddOnce(interval, job)
	})
}

// DelayAddTimes adds a timing job after delay of <interval> duration.
// Also see AddTimes.
func (t *Timer) DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddTimes(interval, times, job)
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

// createJob creates and adds a timing job to the timer.
func (t *Timer) createJob(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Job {
	if times <= 0 {
		times = defaultTimes
	}
	var (
		intervalTicksOfJob = int64(interval / t.options.Interval)
	)
	if intervalTicksOfJob == 0 {
		// If the given interval is lesser than the one of the wheel,
		// then sets it to one tick, which means it will be run in one interval.
		intervalTicksOfJob = 1
	}
	nextTicks := t.ticks.Val() + intervalTicksOfJob
	j := &Job{
		job:       job,
		timer:     t,
		ticks:     intervalTicksOfJob,
		times:     gtype.NewInt(times),
		status:    gtype.NewInt(status),
		singleton: gtype.NewBool(singleton),
		nextTicks: gtype.NewInt64(nextTicks),
	}
	t.queue.Push(j, nextTicks)
	return j
}

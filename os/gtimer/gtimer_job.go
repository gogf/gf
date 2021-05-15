// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer

import (
	"github.com/gogf/gf/container/gtype"
	"math"
)

// Job is the timing job.
type Job struct {
	job       JobFunc      // The job function.
	timer     *Timer       // Belonged timer.
	ticks     int64        // The job runs every ticks.
	times     *gtype.Int   // Limit running times.
	status    *gtype.Int   // Job status.
	singleton *gtype.Bool  // Singleton mode.
	nextTicks *gtype.Int64 // Next run ticks of the job.
}

// JobFunc is the job function.
type JobFunc = func()

// Status returns the status of the job.
func (j *Job) Status() int {
	return j.status.Val()
}

// Run runs the timer job asynchronously.
func (j *Job) Run() {
	leftRunningTimes := j.times.Add(-1)
	if leftRunningTimes < 0 {
		j.status.Set(StatusClosed)
		return
	}
	// This means it does not limit the running times.
	// I know it's ugly, but it is surely high performance for running times limit.
	if leftRunningTimes < 2000000000 && leftRunningTimes > 1000000000 {
		j.times.Set(math.MaxInt32)
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if err != panicExit {
					panic(err)
				} else {
					j.Close()
					return
				}
			}
			if j.Status() == StatusRunning {
				j.SetStatus(StatusReady)
			}
		}()
		j.job()
	}()
}

// doCheckAndRunByTicks checks the if job can run in given timer ticks,
// it runs asynchronously if the given `currentTimerTicks` meets or else
// it increments its ticks and waits for next running check.
func (j *Job) doCheckAndRunByTicks(currentTimerTicks int64) {
	// Ticks check.
	if currentTimerTicks < j.nextTicks.Val() {
		return
	}
	j.nextTicks.Set(currentTimerTicks + j.ticks)
	// Perform job checking.
	switch j.status.Val() {
	case StatusRunning:
		if j.IsSingleton() {
			return
		}
	case StatusReady:
		if !j.status.Cas(StatusReady, StatusRunning) {
			return
		}
	case StatusStopped:
		return
	case StatusClosed:
		return
	}
	// Perform job running.
	j.Run()
}

// SetStatus custom sets the status for the job.
func (j *Job) SetStatus(status int) int {
	return j.status.Set(status)
}

// Start starts the job.
func (j *Job) Start() {
	j.status.Set(StatusReady)
}

// Stop stops the job.
func (j *Job) Stop() {
	j.status.Set(StatusStopped)
}

// Close closes the job, and then it will be removed from the timer.
func (j *Job) Close() {
	j.status.Set(StatusClosed)
}

// Reset reset the job, which resets its ticks for next running.
func (j *Job) Reset() {
	j.nextTicks.Set(j.timer.ticks.Val() + j.ticks)
}

// IsSingleton checks and returns whether the job in singleton mode.
func (j *Job) IsSingleton() bool {
	return j.singleton.Val()
}

// SetSingleton sets the job singleton mode.
func (j *Job) SetSingleton(enabled bool) {
	j.singleton.Set(enabled)
}

// Job returns the job function of this job.
func (j *Job) Job() JobFunc {
	return j.job
}

// SetTimes sets the limit running times for the job.
func (j *Job) SetTimes(times int) {
	j.times.Set(times)
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer

import "time"

// loop starts the ticker using a standalone goroutine.
func (t *Timer) loop() {
	go func() {
		var (
			currentTimerTicks   int64
			timerIntervalTicker = time.NewTicker(t.options.Interval)
		)
		defer timerIntervalTicker.Stop()
		for {
			select {
			case <-timerIntervalTicker.C:
				// Check the timer status.
				switch t.status.Val() {
				case StatusRunning:
					// Timer proceeding.
					if currentTimerTicks = t.ticks.Add(1); currentTimerTicks >= t.queue.NextPriority() {
						t.proceed(currentTimerTicks)
					}

				case StatusStopped:
					// Do nothing.

				case StatusClosed:
					// Timer exits.
					return
				}
			}
		}
	}()
}

// proceed function proceeds the timer job checking and running logic.
func (t *Timer) proceed(currentTimerTicks int64) {
	var (
		value interface{}
	)
	for {
		value = t.queue.Pop()
		if value == nil {
			break
		}
		entry := value.(*Entry)
		// It checks if it meets the ticks' requirement.
		if jobNextTicks := entry.nextTicks.Val(); currentTimerTicks < jobNextTicks {
			// It pushes the job back if current ticks does not meet its running ticks requirement.
			t.queue.Push(entry, entry.nextTicks.Val())
			break
		}
		// It checks the job running requirements and then does asynchronous running.
		entry.doCheckAndRunByTicks(currentTimerTicks)
		// Status check: push back or ignore it.
		if entry.Status() != StatusClosed {
			// It pushes the job back to queue for next running.
			t.queue.Push(entry, entry.nextTicks.Val())
		}
	}
}

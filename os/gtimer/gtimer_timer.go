// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer

import (
	"fmt"
	"time"

	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/container/gtype"
)

// Timer is a Hierarchical Timing Wheel manager for timing jobs.
type Timer struct {
	status     *gtype.Int // Timer status.
	wheels     []*wheel   // The underlying wheels.
	length     int        // Max level of the wheels.
	number     int        // Slot Number of each wheel.
	intervalMs int64      // Interval of the slot in milliseconds.
}

// Wheel is a slot wrapper for timing job install and uninstall.
type wheel struct {
	timer      *Timer        // Belonged timer.
	level      int           // The level in the timer.
	slots      []*glist.List // Slot array.
	number     int64         // Slot Number=len(slots).
	ticks      *gtype.Int64  // Ticked count of the wheel, one tick is one of its interval passed.
	totalMs    int64         // Total duration in milliseconds=number*interval.
	createMs   int64         // Created timestamp in milliseconds.
	intervalMs int64         // Interval in milliseconds, which is the duration of one slot.
}

// New creates and returns a Hierarchical Timing Wheel designed timer.
// The parameter <interval> specifies the interval of the timer.
// The optional parameter <level> specifies the wheels count of the timer,
// which is defaultWheelLevel in default.
func New(slot int, interval time.Duration, level ...int) *Timer {
	if slot <= 0 {
		panic(fmt.Sprintf("invalid slot number: %d", slot))
	}
	length := defaultWheelLevel
	if len(level) > 0 {
		length = level[0]
	}
	t := &Timer{
		status:     gtype.NewInt(StatusRunning),
		wheels:     make([]*wheel, length),
		length:     length,
		number:     slot,
		intervalMs: interval.Nanoseconds() / 1e6,
	}
	for i := 0; i < length; i++ {
		if i > 0 {
			n := time.Duration(t.wheels[i-1].totalMs) * time.Millisecond
			if n <= 0 {
				panic(fmt.Sprintf(`inteval is too large with level: %dms x %d`, interval, length))
			}
			w := t.newWheel(i, slot, n)
			t.wheels[i] = w
			t.wheels[i-1].addEntry(n, w.proceed, false, defaultTimes, StatusReady)
		} else {
			t.wheels[i] = t.newWheel(i, slot, interval)
		}
	}
	t.wheels[0].start()
	return t
}

// newWheel creates and returns a single wheel.
func (t *Timer) newWheel(level int, slot int, interval time.Duration) *wheel {
	w := &wheel{
		timer:      t,
		level:      level,
		slots:      make([]*glist.List, slot),
		number:     int64(slot),
		ticks:      gtype.NewInt64(),
		totalMs:    int64(slot) * interval.Nanoseconds() / 1e6,
		createMs:   time.Now().UnixNano() / 1e6,
		intervalMs: interval.Nanoseconds() / 1e6,
	}
	for i := int64(0); i < w.number; i++ {
		w.slots[i] = glist.New(true)
	}
	return w
}

// Add adds a timing job to the timer, which runs in interval of <interval>.
func (t *Timer) Add(interval time.Duration, job JobFunc) *Entry {
	return t.doAddEntry(interval, job, false, defaultTimes, StatusReady)
}

// AddEntry adds a timing job to the timer with detailed parameters.
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
func (t *Timer) AddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	return t.doAddEntry(interval, job, singleton, times, status)
}

// AddSingleton is a convenience function for add singleton mode job.
func (t *Timer) AddSingleton(interval time.Duration, job JobFunc) *Entry {
	return t.doAddEntry(interval, job, true, defaultTimes, StatusReady)
}

// AddOnce is a convenience function for adding a job which only runs once and then exits.
func (t *Timer) AddOnce(interval time.Duration, job JobFunc) *Entry {
	return t.doAddEntry(interval, job, true, 1, StatusReady)
}

// AddTimes is a convenience function for adding a job which is limited running times.
func (t *Timer) AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
	return t.doAddEntry(interval, job, true, times, StatusReady)
}

// DelayAdd adds a timing job after delay of <interval> duration.
// Also see Add.
func (t *Timer) DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.Add(interval, job)
	})
}

// DelayAddEntry adds a timing job after delay of <interval> duration.
// Also see AddEntry.
func (t *Timer) DelayAddEntry(delay time.Duration, interval time.Duration, job JobFunc, singleton bool, times int, status int) {
	t.AddOnce(delay, func() {
		t.AddEntry(interval, job, singleton, times, status)
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

// doAddEntry adds a timing job to timer for internal usage.
func (t *Timer) doAddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	return t.wheels[t.getLevelByIntervalMs(interval.Nanoseconds()/1e6)].addEntry(interval, job, singleton, times, status)
}

// doAddEntryByParent adds a timing job to timer with parent entry for internal usage.
func (t *Timer) doAddEntryByParent(interval int64, parent *Entry) *Entry {
	return t.wheels[t.getLevelByIntervalMs(interval)].addEntryByParent(interval, parent)
}

// getLevelByIntervalMs calculates and returns the level of timer wheel with given milliseconds.
func (t *Timer) getLevelByIntervalMs(intervalMs int64) int {
	pos, cmp := t.binSearchIndex(intervalMs)
	switch cmp {
	// If equals to the last comparison value, do not add it directly to this wheel,
	// but loop and continue comparison from the index to the first level,
	// and add it to the proper level wheel.
	case 0:
		fallthrough
	// If lesser than the last comparison value,
	// loop and continue comparison from the index to the first level,
	// and add it to the proper level wheel.
	case -1:
		i := pos
		for ; i > 0; i-- {
			if intervalMs > t.wheels[i].intervalMs && intervalMs <= t.wheels[i].totalMs {
				return i
			}
		}
		return i

	// If greater than the last comparison value,
	// loop and continue comparison from the index to the last level,
	// and add it to the proper level wheel.
	case 1:
		i := pos
		for ; i < t.length-1; i++ {
			if intervalMs > t.wheels[i].intervalMs && intervalMs <= t.wheels[i].totalMs {
				return i
			}
		}
		return i
	}
	return 0
}

// binSearchIndex uses binary search algorithm for finding the possible level of the wheel
// for the interval value.
func (t *Timer) binSearchIndex(n int64) (index int, result int) {
	min := 0
	max := t.length - 1
	mid := 0
	cmp := -2
	for min <= max {
		mid = min + int((max-min)/2)
		switch {
		case t.wheels[mid].intervalMs == n:
			cmp = 0
		case t.wheels[mid].intervalMs > n:
			cmp = -1
		case t.wheels[mid].intervalMs < n:
			cmp = 1
		}
		switch cmp {
		case -1:
			max = mid - 1
		case 1:
			min = mid + 1
		case 0:
			return mid, cmp
		}
	}
	return mid, cmp
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron

import (
	"time"
)

// Next returns the next time this schedule is activated, greater than the given
// time.  If no time can be found to satisfy the schedule, return the zero time.
func (s *cronSchedule) Next(lastMeetTime time.Time) time.Time {
	if s.everySeconds != 0 {
		var (
			diff  = lastMeetTime.Unix() - s.createTimestamp
			count = diff/s.everySeconds + 1
		)
		return lastMeetTime.Add(time.Duration(count*s.everySeconds) * time.Second)
	}

	var currentTime = lastMeetTime
	if s.ignoreSeconds {
		// Start at the earliest possible time (the upcoming minute).
		currentTime = currentTime.Add(1*time.Minute - time.Duration(currentTime.Nanosecond())*time.Nanosecond)
	} else {
		// Start at the earliest possible time (the upcoming second).
		currentTime = currentTime.Add(1*time.Second - time.Duration(currentTime.Nanosecond())*time.Nanosecond)
	}

	var (
		loc       = currentTime.Location()
		yearLimit = currentTime.Year() + 5
	)

WRAP:
	if currentTime.Year() > yearLimit {
		return currentTime // who will care the job that run in five years later
	}

	for !s.checkMeetMonth(currentTime) {
		currentTime = currentTime.AddDate(0, 1, 0)
		currentTime = time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, loc)
		if currentTime.Month() == time.January {
			goto WRAP
		}
	}
	for !s.checkMeetWeek(currentTime) || !s.checkMeetDay(currentTime) {
		currentTime = currentTime.AddDate(0, 0, 1)
		currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, loc)
		if currentTime.Day() == 1 {
			goto WRAP
		}
	}
	for !s.checkMeetHour(currentTime) {
		currentTime = currentTime.Add(time.Hour)
		currentTime = currentTime.Truncate(time.Hour)
		if currentTime.Hour() == 0 {
			goto WRAP
		}
	}
	for !s.checkMeetMinute(currentTime) {
		currentTime = currentTime.Add(1 * time.Minute)
		currentTime = currentTime.Truncate(time.Minute)
		if currentTime.Minute() == 0 {
			goto WRAP
		}
	}

	for !s.checkMeetSecond(lastMeetTime, currentTime) {
		currentTime = currentTime.Add(1 * time.Second)
		if currentTime.Second() == 0 {
			goto WRAP
		}
	}
	return currentTime.In(loc)
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron

import "time"

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
		added     = false
		yearLimit = currentTime.Year() + 5
	)

WRAP:
	if currentTime.Year() > yearLimit {
		return currentTime // who will care the job that run in five years later
	}

	for !s.checkMeetMonth(lastMeetTime, currentTime) {
		if !added {
			added = true
			currentTime = time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, loc)
		}
		currentTime = currentTime.AddDate(0, 1, 0)
		// need recheck
		if currentTime.Month() == time.January {
			goto WRAP
		}
	}

	for !s.checkMeetDay(lastMeetTime, currentTime) {
		if !added {
			added = true
			currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, loc)
		}
		currentTime = currentTime.AddDate(0, 0, 1)

		// Notice if the hour is no longer midnight due to DST.
		// Add an hour if it's 23, subtract an hour if it's 1.
		if currentTime.Hour() != 0 {
			if currentTime.Hour() > 12 {
				currentTime = currentTime.Add(time.Duration(24-currentTime.Hour()) * time.Hour)
			} else {
				currentTime = currentTime.Add(time.Duration(-currentTime.Hour()) * time.Hour)
			}
		}
		if currentTime.Day() == 1 {
			goto WRAP
		}
	}
	for !s.checkMeetHour(lastMeetTime, currentTime) {
		if !added {
			added = true
			currentTime = time.Date(
				currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), 0, 0, 0, loc,
			)
		}
		currentTime = currentTime.Add(time.Hour)
		// need recheck
		if currentTime.Hour() == 0 {
			goto WRAP
		}
	}
	for !s.checkMeetMinute(lastMeetTime, currentTime) {
		if !added {
			added = true
			currentTime = currentTime.Truncate(time.Minute)
		}
		currentTime = currentTime.Add(1 * time.Minute)

		if currentTime.Minute() == 0 {
			goto WRAP
		}
	}

	if !s.ignoreSeconds {
		for !s.checkMeetSecond(lastMeetTime, currentTime) {
			if !added {
				added = true
				currentTime = currentTime.Truncate(time.Second)
			}
			currentTime = currentTime.Add(1 * time.Second)
			if currentTime.Second() == 0 {
				goto WRAP
			}
		}
	}
	return currentTime.In(loc)
}

// dayMatches returns true if the schedule's day-of-week and day-of-month
// restrictions are satisfied by the given time.
func (s *cronSchedule) dayMatches(t time.Time) bool {
	_, ok1 := s.dayMap[t.Day()]
	_, ok2 := s.weekMap[int(t.Weekday())]
	return ok1 && ok2
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// checkMeetAndUpdateLastSeconds checks if the given time `t` meets the runnable point for the job.
// This function is called every second.
func (s *cronSchedule) checkMeetAndUpdateLastSeconds(ctx context.Context, currentTime time.Time) (ok bool) {
	var (
		lastCheckTimestamp = s.getAndUpdateLastCheckTimestamp(ctx, currentTime)
		lastCheckTime      = gtime.NewFromTimeStamp(lastCheckTimestamp)
		lastMeetTime       = gtime.NewFromTimeStamp(s.lastMeetTimestamp.Val())
	)
	defer func() {
		if ok {
			s.lastMeetTimestamp.Set(currentTime.Unix())
		}
	}()
	if !s.checkMinIntervalAndItemMapMeet(lastMeetTime.Time, lastCheckTime.Time, currentTime) {
		return false
	}
	return true
}

func (s *cronSchedule) checkMinIntervalAndItemMapMeet(
	lastMeetTime, lastCheckTime, currentTime time.Time,
) (ok bool) {
	if s.everySeconds != 0 {
		// It checks using interval.
		secondsAfterCreated := lastCheckTime.UnixNano()/1e9 - s.createTimestamp
		if secondsAfterCreated > 0 {
			return secondsAfterCreated%s.everySeconds == 0
		}
		return false
	}
	if !s.checkMeetSecond(lastMeetTime, currentTime) {
		return false
	}
	if !s.checkMeetMinute(lastMeetTime, currentTime) {
		return false
	}
	if !s.checkMeetHour(lastMeetTime, currentTime) {
		return false
	}
	if !s.checkMeetDay(lastMeetTime, currentTime) {
		return false
	}
	if !s.checkMeetMonth(lastMeetTime, currentTime) {
		return false
	}
	if !s.checkMeetWeek(lastMeetTime, currentTime) {
		return false
	}
	return true
}

func (s *cronSchedule) checkMeetSecond(lastMeetTime, currentTime time.Time) (ok bool) {
	interval := currentTime.Sub(lastMeetTime)
	if interval.Seconds() < float64(s.minIntervalSeconds) {
		return false
	}
	if s.ignoreSeconds {
		if currentTime.Unix()-lastMeetTime.Unix() < 60 {
			return false
		}
	} else {
		if !s.keyMatch(s.secondMap, currentTime.Second()) {
			return false
		}
	}
	return true
}

func (s *cronSchedule) checkMeetMinute(lastMeetTime, currentTime time.Time) (ok bool) {
	interval := currentTime.Sub(lastMeetTime)
	if interval.Minutes() < float64(s.minIntervalMinutes) {
		return false
	}
	if !s.keyMatch(s.minuteMap, currentTime.Minute()) {
		return false
	}
	return true
}

func (s *cronSchedule) checkMeetHour(lastMeetTime, currentTime time.Time) (ok bool) {
	interval := currentTime.Sub(lastMeetTime)
	if interval.Hours() < float64(s.minIntervalHours) {
		return false
	}
	if !s.keyMatch(s.hourMap, currentTime.Hour()) {
		return false
	}
	return true
}

func (s *cronSchedule) checkMeetDay(lastMeetTime, currentTime time.Time) (ok bool) {
	interval := currentTime.Sub(lastMeetTime)
	if interval.Hours()/24 < float64(s.minIntervalDays) {
		return false
	}
	if !s.keyMatch(s.dayMap, currentTime.Day()) {
		return false
	}
	return true
}

func (s *cronSchedule) checkMeetMonth(lastMeetTime, currentTime time.Time) (ok bool) {
	if s.minIntervalMonths > 0 {
		monthDiff := currentTime.Month() - lastMeetTime.Month()
		if monthDiff < 0 {
			monthDiff += 12
		}
		if int64(monthDiff) < s.minIntervalMonths {
			return false
		}
	}
	if !s.keyMatch(s.monthMap, int(currentTime.Month())) {
		return false
	}
	return true
}

func (s *cronSchedule) checkMeetWeek(lastMeetTime, currentTime time.Time) (ok bool) {
	interval := currentTime.Sub(lastMeetTime)
	if interval.Hours()/24 < float64(s.minIntervalWeeks) {
		return false
	}
	if !s.keyMatch(s.weekMap, int(currentTime.Weekday())) {
		return false
	}
	return true
}

func (s *cronSchedule) keyMatch(m map[int]struct{}, key int) bool {
	_, ok := m[key]
	return ok
}

func (s *cronSchedule) checkMinIntervalMeet(lastMeetTime, currentTime time.Time) (ok bool) {
	interval := currentTime.Sub(lastMeetTime)
	if interval.Seconds() < float64(s.minIntervalSeconds) {
		return false
	}
	if interval.Minutes() < float64(s.minIntervalMinutes) {
		return false
	}
	if interval.Hours() < float64(s.minIntervalHours) {
		return false
	}
	if interval.Hours()/24 < float64(s.minIntervalDays) {
		return false
	}
	if s.minIntervalMonths > 0 {
		monthDiff := currentTime.Month() - lastMeetTime.Month()
		if monthDiff < 0 {
			monthDiff += 12
		}
		if int64(monthDiff) < s.minIntervalMonths {
			return false
		}
	}
	if interval.Hours()/24 < float64(s.minIntervalWeeks) {
		return false
	}
	return true
}

func (s *cronSchedule) checkItemMapMeet(lastMeetTime, currentTime time.Time) (ok bool) {
	// second.
	if s.ignoreSeconds {
		if currentTime.Unix()-lastMeetTime.Unix() < 60 {
			return false
		}
	} else {
		if !s.keyMatch(s.secondMap, currentTime.Second()) {
			return false
		}
	}
	// minute.
	if !s.keyMatch(s.minuteMap, currentTime.Minute()) {
		return false
	}
	// hour.
	if !s.keyMatch(s.hourMap, currentTime.Hour()) {
		return false
	}
	// day.
	if !s.keyMatch(s.dayMap, currentTime.Day()) {
		return false
	}
	// month.
	if !s.keyMatch(s.monthMap, int(currentTime.Month())) {
		return false
	}
	// week.
	if !s.keyMatch(s.weekMap, int(currentTime.Weekday())) {
		return false
	}
	return true
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
)

// cronSchedule is the schedule for cron job.
type cronSchedule struct {
	createTimestamp int64            // Created timestamp in seconds.
	everySeconds    int64            // Running interval in seconds.
	pattern         string           // The raw cron pattern string that is passed in cron job creation.
	ignoreSeconds   bool             // Mark the pattern is standard 5 parts crontab pattern instead 6 parts pattern.
	secondMap       map[int]struct{} // Job can run in these second numbers.
	minuteMap       map[int]struct{} // Job can run in these minute numbers.
	hourMap         map[int]struct{} // Job can run in these hour numbers.
	dayMap          map[int]struct{} // Job can run in these day numbers.
	weekMap         map[int]struct{} // Job can run in these week numbers.
	monthMap        map[int]struct{} // Job can run in these moth numbers.

	// Used in pattern which has interval part like number "2" in pattern "# */2 * * * *".
	// The escaped value should be greater or equal than this interval,
	// unit of which is according to its part type.
	minIntervalSeconds int64
	minIntervalMinutes int64
	minIntervalHours   int64
	minIntervalDays    int64
	minIntervalWeeks   int64
	minIntervalMonths  int64

	// This field stores the timestamp that meets schedule latest.
	lastMeetTimestamp *gtype.Int64

	// Last timestamp number, for timestamp fix in some latency.
	lastCheckTimestamp *gtype.Int64
}

type patternItemType int

const (
	patternItemTypeSecond patternItemType = iota
	patternItemTypeMinute
	patternItemTypeHour
	patternItemTypeDay
	patternItemTypeWeek
	patternItemTypeMonth
)

const (
	// regular expression for cron pattern, which contains 6 parts of time units.
	regexForCron = `^([\-/\d\*,#]+)\s+([\-/\d\*,]+)\s+([\-/\d\*,]+)\s+([\-/\d\*\?,]+)\s+([\-/\d\*,A-Za-z]+)\s+([\-/\d\*\?,A-Za-z]+)$`
)

var (
	// Predefined pattern map.
	predefinedPatternMap = map[string]string{
		"@yearly":   "0 0 0 1 1 *",
		"@annually": "0 0 0 1 1 *",
		"@monthly":  "0 0 0 1 * *",
		"@weekly":   "0 0 0 * * 0",
		"@daily":    "0 0 0 * * *",
		"@midnight": "0 0 0 * * *",
		"@hourly":   "0 0 * * * *",
	}
	// Short month name to its number.
	monthShortNameMap = map[string]int{
		"jan": 1,
		"feb": 2,
		"mar": 3,
		"apr": 4,
		"may": 5,
		"jun": 6,
		"jul": 7,
		"aug": 8,
		"sep": 9,
		"oct": 10,
		"nov": 11,
		"dec": 12,
	}
	// Full month name to its number.
	monthFullNameMap = map[string]int{
		"january":   1,
		"february":  2,
		"march":     3,
		"april":     4,
		"may":       5,
		"june":      6,
		"july":      7,
		"august":    8,
		"september": 9,
		"october":   10,
		"november":  11,
		"december":  12,
	}
	// Short week name to its number.
	weekShortNameMap = map[string]int{
		"sun": 0,
		"mon": 1,
		"tue": 2,
		"wed": 3,
		"thu": 4,
		"fri": 5,
		"sat": 6,
	}
	// Full week name to its number.
	weekFullNameMap = map[string]int{
		"sunday":    0,
		"monday":    1,
		"tuesday":   2,
		"wednesday": 3,
		"thursday":  4,
		"friday":    5,
		"saturday":  6,
	}
)

// newSchedule creates and returns a schedule object for given cron pattern.
func newSchedule(pattern string) (*cronSchedule, error) {
	var currentTimestamp = time.Now().Unix()
	// Check given `pattern` if the predefined patterns.
	if match, _ := gregex.MatchString(`(@\w+)\s*(\w*)\s*`, pattern); len(match) > 0 {
		key := strings.ToLower(match[1])
		if v, ok := predefinedPatternMap[key]; ok {
			pattern = v
		} else if strings.Compare(key, "@every") == 0 {
			d, err := gtime.ParseDuration(match[2])
			if err != nil {
				return nil, err
			}
			return &cronSchedule{
				createTimestamp:    currentTimestamp,
				everySeconds:       int64(d.Seconds()),
				pattern:            pattern,
				lastMeetTimestamp:  gtype.NewInt64(currentTimestamp),
				lastCheckTimestamp: gtype.NewInt64(currentTimestamp),
			}, nil
		} else {
			return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid pattern: "%s"`, pattern)
		}
	}
	// Handle given `pattern` as common 6 parts pattern.
	match, _ := gregex.MatchString(regexForCron, pattern)
	if len(match) != 7 {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid pattern: "%s"`, pattern)
	}
	var (
		err error
		cs  = &cronSchedule{
			createTimestamp:    currentTimestamp,
			everySeconds:       0,
			pattern:            pattern,
			lastMeetTimestamp:  gtype.NewInt64(currentTimestamp),
			lastCheckTimestamp: gtype.NewInt64(currentTimestamp),
		}
	)

	// Second.
	if match[1] == "#" {
		cs.ignoreSeconds = true
	} else {
		cs.secondMap, cs.minIntervalSeconds, err = parsePatternItem(match[1], 0, 59, false, patternItemTypeSecond)
		if err != nil {
			return nil, err
		}
	}
	// Minute.
	cs.minuteMap, cs.minIntervalMinutes, err = parsePatternItem(match[2], 0, 59, false, patternItemTypeMinute)
	if err != nil {
		return nil, err
	}
	// Hour.
	cs.hourMap, cs.minIntervalHours, err = parsePatternItem(match[3], 0, 23, false, patternItemTypeHour)
	if err != nil {
		return nil, err
	}
	// Day.
	cs.dayMap, cs.minIntervalDays, err = parsePatternItem(match[4], 1, 31, true, patternItemTypeDay)
	if err != nil {
		return nil, err
	}
	// Month.
	cs.monthMap, cs.minIntervalMonths, err = parsePatternItem(match[5], 1, 12, false, patternItemTypeMonth)
	if err != nil {
		return nil, err
	}
	// Week.
	cs.weekMap, cs.minIntervalWeeks, err = parsePatternItem(match[6], 0, 6, true, patternItemTypeWeek)
	if err != nil {
		return nil, err
	}
	return cs, nil

}

// parsePatternItem parses every item in the pattern and returns the result as map, which is used for indexing.
func parsePatternItem(
	item string, min int, max int,
	allowQuestionMark bool, itemType patternItemType,
) (itemMap map[int]struct{}, itemInterval int64, err error) {
	itemMap = make(map[int]struct{}, max-min+1)
	if item == "*" || (allowQuestionMark && item == "?") {
		for i := min; i <= max; i++ {
			itemMap[i] = struct{}{}
		}
		return itemMap, 0, nil
	}
	// Example: 1-10/2,11-30/3
	var number int
	for _, itemElem := range strings.Split(item, ",") {
		var (
			interval      = 1
			intervalArray = strings.Split(itemElem, "/")
		)
		if len(intervalArray) == 2 {
			if number, err = strconv.Atoi(intervalArray[1]); err != nil {
				return nil, 0, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid pattern item: "%s"`, itemElem)
			} else {
				itemInterval = int64(number)
			}
		}
		var (
			rangeMin   = min
			rangeMax   = max
			rangeArray = strings.Split(intervalArray[0], "-") // Example: 1-30, JAN-DEC
		)
		// Example: 1-30/2
		if rangeArray[0] != "*" {
			if number, err = parseWeekAndMonthNameToInt(rangeArray[0], itemType); err != nil {
				return nil, 0, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid pattern item: "%s"`, itemElem)
			} else {
				rangeMin = number
				if len(intervalArray) == 1 {
					rangeMax = number
				}
			}
			interval = int(itemInterval)
		}
		// Example: 1-30/2
		if len(rangeArray) == 2 {
			if number, err = parseWeekAndMonthNameToInt(rangeArray[1], itemType); err != nil {
				return nil, 0, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid pattern item: "%s"`, itemElem)
			} else {
				rangeMax = number
			}
		}
		for i := rangeMin; i <= rangeMax; i += interval {
			itemMap[i] = struct{}{}
		}
	}
	return
}

// parseWeekAndMonthNameToInt parses the field value to a number according to its field type.
func parseWeekAndMonthNameToInt(value string, itemType patternItemType) (int, error) {
	if gregex.IsMatchString(`^\d+$`, value) {
		// It is pure number.
		if number, err := strconv.Atoi(value); err == nil {
			return number, nil
		}
	} else {
		// Check if it contains letter,
		// it converts the value to number according to predefined map.
		switch itemType {
		case patternItemTypeWeek:
			if number, ok := weekShortNameMap[strings.ToLower(value)]; ok {
				return number, nil
			}
			if number, ok := weekFullNameMap[strings.ToLower(value)]; ok {
				return number, nil
			}
		case patternItemTypeMonth:
			if number, ok := monthShortNameMap[strings.ToLower(value)]; ok {
				return number, nil
			}
			if number, ok := monthFullNameMap[strings.ToLower(value)]; ok {
				return number, nil
			}
		}
	}
	return 0, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid pattern value: "%s"`, value)
}

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

	if s.everySeconds != 0 {
		// It checks using interval.
		secondsAfterCreated := lastCheckTime.Timestamp() - s.createTimestamp
		if secondsAfterCreated > 0 {
			return secondsAfterCreated%s.everySeconds == 0
		}
		return false
	}
	if !s.checkMinIntervalMeet(lastMeetTime.Time, currentTime) {
		return false
	}
	if !s.checkItemMapMeet(lastCheckTime.Time, currentTime) {
		return false
	}
	return true
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

func (s *cronSchedule) checkItemMapMeet(lastCheckTime, currentTime time.Time) (ok bool) {
	// second.
	if s.ignoreSeconds {
		if currentTime.Unix()-s.lastMeetTimestamp.Val() < 60 {
			return false
		}
	} else {
		if !s.keyMatch(s.secondMap, lastCheckTime.Second()) {
			return false
		}
	}
	// minute.
	if !s.keyMatch(s.minuteMap, lastCheckTime.Minute()) {
		return false
	}
	// hour.
	if !s.keyMatch(s.hourMap, lastCheckTime.Hour()) {
		return false
	}
	// day.
	if !s.keyMatch(s.dayMap, lastCheckTime.Day()) {
		return false
	}
	// month.
	if !s.keyMatch(s.monthMap, int(lastCheckTime.Month())) {
		return false
	}
	// week.
	if !s.keyMatch(s.weekMap, int(lastCheckTime.Weekday())) {
		return false
	}
	return true
}

// Next returns the next time this schedule is activated, greater than the given
// time.  If no time can be found to satisfy the schedule, return the zero time.
func (s *cronSchedule) Next(t time.Time) time.Time {
	if s.everySeconds != 0 {
		var (
			diff  = t.Unix() - s.createTimestamp
			count = diff/s.everySeconds + 1
		)
		return t.Add(time.Duration(count*s.everySeconds) * time.Second)
	}

	if s.ignoreSeconds {
		// Start at the earliest possible time (the upcoming minute).
		t = t.Add(1*time.Minute - time.Duration(t.Nanosecond())*time.Nanosecond)
	} else {
		// Start at the earliest possible time (the upcoming second).
		t = t.Add(1*time.Second - time.Duration(t.Nanosecond())*time.Nanosecond)
	}

	var (
		loc       = t.Location()
		added     = false
		yearLimit = t.Year() + 5
	)

WRAP:
	if t.Year() > yearLimit {
		return t // who will care the job that run in five years later
	}

	for !s.keyMatch(s.monthMap, int(t.Month())) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
		}
		t = t.AddDate(0, 1, 0)
		// need recheck
		if t.Month() == time.January {
			goto WRAP
		}
	}

	for !s.dayMatches(t) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		}
		t = t.AddDate(0, 0, 1)

		// Notice if the hour is no longer midnight due to DST.
		// Add an hour if it's 23, subtract an hour if it's 1.
		if t.Hour() != 0 {
			if t.Hour() > 12 {
				t = t.Add(time.Duration(24-t.Hour()) * time.Hour)
			} else {
				t = t.Add(time.Duration(-t.Hour()) * time.Hour)
			}
		}
		if t.Day() == 1 {
			goto WRAP
		}
	}
	for !s.keyMatch(s.hourMap, t.Hour()) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, loc)
		}
		t = t.Add(time.Hour)
		// need recheck
		if t.Hour() == 0 {
			goto WRAP
		}
	}
	for !s.keyMatch(s.minuteMap, t.Minute()) {
		if !added {
			added = true
			t = t.Truncate(time.Minute)
		}
		t = t.Add(1 * time.Minute)

		if t.Minute() == 0 {
			goto WRAP
		}
	}

	if !s.ignoreSeconds {
		for !s.keyMatch(s.secondMap, t.Second()) {
			if !added {
				added = true
				t = t.Truncate(time.Second)
			}
			t = t.Add(1 * time.Second)
			if t.Second() == 0 {
				goto WRAP
			}
		}
	}
	return t.In(loc)
}

// dayMatches returns true if the schedule's day-of-week and day-of-month
// restrictions are satisfied by the given time.
func (s *cronSchedule) dayMatches(t time.Time) bool {
	_, ok1 := s.dayMap[t.Day()]
	_, ok2 := s.weekMap[int(t.Weekday())]
	return ok1 && ok2
}

func (s *cronSchedule) keyMatch(m map[int]struct{}, key int) bool {
	_, ok := m[key]
	return ok
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron

import (
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
		"@yearly":   "# 0 0 1 1 *",
		"@annually": "# 0 0 1 1 *",
		"@monthly":  "# 0 0 1 * *",
		"@weekly":   "# 0 0 * * 0",
		"@daily":    "# 0 0 * * *",
		"@midnight": "# 0 0 * * *",
		"@hourly":   "# 0 * * * *",
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
		cs.secondMap, err = parsePatternItem(match[1], 0, 59, false, patternItemTypeSecond)
		if err != nil {
			return nil, err
		}
	}
	// Minute.
	cs.minuteMap, err = parsePatternItem(match[2], 0, 59, false, patternItemTypeMinute)
	if err != nil {
		return nil, err
	}
	// Hour.
	cs.hourMap, err = parsePatternItem(match[3], 0, 23, false, patternItemTypeHour)
	if err != nil {
		return nil, err
	}
	// Day.
	cs.dayMap, err = parsePatternItem(match[4], 1, 31, true, patternItemTypeDay)
	if err != nil {
		return nil, err
	}
	// Month.
	cs.monthMap, err = parsePatternItem(match[5], 1, 12, false, patternItemTypeMonth)
	if err != nil {
		return nil, err
	}
	// Week.
	cs.weekMap, err = parsePatternItem(match[6], 0, 6, true, patternItemTypeWeek)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

// parsePatternItem parses every item in the pattern and returns the result as map, which is used for indexing.
func parsePatternItem(
	item string, min int, max int,
	allowQuestionMark bool, itemType patternItemType,
) (itemMap map[int]struct{}, err error) {
	itemMap = make(map[int]struct{}, max-min+1)
	if item == "*" || (allowQuestionMark && item == "?") {
		for i := min; i <= max; i++ {
			itemMap[i] = struct{}{}
		}
		return itemMap, nil
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
				return nil, gerror.NewCodef(
					gcode.CodeInvalidParameter, `invalid pattern item: "%s"`, itemElem,
				)
			} else {
				interval = number
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
				return nil, gerror.NewCodef(
					gcode.CodeInvalidParameter, `invalid pattern item: "%s"`, itemElem,
				)
			} else {
				rangeMin = number
				if len(intervalArray) == 1 {
					rangeMax = number
				}
			}
		}
		// Example: 1-30/2
		if len(rangeArray) == 2 {
			if number, err = parseWeekAndMonthNameToInt(rangeArray[1], itemType); err != nil {
				return nil, gerror.NewCodef(
					gcode.CodeInvalidParameter, `invalid pattern item: "%s"`, itemElem,
				)
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

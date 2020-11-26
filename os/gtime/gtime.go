// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gtime provides functionality for measuring and displaying time.
//
// This package should keep much less dependencies with other packages.
package gtime

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/utils"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/text/gregex"
)

const (
	// Short writes for common usage durations.
	D  = 24 * time.Hour
	H  = time.Hour
	M  = time.Minute
	S  = time.Second
	MS = time.Millisecond
	US = time.Microsecond
	NS = time.Nanosecond

	// Regular expression1(datetime separator supports '-', '/', '.').
	// Eg:
	// "2017-12-14 04:51:34 +0805 LMT",
	// "2017-12-14 04:51:34 +0805 LMT",
	// "2006-01-02T15:04:05Z07:00",
	// "2014-01-17T01:19:15+08:00",
	// "2018-02-09T20:46:17.897Z",
	// "2018-02-09 20:46:17.897",
	// "2018-02-09T20:46:17Z",
	// "2018-02-09 20:46:17",
	// "2018/10/31 - 16:38:46"
	// "2018-02-09",
	// "2018.02.09",
	TIME_REAGEX_PATTERN1 = `(\d{4}[-/\.]\d{2}[-/\.]\d{2})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`

	// Regular expression2(datetime separator supports '-', '/', '.').
	// Eg:
	// 01-Nov-2018 11:50:28
	// 01/Nov/2018 11:50:28
	// 01.Nov.2018 11:50:28
	// 01.Nov.2018:11:50:28
	TIME_REAGEX_PATTERN2 = `(\d{1,2}[-/\.][A-Za-z]{3,}[-/\.]\d{4})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`
)

var (
	// It's more high performance using regular expression
	// than time.ParseInLocation to parse the datetime string.
	timeRegex1, _ = regexp.Compile(TIME_REAGEX_PATTERN1)
	timeRegex2, _ = regexp.Compile(TIME_REAGEX_PATTERN2)

	// Month words to arabic numerals mapping.
	monthMap = map[string]int{
		"jan":       1,
		"feb":       2,
		"mar":       3,
		"apr":       4,
		"may":       5,
		"jun":       6,
		"jul":       7,
		"aug":       8,
		"sep":       9,
		"sept":      9,
		"oct":       10,
		"nov":       11,
		"dec":       12,
		"january":   1,
		"february":  2,
		"march":     3,
		"april":     4,
		"june":      6,
		"july":      7,
		"august":    8,
		"september": 9,
		"october":   10,
		"november":  11,
		"december":  12,
	}
)

// SetTimeZone sets the time zone for current whole process.
// The parameter <zone> is an area string specifying corresponding time zone,
// eg: Asia/Shanghai.
//
// Note that the time zone database needed by LoadLocation may not be
// present on all systems, especially non-Unix systems.
// LoadLocation looks in the directory or uncompressed zip file
// named by the ZONEINFO environment variable, if any, then looks in
// known installation locations on Unix systems,
// and finally looks in $GOROOT/lib/time/zoneinfo.zip.
func SetTimeZone(zone string) error {
	location, err := time.LoadLocation(zone)
	if err == nil {
		time.Local = location
	}
	return err
}

// Timestamp retrieves and returns the timestamp in seconds.
func Timestamp() int64 {
	return Now().Timestamp()
}

// TimestampMilli retrieves and returns the timestamp in milliseconds.
func TimestampMilli() int64 {
	return Now().TimestampMilli()
}

// TimestampMicro retrieves and returns the timestamp in microseconds.
func TimestampMicro() int64 {
	return Now().TimestampMicro()
}

// TimestampNano retrieves and returns the timestamp in nanoseconds.
func TimestampNano() int64 {
	return Now().TimestampNano()
}

// TimestampStr is a convenience method which retrieves and returns
// the timestamp in seconds as string.
func TimestampStr() string {
	return Now().TimestampStr()
}

// TimestampMilliStr is a convenience method which retrieves and returns
// the timestamp in milliseconds as string.
func TimestampMilliStr() string {
	return Now().TimestampMilliStr()
}

// TimestampMicroStr is a convenience method which retrieves and returns
// the timestamp in microseconds as string.
func TimestampMicroStr() string {
	return Now().TimestampMicroStr()
}

// TimestampNanoStr is a convenience method which retrieves and returns
// the timestamp in nanoseconds as string.
func TimestampNanoStr() string {
	return Now().TimestampNanoStr()
}

// Second returns the timestamp in seconds.
// Deprecated, use Timestamp instead.
func Second() int64 {
	return Timestamp()
}

// Millisecond returns the timestamp in milliseconds.
// Deprecated, use TimestampMilli instead.
func Millisecond() int64 {
	return TimestampMilli()
}

// Microsecond returns the timestamp in microseconds.
// Deprecated, use TimestampMicro instead.
func Microsecond() int64 {
	return TimestampMicro()
}

// Nanosecond returns the timestamp in nanoseconds.
// Deprecated, use TimestampNano instead.
func Nanosecond() int64 {
	return TimestampNano()
}

// Date returns current date in string like "2006-01-02".
func Date() string {
	return time.Now().Format("2006-01-02")
}

// Datetime returns current datetime in string like "2006-01-02 15:04:05".
func Datetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// ISO8601 returns current datetime in ISO8601 format like "2006-01-02T15:04:05-07:00".
func ISO8601() string {
	return time.Now().Format("2006-01-02T15:04:05-07:00")
}

// ISO8601 returns current datetime in RFC822 format like "Mon, 02 Jan 06 15:04 MST".
func RFC822() string {
	return time.Now().Format("Mon, 02 Jan 06 15:04 MST")
}

// parseDateStr parses the string to year, month and day numbers.
func parseDateStr(s string) (year, month, day int) {
	array := strings.Split(s, "-")
	if len(array) < 3 {
		array = strings.Split(s, "/")
	}
	if len(array) < 3 {
		array = strings.Split(s, ".")
	}
	// Parsing failed.
	if len(array) < 3 {
		return
	}
	// Checking the year in head or tail.
	if utils.IsNumeric(array[1]) {
		year, _ = strconv.Atoi(array[0])
		month, _ = strconv.Atoi(array[1])
		day, _ = strconv.Atoi(array[2])
	} else {
		if v, ok := monthMap[strings.ToLower(array[1])]; ok {
			month = v
		} else {
			return
		}
		year, _ = strconv.Atoi(array[2])
		day, _ = strconv.Atoi(array[0])
	}
	return
}

// StrToTime converts string to *Time object. It also supports timestamp string.
// The parameter <format> is unnecessary, which specifies the format for converting like "Y-m-d H:i:s".
// If <format> is given, it acts as same as function StrToTimeFormat.
// If <format> is not given, it converts string as a "standard" datetime string.
// Note that, it fails and returns error if there's no date string in <str>.
func StrToTime(str string, format ...string) (*Time, error) {
	if len(format) > 0 {
		return StrToTimeFormat(str, format[0])
	}
	if isTimestampStr(str) {
		timestamp, _ := strconv.ParseInt(str, 10, 64)
		return NewFromTimeStamp(timestamp), nil
	}
	var (
		year, month, day     int
		hour, min, sec, nsec int
		match                []string
		local                = time.Local
	)
	if match = timeRegex1.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
		for k, v := range match {
			match[k] = strings.TrimSpace(v)
		}
		year, month, day = parseDateStr(match[1])
	} else if match = timeRegex2.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
		for k, v := range match {
			match[k] = strings.TrimSpace(v)
		}
		year, month, day = parseDateStr(match[1])
	} else {
		return nil, errors.New("unsupported time format")
	}

	// Time
	if len(match[2]) > 0 {
		s := strings.Replace(match[2], ":", "", -1)
		if len(s) < 6 {
			s += strings.Repeat("0", 6-len(s))
		}
		hour, _ = strconv.Atoi(s[0:2])
		min, _ = strconv.Atoi(s[2:4])
		sec, _ = strconv.Atoi(s[4:6])
	}
	// Nanoseconds, check and perform bit filling
	if len(match[3]) > 0 {
		nsec, _ = strconv.Atoi(match[3])
		for i := 0; i < 9-len(match[3]); i++ {
			nsec *= 10
		}
	}
	// If there's zone information in the string,
	// it then performs time zone conversion, which converts the time zone to UTC.
	if match[4] != "" && match[6] == "" {
		match[6] = "000000"
	}
	// If there's offset in the string, it then firstly processes the offset.
	if match[6] != "" {
		zone := strings.Replace(match[6], ":", "", -1)
		zone = strings.TrimLeft(zone, "+-")
		if len(zone) <= 6 {
			zone += strings.Repeat("0", 6-len(zone))
			h, _ := strconv.Atoi(zone[0:2])
			m, _ := strconv.Atoi(zone[2:4])
			s, _ := strconv.Atoi(zone[4:6])
			if h > 24 || m > 59 || s > 59 {
				return nil, gerror.Newf("invalid zone string: %s", match[6])
			}
			// Comparing the given time zone whether equals to current time zone,
			// it converts it to UTC if they does not equal.
			_, localOffset := time.Now().Zone()
			// Comparing in seconds.
			if (h*3600 + m*60 + s) != localOffset {
				local = time.UTC
				// UTC conversion.
				operation := match[5]
				if operation != "+" && operation != "-" {
					operation = "-"
				}
				switch operation {
				case "+":
					if h > 0 {
						hour -= h
					}
					if m > 0 {
						min -= m
					}
					if s > 0 {
						sec -= s
					}
				case "-":
					if h > 0 {
						hour += h
					}
					if m > 0 {
						min += m
					}
					if s > 0 {
						sec += s
					}
				}
			}
		}
	}
	if year <= 0 {
		return nil, errors.New("invalid time string:" + str)
	}
	return NewFromTime(time.Date(year, time.Month(month), day, hour, min, sec, nsec, local)), nil
}

// ConvertZone converts time in string <strTime> from <fromZone> to <toZone>.
// The parameter <fromZone> is unnecessary, it is current time zone in default.
func ConvertZone(strTime string, toZone string, fromZone ...string) (*Time, error) {
	t, err := StrToTime(strTime)
	if err != nil {
		return nil, err
	}
	if len(fromZone) > 0 {
		if l, err := time.LoadLocation(fromZone[0]); err != nil {
			return nil, err
		} else {
			t.Time = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Time.Second(), t.Time.Nanosecond(), l)
		}
	}
	if l, err := time.LoadLocation(toZone); err != nil {
		return nil, err
	} else {
		return t.ToLocation(l), nil
	}
}

// StrToTimeFormat parses string <str> to *Time object with given format <format>.
// The parameter <format> is like "Y-m-d H:i:s".
func StrToTimeFormat(str string, format string) (*Time, error) {
	return StrToTimeLayout(str, formatToStdLayout(format))
}

// StrToTimeLayout parses string <str> to *Time object with given format <layout>.
// The parameter <layout> is in stdlib format like "2006-01-02 15:04:05".
func StrToTimeLayout(str string, layout string) (*Time, error) {
	if t, err := time.ParseInLocation(layout, str, time.Local); err == nil {
		return NewFromTime(t), nil
	} else {
		return nil, err
	}
}

// ParseTimeFromContent retrieves time information for content string, it then parses and returns it
// as *Time object.
// It returns the first time information if there're more than one time string in the content.
// It only retrieves and parses the time information with given <format> if it's passed.
func ParseTimeFromContent(content string, format ...string) *Time {
	if len(format) > 0 {
		if match, err := gregex.MatchString(formatToRegexPattern(format[0]), content); err == nil && len(match) > 0 {
			return NewFromStrFormat(match[0], format[0])
		}
	} else {
		if match := timeRegex1.FindStringSubmatch(content); len(match) >= 1 {
			return NewFromStr(strings.Trim(match[0], "./_- \n\r"))
		} else if match := timeRegex2.FindStringSubmatch(content); len(match) >= 1 {
			return NewFromStr(strings.Trim(match[0], "./_- \n\r"))
		}
	}
	return nil
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h", "1d" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h", "d".
//
// Very note that it supports unit "d" more than function time.ParseDuration.
func ParseDuration(s string) (time.Duration, error) {
	if utils.IsNumeric(s) {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return time.Duration(v), nil
	}
	match, err := gregex.MatchString(`^([\-\d]+)[dD](.*)$`, s)
	if err != nil {
		return 0, err
	}
	if len(match) == 3 {
		v, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0, err
		}
		return time.ParseDuration(fmt.Sprintf(`%dh%s`, v*24, match[2]))
	}
	return time.ParseDuration(s)
}

// FuncCost calculates the cost time of function <f> in nanoseconds.
func FuncCost(f func()) int64 {
	t := TimestampNano()
	f()
	return TimestampNano() - t
}

// isTimestampStr checks and returns whether given string a timestamp string.
func isTimestampStr(s string) bool {
	length := len(s)
	if length == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

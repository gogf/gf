// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogf/gf/v3/text/gregex"
)

var (
	// Refer: http://php.net/manual/en/function.date.php
	layouts = map[byte]string{
		'd': "02",                        // Day: Day of the month, 2 digits with leading zeros. Eg: 01 to 31.
		'D': "Mon",                       // Day: A textual representation of a day, three letters. Eg: Mon through Sun.
		'w': "Monday",                    // Day: Numeric representation of the day of the week. Eg: 0 (for Sunday) through 6 (for Saturday).
		'N': "Monday",                    // Day: ISO-8601 numeric representation of the day of the week. Eg: 1 (for Monday) through 7 (for Sunday).
		'j': "=j=02",                     // Day: Day of the month without leading zeros. Eg: 1 to 31.
		'S': "02",                        // Day: English ordinal suffix for the day of the month, 2 characters. Eg: st, nd, rd or th. Works well with j.
		'l': "Monday",                    // Day: A full textual representation of the day of the week. Eg: Sunday through Saturday.
		'z': "",                          // Day: The day of the year (starting from 0). Eg: 0 through 365.
		'W': "",                          // Week: ISO-8601 week number of year, weeks starting on Monday. Eg: 42 (the 42nd week in the year).
		'F': "January",                   // Month: A full textual representation of a month, such as January or March. Eg: January through December.
		'm': "01",                        // Month: Numeric representation of a month, with leading zeros. Eg: 01 through 12.
		'M': "Jan",                       // Month: A short textual representation of a month, three letters. Eg: Jan through Dec.
		'n': "1",                         // Month: Numeric representation of a month, without leading zeros. Eg: 1 through 12.
		't': "",                          // Month: Number of days in the given month. Eg: 28 through 31.
		'Y': "2006",                      // Year: A full numeric representation of a year, 4 digits. Eg: 1999 or 2003.
		'y': "06",                        // Year: A two-digit representation of a year. Eg: 99 or 03.
		'a': "pm",                        // Time: Lowercase Ante meridiem and Post meridiem. Eg: am or pm.
		'A': "PM",                        // Time: Uppercase Ante meridiem and Post meridiem. Eg: AM or PM.
		'g': "3",                         // Time: 12-hour layout of an hour without leading zeros. Eg: 1 through 12.
		'G': "=G=15",                     // Time: 24-hour layout of an hour without leading zeros. Eg: 0 through 23.
		'h': "03",                        // Time: 12-hour layout of an hour with leading zeros. Eg: 01 through 12.
		'H': "15",                        // Time: 24-hour layout of an hour with leading zeros. Eg: 00 through 23.
		'i': "04",                        // Time: Minutes with leading zeros. Eg: 00 to 59.
		's': "05",                        // Time: Seconds with leading zeros. Eg: 00 through 59.
		'u': "=u=.000",                   // Time: Milliseconds. Eg: 234, 678.
		'U': "",                          // Time: Seconds since the Unix Epoch (January 1 1970 00:00:00 GMT).
		'O': "-0700",                     // Zone: Difference to Greenwich time (GMT) in hours. Eg: +0200.
		'P': "-07:00",                    // Zone: Difference to Greenwich time (GMT) with colon between hours and minutes. Eg: +02:00.
		'T': "MST",                       // Zone: Timezone abbreviation. Eg: UTC, EST, MDT ...
		'c': "2006-01-02T15:04:05-07:00", // Layout: ISO 8601 date. Eg: 2004-02-12T15:19:21+00:00.
		'r': "Mon, 02 Jan 06 15:04 MST",  // Layout: RFC 2822 layoutted date. Eg: Thu, 21 Dec 2000 16:01:07 +0200.
	}

	// Week to number mapping.
	weekMap = map[string]string{
		"Sunday":    "0",
		"Monday":    "1",
		"Tuesday":   "2",
		"Wednesday": "3",
		"Thursday":  "4",
		"Friday":    "5",
		"Saturday":  "6",
	}

	// Day count of each month which is not in leap year.
	dayOfMonth = []int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}
)

// Layout layouts and returns the layoutted result with custom `layout`.
// Refer method Format if you want to follow stdlib format.
func (t *Time) Layout(layout string) string {
	if t == nil {
		return ""
	}
	runes := []rune(layout)
	buffer := bytes.NewBuffer(nil)
	for i := 0; i < len(runes); {
		switch runes[i] {
		case '\\':
			if i < len(runes)-1 {
				buffer.WriteRune(runes[i+1])
				i += 2
				continue
			} else {
				return buffer.String()
			}
		case 'W':
			buffer.WriteString(strconv.Itoa(t.WeeksOfYear()))
		case 'z':
			buffer.WriteString(strconv.Itoa(t.DayOfYear()))
		case 't':
			buffer.WriteString(strconv.Itoa(t.DaysInMonth()))
		case 'U':
			buffer.WriteString(strconv.FormatInt(t.Unix(), 10))
		default:
			if runes[i] > 255 {
				buffer.WriteRune(runes[i])
				break
			}
			if f, ok := layouts[byte(runes[i])]; ok {
				result := t.Time.Format(f)
				// Particular chars should be handled here.
				switch runes[i] {
				case 'j':
					for _, s := range []string{"=j=0", "=j="} {
						result = strings.ReplaceAll(result, s, "")
					}
					buffer.WriteString(result)
				case 'G':
					for _, s := range []string{"=G=0", "=G="} {
						result = strings.ReplaceAll(result, s, "")
					}
					buffer.WriteString(result)
				case 'u':
					buffer.WriteString(strings.ReplaceAll(result, "=u=.", ""))
				case 'w':
					buffer.WriteString(weekMap[result])
				case 'N':
					buffer.WriteString(strings.ReplaceAll(weekMap[result], "0", "7"))
				case 'S':
					buffer.WriteString(layoutMonthDaySuffixMap(result))
				default:
					buffer.WriteString(result)
				}
			} else {
				buffer.WriteRune(runes[i])
			}
		}
		i++
	}
	return buffer.String()
}

// LayoutNew layouts and returns a new Time object with given custom `layout`.
func (t *Time) LayoutNew(layout string) *Time {
	if t == nil {
		return nil
	}
	return NewFromStr(t.Layout(layout))
}

// LayoutTo layouts `t` with given custom `layout`.
func (t *Time) LayoutTo(layout string) *Time {
	if t == nil {
		return nil
	}
	t.Time = NewFromStr(t.Layout(layout)).Time
	return t
}

// Format layouts the time with stdlib format and returns the layoutted result.
func (t *Time) Format(format string) string {
	if t == nil {
		return ""
	}
	return t.Time.Format(format)
}

// FormatNew layouts the time with stdlib format and returns the new Time object.
func (t *Time) FormatNew(format string) *Time {
	if t == nil {
		return nil
	}
	newTime, err := StrToTimeFormat(t.Format(format), format)
	if err != nil {
		panic(err)
	}
	return newTime
}

// FormatTo layouts `t` with stdlib format.
func (t *Time) FormatTo(format string) *Time {
	if t == nil {
		return nil
	}
	newTime, err := StrToTimeFormat(t.Format(format), format)
	if err != nil {
		panic(err)
	}
	t.Time = newTime.Time
	return t
}

// IsLeapYear checks whether the time is leap year.
func (t *Time) IsLeapYear() bool {
	year := t.Year()
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return true
	}
	return false
}

// DayOfYear checks and returns the position of the day for the year.
func (t *Time) DayOfYear() int {
	var (
		day   = t.Day()
		month = t.Month()
	)
	if t.IsLeapYear() {
		if month > 2 {
			return dayOfMonth[month-1] + day
		}
		return dayOfMonth[month-1] + day - 1
	}
	return dayOfMonth[month-1] + day - 1
}

// DaysInMonth returns the day count of the current month.
func (t *Time) DaysInMonth() int {
	switch t.Month() {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	}
	if t.IsLeapYear() {
		return 29
	}
	return 28
}

// WeeksOfYear returns the point of current week for the year.
func (t *Time) WeeksOfYear() int {
	_, week := t.ISOWeek()
	return week
}

// layoutToStdFormat converts the custom layout to stdlib format.
func layoutToStdFormat(layout string) string {
	b := bytes.NewBuffer(nil)
	for i := 0; i < len(layout); {
		switch layout[i] {
		case '\\':
			if i < len(layout)-1 {
				b.WriteByte(layout[i+1])
				i += 2
				continue
			} else {
				return b.String()
			}

		default:
			if f, ok := layouts[layout[i]]; ok {
				// Handle particular chars.
				switch layout[i] {
				case 'j':
					b.WriteString("2")
				case 'G':
					b.WriteString("15")
				case 'u':
					if i > 0 && layout[i-1] == '.' {
						b.WriteString("000")
					} else {
						b.WriteString(".000")
					}

				default:
					b.WriteString(f)
				}
			} else {
				b.WriteByte(layout[i])
			}
			i++
		}
	}
	return b.String()
}

// layoutToRegexPattern converts the custom layout to its corresponding regular expression.
func layoutToRegexPattern(layout string) string {
	s := regexp.QuoteMeta(layoutToStdFormat(layout))
	s, _ = gregex.ReplaceString(`[0-9]`, `[0-9]`, s)
	s, _ = gregex.ReplaceString(`[A-Za-z]`, `[A-Za-z]`, s)
	s, _ = gregex.ReplaceString(`\s+`, `\s+`, s)
	return s
}

// layoutMonthDaySuffixMap returns the short english word for current day.
func layoutMonthDaySuffixMap(day string) string {
	switch day {
	case "01", "21", "31":
		return "st"
	case "02", "22":
		return "nd"
	case "03", "23":
		return "rd"
	default:
		return "th"
	}
}

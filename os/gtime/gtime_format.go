// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/gogf/gf/text/gregex"
)

var (
	// Please refer: http://php.net/manual/en/function.date.php
	formats = map[byte]string{
		// ================== Day ==================
		'd': "02",     // Day of the month, 2 digits with leading zeros. Eg: 01 to 31
		'D': "Mon",    // A textual representation of a day, three letters. Eg: Mon through Sun
		'w': "Monday", // Numeric representation of the day of the week. Eg: 0 (for Sunday) through 6 (for Saturday)
		'N': "Monday", // ISO-8601 numeric representation of the day of the week. Eg: 1 (for Monday) through 7 (for Sunday)
		'j': "=j=02",  // Day of the month without leading zeros. Eg: 1 to 31
		'S': "02",     // English ordinal suffix for the day of the month, 2 characters. Eg: st, nd, rd or th. Works well with j
		'l': "Monday", // A full textual representation of the day of the week. Eg: Sunday through Saturday
		'z': "",       // The day of the year (starting from 0). Eg: 0 through 365

		// ================== Week ==================
		'W': "", // ISO-8601 week number of year, weeks starting on Monday. Eg: 42 (the 42nd week in the year)

		// ================== Month ==================
		'F': "January", // A full textual representation of a month, such as January or March. Eg: January through December
		'm': "01",      // Numeric representation of a month, with leading zeros. Eg: 01 through 12
		'M': "Jan",     // A short textual representation of a month, three letters. Eg: Jan through Dec
		'n': "1",       // Numeric representation of a month, without leading zeros. Eg: 1 through 12
		't': "",        // Number of days in the given month. Eg: 28 through 31

		// ================== Year ==================
		'Y': "2006", // A full numeric representation of a year, 4 digits. Eg: 1999 or 2003
		'y': "06",   // A two digit representation of a year. Eg: 99 or 03

		// ================== Time ==================
		'a': "pm",      // Lowercase Ante meridiem and Post meridiem. Eg: am or pm
		'A': "PM",      // Uppercase Ante meridiem and Post meridiem. Eg: AM or PM
		'g': "3",       // 12-hour format of an hour without leading zeros. Eg: 1 through 12
		'G': "=G=15",   // 24-hour format of an hour without leading zeros. Eg: 0 through 23
		'h': "03",      // 12-hour format of an hour with leading zeros. Eg: 01 through 12
		'H': "15",      // 24-hour format of an hour with leading zeros. Eg: 00 through 23
		'i': "04",      // Minutes with leading zeros. Eg: 00 to 59
		's': "05",      // Seconds with leading zeros. Eg: 00 through 59
		'u': "=u=.000", // Milliseconds. Eg: 234, 678
		'U': "",        // Seconds since the Unix Epoch (January 1 1970 00:00:00 GMT).

		// ================== Zone ==================
		'O': "-0700",  // Difference to Greenwich time (GMT) in hours. Eg: +0200
		'P': "-07:00", // Difference to Greenwich time (GMT) with colon between hours and minutes. Eg: +02:00
		'T': "MST",    // Timezone abbreviation. Eg: UTC, EST, MDT ...

		// ================== Format ==================
		'c': "2006-01-02T15:04:05-07:00", // ISO 8601 date. Eg: 2004-02-12T15:19:21+00:00
		'r': "Mon, 02 Jan 06 15:04 MST",  // RFC 2822 formatted date. Eg: Thu, 21 Dec 2000 16:01:07 +0200
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

// Format formats and returns the formatted result with custom <format>.
func (t *Time) Format(format string) string {
	runes := []rune(format)
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
			if f, ok := formats[byte(runes[i])]; ok {
				result := t.Time.Format(f)
				// Particular chars should be handled here.
				switch runes[i] {
				case 'j':
					for _, s := range []string{"=j=0", "=j="} {
						result = strings.Replace(result, s, "", -1)
					}
					buffer.WriteString(result)
				case 'G':
					for _, s := range []string{"=G=0", "=G="} {
						result = strings.Replace(result, s, "", -1)
					}
					buffer.WriteString(result)
				case 'u':
					buffer.WriteString(strings.Replace(result, "=u=.", "", -1))
				case 'w':
					buffer.WriteString(weekMap[result])
				case 'N':
					buffer.WriteString(strings.Replace(weekMap[result], "0", "7", -1))
				case 'S':
					buffer.WriteString(formatMonthDaySuffixMap(result))
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

// FormatTo formats and returns a new Time object with given custom <format>.
func (t *Time) FormatTo(format string) *Time {
	t.Time = NewFromStr(t.Format(format)).Time
	return t
}

// Layout formats the time with stdlib layout and returns the formatted result.
func (t *Time) Layout(layout string) string {
	return t.Time.Format(layout)
}

// Layout formats the time with stdlib layout and returns the new Time object.
func (t *Time) LayoutTo(layout string) *Time {
	t.Time = NewFromStr(t.Layout(layout)).Time
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
	day := t.Day()
	month := int(t.Month())
	if t.IsLeapYear() {
		if month > 2 {
			return dayOfMonth[month-1] + day
		}
		return dayOfMonth[month-1] + day - 1
	}
	return dayOfMonth[month-1] + day - 1
}

// DaysInMonth returns the day count of current month.
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

// formatToStdLayout converts custom format to stdlib layout.
func formatToStdLayout(format string) string {
	b := bytes.NewBuffer(nil)
	for i := 0; i < len(format); {
		switch format[i] {
		case '\\':
			if i < len(format)-1 {
				b.WriteByte(format[i+1])
				i += 2
				continue
			} else {
				return b.String()
			}

		default:
			if f, ok := formats[format[i]]; ok {
				// Handle particular chars.
				switch format[i] {
				case 'j':
					b.WriteString("2")
				case 'G':
					b.WriteString("15")
				case 'u':
					if i > 0 && format[i-1] == '.' {
						b.WriteString("000")
					} else {
						b.WriteString(".000")
					}

				default:
					b.WriteString(f)
				}
			} else {
				b.WriteByte(format[i])
			}
			i++
		}
	}
	return b.String()
}

// formatToRegexPattern converts the custom format to its corresponding regular expression.
func formatToRegexPattern(format string) string {
	s := gregex.Quote(formatToStdLayout(format))
	s, _ = gregex.ReplaceString(`[0-9]`, `[0-9]`, s)
	s, _ = gregex.ReplaceString(`[A-Za-z]`, `[A-Za-z]`, s)
	return s
}

// formatMonthDaySuffixMap returns the short english word for current day.
func formatMonthDaySuffixMap(day string) string {
	switch day {
	case "01":
		return "st"
	case "02":
		return "nd"
	case "03":
		return "rd"
	default:
		return "th"
	}
}

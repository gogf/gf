// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime_test

import (
    "fmt"
    "github.com/gogf/gf/v2/os/gtime"
)

// New creates and returns a Time object with given parameter.
// The optional parameter can be type of: time.Time/*time.Time, string or integer.
func ExampleSetTimeZone() {
    gtime.SetTimeZone("Asia/Shanghai")
    fmt.Println(gtime.Datetime())

    gtime.SetTimeZone("Asia/Tokyo")
    fmt.Println(gtime.Datetime())
    // May Output:
    // 2018-08-08 08:08:08
    // 2018-08-08 09:08:08
}

func ExampleTimestamp() {
    fmt.Println(gtime.Timestamp())

    // May Output:
    // 1636359252
}

func ExampleTimestampMilli() {
    fmt.Println(gtime.TimestampMilli())

    // May Output:
    // 1636359252000
}

func ExampleTimestampMicro() {
    fmt.Println(gtime.TimestampMicro())

    // May Output:
    // 1636359252000000
}

func ExampleTimestampNano() {
    fmt.Println(gtime.TimestampNano())

    // May Output:
    // 1636359252000000000
}

func ExampleTimestampStr() {
    fmt.Println(gtime.TimestampStr())

    // May Output:
    // 1636359252
}

func ExampleDate() {
    fmt.Println(gtime.Date())

    // May Output:
    // 2006-01-02
}

func ExampleDatetime() {
    fmt.Println(gtime.Datetime())

    // May Output:
    // 2006-01-02 15:04:05
}

func ExampleISO8601() {
    fmt.Println(gtime.ISO8601())

    // May Output:
    // 2006-01-02T15:04:05-07:00
}

func ExampleRFC822() {
    fmt.Println(gtime.RFC822())

    // May Output:
    // Mon, 02 Jan 06 15:04 MST
}

func ExampleStrToTime() {
    res, _ := gtime.StrToTime("2006-01-02T15:04:05-07:00", "Y-m-d H:i:s")
    fmt.Println(res)

    // May Output:
    // 2006-01-02 15:04:05
}

func ExampleConvertZone() {
    res, _ := gtime.ConvertZone("2006-01-02 15:04:05", "Asia/Tokyo", "Asia/Shanghai")
    fmt.Println(res)

    // Output:
    // 2006-01-02 16:04:05
}

func ExampleStrToTimeFormat() {
    res, _ := gtime.StrToTimeFormat("2006-01-02 15:04:05", "Y-m-d H:i:s")
    fmt.Println(res)

    // Output:
    // 2006-01-02 15:04:05
}

func ExampleStrToTimeLayout() {
    res, _ := gtime.StrToTimeLayout("2018-08-08", "2006-01-02")
    fmt.Println(res)

    // Output:
    // 2018-08-08 00:00:00
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h", "1d" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h", "d".
//
// Very note that it supports unit "d" more than function time.ParseDuration.
func ExampleParseDuration() {
    res, _ := gtime.ParseDuration("+10h")
    fmt.Println(res)

    // Output:
    // 10h0m0s
}

func ExampleTime_Format() {
    gt1 := gtime.New("2018-08-08 08:08:08")

    fmt.Println(gt1.Format("Y-m-d"))
    fmt.Println(gt1.Format("l"))
    fmt.Println(gt1.Format("F j, Y, g:i a"))
    fmt.Println(gt1.Format("j, n, Y"))
    fmt.Println(gt1.Format("h-i-s, j-m-y, it is w Day z"))
    fmt.Println(gt1.Format("D M j G:i:s T Y"))

    // Output:
    // 2018-08-08
    // Wednesday
    // August 8, 2018, 8:08 am
    // 8, 8, 2018
    // 08-08-08, 8-08-18, 0831 0808 3 Wedam18 219
    // Wed Aug 8 8:08:08 CST 2018
}

func ExampleTime_FormatNew() {
    gt1 := gtime.New("2018-08-08 08:08:08")

    fmt.Println(gt1.FormatNew("Y-m-d"))
    fmt.Println(gt1.FormatNew("Y-m-d H:i"))

    // Output:
    // 2018-08-08 00:00:00
    // 2018-08-08 08:08:00
}

func ExampleTime_FormatTo() {
    gt1 := gtime.New("2018-08-08 08:08:08")

    fmt.Println(gt1.FormatTo("Y-m-d"))

    // Output:
    // 2018-08-08 00:00:00
}

func ExampleTime_Layout() {
    gt1 := gtime.New("2018-08-08 08:08:08")

    fmt.Println(gt1.Layout("2006-01-02"))

    // Output:
    // 2018-08-08
}

func ExampleTime_LayoutNew() {
    gt1 := gtime.New("2018-08-08 08:08:08")

    fmt.Println(gt1.LayoutNew("2006-01-02"))

    // Output:
    // 2018-08-08 00:00:00
}

func ExampleTime_LayoutTo() {
    gt1 := gtime.New("2018-08-08 08:08:08")

    fmt.Println(gt1.LayoutTo("2006-01-02"))

    // Output:
    // 2018-08-08 00:00:00
}

func ExampleTime_IsLeapYear() {
    gt1 := gtime.New("2018-08-08 08:08:08")

    fmt.Println(gt1.IsLeapYear())

    // Output:
    // false
}

func ExampleTime_DayOfYear() {
    gt1 := gtime.New("2018-01-08 08:08:08")

    fmt.Println(gt1.DayOfYear())

    // Output:
    // 7
}

// DaysInMonth returns the day count of current month.
func ExampleTime_DaysInMonth() {
    gt1 := gtime.New("2018-08-08 08:08:08")

    fmt.Println(gt1.DaysInMonth())

    // Output:
    // 31
}

// WeeksOfYear returns the point of current week for the year.
func ExampleTime_WeeksOfYear() {
    gt1 := gtime.New("2018-01-08 08:08:08")

    fmt.Println(gt1.WeeksOfYear())

    // Output:
    // 2
}

func ExampleTime_ToZone() {
    gt1 := gtime.Now()
    gt2, _ := gt1.ToZone("Asia/Shanghai")
    gt3, _ := gt1.ToZone("Asia/Tokyo")

    fmt.Println(gt2)
    fmt.Println(gt3)

    // May Output:
    // 2021-11-11 17:10:10
    // 2021-11-11 18:10:10
}


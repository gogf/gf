// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime_test

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"reflect"
	"time"
)

// New creates and returns a Time object with given parameter.
// The optional parameter can be type of: time.Time/*time.Time, string or integer.
func ExampleNew() {
	timeLayout := "2006-01-02 15:04:05"
	curTime := "2018-08-08 08:08:08"
	timer, _ := time.Parse(timeLayout, curTime)
	t1 := gtime.New(timer)
	t2 := gtime.New(&timer)
	t3 := gtime.New(curTime, "Y-m-d H:i:s")
	t4 := gtime.New(curTime)
	t5 := gtime.New(1533686888)

	fmt.Println(t1)
	fmt.Println(t2)
	fmt.Println(t3)
	fmt.Println(t4)
	fmt.Println(t5)

	// Output:
	// 2018-08-08 08:08:08
	// 2018-08-08 08:08:08
	// 2018-08-08 08:08:08
	// 2018-08-08 08:08:08
	// 2018-08-08 08:08:08
}

// Now creates and returns a time object of now.
func ExampleNow() {
	t := gtime.Now()
	fmt.Println(t)

	// May Output:
	// 2021-11-06 13:41:08
}

// NewFromTime creates and returns a Time object with given time.Time object.
func ExampleNewFromTime() {
	timeLayout := "2006-01-02 15:04:05"
	curTime := "2018-08-08 08:08:08"
	timer, _ := time.Parse(timeLayout, curTime)
	nTime := gtime.NewFromTime(timer)

	fmt.Println(nTime)

	// Output:
	// 2018-08-08 08:08:08
}

// NewFromStr creates and returns a Time object with given string.
// Note that it returns nil if there's error occurs.
func ExampleNewFromStr() {
	t := gtime.NewFromStr("2018-08-08 08:08:08")

	fmt.Println(t)

	// Output:
	// 2018-08-08 08:08:08
}

// NewFromStrFormat creates and returns a Time object with given string and
// custom format like: Y-m-d H:i:s.
// Note that it returns nil if there's error occurs.
func ExampleNewFromStrFormat() {
	t := gtime.NewFromStrFormat("2018-08-08 08:08:08", "Y-m-d H:i:s")
	fmt.Println(t)

	// Output:
	// 2018-08-08 08:08:08
}

// NewFromStrLayout creates and returns a Time object with given string and
// stdlib layout like: 2006-01-02 15:04:05.
// Note that it returns nil if there's error occurs.
func ExampleNewFromStrLayout() {
	t := gtime.NewFromStrLayout("2018-08-08 08:08:08", "2006-01-02 15:04:05")
	fmt.Println(t)

	// Output:
	// 2018-08-08 08:08:08
}

// NewFromTimeStamp creates and returns a Time object with given timestamp,
// which can be in seconds to nanoseconds.
// Eg: 1600443866 and 1600443866199266000 are both considered as valid timestamp number.
func ExampleNewFromTimeStamp() {
	t1 := gtime.NewFromTimeStamp(1533686888)
	t2 := gtime.NewFromTimeStamp(1533686888000)

	fmt.Println(t1.String() == t2.String())
	fmt.Println(t1)

	// Output:
	// true
	// 2018-08-08 08:08:08
}

// Timestamp returns the timestamp in seconds.
func ExampleTime_Timestamp() {
	t := gtime.Timestamp()

	fmt.Println(t)

	// May output:
	// 1533686888
}

// Timestamp returns the timestamp in milliseconds.
func ExampleTime_TimestampMilli() {
	t := gtime.TimestampMilli()

	fmt.Println(t)

	// May output:
	// 1533686888000
}

// Timestamp returns the timestamp in microseconds.
func ExampleTime_TimestampMicro() {
	t := gtime.TimestampMicro()

	fmt.Println(t)

	// May output:
	// 1533686888000000
}

// Timestamp returns the timestamp in nanoseconds.
func ExampleTime_TimestampNano() {
	t := gtime.TimestampNano()

	fmt.Println(t)

	// May output:
	// 1533686888000000
}

// TimestampStr is a convenience method which retrieves and returns
// the timestamp in seconds as string.
func ExampleTime_TimestampStr() {
	t := gtime.TimestampStr()

	fmt.Println(reflect.TypeOf(t))

	// Output:
	// string
}

// Month returns the month of the year specified by t.
func ExampleTime_Month() {
	timeLayout := "2006-01-02 15:04:05"
	curTime := "2018-08-08 08:08:08"
	timer, _ := time.Parse(timeLayout, curTime)
	gt := gtime.New(timer)

	t1 := gt.Month()

	fmt.Println(t1)

	// Output:
	// 8
}

// Second returns the second offset within the minute specified by t,
// in the range [0, 59].
func ExampleTime_Second() {
	timeLayout := "2006-01-02 15:04:05"
	curTime := "2018-08-08 08:08:08"
	timer, _ := time.Parse(timeLayout, curTime)
	gt := gtime.New(timer)

	t1 := gt.Second()

	fmt.Println(t1)

	// Output:
	// 8
}

// String returns current time object as string.
func ExampleTime_String() {
	timeLayout := "2006-01-02 15:04:05"
	curTime := "2018-08-08 08:08:08"
	timer, _ := time.Parse(timeLayout, curTime)
	gt := gtime.New(timer)

	t1 := gt.String()

	fmt.Println(t1)
	fmt.Println(reflect.TypeOf(t1))

	// Output:
	// 2018-08-08 08:08:08
	// string
}

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func ExampleTime_IsZero() {
	timeLayout := "2006-01-02 15:04:05"
	curTime := "2018-08-08 08:08:08"
	timer, _ := time.Parse(timeLayout, curTime)
	gt := gtime.New(timer)

	fmt.Println(gt.IsZero())

	// Output:
	// false
}

// Add adds the duration to current time.
func ExampleTime_Add() {
	timeLayout := "2006-01-02 15:04:05"
	curTime := "2018-08-08 08:08:08"
	timer, _ := time.Parse(timeLayout, curTime)
	gt := gtime.New(timer)
	gt1 := gt.Add(time.Duration(10) * time.Second)

	fmt.Println(gt1)

	// Output:
	// 2018-08-08 08:08:18
}

// AddStr parses the given duration as string and adds it to current time.
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ExampleTime_AddStr() {
	timeLayout := "2006-01-02 15:04:05"
	curTime := "2018-08-08 08:08:08"
	timer, _ := time.Parse(timeLayout, curTime)
	gt := gtime.New(timer)
	gt1, _ := gt.AddStr("10s")

	fmt.Println(gt1)

	// Output:
	// 2018-08-08 08:08:18
}

// AddDate adds year, month and day to the time.
func ExampleTime_AddDate() {
	timeLayout := "2006-01-02"
	curTime := "2018-08-08"
	timer, _ := time.Parse(timeLayout, curTime)
	gt := gtime.New(timer)
	gt = gt.AddDate(1, 2, 3)

	fmt.Println(gt)

	// Output:
	// 2019-10-11 00:00:00
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Round operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Round(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func ExampleTime_Round() {
	timeLayout := "2006-01-02"
	curTime := "2018-08-08"
	timer, _ := time.Parse(timeLayout, curTime)
	gt := gtime.New(timer)
	t := gt.Round(time.Duration(10) * time.Second)

	fmt.Println(t)

	// Output:
	// 2018-08-08 00:00:00
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func ExampleTime_Truncate() {
	timeLayout := "2006-01-02 15:04:05"
	curTime := "2018-08-08 08:08:08"
	timer, _ := time.Parse(timeLayout, curTime)
	gt := gtime.New(timer)
	t := gt.Truncate(time.Duration(10) * time.Second)

	fmt.Println(t)

	// Output:
	// 2018-08-08 08:08:00
}

// Equal reports whether t and u represent the same time instant.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
// See the documentation on the Time type for the pitfalls of using == with
// Time values; most code should use Equal instead.
func ExampleTime_Equal() {
	timeLayout := "2006-01-02 15:04:05"
	curTime := "2018-08-08 08:08:08"
	timer, _ := time.Parse(timeLayout, curTime)
	gt1 := gtime.New(timer)
	gt2 := gtime.New(timer)

	fmt.Println(gt1.Equal(gt2))

	// Output:
	// true
}

// Before reports whether the time instant t is before u.
func ExampleTime_Before() {
	timeLayout := "2006-01-02"
	timer1, _ := time.Parse(timeLayout, "2018-08-07")
	timer2, _ := time.Parse(timeLayout, "2018-08-08")
	gt1 := gtime.New(timer1)
	gt2 := gtime.New(timer2)

	fmt.Println(gt1.Before(gt2))

	// Output:
	// true
}

// After reports whether the time instant t is after u.
func ExampleTime_After() {
	timeLayout := "2006-01-02"
	timer1, _ := time.Parse(timeLayout, "2018-08-07")
	timer2, _ := time.Parse(timeLayout, "2018-08-08")
	gt1 := gtime.New(timer1)
	gt2 := gtime.New(timer2)

	fmt.Println(gt1.After(gt2))

	// Output:
	// false
}

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned.
// To compute t-d for a duration d, use t.Add(-d).
func ExampleTime_Sub() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	timer2, _ := time.Parse(timeLayout, "2018-08-08 08:08:10")
	gt1 := gtime.New(timer1)
	gt2 := gtime.New(timer2)

	fmt.Println(gt2.Sub(gt1))

	// Output:
	// 2s
}

// StartOfMinute clones and returns a new time of which the seconds is set to 0.
func ExampleTime_StartOfMinute() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.StartOfMinute())

	// Output:
	// 2018-08-08 08:08:00
}

func ExampleTime_StartOfHour() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.StartOfHour())

	// Output:
	// 2018-08-08 08:00:00
}

func ExampleTime_StartOfDay() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.StartOfDay())

	// Output:
	// 2018-08-08 00:00:00
}

func ExampleTime_StartOfWeek() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.StartOfWeek())

	// Output:
	// 2018-08-05 00:00:00
}

func ExampleTime_StartOfQuarter() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.StartOfQuarter())

	// Output:
	// 2018-07-01 00:00:00
}

func ExampleTime_StartOfHalf() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.StartOfHalf())

	// Output:
	// 2018-07-01 00:00:00
}

func ExampleTime_StartOfYear() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.StartOfYear())

	// Output:
	// 2018-01-01 00:00:00
}

func ExampleTime_EndOfMinute() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.EndOfMinute())

	// Output:
	// 2018-08-08 08:08:59
}

func ExampleTime_EndOfHour() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.EndOfHour())

	// Output:
	// 2018-08-08 08:59:59
}

func ExampleTime_EndOfDay() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.EndOfDay())

	// Output:
	// 2018-08-08 23:59:59
}

func ExampleTime_EndOfWeek() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.EndOfWeek())

	// Output:
	// 2018-08-11 23:59:59
}

func ExampleTime_EndOfMonth() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.EndOfMonth())

	// Output:
	// 2018-08-31 23:59:59
}

func ExampleTime_EndOfQuarter() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.EndOfQuarter())

	// Output:
	// 2018-09-30 23:59:59
}

func ExampleTime_EndOfHalf() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.EndOfHalf())

	// Output:
	// 2018-12-31 23:59:59
}

func ExampleTime_EndOfYear() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	fmt.Println(gt1.EndOfYear())

	// Output:
	// 2018-12-31 23:59:59
}

func ExampleTime_MarshalJSON() {
	timeLayout := "2006-01-02 15:04:05"
	timer1, _ := time.Parse(timeLayout, "2018-08-08 08:08:08")
	gt1 := gtime.New(timer1)

	json, _ := gt1.MarshalJSON()
	fmt.Println(string(json))

	// Output:
	// "2018-08-08 08:08:08"
}

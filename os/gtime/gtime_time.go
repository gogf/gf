// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime

import (
	"bytes"
	"strconv"
	"time"
)

// Time is a wrapper for time.Time for additional features.
type Time struct {
	wrapper
}

// apiUnixNano is an interface definition commonly for custom time.Time wrapper.
type apiUnixNano interface {
	UnixNano() int64
}

// New creates and returns a Time object with given parameter.
// The optional parameter can be type of: time.Time/*time.Time, string or integer.
func New(param ...interface{}) *Time {
	if len(param) > 0 {
		switch r := param[0].(type) {
		case time.Time:
			return NewFromTime(r)
		case *time.Time:
			return NewFromTime(*r)
		case Time:
			return &r
		case *Time:
			return r
		case string:
			return NewFromStr(r)
		case []byte:
			return NewFromStr(string(r))
		case int:
			return NewFromTimeStamp(int64(r))
		case int64:
			return NewFromTimeStamp(r)
		default:
			if v, ok := r.(apiUnixNano); ok {
				return NewFromTimeStamp(v.UnixNano())
			}
		}
	}
	return &Time{
		wrapper{time.Time{}},
	}
}

// Now creates and returns a time object of now.
func Now() *Time {
	return &Time{
		wrapper{time.Now()},
	}
}

// NewFromTime creates and returns a Time object with given time.Time object.
func NewFromTime(t time.Time) *Time {
	return &Time{
		wrapper{t},
	}
}

// NewFromStr creates and returns a Time object with given string.
// Note that it returns nil if there's error occurs.
func NewFromStr(str string) *Time {
	if t, err := StrToTime(str); err == nil {
		return t
	}
	return nil
}

// NewFromStrFormat creates and returns a Time object with given string and
// custom format like: Y-m-d H:i:s.
// Note that it returns nil if there's error occurs.
func NewFromStrFormat(str string, format string) *Time {
	if t, err := StrToTimeFormat(str, format); err == nil {
		return t
	}
	return nil
}

// NewFromStrLayout creates and returns a Time object with given string and
// stdlib layout like: 2006-01-02 15:04:05.
// Note that it returns nil if there's error occurs.
func NewFromStrLayout(str string, layout string) *Time {
	if t, err := StrToTimeLayout(str, layout); err == nil {
		return t
	}
	return nil
}

// NewFromTimeStamp creates and returns a Time object with given timestamp,
// which can be in seconds to nanoseconds.
// Eg: 1600443866 and 1600443866199266000 are both considered as valid timestamp number.
func NewFromTimeStamp(timestamp int64) *Time {
	if timestamp == 0 {
		return &Time{}
	}
	var sec, nano int64
	if timestamp > 1e9 {
		for timestamp < 1e18 {
			timestamp *= 10
		}
		sec = timestamp / 1e9
		nano = timestamp % 1e9
	} else {
		sec = timestamp
	}
	return &Time{
		wrapper{time.Unix(sec, nano)},
	}
}

// Timestamp returns the timestamp in seconds.
func (t *Time) Timestamp() int64 {
	return t.UnixNano() / 1e9
}

// TimestampMilli returns the timestamp in milliseconds.
func (t *Time) TimestampMilli() int64 {
	return t.UnixNano() / 1e6
}

// TimestampMicro returns the timestamp in microseconds.
func (t *Time) TimestampMicro() int64 {
	return t.UnixNano() / 1e3
}

// TimestampNano returns the timestamp in nanoseconds.
func (t *Time) TimestampNano() int64 {
	return t.UnixNano()
}

// TimestampStr is a convenience method which retrieves and returns
// the timestamp in seconds as string.
func (t *Time) TimestampStr() string {
	return strconv.FormatInt(t.Timestamp(), 10)
}

// TimestampMilliStr is a convenience method which retrieves and returns
// the timestamp in milliseconds as string.
func (t *Time) TimestampMilliStr() string {
	return strconv.FormatInt(t.TimestampMilli(), 10)
}

// TimestampMicroStr is a convenience method which retrieves and returns
// the timestamp in microseconds as string.
func (t *Time) TimestampMicroStr() string {
	return strconv.FormatInt(t.TimestampMicro(), 10)
}

// TimestampNanoStr is a convenience method which retrieves and returns
// the timestamp in nanoseconds as string.
func (t *Time) TimestampNanoStr() string {
	return strconv.FormatInt(t.TimestampNano(), 10)
}

// Second returns the second offset within the minute specified by t,
// in the range [0, 59].
func (t *Time) Second() int {
	return t.Time.Second()
}

// Millisecond returns the millisecond offset within the second specified by t,
// in the range [0, 999].
func (t *Time) Millisecond() int {
	return t.Time.Nanosecond() / 1e6
}

// Microsecond returns the microsecond offset within the second specified by t,
// in the range [0, 999999].
func (t *Time) Microsecond() int {
	return t.Time.Nanosecond() / 1e3
}

// Nanosecond returns the nanosecond offset within the second specified by t,
// in the range [0, 999999999].
func (t *Time) Nanosecond() int {
	return t.Time.Nanosecond()
}

// String returns current time object as string.
func (t *Time) String() string {
	if t == nil {
		return ""
	}
	if t.IsZero() {
		return ""
	}
	return t.Format("Y-m-d H:i:s")
}

// Clone returns a new Time object which is a clone of current time object.
func (t *Time) Clone() *Time {
	return New(t.Time)
}

// Add adds the duration to current time.
func (t *Time) Add(d time.Duration) *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.Add(d)
	return newTime
}

// AddStr parses the given duration as string and adds it to current time.
func (t *Time) AddStr(duration string) (*Time, error) {
	if d, err := time.ParseDuration(duration); err != nil {
		return nil, err
	} else {
		return t.Add(d), nil
	}
}

// ToLocation converts current time to specified location.
func (t *Time) ToLocation(location *time.Location) *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.In(location)
	return newTime
}

// ToZone converts current time to specified zone like: Asia/Shanghai.
func (t *Time) ToZone(zone string) (*Time, error) {
	if l, err := time.LoadLocation(zone); err == nil {
		return t.ToLocation(l), nil
	} else {
		return nil, err
	}
}

// UTC converts current time to UTC timezone.
func (t *Time) UTC() *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.UTC()
	return newTime
}

// ISO8601 formats the time as ISO8601 and returns it as string.
func (t *Time) ISO8601() string {
	return t.Layout("2006-01-02T15:04:05-07:00")
}

// RFC822 formats the time as RFC822 and returns it as string.
func (t *Time) RFC822() string {
	return t.Layout("Mon, 02 Jan 06 15:04 MST")
}

// Local converts the time to local timezone.
func (t *Time) Local() *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.Local()
	return newTime
}

// AddDate adds year, month and day to the time.
func (t *Time) AddDate(years int, months int, days int) *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.AddDate(years, months, days)
	return newTime
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Round operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Round(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t *Time) Round(d time.Duration) *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.Round(d)
	return newTime
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t *Time) Truncate(d time.Duration) *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.Truncate(d)
	return newTime
}

// Equal reports whether t and u represent the same time instant.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
// See the documentation on the Time type for the pitfalls of using == with
// Time values; most code should use Equal instead.
func (t *Time) Equal(u *Time) bool {
	return t.Time.Equal(u.Time)
}

// Before reports whether the time instant t is before u.
func (t *Time) Before(u *Time) bool {
	return t.Time.Before(u.Time)
}

// After reports whether the time instant t is after u.
func (t *Time) After(u *Time) bool {
	return t.Time.After(u.Time)
}

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned.
// To compute t-d for a duration d, use t.Add(-d).
func (t *Time) Sub(u *Time) time.Duration {
	return t.Time.Sub(u.Time)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		t.Time = time.Time{}
		return nil
	}
	newTime, err := StrToTime(string(bytes.Trim(b, `"`)))
	if err != nil {
		return err
	}
	t.Time = newTime.Time
	return nil
}

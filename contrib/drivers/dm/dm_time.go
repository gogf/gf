// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

const (
	// minTimeZoneOption and maxTimeZoneOption are the offsets in minutes that the underlying driver
	// accepts for its timeZone option.
	minTimeZoneOption = -779
	maxTimeZoneOption = 840
)

// zonelessTimeFieldTypes are the field types of the time columns without a time zone, which keep
// the wall clock of a time.
var zonelessTimeFieldTypes = map[string]bool{
	"DATE":      true,
	"TIME":      true,
	"DATETIME":  true,
	"TIMESTAMP": true,
}

// zonedTimeFieldTypes are the field types of the time columns with a time zone, which keep the
// instant of a time, in the offset of the session that writes it.
var zonedTimeFieldTypes = map[string]bool{
	"TIME WITH TIME ZONE":      true,
	"DATETIME WITH TIME ZONE":  true,
	"TIMESTAMP WITH TIME ZONE": true,
}

// localTimeFieldTypes are the field types of the time columns with the local time zone. The
// underlying driver converts a time bound for them to the time zone of the session and then again
// to the time zone of the database, so that a time is bound as the text of its wall clock in the
// time zone of the session, which the server reads as a time of the session.
var localTimeFieldTypes = map[string]bool{
	"TIMESTAMP WITH LOCAL TIME ZONE": true,
}

// locations caches the time zones loaded by name.
var locations sync.Map

// loadLocation returns the time zone `name`, which it loads once.
func loadLocation(name string) (*time.Location, error) {
	if location, ok := locations.Load(name); ok {
		return location.(*time.Location), nil
	}
	location, err := time.LoadLocation(name)
	if err != nil {
		return nil, err
	}
	locations.Store(name, location)
	return location, nil
}

// timeZoneOption returns the offset in minutes that the timeZone option of the data source name
// options `extra` gives, and false if it gives none that the underlying driver accepts. The
// underlying driver reads the last option of the name in any letter case.
func timeZoneOption(extra string) (minutes int, ok bool) {
	if _, query, found := strings.Cut(extra, "?"); found {
		extra = query
	}
	for _, option := range strings.Split(extra, "&") {
		key, value, found := strings.Cut(option, "=")
		if !found || !strings.EqualFold(key, "timeZone") {
			continue
		}
		n, err := strconv.Atoi(value)
		ok = err == nil && n >= minTimeZoneOption && n <= maxTimeZoneOption
		minutes = n
	}
	if !ok {
		minutes = 0
	}
	return
}

// offsetMinutes returns the offset of the time zone `location` at `t` in whole minutes, which is
// the time zone of a session that the underlying driver opens at `t` for it.
func offsetMinutes(location *time.Location, t time.Time) int {
	_, offset := t.In(location).Zone()
	return offset / 60
}

// wallClockLocation returns the time zone of the wall clocks that the columns without a time zone
// keep: the configured Timezone, or else the local time zone. An invalid Timezone fails to open
// the database, see buildDataSourceName.
func (d *Driver) wallClockLocation() *time.Location {
	if config := d.GetConfig(); config != nil && config.Timezone != "" {
		if location, err := loadLocation(config.Timezone); err == nil {
			return location
		}
	}
	return time.Local
}

// relabelWallClock returns `t`, a time in the fixed time zone of a session, as the same wall clock
// in the time zone `location`, or `t` as it is if `location` has no instant of that wall clock, as
// within the hour that a daylight saving time skips.
func relabelWallClock(t time.Time, location *time.Location) time.Time {
	relabelled := time.Date(
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), location,
	)
	if y, mo, d := relabelled.Date(); y != t.Year() || mo != t.Month() || d != t.Day() {
		return t
	}
	if h, mi, s := relabelled.Clock(); h != t.Hour() || mi != t.Minute() || s != t.Second() {
		return t
	}
	return relabelled
}

// timeOfValue returns the time that `value` is, or false if it is not a time or is a zero time.
func timeOfValue(value any) (time.Time, bool) {
	var t time.Time
	switch v := value.(type) {
	case time.Time:
		t = v
	case *time.Time:
		if v != nil {
			t = *v
		}
	case gtime.Time:
		t = v.Time
	case *gtime.Time:
		if v != nil {
			t = v.Time
		}
	case sql.NullTime:
		if v.Valid {
			t = v.Time
		}
	case *sql.NullTime:
		if v != nil && v.Valid {
			t = v.Time
		}
	}
	return t, !t.IsZero()
}

// instantTime is a time written to a column with a time zone, which a connection binds as its
// instant, see sessionTimes.convert. It does not embed the time, so that the core passes it to the
// connection as it is instead of as the text of its String method.
type instantTime struct {
	t time.Time
}

// sessionTextTime is a time written to a column with the local time zone, which a connection binds
// as the text of its wall clock in the time zone of its session, see sessionTimes.convert.
type sessionTextTime struct {
	t time.Time
}

// convertTimeForField returns the value for `t`, a time written to a column of the type
// `fieldType`, and false for a column of another type than a time or a character type: the text
// of `t` for a character column, `t` for a column without a time zone, and `t` marked as an
// instantTime or a sessionTextTime for a column with a time zone or with the local time zone, for
// the connection to bind, see sessionTimes.
func convertTimeForField(fieldType string, t time.Time) (any, bool) {
	fieldType = strings.ToUpper(fieldType)
	switch {
	case characterFieldTypes[fieldType]:
		return gtime.NewFromTime(t).String(), true
	case zonelessTimeFieldTypes[fieldType]:
		return t, true
	case zonedTimeFieldTypes[fieldType]:
		return instantTime{t: t}, true
	case localTimeFieldTypes[fieldType]:
		return sessionTextTime{t: t}, true
	}
	return nil, false
}

// sessionTimes converts the time arguments of the statements of a connection for the time zone of
// its session, which the underlying driver converts a time from the offset of its time zone to, in
// whole minutes. A connection opens its session with the offset of the time zone of the wall
// clocks when it opens, see connector.
type sessionTimes struct {
	// wallClock is the time zone of the wall clocks that the columns without a time zone keep.
	wallClock *time.Location
	// session is the fixed time zone of the session.
	session *time.Location
}

// convert converts the value of the argument `nv`, see CheckNamedValue:
//   - an instantTime is passed as its instant, with its offset in whole minutes, which the
//     underlying driver binds;
//   - a sessionTextTime is passed as the text of its wall clock in the time zone of the session;
//   - any other time is passed as its wall clock in the time zone of the wall clocks, in the time
//     zone of the session, so that a time of another offset than the session, like a local mean
//     time before 1901 or a daylight saving time, is written to a column without a time zone as its
//     own wall clock. A time compared with or written to a column with a time zone by a raw argument
//     is so passed as its wall clock too, as the type of the column is not known.
//
// It returns false for a value that is not a time.
func (s sessionTimes) convert(nv *driver.NamedValue) bool {
	switch v := nv.Value.(type) {
	case instantTime:
		_, offset := v.t.Zone()
		nv.Value = v.t.In(time.FixedZone("", offset/60*60))
	case sessionTextTime:
		nv.Value = v.t.In(s.session).Format("2006-01-02 15:04:05.999999999")
	default:
		t, ok := timeOfValue(nv.Value)
		if !ok {
			return false
		}
		t = t.In(s.wallClock)
		nv.Value = time.Date(
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), s.session,
		)
	}
	return true
}

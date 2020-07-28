// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"time"

	"github.com/gogf/gf/internal/utils"
	"github.com/gogf/gf/os/gtime"
)

// Time converts <i> to time.Time.
func Time(i interface{}, format ...string) time.Time {
	// It's already this type.
	if len(format) == 0 {
		if v, ok := i.(time.Time); ok {
			return v
		}
	}
	if t := GTime(i, format...); t != nil {
		return t.Time
	}
	return time.Time{}
}

// Duration converts <i> to time.Duration.
// If <i> is string, then it uses time.ParseDuration to convert it.
// If <i> is numeric, then it converts <i> as nanoseconds.
func Duration(i interface{}) time.Duration {
	// It's already this type.
	if v, ok := i.(time.Duration); ok {
		return v
	}
	s := String(i)
	if !utils.IsNumeric(s) {
		d, _ := gtime.ParseDuration(s)
		return d
	}
	return time.Duration(Int64(i))
}

// GTime converts <i> to *gtime.Time.
// The parameter <format> can be used to specify the format of <i>.
// If no <format> given, it converts <i> using gtime.NewFromTimeStamp if <i> is numeric,
// or using gtime.StrToTime if <i> is string.
func GTime(i interface{}, format ...string) *gtime.Time {
	if i == nil {
		return nil
	}
	// It's already this type.
	if len(format) == 0 {
		if v, ok := i.(*gtime.Time); ok {
			return v
		}
	}
	s := String(i)
	if len(s) == 0 {
		return gtime.New()
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		t, _ := gtime.StrToTimeFormat(s, format[0])
		return t
	}
	if utils.IsNumeric(s) {
		return gtime.NewFromTimeStamp(Int64(s))
	} else {
		t, _ := gtime.StrToTime(s)
		return t
	}
}

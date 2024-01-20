// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"time"

	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gtime"
)

// Time converts `any` to time.Time.
func Time(any interface{}, format ...string) time.Time {
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(time.Time); ok {
			return v
		}
	}
	if t := GTime(any, format...); t != nil {
		return t.Time
	}
	return time.Time{}
}

// Duration converts `any` to time.Duration.
// If `any` is string, then it uses time.ParseDuration to convert it.
// If `any` is numeric, then it converts `any` as nanoseconds.
func Duration(any interface{}) time.Duration {
	// It's already this type.
	if v, ok := any.(time.Duration); ok {
		return v
	}
	s := String(any)
	if !utils.IsNumeric(s) {
		d, _ := gtime.ParseDuration(s)
		return d
	}
	return time.Duration(Int64(any))
}

// GTime converts `any` to *gtime.Time.
// The parameter `format` can be used to specify the format of `any`.
// It returns the converted value that matched the first format of the formats slice.
// If no `format` given, it converts `any` using gtime.NewFromTimeStamp if `any` is numeric,
// or using gtime.StrToTime if `any` is string.
func GTime(any interface{}, format ...string) *gtime.Time {
	if any == nil {
		return nil
	}
	if v, ok := any.(iGTime); ok {
		return v.GTime(format...)
	}
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(*gtime.Time); ok {
			return v
		}
		if t, ok := any.(time.Time); ok {
			return gtime.New(t)
		}
		if t, ok := any.(*time.Time); ok {
			return gtime.New(t)
		}
	}
	s := String(any)
	if len(s) == 0 {
		return gtime.New()
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		for _, item := range format {
			t, err := gtime.StrToTimeFormat(s, item)
			if t != nil && err == nil {
				return t
			}
		}
		return nil
	}
	if utils.IsNumeric(s) {
		return gtime.NewFromTimeStamp(Int64(s))
	} else {
		t, _ := gtime.StrToTime(s)
		return t
	}
}

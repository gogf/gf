// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"time"

	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Time converts `any` to time.Time.
func (c *Converter) Time(any interface{}, format ...string) (time.Time, error) {
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(time.Time); ok {
			return v, nil
		}
	}
	t, err := c.GTime(any, format...)
	if err != nil {
		return time.Time{}, err
	}
	if t != nil {
		return t.Time, nil
	}
	return time.Time{}, nil
}

// Duration converts `any` to time.Duration.
// If `any` is string, then it uses time.ParseDuration to convert it.
// If `any` is numeric, then it converts `any` as nanoseconds.
func (c *Converter) Duration(any interface{}) (time.Duration, error) {
	// It's already this type.
	if v, ok := any.(time.Duration); ok {
		return v, nil
	}
	s, err := c.String(any)
	if err != nil {
		return 0, err
	}
	if !utils.IsNumeric(s) {
		return gtime.ParseDuration(s)
	}
	i, err := c.Int64(any)
	if err != nil {
		return 0, err
	}
	return time.Duration(i), nil
}

// GTime converts `any` to *gtime.Time.
// The parameter `format` can be used to specify the format of `any`.
// It returns the converted value that matched the first format of the formats slice.
// If no `format` given, it converts `any` using gtime.NewFromTimeStamp if `any` is numeric,
// or using gtime.StrToTime if `any` is string.
func (c *Converter) GTime(any interface{}, format ...string) (*gtime.Time, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	if v, ok := any.(localinterface.IGTime); ok {
		return v.GTime(format...), nil
	}
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(*gtime.Time); ok {
			return v, nil
		}
		if t, ok := any.(time.Time); ok {
			return gtime.New(t), nil
		}
		if t, ok := any.(*time.Time); ok {
			return gtime.New(t), nil
		}
	}
	s, err := c.String(any)
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return gtime.New(), nil
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		for _, item := range format {
			t, err := gtime.StrToTimeFormat(s, item)
			if err != nil {
				return nil, err
			}
			if t != nil {
				return t, nil
			}
		}
		return nil, nil
	}
	if utils.IsNumeric(s) {
		i, err := c.Int64(s)
		if err != nil {
			return nil, err
		}
		return gtime.NewFromTimeStamp(i), nil
	} else {
		return gtime.StrToTime(s)
	}
}

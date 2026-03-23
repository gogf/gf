// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// RuleDate implements `date` rule:
// Standard date, like: 2006-01-02, 20060102, 2006.01.02.
//
// Format: date
type RuleDate struct{}

func init() {
	Register(RuleDate{})
}

func (r RuleDate) Name() string {
	return "date"
}

func (r RuleDate) Message() string {
	return "The {field} value `{value}` is not a valid date"
}

func (r RuleDate) Run(in RunInput) error {
	type iTime interface {
		Date() (year int, month time.Month, day int)
		IsZero() bool
	}
	// support for time value, eg: gtime.Time/*gtime.Time, time.Time/*time.Time.
	if obj, ok := in.Value.Val().(iTime); ok {
		if obj.IsZero() {
			return errors.New(in.Message)
		}
		return nil
	}
	// Try direct time conversion for validation, which handles both format and date validity.
	// Support common date formats: 2006-01-02, 20060102, 2006.01.02, 2006/01/02
	if _, err := gtime.StrToTimeFormat(in.Value.String(), "Ymd"); err != nil {
		// Try with different separator formats
		if _, err := gtime.StrToTimeFormat(in.Value.String(), "Y-m-d"); err != nil {
			if _, err := gtime.StrToTimeFormat(in.Value.String(), "Y.m.d"); err != nil {
				if _, err := gtime.StrToTimeFormat(in.Value.String(), "Y/m/d"); err != nil {
					return errors.New(in.Message)
				}
			}
		}
	}
	return nil
}

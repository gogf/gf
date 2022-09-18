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

// RuleDatetime implements `datetime` rule:
// Standard datetime, like: 2006-01-02 12:00:00.
//
// Format: datetime
type RuleDatetime struct{}

func init() {
	Register(&RuleDatetime{})
}

func (r *RuleDatetime) Name() string {
	return "datetime"
}

func (r *RuleDatetime) Message() string {
	return "The {attribute} value `{value}` is not a valid datetime"
}

func (r *RuleDatetime) Run(in RunInput) error {
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
	if _, err := gtime.StrToTimeFormat(in.Value.String(), `Y-m-d H:i:s`); err != nil {
		return errors.New(in.Message)
	}
	return nil
}

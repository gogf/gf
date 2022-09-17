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

type RuleDateFormat struct{}

func init() {
	Register(&RuleDateFormat{})
}

func (r *RuleDateFormat) Name() string {
	return "date-format"
}

func (r *RuleDateFormat) Message() string {
	return "The {attribute} value `{value}` does not match the format: {pattern}"
}

func (r *RuleDateFormat) Run(in RunInput) error {
	type iTime interface {
		Date() (year int, month time.Month, day int)
		IsZero() bool
	}
	// support for time value, eg: gtime.Time/*gtime.Time, time.Time/*time.Time.
	if obj, ok := in.Value.Val().(iTime); ok {
		if obj.IsZero() {
			return errors.New(in.Message)
		}
	}
	if _, err := gtime.StrToTimeFormat(in.Value.String(), in.RulePattern); err != nil {
		return errors.New(in.Message)
	}
	return nil
}
